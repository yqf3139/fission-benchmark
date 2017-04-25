package main

import (
	"github.com/fission/fission"
	controller "github.com/fission/fission/controller/client"
	"github.com/urfave/cli"
	"github.com/yqf3139/fission-benchmark/requester"
	"github.com/yqf3139/fission-benchmark/tpr"
	"k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/pkg/api/unversioned"
	"k8s.io/client-go/1.5/pkg/labels"
	"k8s.io/client-go/1.5/pkg/watch"
	"log"
	"time"
)

const BENCHMARK_NAMESPACE = "fission-benchmark"

type Worker struct {
	configApi   tpr.ConfigInterface
	instanceApi tpr.InstanceInterface
	controller  *controller.Client
	routerUrl   string
	instanceMap map[string]*InstanceRunner
}

type InstanceRunner struct {
	instance    *tpr.Instance
	configApi   tpr.ConfigInterface
	instanceApi tpr.InstanceInterface
	controller  *controller.Client
	routerUrl   string

	config            *tpr.Config
	functionMap       map[string]tpr.Function
	workloadMap       map[string]tpr.Workload
	fissionFuncMap    map[string]*fission.Function
	workloadRunnerMap map[string]WorkloadRunner
	reportNumber      int
	interruptChan     chan struct{}
	stopped           bool
}

func (runner *InstanceRunner) updateReport(instance *tpr.Instance, ch chan *requester.ReportWrapper) {
	finished := 0
	var ins *tpr.Instance = nil
	for {
		wrapper, more := <-ch
		if !more {
			break
		}
		err := wrapper.Error
		if err != nil {
			logErr(err, "Workload error occured")
			e := tpr.Error{"WORKLOAD_ERROR", err.Error()}
			instance.Spec.Errors = append(instance.Spec.Errors, e)
		}

		finished += 1
		progress := float32(finished) / float32(runner.reportNumber)
		instance.Spec.Progress = progress

		if wrapper.Report != nil {
			idx := instance.Spec.ReportIndex[finished-1]
			instance.Spec.Reports[idx.X][idx.Y] = *wrapper.Report
		}

		if finished == runner.reportNumber {
			break
		}
		ins, err = runner.instanceApi.Update(instance)
		*instance = *ins
		logErr(err, "Instance update error occured")
	}

	if finished < runner.reportNumber {
		return
	}
	instance.Spec.EndTimeStamp = time.Now().Unix()
	instance.Spec.Progress = 1.0
	instance.Spec.Status = "finished"

	ins, err := runner.instanceApi.Update(instance)
	logErr(err, "Instance update error occured")
	*instance = *ins
	if instance != nil && ins != nil {
		*instance = *ins
	}
}

func (runner *InstanceRunner) run() {
	conf := runner.config
	instance := runner.instance
	reportChan := make(chan *requester.ReportWrapper, 3)

	// update the instance tpr if new report received
	go runner.updateReport(instance, reportChan)
	defer close(reportChan)

	for i, p := range conf.Spec.Pairs {
		log.Printf("running pair #%v", i)
		function := runner.functionMap[p.Function]
		fissionFunc := runner.fissionFuncMap[p.Function]
		workload := runner.workloadMap[p.Workload]
		run := runner.workloadRunnerMap[p.Workload]

		runner.interruptChan = make(chan struct{}, 1)
		run(&function, fissionFunc, &workload, runner.controller, runner.routerUrl,
			reportChan, runner.interruptChan)
		close(runner.interruptChan)
		log.Printf("running pair #%v done", i)
		if runner.stopped {
			break
		}
	}
}

func (runner *InstanceRunner) interrupt() {
	if runner.interruptChan != nil {
		runner.stopped = true
		runner.interruptChan <- struct{}{}
	}
}

func (runner *InstanceRunner) validate() bool {
	// fetch conf from label
	instance := runner.instance
	configName := runner.instance.Metadata.Labels["config"]
	instance.Spec.Errors = make([]tpr.Error, 0)
	if len(configName) == 0 {
		e := tpr.Error{"CONFIG_NOT_FOND", ""}
		instance.Spec.Errors = append(instance.Spec.Errors, e)
		return false
	}

	conf, err := runner.configApi.Get(configName)
	if err != nil {
		e := tpr.Error{"CONFIG_NOT_FOND", err.Error()}
		instance.Spec.Errors = append(instance.Spec.Errors, e)
		return false
	}
	runner.config = conf

	runner.fissionFuncMap = make(map[string]*fission.Function)
	runner.functionMap = make(map[string]tpr.Function)
	for _, f := range conf.Spec.Functions {
		function, err := runner.controller.FunctionGet(&fission.Metadata{Name: f.Name})
		if err != nil {
			e := tpr.Error{"FISSION_FUNCTION_NOT_FOND", ""}
			instance.Spec.Errors = append(instance.Spec.Errors, e)
			return false
		}
		runner.fissionFuncMap[f.Name] = function
		runner.functionMap[f.Name] = f
	}

	runner.workloadMap = make(map[string]tpr.Workload)
	runner.workloadRunnerMap = make(map[string]WorkloadRunner)
	for _, w := range conf.Spec.Workloads {
		marshaller := triggerKind2TriggerSpecMarshaler[w.Trigger.Kind]
		if marshaller == nil {
			e := tpr.Error{"TRIGGER_MARSHALLER_NOT_FOND", ""}
			instance.Spec.Errors = append(instance.Spec.Errors, e)
			return false
		}

		w.Trigger.ParsedSpec, err = marshaller(w.Trigger.Spec)
		if err != nil {
			e := tpr.Error{"PARSE_TRIGGER_ERROR", ""}
			instance.Spec.Errors = append(instance.Spec.Errors, e)
			return false
		}

		workloadRunner := triggerKind2WorkloadRunner[w.Trigger.Kind]
		if workloadRunner == nil {
			e := tpr.Error{"WORKLOAD_RUNNER_NOT_FOND", ""}
			instance.Spec.Errors = append(instance.Spec.Errors, e)
			return false
		}

		runner.workloadRunnerMap[w.Name] = workloadRunner
		runner.workloadMap[w.Name] = w
	}

	runner.reportNumber = 0
	instance.Spec.Reports = make([][]requester.Report, len(conf.Spec.Pairs))
	instance.Spec.ReportIndex = make([]struct{ X, Y int }, 0)
	for i, p := range conf.Spec.Pairs {
		f, ok := runner.functionMap[p.Function]
		if !ok {
			e := tpr.Error{"PAIRED_FUNCTION_NOT_FOND", ""}
			instance.Spec.Errors = append(instance.Spec.Errors, e)
			return false
		}
		_, ok = runner.workloadMap[p.Workload]
		if !ok {
			e := tpr.Error{"PAIRED_WORKLOAD_NOT_FOND", ""}
			instance.Spec.Errors = append(instance.Spec.Errors, e)
			return false
		}
		instance.Spec.Reports[i] = make([]requester.Report, len(f.Controls)+1)

		// number of controls + self
		runner.reportNumber += len(f.Controls) + 1
		for j := 0; j < len(f.Controls)+1; j++ {
			instance.Spec.ReportIndex = append(instance.Spec.ReportIndex, struct{ X, Y int }{i, j})
		}
	}
	return true
}

func runService(c *cli.Context) error {
	controllerClient := getController(c.GlobalString("controller"))
	routerUrl := getRouterUrl(c.GlobalString("router"))

	kubeconfig := c.GlobalString("kubeconfig")
	config, clientset, err := tpr.GetKubernetesClient(&kubeconfig)
	if err != nil {
		return err
	}

	err = tpr.EnsureFissionBenchmarkTPRs(clientset)
	if err != nil {
		return err
	}

	// wait for the resource to be registered
	time.Sleep(time.Second * 2)

	client, err := tpr.GetTprClient(config)
	if err != nil {
		return err
	}

	worker := Worker{
		configApi:   tpr.MakeConfigInterface(client, BENCHMARK_NAMESPACE),
		instanceApi: tpr.MakeInstanceInterface(client, BENCHMARK_NAMESPACE),
		controller:  controllerClient,
		routerUrl:   routerUrl,
		instanceMap: make(map[string]*InstanceRunner),
	}

	go worker.instanceService()
	go worker.configService()

	for {
		time.Sleep(time.Second * 30)
	}

	log.Println("exit")
	return nil
}

func (w *Worker) instanceModify(instance *tpr.Instance) {
	// the main status transform routine
	// running-request -> pending -> running ->  finished
	//                                       `-> stop-request -> stopped
	// otherwise report the error

	uid := string(instance.Metadata.UID)
	var oldInstance *tpr.Instance = nil
	runner := w.instanceMap[uid]
	if runner != nil {
		oldInstance = runner.instance
	}

	switch instance.Spec.Status {
	case "running-request":
		if oldInstance != nil {
			e := tpr.Error{"INSTANCE_ONLY_RUN_ONCE", ""}
			oldInstance.Spec.Errors = append(oldInstance.Spec.Errors, e)
			oldInstance.Spec.Status = "stopped"
		} else {
			// validate the instance's config
			// make sure the labeled config exists
			// validate the config file
			oldInstance = instance
			oldInstance.Spec.Status = "running"
			oldInstance.Spec.StartTimeStamp = time.Now().Unix()
			runner = &InstanceRunner{
				instance:    oldInstance,
				instanceApi: w.instanceApi,
				configApi:   w.configApi,
				controller:  w.controller,
				routerUrl:   w.routerUrl,
			}
			if !runner.validate() {
				oldInstance.Spec.Status = "stopped"
				e := tpr.Error{"INSTANCE_NOT_VALIDATED", ""}
				oldInstance.Spec.Errors = append(oldInstance.Spec.Errors, e)
				break
			}
			w.instanceMap[uid] = runner
			go runner.run()
		}
		break
	case "stop-request":
		if oldInstance != nil {
			w.instanceStop(runner)
		} else {
			oldInstance = instance
			oldInstance.Spec.Status = "stopped"
			e := tpr.Error{"INSTANCE_NOT_RUNNING", ""}
			oldInstance.Spec.Errors = append(oldInstance.Spec.Errors, e)
		}
		break
	default:
		return
	}
	// finally update the instance
	oldInstance.Metadata = instance.Metadata
	ins, err := w.instanceApi.Update(oldInstance)
	logErr(err, "Update instance status")
	if oldInstance != nil && ins != nil {
		*oldInstance = *ins
	}
}

func (w *Worker) instanceAdd(instance *tpr.Instance) {
	// transform the status from create-request to created
	instance.Spec.Status = "created"
	ins, err := w.instanceApi.Update(instance)
	*instance = *ins
	logErr(err, "Update instance status to created")
}

func (w *Worker) instanceStop(runner *InstanceRunner) {
	instance := runner.instance
	if instance.Spec.Status == "running" {
		// send stop sig to the worker thread
		log.Println("interrupting runner")
		runner.interrupt()
		instance.Spec.Status = "stopped"
	} else {
		e := tpr.Error{"INSTANCE_NOT_RUNNING", ""}
		instance.Spec.Errors = append(instance.Spec.Errors, e)
	}
}

func (w *Worker) instanceDelete(instance *tpr.Instance) {
	// if the instance is running, stop the instance
	uid := string(instance.Metadata.UID)
	if runner, ok := w.instanceMap[uid]; ok {
		w.instanceStop(runner)
		delete(w.instanceMap, string(uid))
	}
}

func (w *Worker) instanceService() {
	listOptions := api.ListOptions{}

	wi, err := w.instanceApi.Watch(listOptions)
	if err != nil {
		panic(err)
	}
	for {
		log.Println("instance service waiting")
		res := <-wi.ResultChan()
		switch resType := res.Object.(type) {
		case *tpr.Instance:
			instance := res.Object.(*tpr.Instance)
			log.Println(res.Type, instance.Spec.Status)
			switch res.Type {
			case watch.Added:
				w.instanceAdd(instance)
				break
			case watch.Deleted:
				w.instanceDelete(instance)
				break
			case watch.Modified:
				w.instanceModify(instance)
				break
			case watch.Error:
				break
			}
			break
		case *unversioned.Status:
			log.Println("instance service get,", res.Object.(*unversioned.Status))
			break
		default:
			log.Println("instance service unkown type,", resType)
			wi.Stop()
			wi, err = w.instanceApi.Watch(listOptions)
			if err != nil {
				panic(err)
			}
			break
		}

	}
	log.Println("exit")
}

func (w *Worker) configService() {
	listOptions := api.ListOptions{}
	deleOptions := api.DeleteOptions{}
	wi, err := w.configApi.Watch(listOptions)
	if err != nil {
		panic(err)
	}
	for {
		res := <-wi.ResultChan()

		switch resType := res.Object.(type) {
		case *tpr.Config:
			config := res.Object.(*tpr.Config)
			switch res.Type {
			case watch.Deleted:
				// Delete all the instances labeled with the config
				set := labels.Set(map[string]string{
					"config": config.Metadata.Name,
				})
				list, err := w.instanceApi.List(api.ListOptions{LabelSelector: set.AsSelector()})
				log.Printf("delete %v instances", len(list.Items))
				if err != nil {
					logErr(err, "Select instances with config labeled")
					continue
				}

				for _, instance := range list.Items {
					err = w.instanceApi.Delete(instance.Metadata.Name, &deleOptions)
					logErr(err, "Delete instances with config labeled")
				}

				break
			default:
				break

			}
			break
		case *unversioned.Status:
			log.Println("config service get,", res.Object.(*unversioned.Status))
			break
		default:
			log.Println("config service unkown type,", resType)
			wi.Stop()
			wi, err = w.configApi.Watch(listOptions)
			if err != nil {
				panic(err)
			}
			break
		}
	}

}
