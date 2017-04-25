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
	InstanceInterface interface {
		Create(*Instance) (*Instance, error)
		Get(name string) (*Instance, error)
		Update(*Instance) (*Instance, error)
		Delete(name string, options *api.DeleteOptions) error
		List(opts api.ListOptions) (*InstanceList, error)
		Watch(opts api.ListOptions) (watch.Interface, error)
	}

	instanceClient struct {
		client    *rest.RESTClient
		namespace string
	}
)

func MakeInstanceInterface(tprClient *rest.RESTClient, namespace string) InstanceInterface {
	return &instanceClient{
		client:    tprClient,
		namespace: namespace,
	}
}

func (ic *instanceClient) Create(ins *Instance) (*Instance, error) {
	var result Instance
	err := ic.client.Post().
		Resource("instances").
		Namespace(ic.namespace).
		Body(ins).
		Do().Into(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (ic *instanceClient) Get(name string) (*Instance, error) {
	var result Instance
	err := ic.client.Get().
		Resource("instances").
		Namespace(ic.namespace).
		Name(name).
		Do().Into(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (ic *instanceClient) Update(ins *Instance) (*Instance, error) {
	var result Instance
	err := ic.client.Put().
		Resource("instances").
		Namespace(ic.namespace).
		Name(ins.Metadata.Name).
		Body(ins).
		Do().Into(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (ic *instanceClient) Delete(name string, opts *api.DeleteOptions) error {
	return ic.client.Delete().
		Namespace(ic.namespace).
		Resource("instances").
		Name(name).
		Body(opts).
		Do().
		Error()
}

func (ic *instanceClient) List(opts api.ListOptions) (*InstanceList, error) {
	var result InstanceList
	err := ic.client.Get().
		Namespace(ic.namespace).
		Resource("instances").
		VersionedParams(&opts, api.ParameterCodec).
		Do().
		Into(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (ic *instanceClient) Watch(opts api.ListOptions) (watch.Interface, error) {
	return ic.client.Get().
		Prefix("watch").
		Namespace(ic.namespace).
		Resource("instances").
		VersionedParams(&opts, api.ParameterCodec).
		Watch()
}
