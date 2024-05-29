package main

import (
	"context"
	"encoding/json"
	"flag"
	nadv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	clientsetnad "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/client/clientset/versioned"
	"k8s.io/api/admission/v1beta1"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"
)

const (
	port = ":8443"
)

var (
	tlscert, tlskey string

	// outside cluster client
	/*
		config, _       = clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
		clientset, _    = kubernetes.NewForConfig(config)
	*/

	// inside cluster client
	// creates the in-cluster config
	config, _ = rest.InClusterConfig()

	// creates the clientset
	clientset, _ = kubernetes.NewForConfig(config)
)

// Validate handler accepts or rejects based on request contents
func Validate(w http.ResponseWriter, r *http.Request) {

	// var arReview with struct v1beta1.AdmissionReview{}
	arReview := v1beta1.AdmissionReview{}

	// decode arReview to json and check request
	if err := json.NewDecoder(r.Body).Decode(&arReview); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if arReview.Request == nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// https://kubernetes.io/docs/reference/config-api/apiserver-admission.v1/ <--- about admission request and response
	// get namespace name and put it to var
	nsName := arReview.Request.Namespace
	// get requester name and put it to var
	userInfo := arReview.Request.UserInfo.Username
	log.Printf("Requested namespace is %s", nsName)
	log.Printf("Requester for namespace is %s", userInfo)

	// object struct AdmissionResponse
	arReview.Response = &v1beta1.AdmissionResponse{
		UID:     arReview.Request.UID,
		Allowed: true,
	}
	//log.Println("The end of func validate")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&arReview)

	//========================================================================================
	// is it allowed to create quota and limit range?
	//log.Printf("start to check namespace %s", nsName)
	// get list namespaces from configmap
	exceptionNs, err := os.ReadFile("/files/_namespacelist")
	if err != nil {
		log.Println("config not found...")
		log.Println(err)
	}
	// convert it to string
	nsToStr := string(exceptionNs)

	// convert to slice
	stringToSlice := strings.Split(nsToStr, "\n")
	log.Printf("Checking namespace %s", nsName)
	// check namespace name forbidden to create resources
	if slices.Contains(stringToSlice, nsName) {
		log.Printf("create resources for namespace %s is forbidden", nsName)
	} else {
		//log.Println("start to create resources")
		go createObjects(nsName, userInfo)
	}

}

func createObjects(nsName, userInfo string) {

	//========================================================================================
	// apply resources quota, limit range and role binding
	time.Sleep(1 * time.Second)

	// ResourceQuota create
	a, err := os.ReadFile("/files/_resourcequota.yaml")
	if err != nil {
		panic(err)
	}
	// get yaml and convert it to  v1.ResourceQuota
	// to provide it import "k8s.io/apimachinery/pkg/util/yaml"
	quotaData := &v1.ResourceQuota{}
	err = yaml.Unmarshal(a, quotaData)
	if err != nil {
		panic(err)
	}

	// create quota for new namespaces
	_, err = clientset.CoreV1().ResourceQuotas(nsName).Create(context.TODO(), quotaData, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
		log.Printf("Failed to create ResourceQuota for %s ", nsName)

	} else {
		log.Printf("Created ResourceQuota for %s", nsName)
	}

	// LimitRange create
	b, err := os.ReadFile("/files/_limitrange.yaml")
	if err != nil {
		panic(err)
	}
	// get yaml and convert it to  v1.ResourceQuota
	// to provide it import "k8s.io/apimachinery/pkg/util/yaml"
	limitRange := &v1.LimitRange{}
	err = yaml.Unmarshal(b, limitRange)
	if err != nil {
		panic(err)
	}

	// create quota for new namespaces
	_, err = clientset.CoreV1().LimitRanges(nsName).Create(context.TODO(), limitRange, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
		log.Printf("Failed to create LimitRange for %s ", nsName)

	} else {
		log.Printf("Created LimitRange for %s", nsName)
	}

	// ---------------------------------------
	// create net-attach-def for new namespaces

	// nad client
	cfg, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		log.Printf("Error building kubeconfig: %v", err)
	}

	nadClient, err := clientsetnad.NewForConfig(cfg)
	if err != nil {
		log.Printf("Error building example clientset: %v", err)
	}

	nad := &nadv1.NetworkAttachmentDefinition{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "k8s.cni.cncf.io/v1",
			Kind:       "NetworkAttachmentDefinition",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "istio-cni",
		},
	}
	_, err = nadClient.K8sCniCncfIoV1().NetworkAttachmentDefinitions(nsName).Create(context.Background(), nad, metav1.CreateOptions{})
	if err != nil {
		log.Println("Cluster is k8s not OpenShift no need to create NetworkAttachmentDefinition ")
		log.Printf("Error %s", err)
	} else {
		log.Println("Created NetworkAttachmentDefinition ")
	}

	// ---------------------------------------
	// parse requester user or service account
	log.Printf("Requester %s", userInfo)
	// empty slice
	sliceUser := []string{}
	// append to slice
	sliceUser = append(sliceUser, userInfo)
	parsedUser := strings.Split(userInfo, ":")
	if len(parsedUser) > 1 {

		log.Printf("parsed user %s", parsedUser[3])
		saUser := parsedUser[3]
		saNs := parsedUser[2]
		//system:serviceaccount:vlku4:vlku4

		roleBinding := rbacv1.RoleBinding{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "rbac.authorization.k8s.io/v1",
				Kind:       "RoleBinding",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: saUser + "-admin-" + nsName,
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     "admin",
			},
			Subjects: []rbacv1.Subject{{Kind: "ServiceAccount", Name: saUser, Namespace: saNs}},
		}
		// create role binding for service account
		_, err := clientset.RbacV1().RoleBindings(nsName).Create(context.Background(), &roleBinding, metav1.CreateOptions{})
		if err != nil {
			log.Println(err)
			log.Printf("failed to create rolebinding for %s ", saUser)

		} else {
			log.Printf("Created RoleBinding for %s", saUser)
		}

	} else {
		// create role binding for user
		//log.Printf("creating rolebinding for %s", userInfo)
		roleBinding := rbacv1.RoleBinding{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "rbac.authorization.k8s.io/v1",
				Kind:       "RoleBinding",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: userInfo + "-admin-" + nsName,
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     "admin",
			},
			Subjects: []rbacv1.Subject{{Kind: "User", Name: userInfo}},
		}

		_, err = clientset.RbacV1().RoleBindings(nsName).Create(context.Background(), &roleBinding, metav1.CreateOptions{})
		if err != nil {
			log.Println(err)
			log.Printf("failed to create rolebinding for %s ", userInfo)

		} else {
			log.Printf("created rolebinding for %s", userInfo)
		}

	}

}

func main() {
	log.Printf("Application listening port %s\n", port)
	flag.StringVar(&tlscert, "tlsCertFile", "/certs/tls.crt",
		"File containing a certificate for HTTPS.")
	flag.StringVar(&tlskey, "tlsKeyFile", "/certs/tls.key",
		"File containing a private key for HTTPS.")
	flag.Parse()

	http.HandleFunc("/validate", Validate)
	log.Fatal(http.ListenAndServeTLS(port, tlscert, tlskey, nil))
}
