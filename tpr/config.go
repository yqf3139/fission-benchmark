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
	"k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/pkg/watch"
	"k8s.io/client-go/1.5/rest"
)

type (
	ConfigInterface interface {
		Create(*Config) (*Config, error)
		Get(name string) (*Config, error)
		Update(*Config) (*Config, error)
		Delete(name string, options *api.DeleteOptions) error
		List(opts api.ListOptions) (*ConfigList, error)
		Watch(opts api.ListOptions) (watch.Interface, error)
	}

	configClient struct {
		client    *rest.RESTClient
		namespace string
	}
)

func MakeConfigInterface(tprClient *rest.RESTClient, namespace string) ConfigInterface {
	return &configClient{
		client:    tprClient,
		namespace: namespace,
	}
}

func (cc *configClient) Create(f *Config) (*Config, error) {
	var result Config
	err := cc.client.Post().
		Resource("configs").
		Namespace(cc.namespace).
		Body(f).
		Do().Into(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (cc *configClient) Get(name string) (*Config, error) {
	var result Config
	err := cc.client.Get().
		Resource("configs").
		Namespace(cc.namespace).
		Name(name).
		Do().Into(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (cc *configClient) Update(f *Config) (*Config, error) {
	var result Config
	err := cc.client.Put().
		Resource("configs").
		Namespace(cc.namespace).
		Name(f.Metadata.Name).
		Body(f).
		Do().Into(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (cc *configClient) Delete(name string, opts *api.DeleteOptions) error {
	return cc.client.Delete().
		Namespace(cc.namespace).
		Resource("configs").
		Name(name).
		Body(opts).
		Do().
		Error()
}

func (cc *configClient) List(opts api.ListOptions) (*ConfigList, error) {
	var result ConfigList
	err := cc.client.Get().
		Namespace(cc.namespace).
		Resource("configs").
		VersionedParams(&opts, api.ParameterCodec).
		Do().
		Into(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (cc *configClient) Watch(opts api.ListOptions) (watch.Interface, error) {
	return cc.client.Get().
		Prefix("watch").
		Namespace(cc.namespace).
		Resource("configs").
		VersionedParams(&opts, api.ParameterCodec).
		Watch()
}
