package controller

import (
	"context"
	"fmt"
	controllerv1 "github.com/sabbir-hossain70/crd-controller/pkg/apis/crd.com/v1"
	clientset "github.com/sabbir-hossain70/crd-controller/pkg/generated/clientset/versioned"
	informer "github.com/sabbir-hossain70/crd-controller/pkg/generated/informers/externalversions/crd.com/v1"
	lister "github.com/sabbir-hossain70/crd-controller/pkg/generated/listers/crd.com/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	appsinformer "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	appslisters "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"log"
	"math/rand"
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
			fmt.Println("UpdateFunc called")
			ctrl.enqueueSabbir(new)
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println("DeleteFunc called")
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
		go wait.Until(c.runworker, time.Second*2, stopCh)
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
		klog.Infoln("Printing key", key)
		if err := c.syncHandler(key); err != nil {
			c.workQueue.AddRateLimited(key)
			fmt.Println("Error in syncHandler")
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}
		c.workQueue.Forget(obj)
		log.Printf("Successfully synced '%s'", key)
		return nil

	}(obj)
	if err != nil {
		utilruntime.HandleError(err)
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

	//deploymentName := "custom-sabbir"

	fmt.Println("Deployment Name: ", deploymentName)

	if deploymentName == "" {
		utilruntime.HandleError(fmt.Errorf("deployment must be specified: %s", key))
		return nil
	}

	deployment, err := c.deploymentLister.Deployments(sabbir.Namespace).Get(deploymentName)
	//fmt.Println("sabbir.Namespace: ", sabbir.Namespace)
	if errors.IsNotFound(err) {
		fmt.Println("Is not found .......")
		deployment, err = c.kubeclientset.AppsV1().Deployments(sabbir.Namespace).Create(context.TODO(), newDeployment(sabbir), metav1.CreateOptions{})
	}

	if err != nil {
		return err
	}

	if sabbir.Spec.Replicas != nil && *sabbir.Spec.Replicas != *deployment.Spec.Replicas {
		println("*sabbir.Spec.Replicas: ", *sabbir.Spec.Replicas)
		println("*deployment.Spec.Replicas: ", *deployment.Spec.Replicas)
		log.Printf("Sabbir %s replicas : %d, deployment replicas: %d\n", name, *sabbir.Spec.Replicas, *deployment.Spec.Replicas)

		deployment, err = c.kubeclientset.AppsV1().Deployments(namespace).Update(context.TODO(), newDeployment(sabbir), metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}
	fmt.Println("before update sabbir status.....")
	err = c.updateSabbirStatus(sabbir, deployment)
	if err != nil {
		return err
	}

	serviceName := sabbir.Spec.Name + "-service"
	//serviceName := "custom-sabbir" + "-service"
	println("Service Name: ", serviceName)
	service, err := c.kubeclientset.CoreV1().Services(sabbir.Namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		service, err = c.kubeclientset.CoreV1().Services(sabbir.Namespace).Create(context.TODO(), newService(sabbir), metav1.CreateOptions{})
		if err != nil {
			log.Println(err)
			return err
		}
		log.Printf("\nSabbir %s service created\n", serviceName)
	} else if err != nil {
		log.Println(err)
		return err
	}

	_, err = c.kubeclientset.CoreV1().Services(sabbir.Namespace).Update(context.TODO(), service, metav1.UpdateOptions{})
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (c *Controller) updateSabbirStatus(sabbir *controllerv1.Sabbir, deployment *appsv1.Deployment) error {
	time.Sleep(time.Second * 10)
	sb, err := c.sampleclientset.CrdV1().Sabbirs(sabbir.Namespace).Get(context.TODO(), sabbir.Name, metav1.GetOptions{})

	sb.TypeMeta.Kind = "Sabbir"
	sb.TypeMeta.APIVersion = "crd.com/v1"
	sabbirCopy := sb.DeepCopy()
	//sabbirCopy.Spec.Replicas = deployment.Spec.Replicas
	var y int32 = int32(rand.Intn(10))
	sabbirCopy.Spec.Replicas = &y
	klog.Infoln("kind :", sabbir.Kind)

	_, err = c.sampleclientset.CrdV1().Sabbirs(sabbir.Namespace).Update(context.TODO(), sabbirCopy, metav1.UpdateOptions{})
	fmt.Println("sabbirCopy.Spec.Replicas: ", *sabbirCopy.Spec.Replicas)
	fmt.Println("sb.Spec.Replicas: ", *sb.Spec.Replicas)
	fmt.Println("update err: ", err)

	sabbir, err = c.sampleclientset.CrdV1().Sabbirs(sabbir.Namespace).Get(context.TODO(), "sabbir", metav1.GetOptions{})
	if err != nil {
		fmt.Println("Error fetching custom resource:", err)
	}
	fmt.Println("now sabbir.Spec.Replicas: ", *sabbir.Spec.Replicas)

	return err
}

func (c *Controller) updateSabbirStatustest(sabbir *controllerv1.Sabbir, deployment *appsv1.Deployment) error {

	sabbirCopy := sabbir.DeepCopy()
	fmt.Println("SabbirCopy.Status.AvailableReplicas: ", sabbirCopy.Status.AvailableReplicas)
	fmt.Println("deployment.Status.AvailableReplicas: ", deployment.Status.AvailableReplicas)
	sabbirCopy.Status.AvailableReplicas = deployment.Status.AvailableReplicas

	var x int32 = int32(rand.Intn(10))
	var y int32 = int32(rand.Intn(10))
	fmt.Println("x : ", x, " y : ", y)
	deployment.Spec.Replicas = &x

	deployment.Status.Replicas = y + 4

	deploymentCopy := deployment.DeepCopy()
	deploymentCopy.Spec.Replicas = &y // or any new value for replicas
	_, err := c.kubeclientset.AppsV1().Deployments(sabbir.Namespace).Update(context.TODO(), deploymentCopy, metav1.UpdateOptions{})
	if err != nil {
		log.Println("Failed to update Deployment:", err)
		return err
	}

	fmt.Println("deployment.Spec.Replicas: ", *deployment.Spec.Replicas)
	fmt.Println("deployment.Status.Replicas ", deployment.Status.Replicas)

	//fmt.Println("sabbir.Spec.Replicas: ", *sabbir.Spec.Replicas)
	sabbir.Status.AvailableReplicas = y

	//fmt.Println("sabbir.Status.AvailableReplicas: ", sabbir.Status.AvailableReplicas)
	//sabbirCopy.Spec.Replicas = &y
	_, err = c.sampleclientset.CrdV1().Sabbirs(sabbir.Namespace).Update(context.TODO(), sabbirCopy, metav1.UpdateOptions{})

	fmt.Println("sabbir.Status.AvailableReplicas: ", sabbir.Status.AvailableReplicas)
	fmt.Println("sabbirCopy.Status.AvailableReplicas: ", sabbirCopy.Status.AvailableReplicas)

	return err
}

func newDeployment(sabbir *controllerv1.Sabbir) *appsv1.Deployment {
	fmt.Println("..............Inside newDeployment")
	fmt.Println(" sabbir.spec.replicas: ", *sabbir.Spec.Replicas)
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: sabbir.Spec.Name,
			//Name:      "custom-sabbir",
			Namespace: sabbir.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "my-app",
				},
			},
			Replicas: sabbir.Spec.Replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "my-app",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "my-app",
							Image: sabbir.Spec.Container.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

}

func newService(sabbir *controllerv1.Sabbir) *corev1.Service {
	labels := map[string]string{
		"app": "my-app",
	}
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind: "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: sabbir.Spec.Name + "-service",
			//Name: "custom-sabbir-service",
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(sabbir, controllerv1.SchemeGroupVersion.WithKind("Sabbir")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     corev1.ServiceTypeNodePort,
			Ports: []corev1.ServicePort{
				{
					Port:       80,
					TargetPort: intstr.FromInt32(80),
					NodePort:   30007,
				},
			},
		},
	}

}
