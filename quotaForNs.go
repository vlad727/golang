package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"slices"
	"sort"
	"strings"
)

var (
	config, _    = clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	clientset, _ = kubernetes.NewForConfig(config)
)

// get all namespaces
func getNamespace() {

	// time out
	//time.Sleep(10 * time.Second) // sleep 10 sec

	// get namespaces with clientset
	//method Watch, use the watcher to gain access to the event notifications from a Go channel via method ResultChan.
	watchNs, _ := clientset.CoreV1().Namespaces().Watch(context.Background(), metav1.ListOptions{})

	for event := range watchNs.ResultChan() {
		item := event.Object.(*v1.Namespace)

		switch event.Type {
		case watch.Modified:
		case watch.Bookmark:
		case watch.Error:
		case watch.Deleted:
			log.Printf("Namespace %s: has been deleted...\n", item.GetName())
		case watch.Added:
			exceptNs(item.GetName())
		}
	}
}

// exclude namespace, check input namespace and compare it with namespace from config file
func exceptNs(name string) {

	// get list from configmap
	exceptionNs, err := os.ReadFile("config")
	if err != nil {
		log.Println("Config not found...\n")
	}
	// convert it to string
	nsToStr := string(exceptionNs)
	// convert to slice
	stringToSlice := strings.Fields(nsToStr)
	// sort list from config
	sort.Strings(stringToSlice)

	if slices.Contains(stringToSlice, name) {
		log.Printf("Нахуй квота для namespace: %s", name)
	} else {
		createQuota(name)
	}

}

// create quota for namespace
func createQuota(x string) {
	/*
		// declare quota
		quota := &v1.ResourceQuota{
			ObjectMeta: metav1.ObjectMeta{
				Name: "default",
			},
		}

	*/
	b, err := os.ReadFile("quota.yaml")
	if err != nil {
		panic(err)
	}
	// get yaml and convert it to  v1.ResourceQuota
	// to provide it import "k8s.io/apimachinery/pkg/util/yaml"
	quotaData := &v1.ResourceQuota{}
	err = yaml.Unmarshal(b, quotaData)
	if err != nil {
		panic(err)
	}
	// create quota for new namespaces
	_, err = clientset.CoreV1().ResourceQuotas(x).Create(context.TODO(), quotaData, metav1.CreateOptions{})
	if err != nil {
		log.Info(err, " ", "for ", x)
	} else {
		log.Info("Created quota for: ", x)
	}

}

// main func
func main() {

	getNamespace()
}

/*
+ exclude namespaces like: kube-system, kubernetes-, openshift-,box-,default etc.
quota rbac should be configmap
inside cluster solution
time out (may be too much load for api)
*/
