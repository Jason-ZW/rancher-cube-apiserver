package backend

import (
	"time"

	"github.com/Sirupsen/logrus"
	infracs "github.com/rancher/rancher-cube-apiserver/k8s/pkg/client/clientset/versioned"
	infrainformers "github.com/rancher/rancher-cube-apiserver/k8s/pkg/client/informers/externalversions"
	apics "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type ClientGenerator struct {
	clientset            kubernetes.Clientset
	apiclientset         apics.Clientset
	informerFactory      informers.SharedInformerFactory
	infraInformerFactory infrainformers.SharedInformerFactory
}

func NewClientGenerator(kubeConfig string) *ClientGenerator {
	var config *rest.Config
	var err error

	if kubeConfig == "" {
		config, err = rest.InClusterConfig()
		if err != nil {
			logrus.Fatalf("RancherCUBE: generate config failed: %v", err)
		}
	}

	config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		logrus.Fatalf("RancherCUBE: generate config failed: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Fatalf("RancherCUBE: generate clientset failed: %v", err)
	}

	apiclientset, err := apics.NewForConfig(config)
	if err != nil {
		logrus.Fatalf("RancherCUBE: generate extensions clientset failed: %v", err)
	}

	infraclientset, err := infracs.NewForConfig(config)
	if err != nil {
		logrus.Fatalf("RancherCUBE: generate infra clientset failed: %v", err)
	}

	informerFactory := informers.NewSharedInformerFactory(clientset, time.Second*30)
	infraInformerFactory := infrainformers.NewSharedInformerFactory(infraclientset, time.Second*30)

	return &ClientGenerator{
		clientset:            *clientset,
		apiclientset:         *apiclientset,
		informerFactory:      informerFactory,
		infraInformerFactory: infraInformerFactory,
	}
}
