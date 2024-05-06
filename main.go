package main

import (
	"context"
	"flag"
	myv1 "github.com/sabbir-hossain70/crd/pkg/apis/crd.com/v1"
	sbclientset "github.com/sabbir-hossain70/crd/pkg/generated/clientset/versioned"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	crdclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	_ "k8s.io/code-generator"
	"log"
	"os"
	"os/signal"
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
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	crdClient, err := crdclientset.NewForConfig(config)

	if err != nil {
		panic(err.Error())
	}

	customCRD := v1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "sabbirs.crd.com",
		},
		Spec: v1.CustomResourceDefinitionSpec{
			Group: "crd.com",
			Versions: []v1.CustomResourceDefinitionVersion{
				{
					Name:    "v1",
					Served:  true,
					Storage: true,
					Schema: &v1.CustomResourceValidation{
						OpenAPIV3Schema: &v1.JSONSchemaProps{
							Type: "object",

							Properties: map[string]v1.JSONSchemaProps{
								"spec": {
									Type: "object",
									Properties: map[string]v1.JSONSchemaProps{
										"name": {
											Type: "string",
										},
										"replicas": {
											Type: "integer",
										},
										"container": {
											Type: "object",
											Properties: map[string]v1.JSONSchemaProps{
												"image": {
													Type: "string",
												},
												"port": {
													Type: "integer",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			Scope: "Namespaced",
			Names: v1.CustomResourceDefinitionNames{
				Kind:     "Sabbir",
				Plural:   "sabbirs",
				Singular: "sabbir",
				ShortNames: []string{
					"sb",
				},
				Categories: []string{
					"all",
				},
			},
		},
	}

	ctx := context.TODO()
	_ = crdClient.ApiextensionsV1().CustomResourceDefinitions().Delete(ctx, customCRD.Name, metav1.DeleteOptions{})

	_, err = crdClient.ApiextensionsV1().CustomResourceDefinitions().Create(ctx, &customCRD, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
	time.Sleep(2 * time.Second)
	log.Println("CRD created")
	log.Println(customCRD.Name)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	println("press ctrl+c to stop")
	<-signalChan
	log.Println("Creating Sabbir")

	client, err := sbclientset.NewForConfig(config)

	if err != nil {
		print("error here")
		panic(err.Error())
	}

	sbObj := myv1.Sabbir{
		ObjectMeta: metav1.ObjectMeta{
			Name: "sabbir",
		},
		Spec: myv1.SabbirSpec{
			Name:     "CustomSabbir",
			Replicas: intptr(3),
			Container: myv1.ContainerSpec{
				Image: "sabbir70/api-bookserver",
			},
		},
	}

	_, err = client.CrdV1().Sabbirs("default").Create(ctx, &sbObj, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
	time.Sleep(2 * time.Second)
	log.Println("Sabbir created")
	signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	err = client.CrdV1().Sabbirs("default").Delete(ctx, "sabbir", metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
	err = crdClient.ApiextensionsV1().CustomResourceDefinitions().Delete(ctx, customCRD.Name, metav1.DeleteOptions{})

	if err != nil {
		panic(err.Error())
	}
	log.Println("cleaned up")

}

func intptr(i int32) *int32 {
	return &i
}
