/*
Copyright 2016 The Fission Authors.

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

package tpr

import (
	"k8s.io/client-go/1.5/kubernetes"
	"k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/pkg/api/errors"
	"k8s.io/client-go/1.5/pkg/api/unversioned"
	"k8s.io/client-go/1.5/pkg/api/v1"
	"k8s.io/client-go/1.5/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/1.5/pkg/runtime"
	"k8s.io/client-go/1.5/pkg/runtime/serializer"
	"k8s.io/client-go/1.5/rest"
	"k8s.io/client-go/1.5/tools/clientcmd"
)

func buildConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}

// Get a kubernetes client using the pod's service account.
func GetKubernetesClient(kubeconfig *string) (*rest.Config, *kubernetes.Clientset, error) {
	// creates the in-cluster config
	// Create the client config. Use kubeconfig if given, otherwise assume in-cluster.
	config, err := buildConfig(*kubeconfig)
	if err != nil {
		return nil, nil, err
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	return config, clientset, nil
}

// ensureTPR checks if the given TPR type exists, and creates it if
// needed. (Note that this creates the TPR type; it doesn't create any
// _instances_ of that type.)
func ensureTPR(clientset *kubernetes.Clientset, tpr *v1beta1.ThirdPartyResource) error {
	_, err := clientset.Extensions().ThirdPartyResources().Get(tpr.ObjectMeta.Name)
	if err != nil {
		if errors.IsNotFound(err) {
			_, err := clientset.Extensions().ThirdPartyResources().Create(tpr)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func EnsureFissionBenchmarkTPRs(clientset *kubernetes.Clientset) error {
	tprs := []v1beta1.ThirdPartyResource{
		{
			ObjectMeta: v1.ObjectMeta{
				Name: "config.benchmark.fission.io",
			},
			Versions: []v1beta1.APIVersion{
				{Name: "v1"},
			},
			Description: "Configs",
		},
		{
			ObjectMeta: v1.ObjectMeta{
				Name: "instance.benchmark.fission.io",
			},
			Versions: []v1beta1.APIVersion{
				{Name: "v1"},
			},
			Description: "Instances",
		},
	}
	for _, tpr := range tprs {
		err := ensureTPR(clientset, &tpr)
		if err != nil {
			return err
		}
	}
	return nil
}

// This is copied from the client-go TPR example.  (I don't understand
// all of it completely.)  It registers our types with the global API
// "scheme" (api.Scheme), which keeps a directory of types [I guess so
// it can use the string in the Kind field to make a Go object?].  It
// also puts the fission TPR types under a "group version" which we
// create for our TPRs types.
func configureClient(config *rest.Config) {
	groupversion := unversioned.GroupVersion{
		Group:   "benchmark.fission.io",
		Version: "v1",
	}
	config.GroupVersion = &groupversion
	config.APIPath = "/apis"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: api.Codecs}

	schemeBuilder := runtime.NewSchemeBuilder(
		func(scheme *runtime.Scheme) error {
			scheme.AddKnownTypes(
				groupversion,
				&Config{},
				&ConfigList{},
				&api.ListOptions{},
				&api.DeleteOptions{},
			)
			scheme.AddKnownTypes(
				groupversion,
				&Instance{},
				&InstanceList{},
				&api.ListOptions{},
				&api.DeleteOptions{},
			)
			return nil
		})
	schemeBuilder.AddToScheme(api.Scheme)
}

func GetTprClient(config *rest.Config) (*rest.RESTClient, error) {

	// mutate config to add our types
	configureClient(config)

	// make a REST client with that config
	return rest.RESTClientFor(config)
}
