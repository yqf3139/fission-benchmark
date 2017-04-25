Result of fission 01b3991, on a 8-thread server:

```bash

Running Spec #0
Function: sleep, Workload: simple-workload
Function created, waiting for the sync ... done.
Pre request to warm the function up ... done.
Requesting ... 
All requests done.

Summary:
  Total:	5.1066 secs
  Slowest:	0.1136 secs
  Fastest:	0.1022 secs
  Average:	0.1050 secs
  Requests/sec:	9.7913
  Total data:	1000 bytes
  Size/request:	20 bytes

Status code distribution:
  [200]	50 responses

Running #0 control
Requesting ... 
All requests done.

Summary:
  Total:	5.1027 secs
  Slowest:	0.1090 secs
  Fastest:	0.1015 secs
  Average:	0.1038 secs
  Requests/sec:	9.7987
  Total data:	1000 bytes
  Size/request:	20 bytes

Status code distribution:
  [200]	50 responses

Running #1 control
Requesting ... 
All requests done.

Summary:
  Total:	5.1024 secs
  Slowest:	0.1100 secs
  Fastest:	0.1006 secs
  Average:	0.1022 secs
  Requests/sec:	9.7993
  Total data:	1000 bytes
  Size/request:	20 bytes

Status code distribution:
  [200]	50 responses

Running Spec #1
Function: sleep, Workload: medium-workload
Function created, waiting for the sync ... done.
Pre request to warm the function up ... done.
Requesting ... 
All requests done.

Summary:
  Total:	2.1736 secs
  Slowest:	0.1218 secs
  Fastest:	0.1016 secs
  Average:	0.1069 secs
  Requests/sec:	92.0138
  Total data:	4000 bytes
  Size/request:	20 bytes

Status code distribution:
  [200]	200 responses

Running #0 control
Requesting ... 
All requests done.

Summary:
  Total:	2.1092 secs
  Slowest:	0.1126 secs
  Fastest:	0.1005 secs
  Average:	0.1045 secs
  Requests/sec:	94.8247
  Total data:	4000 bytes
  Size/request:	20 bytes

Status code distribution:
  [200]	200 responses

Running #1 control
Requesting ... 
All requests done.

Summary:
  Total:	2.0787 secs
  Slowest:	0.1146 secs
  Fastest:	0.1006 secs
  Average:	0.1028 secs
  Requests/sec:	96.2120
  Total data:	4000 bytes
  Size/request:	20 bytes

Status code distribution:
  [200]	200 responses

Running Spec #2
Function: sleep, Workload: heavy-workload
Function created, waiting for the sync ... done.
Pre request to warm the function up ... done.
Requesting ... 
All requests done.

Summary:
  Total:	3.9268 secs
  Slowest:	0.4744 secs
  Fastest:	0.1016 secs
  Average:	0.1826 secs
  Requests/sec:	1273.3113
  Total data:	100000 bytes
  Size/request:	20 bytes

Status code distribution:
  [200]	5000 responses

Running #0 control
Requesting ... 
All requests done.

Summary:
  Total:	2.7690 secs
  Slowest:	0.2273 secs
  Fastest:	0.1041 secs
  Average:	0.1351 secs
  Requests/sec:	1805.7008
  Total data:	100000 bytes
  Size/request:	20 bytes

Status code distribution:
  [200]	5000 responses

Running #1 control
Requesting ... 
All requests done.

Summary:
  Total:	2.3248 secs
  Slowest:	0.1671 secs
  Fastest:	0.0999 secs
  Average:	0.1089 secs
  Requests/sec:	2150.6859
  Total data:	100000 bytes
  Size/request:	20 bytes

Status code distribution:
  [200]	5000 responses


```
