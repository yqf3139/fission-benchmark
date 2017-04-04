Result of fission 01b3991, on a minikube:

```bash
Running Spec #0
Function: helloworld, Workload: simple-workload
Function created, waiting for the sync ... done.
Pre request to warm the function up ... done.
Requesting ... 
All requests done.

Summary:
  Total:	5.0420 secs
  Slowest:	0.0851 secs
  Fastest:	0.0024 secs
  Average:	0.0248 secs
  Requests/sec:	9.9167
  Total data:	700 bytes
  Size/request:	14 bytes

Status code distribution:
  [200]	50 responses

Running #0 control
Requesting ... 
All requests done.

Summary:
  Total:	5.0059 secs
  Slowest:	0.0175 secs
  Fastest:	0.0017 secs
  Average:	0.0042 secs
  Requests/sec:	9.9882
  Total data:	700 bytes
  Size/request:	14 bytes

Status code distribution:
  [200]	50 responses

Running Spec #1
Function: helloworld, Workload: medium-workload
Function created, waiting for the sync ... done.
Pre request to warm the function up ... done.
Requesting ... 
All requests done.

Summary:
  Total:	2.0654 secs
  Slowest:	1.0927 secs
  Fastest:	0.0053 secs
  Average:	0.0635 secs
  Requests/sec:	96.8319
  Total data:	2800 bytes
  Size/request:	14 bytes

Status code distribution:
  [200]	200 responses

Running #0 control
Requesting ... 
All requests done.

Summary:
  Total:	1.1520 secs
  Slowest:	0.5996 secs
  Fastest:	0.0010 secs
  Average:	0.0560 secs
  Requests/sec:	173.6115
  Total data:	2800 bytes
  Size/request:	14 bytes

Status code distribution:
  [200]	200 responses

Running Spec #2
Function: helloworld, Workload: heavy-workload
Function created, waiting for the sync ... done.
Pre request to warm the function up ... done.
Requesting ... 
All requests done.

Summary:
  Total:	56.5146 secs
  Slowest:	53.3861 secs
  Fastest:	0.0018 secs
  Average:	1.1391 secs
  Requests/sec:	8.4934
  Total data:	6608 bytes
  Size/request:	13 bytes

Status code distribution:
  [200]	472 responses
  [502]	8 responses

Running #0 control
Requesting ... 
All requests done.

Summary:
  Total:	0.7800 secs
  Slowest:	0.1060 secs
  Fastest:	0.0027 secs
  Average:	0.0452 secs
  Requests/sec:	615.3712
  Total data:	6720 bytes
  Size/request:	14 bytes

Status code distribution:
  [200]	480 responses

```