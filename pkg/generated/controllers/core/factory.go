// Code generated by main. DO NOT EDIT.

package core

import (
	"context"
	"time"

	"github.com/rancher/wrangler/pkg/generic"
	"k8s.io/apimachinery/pkg/runtime/schema"
	informers "k8s.io/client-go/informers"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Factory struct {
	synced            bool
	informerFactory   informers.SharedInformerFactory
	clientset         clientset.Interface
	controllerManager *generic.ControllerManager
	threadiness       map[schema.GroupVersionKind]int
}

func NewFactoryFromConfigOrDie(config *rest.Config) *Factory {
	f, err := NewFactoryFromConfig(config)
	if err != nil {
		panic(err)
	}
	return f
}

func NewFactoryFromConfig(config *rest.Config) (*Factory, error) {
	cs, err := clientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	informerFactory := informers.NewSharedInformerFactory(cs, 2*time.Hour)
	return NewFactory(cs, informerFactory), nil
}

func NewFactoryFromConfigWithNamespace(config *rest.Config, namespace string) (*Factory, error) {
	if namespace == "" {
		return NewFactoryFromConfig(config)
	}

	cs, err := clientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	informerFactory := informers.NewSharedInformerFactoryWithOptions(cs, 2*time.Hour, informers.WithNamespace(namespace))
	return NewFactory(cs, informerFactory), nil
}

func NewFactory(clientset clientset.Interface, informerFactory informers.SharedInformerFactory) *Factory {
	return &Factory{
		threadiness:       map[schema.GroupVersionKind]int{},
		controllerManager: &generic.ControllerManager{},
		clientset:         clientset,
		informerFactory:   informerFactory,
	}
}

func (c *Factory) SetThreadiness(gvk schema.GroupVersionKind, threadiness int) {
	c.threadiness[gvk] = threadiness
}

func (c *Factory) Sync(ctx context.Context) error {
	c.informerFactory.Start(ctx.Done())
	c.informerFactory.WaitForCacheSync(ctx.Done())
	return nil
}

func (c *Factory) Start(ctx context.Context, defaultThreadiness int) error {
	if err := c.Sync(ctx); err != nil {
		return err
	}

	return c.controllerManager.Start(ctx, defaultThreadiness, c.threadiness)
}

func (c *Factory) Core() Interface {
	return New(c.controllerManager, c.informerFactory.Core(), c.clientset)
}
