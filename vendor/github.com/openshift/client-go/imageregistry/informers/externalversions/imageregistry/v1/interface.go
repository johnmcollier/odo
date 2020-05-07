// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	internalinterfaces "github.com/openshift/client-go/imageregistry/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// Configs returns a ConfigInformer.
	Configs() ConfigInformer
	// ImagePruners returns a ImagePrunerInformer.
	ImagePruners() ImagePrunerInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// Configs returns a ConfigInformer.
func (v *version) Configs() ConfigInformer {
	return &configInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}

// ImagePruners returns a ImagePrunerInformer.
func (v *version) ImagePruners() ImagePrunerInformer {
	return &imagePrunerInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}
