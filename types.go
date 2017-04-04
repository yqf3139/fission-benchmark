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

package main

import (
	"gopkg.in/yaml.v2"

	controller "github.com/fission/fission/controller/client"
)

type (
	Metadata struct {
		Name string
		Desc string `yaml:"desc,omitempty"`
	}

	Config struct {
		Version string
		Metadata
		Specs     []Spec
		Functions []Function
		Workloads []Workload
	}

	Spec struct {
		Function string
		Workload string
	}

	Function struct {
		Metadata
		Environment     Metadata
		File            string
		code            string
		RunControlGroup bool
		Controls        []Endpoint
	}

	Workload struct {
		Metadata
		Kind        string
		Number      int
		Concurrence int
		Timeout     int
		Qps         int
		Verbose     bool
		Trigger
	}

	Trigger struct {
		Kind   string
		Spec   yaml.MapSlice
		spec   interface{}
		runner WorkloadRunner
	}

	HttpTriggerSpec struct {
		Method string
		Data   string
	}

	Endpoint struct {
		Metadata
		Endpoint string
	}

	WorkloadRunner       func(f Function, w Workload, controller *controller.Client, routerUrl string) error
	TriggerSpecMarshaler func(slice yaml.MapSlice) (interface{}, error)
)

const (
	FUNCTION_PREFIX = "fission-benchmark-"
)
