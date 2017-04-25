Result of fission 01b3991, on a 8-thread server:
```bash

Running Spec #0
Function: big-response, Workload: simple-workload
Function created, waiting for the sync ... done.
Pre request to warm the function up ... done.
Requesting ... 
All requests done.

Summary:
  Total:	5.0221 secs
  Slowest:	0.0477 secs
  Fastest:	0.0112 secs
  Average:	0.0203 secs
  Requests/sec:	9.9561
  Total data:	104857600 bytes
  Size/request:	2097152 bytes

Status code distribution:
  [200]	50 responses

Running #0 control
Requesting ... 
All requests done.

Summary:
  Total:	5.0174 secs
  Slowest:	0.0250 secs
  Fastest:	0.0074 secs
  Average:	0.0150 secs
  Requests/sec:	9.9654
  Total data:	104857600 bytes
  Size/request:	2097152 bytes

Status code distribution:
  [200]	50 responses

Running #1 control
Requesting ... 
All requests done.

Summary:
  Total:	5.0173 secs
  Slowest:	0.0239 secs
  Fastest:	0.0073 secs
  Average:	0.0137 secs
  Requests/sec:	9.9656
  Total data:	104857600 bytes
  Size/request:	2097152 bytes

Status code distribution:
  [200]	50 responses

Running Spec #1
Function: big-response, Workload: medium-workload
Function created, waiting for the sync ... done.
Pre request to warm the function up ... done.
Requesting ... 
All requests done.

Summary:
  Total:	1.6768 secs
  Slowest:	0.6100 secs
  Fastest:	0.0094 secs
  Average:	0.0685 secs
  Requests/sec:	119.2759
  Total data:	419430400 bytes
  Size/request:	2097152 bytes

Status code distribution:
  [200]	200 responses

Running #0 control
Requesting ... 
All requests done.

Summary:
  Total:	1.2418 secs
  Slowest:	0.1223 secs
  Fastest:	0.0060 secs
  Average:	0.0576 secs
  Requests/sec:	161.0519
  Total data:	419430400 bytes
  Size/request:	2097152 bytes

Status code distribution:
  [200]	200 responses

Running #1 control
Requesting ... 
All requests done.

Summary:
  Total:	0.6230 secs
  Slowest:	0.0600 secs
  Fastest:	0.0060 secs
  Average:	0.0226 secs
  Requests/sec:	321.0387
  Total data:	419430400 bytes
  Size/request:	2097152 bytes

Status code distribution:
  [200]	200 responses

Running Spec #2
Function: big-response, Workload: heavy-workload
Function created, waiting for the sync ... done.
Pre request to warm the function up ... done.
Requesting ... 
All requests done.

Summary:
  Total:	65.6733 secs
  Slowest:	51.8742 secs
  Fastest:	0.0128 secs
  Average:	1.8946 secs
  Requests/sec:	76.1345
  Total data:	10267656192 bytes
  Size/request:	2053531 bytes

Status code distribution:
  [200]	4896 responses
  [502]	104 responses

Running #0 control
Requesting ... 
All requests done.

Summary:
  Total:	30.3150 secs
  Slowest:	3.7079 secs
  Fastest:	0.0398 secs
  Average:	1.4411 secs
  Requests/sec:	164.9346
  Total data:	10485760000 bytes
  Size/request:	2097152 bytes

Status code distribution:
  [200]	5000 responses

Running #1 control
Requesting ... 
All requests done.

Summary:
  Total:	7.2993 secs
  Slowest:	0.7258 secs
  Fastest:	0.0062 secs
  Average:	0.3111 secs
  Requests/sec:	684.9962
  Total data:	10485760000 bytes
  Size/request:	2097152 bytes

Status code distribution:
  [200]	5000 responses



```
