// client go; restApi; structSlice;
// slice append; 
package main

import (
	"context"
	"encoding/json"
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"net/http"
	"os"
)

var (
	// get kubeconfig with os.Getenv
	config, _    = clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	clientset, _ = kubernetes.NewForConfig(config)
)

type Article struct {
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

// test struct
type SliceNsStruct struct {
	ControlPlane1 string `json:"Name"`
	ControlPlane2 string `json:"Name1"`
	ControlPlane3 string `json:"Name2"`
	ControlPlane4 string `json:"Name3"`
}

// struct for slice control plane
type TestSliceStruct struct {
	// Field should be with LowerCase
	List []string `json:"list"`
}

var SliceTest []TestSliceStruct
var Articles []Article
var CpNs []SliceNsStruct

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	json.NewEncoder(w).Encode(Articles)
}

func returnSliceTest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	json.NewEncoder(w).Encode(SliceTest)
}
func returnCpNs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnCpNs")
	json.NewEncoder(w).Encode(CpNs)
}

// function that will handle all requests to our root URL
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Synapse Service Mesh Manager!\n")
	// Endpoint hit (попасть в конечную точку)
	fmt.Println("Endpoint Hit: homePage")
}

// match the URL path hit with a defined function
func handleRequests() {
	//return home page
	http.HandleFunc("/", homePage)

	//retrun articles
	http.HandleFunc("/articles", returnAllArticles)

	//retrun test cpns
	http.HandleFunc("/cpns", returnCpNs)

	//retrun slice control plane
	http.HandleFunc("/cpslice", returnSliceTest)

	// listen port 10000 (listener)
	log.Fatal(http.ListenAndServe(":10000", nil))
	//return all articles
}

func main() {

	Articles = []Article{
		Article{Title: "Hello", Desc: "Article Description", Content: "Article Content"},
	}
	// test instance
	CpNs = []SliceNsStruct{
		{"s1", "s2", "s3", "s4"},
	}
	//namespaces which one has label auto-ose-cp=syn is control planes
	labelForIo := "auto-ose-cp=syn"

	// get list namespaces with label "auto-ose-cp=syn"
	list, err := clientset.CoreV1().Namespaces().List(context.Background(), v1.ListOptions{LabelSelector: string(labelForIo)})
	if err != nil {
		err = fmt.Errorf("error getting namespaces: %v\n", err)
	}
	// create empty slice and fill it with func append
	s1 := []string{}
	for _, item := range list.Items {
		//s = append(s, item.Name)
		s1 = append(s1, item.Name)

	}
	// create instance for struct
	SliceTest = []TestSliceStruct{
		//add slice with control plane list
		{s1},
	}

	// run listener func (port 10000)
	handleRequests()

}
