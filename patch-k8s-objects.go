package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

var (
	nsname   = "vlku7"
	username = "kubernetes-admin"

	// outside cluster client

	config, _    = clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	clientset, _ = kubernetes.NewForConfig(config)

	/*
		// inside cluster client
		// creates the in-cluster config
		config, _ = rest.InClusterConfig()

		// creates the clientset
		clientset, _ = kubernetes.NewForConfig(config)

	*/
)

type Y struct {
	Metadata Annotations `json:"metadata"`
}

type Annotations struct {
	Annotations Requester `json:"annotations"`
}
type Requester struct {
	Requester string `json:"requester"`
}

func main() {
	//kubectl patch ns  vlku7   --type='merge' -p='{"metadata": {"annotations": {"test": "empty"}}}'
	//kubectl patch ns a --type='merge' -p='{"metadata": {"labels": {"test": "empty"}}}'

	requesterName := "admin"

	setAnnotation := Y{
		Metadata: Annotations{
			Requester{requesterName},
		},
	}
	bytes, _ := json.Marshal(setAnnotation)

	_, err := clientset.CoreV1().Namespaces().Patch(context.TODO(), nsname, types.MergePatchType, bytes, metav1.PatchOptions{})
	if err != nil {

		fmt.Println(err)
	}
	fmt.Println("Succesed patch namespace", string(bytes))
}
/*
Result
k get ns vlku7 -o json | jq .metadata.annotations
{
  "requester": "admin"
}
*/

//https://stackoverflow.com/questions/69125257/golang-kubernetes-client-patching-an-existing-resource-with-a-label <<< diff merge and json

