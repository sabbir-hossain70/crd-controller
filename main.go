package main

import (
	"flag"
	"fmt"
	"github.com/sabbir-hossain70/crd-controller/controller"
	clientset "github.com/sabbir-hossain70/crd-controller/pkg/generated/clientset/versioned"
	myinformers "github.com/sabbir-hossain70/crd-controller/pkg/generated/informers/externalversions"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	_ "k8s.io/code-generator"
	"log"
	"path/filepath"
	"time"
)

func main() {

	log.Println("Configuring kubeconfig")
	var kubeconfig *string
	home := homedir.HomeDir()
	if home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	fmt.Println("kubeconfig: ", *kubeconfig)
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	var exampleClient clientset.Interface
	exampleClient, err = clientset.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*10)

	exampleInformerFactor := myinformers.NewSharedInformerFactory(exampleClient, time.Second*10)

	ctrl := controller.NewController(kubeClient, exampleClient,
		kubeInformerFactory.Apps().V1().Deployments(),
		exampleInformerFactor.Crd().V1().Sabbirs())

	stopCh := make(chan struct{})
	kubeInformerFactory.Start(stopCh)
	exampleInformerFactor.Start(stopCh)
	//print(ctrl)
	if err = ctrl.Run(2, stopCh); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Main function terminated +++++++++++")

}
