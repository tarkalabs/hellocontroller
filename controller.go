package main

import (
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	cv1Informers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/tools/cache"
)

type Controller struct {
	PodInformer cv1Informers.PodInformer
}

func NewController(pInformer cv1Informers.PodInformer) *Controller {
	return &Controller{PodInformer: pInformer}
}

func (c *Controller) Run() {
	c.PodInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
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
}
