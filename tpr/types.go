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
	"github.com/yqf3139/fission-benchmark/requester"
	"k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/pkg/api/meta"
	"k8s.io/client-go/1.5/pkg/api/unversioned"
)

type (
	ConfigSpec struct {
		Pairs     []FunctionWorkloadPair `json:"pairs"`
		Functions []Function             `json:"functions"`
		Workloads []Workload             `json:"workloads"`
	}

	Config struct {
		unversioned.TypeMeta `json:",inline"`
		Metadata             api.ObjectMeta `json:"metadata"`
		Spec                 ConfigSpec     `json:"spec"`
	}

	ConfigList struct {
		unversioned.TypeMeta `json:",inline"`
		Metadata             unversioned.ListMeta `json:"metadata"`

		Items []Config `json:"items"`
	}

	FunctionWorkloadPair struct {
		Function string `json:"function"`
		Workload string `json:"workload"`
	}

	Function struct {
		Name     string     `json:"name"`
		Desc     string     `json:"desc"`
		Controls []Endpoint `json:"controls"`
	}

	Workload struct {
		Name             string  `json:"name"`
		Desc             string  `json:"desc"`
		Kind             string  `json:"kind"`
		Number           int     `json:"number"`
		Concurrence      int     `json:"concurrence"`
		Timeout          int     `json:"timeout"`
		Qps              int     `json:"qps"`
		Verbose          bool    `json:"verbose"`
		DisableKeepAlive bool    `json:"disablekeepalive"`
		Trigger          Trigger `json:"trigger"`
	}

	Trigger struct {
		Kind       string            `json:"kind"`
		Spec       map[string]string `json:"spec"`
		ParsedSpec interface{}       `json:"parsedspec"` // ignore
	}

	HttpTriggerSpec struct {
		Method string `json:"method"`
		Data   string `json:"data"`
	}

	Endpoint struct {
		Name     string `json:"name"`
		Desc     string `json:"desc"`
		Endpoint string `json:"endpoint"`
	}

	InstanceSpec struct {
		StartTimeStamp int64                `json:"starttimestamp"`
		EndTimeStamp   int64                `json:"endtimestamp"`
		Status         string               `json:"status"`
		Errors         []Error              `json:"errors"`
		Progress       float32              `json:"progress"`
		Reports        [][]requester.Report `json:"reports"`
		ReportIndex    []struct{ X, Y int } `json:"reportindex"`
	}

	Instance struct {
		unversioned.TypeMeta `json:",inline"`
		Metadata             api.ObjectMeta `json:"metadata"`
		Spec                 InstanceSpec   `json:"spec"`
	}

	InstanceList struct {
		unversioned.TypeMeta `json:",inline"`
		Metadata             unversioned.ListMeta `json:"metadata"`

		Items []Instance `json:"items"`
	}

	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
)

func (c *Config) GetObjectKind() unversioned.ObjectKind {
	return &c.TypeMeta
}

func (c *Config) GetObjectMeta() meta.Object {
	return &c.Metadata
}

func (i *Instance) GetObjectKind() unversioned.ObjectKind {
	return &i.TypeMeta
}

func (i *Instance) GetObjectMeta() meta.Object {
	return &i.Metadata
}

func (l *ConfigList) GetObjectKind() unversioned.ObjectKind {
	return &l.TypeMeta
}

func (l *ConfigList) GetListMeta() meta.List {
	return &l.Metadata
}

func (l *InstanceList) GetObjectKind() unversioned.ObjectKind {
	return &l.TypeMeta
}

func (l *InstanceList) GetListMeta() meta.List {
	return &l.Metadata
}
