Result of fission 01b3991, on a 8-thread server:

`workload.yaml` results:
```bash
Running Spec #0
Function: cpu-intensive, Workload: simple-workload
Function created, waiting for the sync ... done.
Pre request to warm the function up ... done.
Requesting ... 
All requests done.

Summary:
  Total:	5.2840 secs
  Slowest:	0.2171 secs
  Fastest:	0.1125 secs
  Average:	0.1782 secs
  Requests/sec:	9.4625
  Total data:	1000 bytes
  Size/request:	20 bytes

Status code distribution:
  [200]	50 responses

Running #0 control
Requesting ... 
All requests done.

Summary:
  Total:	5.2115 secs
  Slowest:	0.2104 secs
  Fastest:	0.1001 secs
  Average:	0.1545 secs
  Requests/sec:	9.5942
  Total data:	1000 bytes
  Size/request:	20 bytes

Status code distribution:
  [200]	50 responses

Running #1 control
Requesting ... 
All requests done.

Summary:
  Total:	5.1011 secs
  Slowest:	0.1360 secs
  Fastest:	0.1007 secs
  Average:	0.1024 secs
  Requests/sec:	9.8019
  Total data:	1000 bytes
  Size/request:	20 bytes

Status code distribution:
  [200]	50 responses

Running Spec #1
Function: cpu-intensive, Workload: medium-workload
Function created, waiting for the sync ... done.
Pre request to warm the function up ... done.
Requesting ... 
All requests done.

Summary:
  Total:	20.3160 secs
  Slowest:	1.9439 secs
  Fastest:	0.1020 secs
  Average:	0.9341 secs
  Requests/sec:	9.8445
  Total data:	4000 bytes
  Size/request:	20 bytes

Status code distribution:
  [200]	200 responses

Running #0 control
Requesting ... 
All requests done.

Summary:
  Total:	20.1515 secs
  Slowest:	1.9219 secs
  Fastest:	0.1017 secs
  Average:	0.9263 secs
  Requests/sec:	9.9248
  Total data:	4000 bytes
  Size/request:	20 bytes

Status code distribution:
  [200]	200 responses

Running #1 control
Requesting ... 
All requests done.

Summary:
  Total:	6.1189 secs
  Slowest:	0.5060 secs
  Fastest:	0.1001 secs
  Average:	0.2297 secs
  Requests/sec:	32.6859
  Total data:	4000 bytes
  Size/request:	20 bytes

Status code distribution:
  [200]	200 responses

Running Spec #2
Function: cpu-intensive, Workload: long-workload
Function created, waiting for the sync ... done.
Pre request to warm the function up ... done.
Requesting ... 
^Bn^BnAll requests done.

Summary:
  Total:	505.2944 secs
  Slowest:	0.9159 secs
  Fastest:	0.1066 secs
  Average:	0.4447 secs
  Requests/sec:	9.8952
  Total data:	100000 bytes
  Size/request:	20 bytes

Status code distribution:
  [200]	5000 responses

Running #0 control
Requesting ... 
All requests done.

Summary:
  Total:	504.0770 secs
  Slowest:	0.9165 secs
  Fastest:	0.1009 secs
  Average:	0.4455 secs
  Requests/sec:	9.9191
  Total data:	100000 bytes
  Size/request:	20 bytes

Status code distribution:
  [200]	5000 responses

Running #1 control
Requesting ... 
All requests done.

Summary:
  Total:	302.8938 secs
  Slowest:	0.5227 secs
  Fastest:	0.0999 secs
  Average:	0.2052 secs
  Requests/sec:	16.5074
  Total data:	100000 bytes
  Size/request:	20 bytes

Status code distribution:
  [200]	5000 responses

```

`hpa-worload.yaml` results:

```bash
Running Spec #0
Function: cpu-intensive, Workload: long-long-workload
Function created, waiting for the sync ... done.
Pre request to warm the function up ... done.
Enter y to continue, otherwise skip ... 
y
Requesting with disable-keep-alive[true] ... 
All requests done.

Summary:
  Total:        303.4516 secs
  Slowest:      12.0493 secs
  Fastest:      0.1030 secs
  Average:      5.9152 secs
  Requests/sec: 9.8863
  Total data:   60000 bytes
  Size/request: 20 bytes

Status code distribution:
  [200] 3000 responses

Running #0 control control1
Enter y to continue, otherwise skip ... 
y
Requesting with disable-keep-alive[true] ... 
All requests done.

Summary:
  Total:        303.8019 secs
  Slowest:      6.1769 secs
  Fastest:      0.1677 secs
  Average:      6.0164 secs
  Requests/sec: 9.8749
  Total data:   60000 bytes
  Size/request: 20 bytes
  
Status code distribution:
  [200] 3000 responses

Running #1 control control2
Enter y to continue, otherwise skip ... 
y
Requesting with disable-keep-alive[true] ... 
All requests done.

Summary:
  Total:        65.8950 secs
  Slowest:      4.9811 secs
  Fastest:      0.1006 secs
  Average:      1.1978 secs
  Requests/sec: 45.5270
  Total data:   60000 bytes
  Size/request: 20 bytes

Status code distribution:
  [200] 3000 responses

Running #2 control control3
Enter y to continue, otherwise skip ... 
y
Requesting with disable-keep-alive[true] ... 
All requests done.

Summary:
  Total:        147.3005 secs
  Slowest:      6.0881 secs
  Fastest:      0.1007 secs
  Average:      2.7982 secs
  Requests/sec: 20.3597
  Total data:   59980 bytes
  Size/request: 20 bytes

Status code distribution:
  [200] 2999 responses

```