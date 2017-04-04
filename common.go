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
	"os"
	"strings"

	controller "github.com/fission/fission/controller/client"
	"io/ioutil"
	"net/http"
)

func fatal(msg string) {
	os.Stderr.WriteString(msg + "\n")
	os.Exit(1)
}

func getController(controllerUrl string) *controller.Client {
	if len(controllerUrl) == 0 {
		fatal("Need --controller or FISSION_CONTROLLER set to your fission server.")
	}

	isHTTPS := strings.Index(controllerUrl, "https://") == 0
	isHTTP := strings.Index(controllerUrl, "http://") == 0

	if !(isHTTP || isHTTPS) {
		controllerUrl = "http://" + controllerUrl
	}

	return controller.MakeClient(controllerUrl)
}

func getRouterUrl(routerUrl string) string {
	if len(routerUrl) == 0 {
		fatal("Need --router or FISSION_ROUTER set to your fission router.")
	}

	isHTTPS := strings.Index(routerUrl, "https://") == 0
	isHTTP := strings.Index(routerUrl, "http://") == 0

	if !(isHTTP || isHTTPS) {
		routerUrl = "http://" + routerUrl
	}

	return routerUrl
}

func checkErr(err error, msg string) {
	if err != nil {
		fatal(fmt.Sprintf("Failed to %v: %v", msg, err))
	}
}

func fetchFile(prefix, filePath string) []byte {
	var file []byte
	var err error

	if strings.HasPrefix(filePath, "http://") || strings.HasPrefix(filePath, "https://") {
		var resp *http.Response
		resp, err = http.Get(filePath)
		if err != nil {
			checkErr(err, fmt.Sprintf("download function"))
		}

		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("%v - HTTP response returned non 200 status", resp.StatusCode)
			checkErr(err, fmt.Sprintf("download function"))
		}

		file, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			checkErr(err, fmt.Sprintf("download function body %v", filePath))
		}
	} else {
		file, err = ioutil.ReadFile(prefix + filePath)
		checkErr(err, fmt.Sprintf("read %v", prefix+filePath))
	}
	return file
}
