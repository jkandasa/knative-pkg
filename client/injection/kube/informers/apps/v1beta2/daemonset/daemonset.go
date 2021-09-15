/*
Copyright 2021 The Knative Authors

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

// Code generated by injection-gen. DO NOT EDIT.

package daemonset

import (
	context "context"

	apiappsv1beta2 "k8s.io/api/apps/v1beta2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	v1beta2 "k8s.io/client-go/informers/apps/v1beta2"
	kubernetes "k8s.io/client-go/kubernetes"
	appsv1beta2 "k8s.io/client-go/listers/apps/v1beta2"
	cache "k8s.io/client-go/tools/cache"
	client "knative.dev/pkg/client/injection/kube/client"
	factory "knative.dev/pkg/client/injection/kube/informers/factory"
	controller "knative.dev/pkg/controller"
	injection "knative.dev/pkg/injection"
	logging "knative.dev/pkg/logging"
)

func init() {
	injection.Default.RegisterInformer(withInformer)
	injection.Dynamic.RegisterDynamicInformer(withDynamicInformer)
}

// Key is used for associating the Informer inside the context.Context.
type Key struct{}

func withInformer(ctx context.Context) (context.Context, controller.Informer) {
	f := factory.Get(ctx)
	inf := f.Apps().V1beta2().DaemonSets()
	return context.WithValue(ctx, Key{}, inf), inf.Informer()
}

func withDynamicInformer(ctx context.Context) context.Context {
	inf := &wrapper{client: client.Get(ctx)}
	return context.WithValue(ctx, Key{}, inf)
}

// Get extracts the typed informer from the context.
func Get(ctx context.Context) v1beta2.DaemonSetInformer {
	untyped := ctx.Value(Key{})
	if untyped == nil {
		logging.FromContext(ctx).Panic(
			"Unable to fetch k8s.io/client-go/informers/apps/v1beta2.DaemonSetInformer from context.")
	}
	return untyped.(v1beta2.DaemonSetInformer)
}

type wrapper struct {
	client kubernetes.Interface

	namespace string
}

var _ v1beta2.DaemonSetInformer = (*wrapper)(nil)
var _ appsv1beta2.DaemonSetLister = (*wrapper)(nil)

func (w *wrapper) Informer() cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(nil, &apiappsv1beta2.DaemonSet{}, 0, nil)
}

func (w *wrapper) Lister() appsv1beta2.DaemonSetLister {
	return w
}

func (w *wrapper) DaemonSets(namespace string) appsv1beta2.DaemonSetNamespaceLister {
	return &wrapper{client: w.client, namespace: namespace}
}

func (w *wrapper) List(selector labels.Selector) (ret []*apiappsv1beta2.DaemonSet, err error) {
	lo, err := w.client.AppsV1beta2().DaemonSets(w.namespace).List(context.TODO(), v1.ListOptions{
		LabelSelector: selector.String(),
		// TODO(mattmoor): Incorporate resourceVersion bounds based on staleness criteria.
	})
	if err != nil {
		return nil, err
	}
	for idx := range lo.Items {
		ret = append(ret, &lo.Items[idx])
	}
	return ret, nil
}

func (w *wrapper) Get(name string) (*apiappsv1beta2.DaemonSet, error) {
	return w.client.AppsV1beta2().DaemonSets(w.namespace).Get(context.TODO(), name, v1.GetOptions{
		// TODO(mattmoor): Incorporate resourceVersion bounds based on staleness criteria.
	})
}