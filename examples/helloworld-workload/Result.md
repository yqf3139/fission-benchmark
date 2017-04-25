Result of fission 01b3991, on a 8 thread server:

```bash

Running Spec #0
Function: helloworld, Workload: simple-workload
Function created, waiting for the sync ... done.
Pre request to warm the function up ... done.
Requesting ... 
All requests done.

Summary:
  Total:	5.0063 secs
  Slowest:	0.0386 secs
  Fastest:	0.0026 secs
  Average:	0.0075 secs
  Requests/sec:	9.9875
  Total data:	41150 bytes
  Size/request:	823 bytes

Status code distribution:
  [200]	50 responses

Running #0 control
Requesting ... 
All requests done.

Summary:
  Total:	5.0046 secs
  Slowest:	0.0080 secs
  Fastest:	0.0014 secs
  Average:	0.0029 secs
  Requests/sec:	9.9909
  Total data:	45650 bytes
  Size/request:	913 bytes

Status code distribution:
  [200]	50 responses

Running #1 control
Requesting ... 
All requests done.

Summary:
  Total:	5.0030 secs
  Slowest:	0.0068 secs
  Fastest:	0.0018 secs
  Average:	0.0035 secs
  Requests/sec:	9.9940
  Total data:	45650 bytes
  Size/request:	913 bytes

Status code distribution:
  [200]	50 responses

Running Spec #1
Function: helloworld, Workload: medium-workload
Function created, waiting for the sync ... done.
Pre request to warm the function up ... done.
Requesting ... 
All requests done.

Summary:
  Total:	0.3853 secs
  Slowest:	0.0491 secs
  Fastest:	0.0016 secs
  Average:	0.0153 secs
  Requests/sec:	519.1352
  Total data:	164600 bytes
  Size/request:	823 bytes

Status code distribution:
  [200]	200 responses

Running #0 control
Requesting ... 
All requests done.

Summary:
  Total:	0.3425 secs
  Slowest:	0.0146 secs
  Fastest:	0.0009 secs
  Average:	0.0057 secs
  Requests/sec:	583.9057
  Total data:	182600 bytes
  Size/request:	913 bytes

Status code distribution:
  [200]	200 responses

Running #1 control
Requesting ... 
All requests done.

Summary:
  Total:	0.3382 secs
  Slowest:	0.0264 secs
  Fastest:	0.0008 secs
  Average:	0.0056 secs
  Requests/sec:	591.4065
  Total data:	182600 bytes
  Size/request:	913 bytes

Status code distribution:
  [200]	200 responses

Running Spec #2
Function: helloworld, Workload: heavy-workload
Function created, waiting for the sync ... done.
Pre request to warm the function up ... done.
Requesting ... 
All requests done.

Summary:
  Total:	103.3595 secs
  Slowest:	51.4605 secs
  Fastest:	0.0010 secs
  Average:	1.5794 secs
  Requests/sec:	95.7822
  Total data:	7923021 bytes
  Size/request:	800 bytes

Status code distribution:
  [200]	9627 responses
  [502]	273 responses

Running #0 control
Requesting ... 
All requests done.

Summary:
  Total:	5.4150 secs
  Slowest:	1.3835 secs
  Fastest:	0.0013 secs
  Average:	0.1527 secs
  Requests/sec:	1828.2483
  Total data:	9038700 bytes
  Size/request:	913 bytes

Status code distribution:
  [200]	9900 responses

Running #1 control
Requesting ... 
All requests done.

Summary:
  Total:	1.8896 secs
  Slowest:	0.2906 secs
  Fastest:	0.0009 secs
  Average:	0.0505 secs
  Requests/sec:	5239.2992
  Total data:	9038700 bytes
  Size/request:	913 bytes

Status code distribution:
  [200]	9900 responses

```