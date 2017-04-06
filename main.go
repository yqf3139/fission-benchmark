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
	"fmt"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"os"
	"path"

	"github.com/fission/fission"
	controller "github.com/fission/fission/controller/client"
	"github.com/yqf3139/fission-benchmark/requester"
	"net/http"
	"runtime"
	"strings"
	"time"
)

var triggerKind2WorkloadRunner = map[string]WorkloadRunner{
	"Http": func(f Function, w Workload, controller *controller.Client, routerUrl string) error {
		fmt.Printf("Function: %v, Workload: %v\n", f.Name, w.Name)

		metadata, err := controller.FunctionCreate(&fission.Function{
			Metadata: fission.Metadata{
				Name: FUNCTION_PREFIX + f.Metadata.Name,
			},
			Environment: fission.Metadata{
				Name: f.Environment.Name,
			},
			Code: f.code,
		})
		if err != nil {
			return err
		}

		// wait for the function info sync in fission
		fmt.Print("Function created, waiting for the sync ... ")
		time.Sleep(4 * time.Second)
		fmt.Println("done.")

		spec := w.Trigger.spec.(HttpTriggerSpec)

		req, err := http.NewRequest(strings.ToUpper(spec.Method),
			fmt.Sprintf("%v/fission-function/%v%v", routerUrl, FUNCTION_PREFIX, f.Name), nil)
		if err != nil {
			return err
		}
		req.Header.Add("Content-Type", "application/json")

		if w.Kind == "warm" {
			fmt.Print("Pre request to warm the function up ... ")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return err
			}
			if resp.StatusCode >= 300 {
				return fmt.Errorf("Pre request for %v not success", f.Name)
			}
			fmt.Println("done.")
		} else {
			fmt.Println("Pre request is disabled, testing cold fission")
		}

		fmt.Println("Requesting ... ")
		report := (&requester.Work{
			Request:     req,
			RequestBody: []byte(spec.Data),
			N:           w.Number,
			C:           w.Concurrence,
			QPS:         w.Qps,
			Timeout:     w.Timeout,
			Output:      "",
			EnableTrace: false,
		}).Run(w.Verbose)
		report.Finalize()
		report.Print(w.Verbose, w.Verbose)

		err = controller.FunctionDelete(metadata)
		if err != nil {
			return err
		}

		for idx, control := range f.Controls {
			fmt.Printf("\nRunning #%v control\n", idx)
			req, err := http.NewRequest(
				strings.ToUpper(spec.Method), control.Endpoint, nil)
			if err != nil {
				return err
			}
			req.Header.Add("Content-Type", "application/json")

			fmt.Println("Requesting ... ")
			report = (&requester.Work{
				Request:     req,
				RequestBody: []byte(spec.Data),
				N:           w.Number,
				C:           w.Concurrence,
				QPS:         w.Qps,
				Timeout:     w.Timeout,
				Output:      "",
				EnableTrace: false,
			}).Run(w.Verbose)
			report.Finalize()
			report.Print(w.Verbose, w.Verbose)
		}
		fmt.Println()
		time.Sleep(1 * time.Second)
		return nil
	},
}

var triggerKind2TriggerSpecMarshaler = map[string]TriggerSpecMarshaler{
	"Http": func(slice yaml.MapSlice) (interface{}, error) {
		httpTriggerSpec := HttpTriggerSpec{}
		data, err := yaml.Marshal(slice)
		if err != nil {
			return nil, err
		}
		yaml.Unmarshal(data, &httpTriggerSpec)
		return httpTriggerSpec, nil
	},
}

func fetchConfig(filePath string) (*Config, error) {
	file := fetchFile("", filePath)
	config := Config{}
	yaml.Unmarshal(file, &config)
	for idx := range config.Workloads {
		w := &config.Workloads[idx]
		var err error

		marshaller := triggerKind2TriggerSpecMarshaler[w.Trigger.Kind]
		if marshaller == nil {
			return nil, fmt.Errorf("No marshaller for %v", w)
		}
		w.Trigger.spec, err = marshaller(w.Trigger.Spec)
		if err != nil {
			return nil, fmt.Errorf("Marshaler error for %v", w)
		}

		runner := triggerKind2WorkloadRunner[w.Trigger.Kind]
		if runner == nil {
			return nil, fmt.Errorf("No runner for %v", w)
		}
		w.Trigger.runner = runner
	}
	for idx := range config.Functions {
		f := &config.Functions[idx]
		prefix, _ := path.Split(filePath)
		f.code = string(fetchFile(prefix, f.File))
	}
	return &config, nil
}

func runWorkloads(c *cli.Context) error {
	controllerClient := getController(c.GlobalString("controller"))
	routerUrl := getRouterUrl(c.GlobalString("router"))

	filePath := c.String("f")
	if len(filePath) == 0 {
		fatal("Need a workload config, use -f.")
	}

	config, err := fetchConfig(filePath)
	if err != nil {
		return err
	}

	functionMap := make(map[string]Function)
	for _, f := range config.Functions {
		functionMap[f.Metadata.Name] = f
	}

	workloadMap := make(map[string]Workload)
	for _, w := range config.Workloads {
		workloadMap[w.Metadata.Name] = w
	}

	for i, s := range config.Specs {
		f, ok := functionMap[s.Function]
		if !ok {
			return fmt.Errorf("Function %v not found", s.Function)
		}
		w, ok := workloadMap[s.Workload]
		if !ok {
			return fmt.Errorf("Workload %v not found", s.Function)
		}
		fmt.Printf("Running Spec #%v\n", i)
		err = w.runner(f, w, controllerClient, routerUrl)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	runtime.GOMAXPROCS(-1)

	app := cli.NewApp()
	app.Name = "fission-benchmark"
	app.Usage = "Benchmark tools and workloads for Fission"
	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "controller", Usage: "Fission controller URL", EnvVar: "FISSION_CONTROLLER"},
		cli.StringFlag{Name: "router", Usage: "Fission router URL", EnvVar: "FISSION_ROUTER"},
	}

	fileFlag := cli.StringFlag{Name: "f", Usage: "workload config name"}

	app.Commands = []cli.Command{
		{Name: "run", Usage: "Run workload", Flags: []cli.Flag{fileFlag}, Action: runWorkloads},
	}

	app.Run(os.Args)
}
