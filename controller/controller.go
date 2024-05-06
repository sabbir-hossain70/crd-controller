package controller

import (
	"context"
	"fmt"
	controllerv1 "github.com/sabbir-hossain70/crd-controller/pkg/apis/crd.com/v1"
	clientset "github.com/sabbir-hossain70/crd-controller/pkg/generated/clientset/versioned"
	informer "github.com/sabbir-hossain70/crd-controller/pkg/generated/informers/externalversions/crd.com/v1"
	lister "github.com/sabbir-hossain70/crd-controller/pkg/generated/listers/crd.com/v1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	appsinformer "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	appslisters "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"log"
	"time"
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
	println("key added in workqueue ", key)
	c.workQueue.Add(key)
}

func (c *Controller) Run(workers int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workQueue.ShutDown()

	log.Println("Starting Sabbir controller")

	if ok := cache.WaitForCacheSync(stopCh, c.deploymentSynced, c.sabbirSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	log.Println("Starting workers")

	for i := 0; i < workers; i++ {
		go wait.Until(c.runworker, time.Second, stopCh)
	}

	log.Println("Started workers")
	<-stopCh

	log.Println("Shutting down workers")

	return nil
}

func (c *Controller) runworker() {
	for c.ProcessNextItem() {

	}
}

func (c *Controller) ProcessNextItem() bool {

	obj, shutdown := c.workQueue.Get()
	if shutdown {
		return false
	}

	err := func(obj interface{}) error {
		defer c.workQueue.Done(obj)
		var key string
		var ok bool

		if key, ok = obj.(string); !ok {
			c.workQueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}

		if err := c.syncHandler(key); err != nil {
			c.workQueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}
		c.workQueue.Forget(obj)
		log.Printf("Successfully synced '%s'", key)
		return nil

	}

	return true
}

func (c *Controller) syncHandler(key string) error {

	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	sabbir, err := c.sabbirLister.Sabbirs(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("sabbir '%s' in work queue no longer exists", key))
			return nil
		}
		return err
	}

	deploymentName := sabbir.Spec.Name

	if deploymentName == "" {
		utilruntime.HandleError(fmt.Errorf("deployment must be specified: %s", key))
		return nil
	}

	deployment, err := c.deploymentLister.Deployments(sabbir.Namespace).Get(deploymentName)
	if errors.IsNotFound(err) {
		deployment, err := c.kubeclientset.AppsV1().Deployments(sabbir.Namespace).Create(context.TODO(), newDeployment(sabbir), metav1.CreateOptions{})
	}

	if err != nil {
		return err
	}

	if sabbir.Spec.Replicas != nil && *sabbir.Spec.Replicas != *deployment.Spec.Replicas {
		log.Printf("Sabbir %s replicas : %d, deployment replicas: %d\n", name, *sabbir.Spec.Replicas, *deployment.Spec.Replicas)

		deployment, err = c.kubeclientset.AppsV1().Deployments(namespace).Update(context.TODO(), newDeployment(sabbir), metav1.UpdateOptions{})

		if err != nil {
			return err
		}
	}

	return nil
}

func newDeployment(sabbir *controllerv1.Sabbir) *appsv1.Deployment {

}
