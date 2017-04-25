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
	"github.com/urfave/cli"
	"os"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(-1)

	app := cli.NewApp()
	app.Name = "fission-benchmark"
	app.Usage = "Benchmark tools and workloads for Fission"
	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "controller", Usage: "Fission controller URL", EnvVar: "FISSION_CONTROLLER"},
		cli.StringFlag{Name: "router", Usage: "Fission router URL", EnvVar: "FISSION_ROUTER"},
		cli.StringFlag{Name: "kubeconfig", Usage: "Path to a kube config.", EnvVar: "KUBE_CONFIG"},
	}

	fileFlag := cli.StringFlag{Name: "f", Usage: "workload config name"}
	reportFlag := cli.StringFlag{Name: "r", Usage: "workload report name"}

	app.Commands = []cli.Command{
		{Name: "run", Usage: "Run workload locally", Flags: []cli.Flag{fileFlag, reportFlag},
			Action: runWorkloads},
		{Name: "service", Usage: "Run as a service", Flags: []cli.Flag{}, Action: runService},
	}

	app.Run(os.Args)
}
