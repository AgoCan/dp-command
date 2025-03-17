package k8s

import (
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Client struct {
	ClientSet *kubernetes.Clientset
}

func New(config *rest.Config) *Client {
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		// 不应该panic，因为k8s挂了，不能影程序运行
		zap.S().Error("new clientSet err: ", err)

	}
	return &Client{
		ClientSet: clientSet,
	}
}
