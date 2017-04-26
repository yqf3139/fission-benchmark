# fission-benchmark
Benchmark tools and workloads for Fission.

# Compile
```bash
# Get dependencies
$ glide install
# Build benchmark-bundle binary
$ cd benchmark-bundle && go build
# For local use, just run the benchmark-bundle binary
# To use fission-benchmark as a service, a docker images should be built and pushed
$ ./push $TAG
```

# Concepts
- Function: information of a fission function, so benchmark can create functions in need
    - Function name: the function should be created in fission in advance
    - Control groups: some services endpoints to be compared with
- Workload: how to benchmark the Function and all the parameters
    - Trigger: currently only http trigger
    - Kind: Run the function on a cold or warm container
    - Other parameters: number, concurrence, timeout, qps
- Pairs: a list of Function-Workload pair, the workload will be applied to the related function

# Usage

## Add control endpoints
If we need to compare fission with raw k8s service and deployment, create an control endpoint:
```bash
# Build an image with the function built in
$ cd workloads/helloworld-workload && ./build.sh

# Modify the k8s yaml file
$ $EDITOR helloworld-k8s-control-group.yaml

# Create k8s service and endpoint
$ kubectl create -f helloworld-k8s-control-group.yaml

# Test if the endpoint is created
$ curl http://the-endpoint:port
```
And then add the control endpoint into the config file.

## Run the benchmark locally
Modify the benchmark yaml file and run the benchmark
```bash
$ $EDITOR helloworld-workload.yaml
$ export FISSION_CONTROLLER=$ip:port
$ export FISSION_ROUTER=$ip:port
$ fission-benchmark run -f helloworld-workload.yaml -r report.json

Running Spec #0
Function: helloworld, Workload: simple-workload
Function created, waiting for the sync ... done.
Pre request to warm the function up ... done.
All requests done.

Summary:
  Total:        5.0420 secs
  Slowest:      0.0811 secs
  Fastest:      0.0024 secs
  Average:      0.0278 secs
  Requests/sec: 9.9166
  Total data:   700 bytes
  Size/request: 14 bytes

Status code distribution:
  [200] 50 responses

Running #0 control
All requests done.

Summary:
  Total:        5.0027 secs
  Slowest:      0.0059 secs
  Fastest:      0.0012 secs
  Average:      0.0027 secs
  Requests/sec: 9.9946
  Total data:   700 bytes
  Size/request: 14 bytes

Status code distribution:
  [200] 50 responses
```

## Run the benchmark as a service in k8s

Fission-benchmark can also work as service. It watches the config and instance TPRs. 
Once there is a valid benchmark config and instance, the service will run the instance and update report into TPR.
The process can be controlled by creating and modifying k8s TPRs.

Create the fission-benchmark namespace and the deployment:
```bash
$ cd deploy
$ kubectl create -f k8s-namespace.yaml -f fission-benchmark.yaml
```

Create the workload config:
```bash
$ cd any-workload
$ kubectl create -f config.yaml
```

Create the config's instance:
```bash
# make sure the instance is labelled with the config's name
# to create an instance, leave the status as create-request
$ kubectl create -f instance.yaml
```

Run, stop, watch the instance:
```bash
# change the instance status
# running-request to run the instance
# stop-request to stop the running instance
# running-request -> running ->  finished
#                            `-> stop-request -> stopped
$ kubectl replace -f changed-instance.yaml

# to watch current instance
$ watch kubectl get instances hpa-workload-instance --namespace=fission-benchmark -o json
```

The report will be filled into instance and progress will be updated during running.