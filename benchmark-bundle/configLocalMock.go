package main

import (
	"fmt"
	"github.com/yqf3139/fission-benchmark/tpr"
	"k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/pkg/watch"
)

type (
	configLocalClient struct {
		store map[string]*tpr.Config
	}
)

func MakeConfigLocalInterface() tpr.ConfigInterface {
	return &configLocalClient{
		store: make(map[string]*tpr.Config),
	}
}

func (cc *configLocalClient) Create(f *tpr.Config) (*tpr.Config, error) {
	cc.store[f.Metadata.Name] = f
	return f, nil
}

func (cc *configLocalClient) Get(name string) (*tpr.Config, error) {
	result := cc.store[name]
	if result == nil {
		return nil, fmt.Errorf("Config not found")
	}
	return result, nil
}

func (cc *configLocalClient) Update(f *tpr.Config) (*tpr.Config, error) {
	return cc.Create(f)
}

func (cc *configLocalClient) Delete(name string, opts *api.DeleteOptions) error {
	delete(cc.store, name)
	return nil
}

func (cc *configLocalClient) List(opts api.ListOptions) (*tpr.ConfigList, error) {
	items := make([]tpr.Config, len(cc.store))
	idx := 0
	for _, value := range cc.store {
		items[idx] = *value
		idx++
	}
	list := &tpr.ConfigList{Items: items}
	return list, nil
}

func (cc *configLocalClient) Watch(opts api.ListOptions) (watch.Interface, error) {
	return nil, fmt.Errorf("Not implemented")
}
