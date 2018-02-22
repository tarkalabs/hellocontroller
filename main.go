package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	v1alpha "github.com/tarkalabs/hellocontroller/pkg/apis/hellocontroller/v1alpha"
	clientset "github.com/tarkalabs/hellocontroller/pkg/client/clientset/versioned"
	hcInformers "github.com/tarkalabs/hellocontroller/pkg/client/informers/externalversions"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig string

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "You only need this when running outside the cluster")
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})
}
func main() {
	stop := make(chan struct{})
	signalChan := make(chan os.Signal, 2)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		close(stop)
		<-signalChan
		os.Exit(1)
	}()
	flag.Parse()
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Errorf("Unable to initialize config : %s", err.Error())
	}
	clientSet, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Errorf("Unable to create client : %s", err.Error())
	}
	hcClientSet, err := clientset.NewForConfig(cfg)
	if err != nil {
		log.Error("Unable to create client set for hello controller %s", err.Error())
	}
	hcInformerFactory := hcInformers.NewSharedInformerFactory(hcClientSet, 10*time.Second)
	informerFactory := informers.NewSharedInformerFactory(clientSet, 10*time.Second)
	go hcInformerFactory.Start(stop)
	dbInformer := hcInformerFactory.Hellocontroller().V1alpha().Databases().Informer()
	dbInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			if db, ok := obj.(*v1alpha.Database); ok {
				log.Infof("Added database %s ", db.Spec.DatabaseName)
			}
		},
	})
	podInformer := informerFactory.Core().V1().Pods()
	NewController(podInformer).Run()
	log.Infof("Successfully connected to kubernetes %v", clientSet)
	go informerFactory.Start(stop)
	<-stop
}
