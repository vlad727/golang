package main

import (
    "os"
    "os/exec"
    "os/signal"
    "strings"
    "syscall"
    "time"

    "github.com/golang/glog"
    v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/util/runtime"
    "k8s.io/client-go/informers"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/cache"
    "k8s.io/client-go/tools/clientcmd"
)

func newNamespace(obj interface{}) {
    ns := obj.(v1.Object)
    glog.Error("New Namespace ", ns.GetName())
}

func modNamespace(objOld interface{}, objNew interface{}) {
}

func delNamespace(obj interface{}) {
    ns := obj.(v1.Object)
    glog.Error("Del Namespace ", ns.GetName())
}

func watchNamespace(k8s *kubernetes.Clientset) {
    // Add watcher for the Namespace.
    factory := informers.NewSharedInformerFactory(k8s, 5*time.Second)
    nsInformer := factory.Core().V1().Namespaces().Informer()
    nsInformerChan := make(chan struct{})
    //defer close(nsInformerChan)
    defer runtime.HandleCrash()

    // Namespace informer state change handler
    nsInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
        // When a new namespace gets created
        AddFunc: func(obj interface{}) {
            newNamespace(obj)
        },
        // When a namespace gets updated
        UpdateFunc: func(oldObj interface{}, newObj interface{}) {
            modNamespace(oldObj, newObj)
        },
        // When a namespace gets deleted
        DeleteFunc: func(obj interface{}) {
            delNamespace(obj)
        },
    })
    factory.Start(nsInformerChan)
    //go nsInformer.GetController().Run(nsInformerChan)
    go nsInformer.Run(nsInformerChan)
}

func main() {
    kconfig := os.Getenv("KUBECONFIG")
    glog.Error("KCONFIG", kconfig)
    var config *rest.Config
    var clientset *kubernetes.Clientset
    var err error
    for {
        if config == nil {
            config, err = clientcmd.BuildConfigFromFlags("", kconfig)
            if err != nil {
                glog.Error("Cant create  kubernetes config")
                time.Sleep(time.Second)
                continue
            }
        }
        // creates the clientset
        clientset, err = kubernetes.NewForConfig(config)
        if err != nil {
            glog.Error("Cannot create kubernetes client")
            time.Sleep(time.Second)
            continue
        }
        break
    }
    watchNamespace(clientset)
    glog.Error("Watch started")

    term := make(chan os.Signal, 1)
    signal.Notify(term, os.Interrupt)
    signal.Notify(term, syscall.SIGTERM)
    select {
    case <-term:
    }
}
