// Code generated by main. DO NOT EDIT.

package v1

import (
	"context"

	"github.com/rancher/wrangler/pkg/generic"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	informers "k8s.io/client-go/informers/core/v1"
	clientset "k8s.io/client-go/kubernetes/typed/core/v1"
	listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
)

type ConfigMapHandler func(string, *v1.ConfigMap) (*v1.ConfigMap, error)

type ConfigMapController interface {
	ConfigMapClient

	OnChange(ctx context.Context, name string, sync ConfigMapHandler)
	OnRemove(ctx context.Context, name string, sync ConfigMapHandler)
	Enqueue(namespace, name string)

	Cache() ConfigMapCache

	Informer() cache.SharedIndexInformer
	GroupVersionKind() schema.GroupVersionKind

	AddGenericHandler(ctx context.Context, name string, handler generic.Handler)
	AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler)
	Updater() generic.Updater
}

type ConfigMapClient interface {
	Create(*v1.ConfigMap) (*v1.ConfigMap, error)
	Update(*v1.ConfigMap) (*v1.ConfigMap, error)

	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1.ConfigMap, error)
	List(namespace string, opts metav1.ListOptions) (*v1.ConfigMapList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.ConfigMap, err error)
}

type ConfigMapCache interface {
	Get(namespace, name string) (*v1.ConfigMap, error)
	List(namespace string, selector labels.Selector) ([]*v1.ConfigMap, error)

	AddIndexer(indexName string, indexer ConfigMapIndexer)
	GetByIndex(indexName, key string) ([]*v1.ConfigMap, error)
}

type ConfigMapIndexer func(obj *v1.ConfigMap) ([]string, error)

type configMapController struct {
	controllerManager *generic.ControllerManager
	clientGetter      clientset.ConfigMapsGetter
	informer          informers.ConfigMapInformer
	gvk               schema.GroupVersionKind
}

func NewConfigMapController(gvk schema.GroupVersionKind, controllerManager *generic.ControllerManager, clientGetter clientset.ConfigMapsGetter, informer informers.ConfigMapInformer) ConfigMapController {
	return &configMapController{
		controllerManager: controllerManager,
		clientGetter:      clientGetter,
		informer:          informer,
		gvk:               gvk,
	}
}

func FromConfigMapHandlerToHandler(sync ConfigMapHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1.ConfigMap
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1.ConfigMap))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *configMapController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1.ConfigMap))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateConfigMapOnChange(updater generic.Updater, handler ConfigMapHandler) ConfigMapHandler {
	return func(key string, obj *v1.ConfigMap) (*v1.ConfigMap, error) {
		if obj == nil {
			return handler(key, nil)
		}

		copyObj := obj.DeepCopy()
		newObj, err := handler(key, copyObj)
		if newObj != nil {
			copyObj = newObj
		}
		if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
			newObj, err := updater(copyObj)
			if newObj != nil && err == nil {
				copyObj = newObj.(*v1.ConfigMap)
			}
		}

		return copyObj, err
	}
}

func (c *configMapController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, handler)
}

func (c *configMapController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), handler)
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, removeHandler)
}

func (c *configMapController) OnChange(ctx context.Context, name string, sync ConfigMapHandler) {
	c.AddGenericHandler(ctx, name, FromConfigMapHandlerToHandler(sync))
}

func (c *configMapController) OnRemove(ctx context.Context, name string, sync ConfigMapHandler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), FromConfigMapHandlerToHandler(sync))
	c.AddGenericHandler(ctx, name, removeHandler)
}

func (c *configMapController) Enqueue(namespace, name string) {
	c.controllerManager.Enqueue(c.gvk, namespace, name)
}

func (c *configMapController) Informer() cache.SharedIndexInformer {
	return c.informer.Informer()
}

func (c *configMapController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *configMapController) Cache() ConfigMapCache {
	return &configMapCache{
		lister:  c.informer.Lister(),
		indexer: c.informer.Informer().GetIndexer(),
	}
}

func (c *configMapController) Create(obj *v1.ConfigMap) (*v1.ConfigMap, error) {
	return c.clientGetter.ConfigMaps(obj.Namespace).Create(obj)
}

func (c *configMapController) Update(obj *v1.ConfigMap) (*v1.ConfigMap, error) {
	return c.clientGetter.ConfigMaps(obj.Namespace).Update(obj)
}

func (c *configMapController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return c.clientGetter.ConfigMaps(namespace).Delete(name, options)
}

func (c *configMapController) Get(namespace, name string, options metav1.GetOptions) (*v1.ConfigMap, error) {
	return c.clientGetter.ConfigMaps(namespace).Get(name, options)
}

func (c *configMapController) List(namespace string, opts metav1.ListOptions) (*v1.ConfigMapList, error) {
	return c.clientGetter.ConfigMaps(namespace).List(opts)
}

func (c *configMapController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientGetter.ConfigMaps(namespace).Watch(opts)
}

func (c *configMapController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.ConfigMap, err error) {
	return c.clientGetter.ConfigMaps(namespace).Patch(name, pt, data, subresources...)
}

type configMapCache struct {
	lister  listers.ConfigMapLister
	indexer cache.Indexer
}

func (c *configMapCache) Get(namespace, name string) (*v1.ConfigMap, error) {
	return c.lister.ConfigMaps(namespace).Get(name)
}

func (c *configMapCache) List(namespace string, selector labels.Selector) ([]*v1.ConfigMap, error) {
	return c.lister.ConfigMaps(namespace).List(selector)
}

func (c *configMapCache) AddIndexer(indexName string, indexer ConfigMapIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1.ConfigMap))
		},
	}))
}

func (c *configMapCache) GetByIndex(indexName, key string) (result []*v1.ConfigMap, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		result = append(result, obj.(*v1.ConfigMap))
	}
	return result, nil
}
