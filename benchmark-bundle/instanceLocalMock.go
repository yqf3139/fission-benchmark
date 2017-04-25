package main

import (
	"fmt"
	"github.com/yqf3139/fission-benchmark/tpr"
	"k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/pkg/watch"
)

type (
	instanceLocalClient struct {
		store map[string]*tpr.Instance
	}
)

func MakeInstanceLocalInterface() tpr.InstanceInterface {
	return &instanceLocalClient{
		store: make(map[string]*tpr.Instance),
	}
}

func (ic *instanceLocalClient) Create(f *tpr.Instance) (*tpr.Instance, error) {
	ic.store[f.Metadata.Name] = f
	return f, nil
}

func (ic *instanceLocalClient) Get(name string) (*tpr.Instance, error) {
	result := ic.store[name]
	if result == nil {
		return nil, fmt.Errorf("Config not found")
	}
	return result, nil
}

func (ic *instanceLocalClient) Update(f *tpr.Instance) (*tpr.Instance, error) {
	return ic.Create(f)
}

func (ic *instanceLocalClient) Delete(name string, opts *api.DeleteOptions) error {
	delete(ic.store, name)
	return nil
}

func (ic *instanceLocalClient) List(opts api.ListOptions) (*tpr.InstanceList, error) {
	items := make([]tpr.Instance, len(ic.store))
	idx := 0
	for _, value := range ic.store {
		items[idx] = *value
		idx++
	}
	list := &tpr.InstanceList{Items: items}
	return list, nil
}

func (ic *instanceLocalClient) Watch(opts api.ListOptions) (watch.Interface, error) {
	return nil, fmt.Errorf("Not implemented")
}
