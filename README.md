# fission-benchmark
Benchmark tools and workloads for Fission

# Compile
```bash
# Get dependencies
$ glide install
# Build fission-benchmark binary
$ go build
```

# Concepts
- Function: information of a fission function, so benchmark can create functions in need
    - Function file and Environment
    - Control groups: some services endpoints to be compared with
- Workload: how to benchmark the Function and all the parameters
    - Trigger: currently only http trigger
    - Kind: Run the function on a cold or warm container
    - Other parameters: number, concurrence, timeout, qps
- Spec: a Function-Workload pair, the workload will be applied to the function

# Usage
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

Then modify the benchmark yaml file and run the benchmark
```bash
$ $EDITOR helloworld-workload.yaml
$ export FISSION_CONTROLLER=$ip:port
$ export FISSION_ROUTER=$ip:port
$ fission-benchmark run -f helloworld-workload.yaml

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