package main

import (
	"encoding/json"
	"fmt"
	"github.com/fission/fission"
	controller "github.com/fission/fission/controller/client"
	"github.com/urfave/cli"
	"github.com/yqf3139/fission-benchmark/requester"
	"github.com/yqf3139/fission-benchmark/tpr"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/1.5/pkg/api"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

type WorkloadRunner func(function *tpr.Function, fissionFunc *fission.Function,
	workload *tpr.Workload, controller *controller.Client, routerUrl string,
	ch chan *requester.ReportWrapper, interrupt chan struct{})

type TriggerSpecMarshaler func(slice map[string]string) (interface{}, error)

func runHttpTriggerFunction(f *tpr.Function, fissionFunc *fission.Function, w *tpr.Workload,
	controller *controller.Client, routerUrl string, interrupt chan struct{}) (*requester.Report, error) {
	log.Printf("Function: %v, Workload: %v\n", f.Name, w.Name)
	metadata, err := controller.FunctionCreate(&fission.Function{
		Metadata: fission.Metadata{
			Name: FUNCTION_PREFIX + fissionFunc.Name,
		},
		Environment: fission.Metadata{
			Name: fissionFunc.Environment.Name,
		},
		Code: fissionFunc.Code,
	})
	if err != nil {
		return nil, err
	}

	// wait for the function info sync in fission
	log.Print("Function created, waiting for the sync ... ")
	time.Sleep(4 * time.Second)
	log.Println("done.")

	httpTriggerSpec := w.Trigger.ParsedSpec.(tpr.HttpTriggerSpec)

	req, err := http.NewRequest(strings.ToUpper(httpTriggerSpec.Method),
		fmt.Sprintf("%v/fission-function/%v%v", routerUrl, FUNCTION_PREFIX, f.Name), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	if w.Kind == "warm" {
		log.Print("Pre request to warm the function up ... ")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode >= 300 {
			return nil, fmt.Errorf("Pre request for %v not success", f.Name)
		}
		log.Println("done.")
	} else {
		log.Println("Pre request is disabled, testing cold fission")
	}

	log.Printf("Requesting with disable-keep-alive[%v] ... \n", w.DisableKeepAlive)
	report := (&requester.Work{
		Request:           req,
		RequestBody:       []byte(httpTriggerSpec.Data),
		N:                 w.Number,
		C:                 w.Concurrence,
		QPS:               w.Qps,
		Timeout:           w.Timeout,
		Output:            "",
		EnableTrace:       false,
		DisableKeepAlives: w.DisableKeepAlive,
		Stopped:           false,
		Interrupt:         interrupt,
	}).Run(w.Verbose)
	report.Finalize()
	report.Print(w.Verbose, w.Verbose)
	err = controller.FunctionDelete(metadata)
	return report, err
}

func runHttpTriggerControl(w *tpr.Workload, control tpr.Endpoint, interrupt chan struct{}) (*requester.Report, error) {
	httpTriggerSpec := w.Trigger.ParsedSpec.(tpr.HttpTriggerSpec)
	req, err := http.NewRequest(
		strings.ToUpper(httpTriggerSpec.Method), control.Endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	log.Printf("Requesting with disable-keep-alive[%v] ... \n", w.DisableKeepAlive)
	report := (&requester.Work{
		Request:           req,
		RequestBody:       []byte(httpTriggerSpec.Data),
		N:                 w.Number,
		C:                 w.Concurrence,
		QPS:               w.Qps,
		Timeout:           w.Timeout,
		Output:            "",
		EnableTrace:       false,
		DisableKeepAlives: w.DisableKeepAlive,
		Stopped:           false,
		Interrupt:         interrupt,
	}).Run(w.Verbose)
	report.Finalize()
	report.Print(w.Verbose, w.Verbose)
	return report, nil
}

var triggerKind2WorkloadRunner = map[string]WorkloadRunner{
	"Http": func(f *tpr.Function, fissionFunc *fission.Function, w *tpr.Workload,
		controller *controller.Client, routerUrl string,
		ch chan *requester.ReportWrapper, interrupt chan struct{}) {

		stopped := false
		interruptSubChan := make(chan struct{}, 1)

		go func() {
			<-interrupt
			stopped = true
			if interruptSubChan != nil {
				interruptSubChan <- struct{}{}
			}
		}()
		report, err := runHttpTriggerFunction(f, fissionFunc, w, controller, routerUrl, interruptSubChan)
		time.Sleep(time.Second)
		close(interruptSubChan)
		interruptSubChan = nil
		ch <- &requester.ReportWrapper{Report: report, Error: err}

		for _, control := range f.Controls {
			if stopped {
				break
			}
			interruptSubChan = make(chan struct{}, 1)
			report, err := runHttpTriggerControl(w, control, interruptSubChan)
			close(interruptSubChan)
			interruptSubChan = nil
			ch <- &requester.ReportWrapper{Report: report, Error: err}
		}
	},
}

var triggerKind2TriggerSpecMarshaler = map[string]TriggerSpecMarshaler{
	"Http": func(slice map[string]string) (interface{}, error) {
		httpTriggerSpec := tpr.HttpTriggerSpec{}
		data, err := yaml.Marshal(slice)
		if err != nil {
			return nil, err
		}
		yaml.Unmarshal(data, &httpTriggerSpec)
		return httpTriggerSpec, nil
	},
}

func fetchConfig(filePath string) (*tpr.Config, error) {
	data := fetchFile("", filePath)
	config := tpr.Config{}

	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func runWorkloads(c *cli.Context) error {
	controller := getController(c.GlobalString("controller"))
	routerUrl := getRouterUrl(c.GlobalString("router"))

	filePath := c.String("f")
	if len(filePath) == 0 {
		fatal("Need a workload config, use -f.")
	}

	reportPath := c.String("r")
	if len(reportPath) == 0 {
		fatal("Need a workload report path, use -r.")
	}

	config, err := fetchConfig(filePath)
	if err != nil {
		panic(err)
	}

	instanceApi := MakeInstanceLocalInterface()
	configApi := MakeConfigLocalInterface()

	instance := &tpr.Instance{
		Metadata: api.ObjectMeta{
			Name: "instance",
			Labels: map[string]string{
				"config": config.Metadata.Name,
			},
		},
	}
	// create the related config
	configApi.Create(config)

	runner := &InstanceRunner{
		instance:    instance,
		instanceApi: instanceApi,
		configApi:   configApi,
		controller:  controller,
		routerUrl:   routerUrl,
	}

	if !runner.validate() {
		return fmt.Errorf("Config is not validated")
	}

	// receive int sig
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch
		log.Println("try to interrupt runner")
		runner.interrupt()
	}()

	// blocking on the runner thread
	runner.run()

	// print the report file in json
	// only one instance exists
	list, err := instanceApi.List(api.ListOptions{})
	if err != nil {
		panic(err)
	}
	if len(list.Items) != 1 {
		panic("Instance list length is not 1")
	}
	return writeInstanceReports(reportPath, list.Items[0])
}

func writeInstanceReports(path string, instance tpr.Instance) error {
	content, err := json.Marshal(instance.Spec.Reports)
	if err != nil {
		return err
	}
	return writeFile(path, content)
}
