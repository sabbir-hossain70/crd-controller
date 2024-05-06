package controller

import (
	clientset "github.com/sabbir-hossain70/crd-controller/pkg/generated/clientset/versioned"
	informer "github.com/sabbir-hossain70/crd-controller/pkg/generated/informers/externalversions/crd.com/v1"
	lister "github.com/sabbir-hossain70/crd-controller/pkg/generated/listers/crd.com/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	appsinformer "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	appslisters "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"log"
)

type Controller struct {
	kubeclientset    kubernetes.Interface
	sampleclientset  clientset.Interface
	deploymentLister appslisters.DeploymentLister
	deploymentSynced cache.InformerSynced
	sabbirLister     lister.SabbirLister
	sabbirSynced     cache.InformerSynced
	workQueue        workqueue.RateLimitingInterface
}

func NewController(kubeclientset kubernetes.Interface, sampleClientset clientset.Interface,
	deploymentInformer appsinformer.DeploymentInformer,
	sabbirInformer informer.SabbirInformer) *Controller {

	ctrl := &Controller{
		kubeclientset:    kubeclientset,
		sampleclientset:  sampleClientset,
		deploymentLister: deploymentInformer.Lister(),
		deploymentSynced: deploymentInformer.Informer().HasSynced,
		sabbirLister:     sabbirInformer.Lister(),
		sabbirSynced:     sabbirInformer.Informer().HasSynced,
		workQueue:        workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
	}
	log.Println("Setting up event handlers")

	sabbirInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: ctrl.enqueueSabbir,
		UpdateFunc: func(old, new interface{}) {
			ctrl.enqueueSabbir(new)
		},
		DeleteFunc: func(obj interface{}) {
			ctrl.enqueueSabbir(obj)
		},
	})

	return ctrl
}

func (c *Controller) enqueueSabbir(obj interface{}) {
	log.Println("Sabbir Controller: Enqueuing Sabbir")
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.workQueue.Add(key)
}
