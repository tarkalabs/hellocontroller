package main

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig string

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "You only need this when running outside the cluster")
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})
}
func main() {
	flag.Parse()
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Errorf("Unable to initialize config : %s", err.Error())
	}
	clientSet, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Errorf("Unable to create client : %s", err.Error())
	}
	log.Infof("Successfully connected to kubernetes %v", clientSet)

}
