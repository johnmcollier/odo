/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by informer-gen. DO NOT EDIT.

package internalversion

import (
	time "time"

	servicecatalog "github.com/kubernetes-sigs/service-catalog/pkg/apis/servicecatalog"
	internalclientset "github.com/kubernetes-sigs/service-catalog/pkg/client/clientset_generated/internalclientset"
	internalinterfaces "github.com/kubernetes-sigs/service-catalog/pkg/client/informers_generated/internalversion/internalinterfaces"
	internalversion "github.com/kubernetes-sigs/service-catalog/pkg/client/listers_generated/servicecatalog/internalversion"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// ClusterServicePlanInformer provides access to a shared informer and lister for
// ClusterServicePlans.
type ClusterServicePlanInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() internalversion.ClusterServicePlanLister
}

type clusterServicePlanInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewClusterServicePlanInformer constructs a new informer for ClusterServicePlan type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewClusterServicePlanInformer(client internalclientset.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredClusterServicePlanInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredClusterServicePlanInformer constructs a new informer for ClusterServicePlan type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredClusterServicePlanInformer(client internalclientset.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.Servicecatalog().ClusterServicePlans().List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.Servicecatalog().ClusterServicePlans().Watch(options)
			},
		},
		&servicecatalog.ClusterServicePlan{},
		resyncPeriod,
		indexers,
	)
}

func (f *clusterServicePlanInformer) defaultInformer(client internalclientset.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredClusterServicePlanInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *clusterServicePlanInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&servicecatalog.ClusterServicePlan{}, f.defaultInformer)
}

func (f *clusterServicePlanInformer) Lister() internalversion.ClusterServicePlanLister {
	return internalversion.NewClusterServicePlanLister(f.Informer().GetIndexer())
}
