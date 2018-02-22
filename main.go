package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
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
	informerFactory := informers.NewSharedInformerFactory(clientSet, 10*time.Second)

	podInformer := informerFactory.Core().V1().Pods().Informer()
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			if pod, ok := obj.(*corev1.Pod); ok {
				log.Infof("Added a pod : %s", pod.Name)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			if pod, ok := newObj.(*corev1.Pod); ok {
				log.Infof("Updated a pod : %s", pod.Name)
			}
		},
		DeleteFunc: func(obj interface{}) {
			if pod, ok := obj.(*corev1.Pod); ok {
				log.Infof("deleted a pod : %s", pod.Name)
			} else {
				log.Infof("Could not type cast... probably go deleted")
			}
		},
	})

	log.Infof("Successfully connected to kubernetes %v", clientSet)
	go informerFactory.Start(stop)
	<-stop
}
