# Performance tuning


In this blog post we will see how we can se the performance analyzer tool for an ensemble model. The goal is to understand the models performance in find which part of the ensemble is the bottle neck for the throughput and see if we can find any potential solutions to move the bottleneck. Since performance improvements is a iterative process and we are more interested in showing the workflow in this blog post we will stop the performance hunt after the first iteration and instead leave it to the reader to continue the hunt. 

This will required that it is run on a linux machine since `--net=host` [is not supported using mac](https://docs.docker.com/network/host/)

> "The host networking driver only works on Linux hosts, and is not supported on Docker for Mac, Docker for Windows, or Docker EE for Windows Server."

```
docker run --runtime=nvidia -it --shm-size=1gb --rm --net=host -v$(pwd):/workspace/ -v/$(pwd)/model_repository:/models nvcr.io/nvidia/tritonserver:22.06-py3 bash
```

start the container with the performance analyzer aviliable. 

```
docker run -it --rm --net=host -v$(pwd):/workspace/  nvcr.io/nvidia/tritonserver:22.05-py3-sdk /bin/bash
```

Now it is time to run the performance analyzer. 

```
perf_analyzer -m ensemble_python_resnet50 -u localhost:8001 -i gRPC
```

however this returns

```
error: failed to create concurrency manager: input INPUT contains dynamic shape, provide shapes to send along with the request
```

due to the dynamic batching. The docs for the [perf_analyzer](https://github.com/triton-inference-server/server/blob/main/docs/perf_analyzer.md) tool mentiones the required parameter. You can also run `perf_analyzer --help` to get a full description oc the cli. 


The correct command is:

```
perf_analyzer -m ensemble_python_resnet50 -u localhost:8001 -i gRPC --concurrency-range 1:8 --shape INPUT:1 --input-data sample_input.json
```
Since we just send an array of `UINT8`, however it seems like perf_analyzer can not work with dynamic input sizes. And when we set it we seen to get an issue with the Python Pil library. 

> Failed to maintain requested inference load. Worker thread(s) failed to generate concurrent requests.
Thread [0] had error: in ensemble 'ensemble_python_resnet50', Failed to process the request(s) for model instance 'preprocessing_0', message: UnidentifiedImageError: cannot identify image file <_io.BytesIO object at 0x7f23b02dbd60> <br>  <br> At:
/usr/local/lib/python3.8/dist-packages/PIL/Image.py(3147): open
/models/preprocessing/1/model.py(125): execute

In this case it seems like the best option would be to [set the input data](https://github.com/triton-inference-server/server/blob/main/docs/perf_analyzer.md#real-input-data). In order to get the input data, we will deconstruct the client and check what do we exactly send(might feel a bit backwards since we SHOULD know what we send but it is a good hack if you inherit the model from a teammate). 


Lets run the python code to check the output

Get the image: 

```
wget https://raw.githubusercontent.com/triton-inference-server/server/main/qa/images/mug.jpg -O "mug.jpg"
```
from the same folder run: 
```
docker run -it --rm --net=host -v $(pwd):/workspace/ nvcr.io/nvidia/tritonserver:22.06-py3-sdk /bin/bash
```

and then start a python shell and execute the python code to check the input to the model to get some real example data. 
```bash
 $ python
 > import numpy as np
 > image_data = np.fromfile("mug.jpg",dtype='uint8')
 > image_data = np.expand_dims(image_data, axis=0)
 > image_data.shape
 > np.savetxt('array.txt',image_data, delimiter=', ', fmt='%.f')
```
Then we can exit the python shell copy the matrix to the input data file instead. Example matrix. The sample input is available [here](./sample_input.json). The shape will be important later on for the input to the perf_analyzer since it is still required even though we have the set input data due to the dynamic input axis.


from the folder where we have the `sample_input.json`
```
docker run -it --rm --net=host -v$(pwd):/workspace/  nvcr.io/nvidia/tritonserver:22.05-py3-sdk /bin/bash
```

THis will still complain on that the 

```bash
perf_analyzer -m ensemble_python_resnet50 -u localhost:8001 -i gRPC --concurrency-range 1:8 --input-data sample_input.json --shape INPUT:1005970
```

output:

```bash
 Successfully read data for 1 stream/streams with 1 step/steps.
*** Measurement Settings ***
  Batch size: 1
  Using "time_windows" mode for stabilization
  Measurement window: 5000 msec
  Latency limit: 0 msec
  Concurrency limit: 8 concurrent requests
  Using synchronous calls for inference
  Stabilizing using average latency

Request concurrency: 1
  Client:
    Request count: 83
    Throughput: 16.6 infer/sec
    Avg latency: 59984 usec (standard deviation 538 usec)
    p50 latency: 59971 usec
    p90 latency: 60651 usec
    p95 latency: 60954 usec
    p99 latency: 61279 usec
    Avg gRPC time: 60155 usec ((un)marshal request/response 137 usec + response wait 60018 usec)
  Server:
    Inference count: 100
    Execution count: 100
    Successful request count: 100
    Avg request latency: 59299 usec (overhead 43 usec + queue 128 usec + compute 59128 usec)

  Composing models:
  postprocessing, version:
      Inference count: 100
      Execution count: 100
      Successful request count: 100
      Avg request latency: 447 usec (overhead 6 usec + queue 57 usec + compute input 41 usec + compute infer 225 usec + compute output 117 usec)

  preprocessing, version:
      Inference count: 100
      Execution count: 100
      Successful request count: 100
      Avg request latency: 56972 usec (overhead 14 usec + queue 20 usec + compute input 102 usec + compute infer 56698 usec + compute output 138 usec)

  resnet50_trt, version:
      Inference count: 100
      Execution count: 100
      Successful request count: 100
      Avg request latency: 1894 usec (overhead 37 usec + queue 51 usec + compute input 250 usec + compute infer 1498 usec + compute output 57 usec)

Request concurrency: 2
  Client:
    Request count: 86
    Throughput: 17.2 infer/sec
    Avg latency: 117130 usec (standard deviation 2486 usec)
    p50 latency: 116705 usec
    p90 latency: 118674 usec
    p95 latency: 120111 usec
    p99 latency: 129282 usec
    Avg gRPC time: 117020 usec ((un)marshal request/response 206 usec + response wait 116814 usec)
  Server:
    Inference count: 102
    Execution count: 102
    Successful request count: 102
    Avg request latency: 115949 usec (overhead 72 usec + queue 55185 usec + compute 60692 usec)

  Composing models:
  postprocessing, version:
      Inference count: 102
      Execution count: 102
      Successful request count: 102
      Avg request latency: 472 usec (overhead 25 usec + queue 26 usec + compute input 62 usec + compute infer 213 usec + compute output 145 usec)

  preprocessing, version:
      Inference count: 103
      Execution count: 103
      Successful request count: 103
      Avg request latency: 113601 usec (overhead 22 usec + queue 55109 usec + compute input 211 usec + compute infer 58121 usec + compute output 138 usec)

  resnet50_trt, version:
      Inference count: 102
      Execution count: 102
      Successful request count: 102
      Avg request latency: 1886 usec (overhead 35 usec + queue 50 usec + compute input 262 usec + compute infer 1518 usec + compute output 20 usec)

Request concurrency: 3
  Client:
    Request count: 86
    Throughput: 17.2 infer/sec
    Avg latency: 174654 usec (standard deviation 1266 usec)
    p50 latency: 174577 usec
    p90 latency: 176394 usec
    p95 latency: 177356 usec
    p99 latency: 178063 usec
    Avg gRPC time: 174725 usec ((un)marshal request/response 164 usec + response wait 174561 usec)
  Server:
    Inference count: 103
    Execution count: 103
    Successful request count: 103
    Avg request latency: 173803 usec (overhead 61 usec + queue 113315 usec + compute 60427 usec)

  Composing models:
  postprocessing, version:
      Inference count: 103
      Execution count: 103
      Successful request count: 103
      Avg request latency: 466 usec (overhead 24 usec + queue 23 usec + compute input 63 usec + compute infer 209 usec + compute output 146 usec)

  preprocessing, version:
      Inference count: 103
      Execution count: 103
      Successful request count: 103
      Avg request latency: 171473 usec (overhead 21 usec + queue 113238 usec + compute input 205 usec + compute infer 57871 usec + compute output 138 usec)

  resnet50_trt, version:
      Inference count: 103
      Execution count: 103
      Successful request count: 103
      Avg request latency: 1880 usec (overhead 32 usec + queue 54 usec + compute input 329 usec + compute infer 1444 usec + compute output 20 usec)

Request concurrency: 4
  Client:
    Request count: 89
    Throughput: 17.8 infer/sec
    Avg latency: 223871 usec (standard deviation 5000 usec)
    p50 latency: 225257 usec
    p90 latency: 228915 usec
    p95 latency: 229557 usec
    p99 latency: 230342 usec
    Avg gRPC time: 224461 usec ((un)marshal request/response 175 usec + response wait 224286 usec)
  Server:
    Inference count: 107
    Execution count: 107
    Successful request count: 107
    Avg request latency: 223490 usec (overhead 53 usec + queue 165130 usec + compute 58307 usec)

  Composing models:
  postprocessing, version:
      Inference count: 107
      Execution count: 107
      Successful request count: 107
      Avg request latency: 485 usec (overhead 18 usec + queue 33 usec + compute input 63 usec + compute infer 225 usec + compute output 145 usec)

  preprocessing, version:
      Inference count: 107
      Execution count: 107
      Successful request count: 107
      Avg request latency: 221138 usec (overhead 26 usec + queue 165035 usec + compute input 197 usec + compute infer 55743 usec + compute output 136 usec)

  resnet50_trt, version:
      Inference count: 107
      Execution count: 107
      Successful request count: 107
      Avg request latency: 1891 usec (overhead 33 usec + queue 62 usec + compute input 283 usec + compute infer 1482 usec + compute output 29 usec)

Request concurrency: 5
  Client:
    Request count: 88
    Throughput: 17.6 infer/sec
    Avg latency: 284336 usec (standard deviation 1334 usec)
    p50 latency: 284323 usec
    p90 latency: 285979 usec
    p95 latency: 286954 usec
    p99 latency: 287328 usec
    Avg gRPC time: 284356 usec ((un)marshal request/response 161 usec + response wait 284195 usec)
  Server:
    Inference count: 106
    Execution count: 106
    Successful request count: 106
    Avg request latency: 283366 usec (overhead 51 usec + queue 224295 usec + compute 59020 usec)

  Composing models:
  postprocessing, version:
      Inference count: 106
      Execution count: 106
      Successful request count: 106
      Avg request latency: 454 usec (overhead 14 usec + queue 20 usec + compute input 61 usec + compute infer 214 usec + compute output 144 usec)

  preprocessing, version:
      Inference count: 106
      Execution count: 106
      Successful request count: 106
      Avg request latency: 281119 usec (overhead 28 usec + queue 224213 usec + compute input 194 usec + compute infer 56545 usec + compute output 138 usec)

  resnet50_trt, version:
      Inference count: 106
      Execution count: 106
      Successful request count: 106
      Avg request latency: 1814 usec (overhead 30 usec + queue 62 usec + compute input 294 usec + compute infer 1408 usec + compute output 19 usec)

Request concurrency: 6
  Client:
    Request count: 88
    Throughput: 17.6 infer/sec
    Avg latency: 341010 usec (standard deviation 1716 usec)
    p50 latency: 340318 usec
    p90 latency: 343443 usec
    p95 latency: 344438 usec
    p99 latency: 345423 usec
    Avg gRPC time: 341357 usec ((un)marshal request/response 164 usec + response wait 341193 usec)
  Server:
    Inference count: 106
    Execution count: 106
    Successful request count: 106
    Avg request latency: 340402 usec (overhead 59 usec + queue 281335 usec + compute 59008 usec)

  Composing models:
  postprocessing, version:
      Inference count: 106
      Execution count: 106
      Successful request count: 106
      Avg request latency: 481 usec (overhead 24 usec + queue 27 usec + compute input 58 usec + compute infer 213 usec + compute output 157 usec)

  preprocessing, version:
      Inference count: 106
      Execution count: 106
      Successful request count: 106
      Avg request latency: 338108 usec (overhead 24 usec + queue 281251 usec + compute input 197 usec + compute infer 56496 usec + compute output 139 usec)

  resnet50_trt, version:
      Inference count: 106
      Execution count: 106
      Successful request count: 106
      Avg request latency: 1834 usec (overhead 32 usec + queue 57 usec + compute input 294 usec + compute infer 1433 usec + compute output 17 usec)

Request concurrency: 7
  Client:
    Request count: 87
    Throughput: 17.4 infer/sec
    Avg latency: 400005 usec (standard deviation 5585 usec)
    p50 latency: 397898 usec
    p90 latency: 411286 usec
    p95 latency: 413368 usec
    p99 latency: 414352 usec
    Avg gRPC time: 399491 usec ((un)marshal request/response 167 usec + response wait 399324 usec)
  Server:
    Inference count: 106
    Execution count: 106
    Successful request count: 106
    Avg request latency: 398493 usec (overhead 50 usec + queue 339226 usec + compute 59217 usec)

  Composing models:
  postprocessing, version:
      Inference count: 106
      Execution count: 106
      Successful request count: 106
      Avg request latency: 461 usec (overhead 17 usec + queue 24 usec + compute input 62 usec + compute infer 211 usec + compute output 146 usec)

  preprocessing, version:
      Inference count: 106
      Execution count: 106
      Successful request count: 106
      Avg request latency: 396202 usec (overhead 22 usec + queue 339150 usec + compute input 198 usec + compute infer 56692 usec + compute output 139 usec)

  resnet50_trt, version:
      Inference count: 106
      Execution count: 106
      Successful request count: 106
      Avg request latency: 1846 usec (overhead 27 usec + queue 52 usec + compute input 264 usec + compute infer 1489 usec + compute output 13 usec)

Request concurrency: 8
  Client:
    Request count: 88
    Throughput: 17.6 infer/sec
    Avg latency: 453741 usec (standard deviation 3410 usec)
    p50 latency: 454444 usec
    p90 latency: 456791 usec
    p95 latency: 458959 usec
    p99 latency: 460669 usec
    Avg gRPC time: 454197 usec ((un)marshal request/response 165 usec + response wait 454032 usec)
  Server:
    Inference count: 106
    Execution count: 106
    Successful request count: 106
    Avg request latency: 453227 usec (overhead 46 usec + queue 394305 usec + compute 58876 usec)

  Composing models:
  postprocessing, version:
      Inference count: 106
      Execution count: 106
      Successful request count: 106
      Avg request latency: 453 usec (overhead 12 usec + queue 28 usec + compute input 55 usec + compute infer 216 usec + compute output 141 usec)

  preprocessing, version:
      Inference count: 106
      Execution count: 106
      Successful request count: 106
      Avg request latency: 450978 usec (overhead 26 usec + queue 394213 usec + compute input 192 usec + compute infer 56406 usec + compute output 141 usec)

  resnet50_trt, version:
      Inference count: 106
      Execution count: 106
      Successful request count: 106
      Avg request latency: 1817 usec (overhead 29 usec + queue 64 usec + compute input 286 usec + compute infer 1419 usec + compute output 17 usec)

Inferences/Second vs. Client Average Batch Latency
Concurrency: 1, throughput: 16.6 infer/sec, latency 59984 usec
Concurrency: 2, throughput: 17.2 infer/sec, latency 117130 usec
Concurrency: 3, throughput: 17.2 infer/sec, latency 174654 usec
Concurrency: 4, throughput: 17.8 infer/sec, latency 223871 usec
Concurrency: 5, throughput: 17.6 infer/sec, latency 284336 usec
Concurrency: 6, throughput: 17.6 infer/sec, latency 341010 usec
Concurrency: 7, throughput: 17.4 infer/sec, latency 400005 usec
Concurrency: 8, throughput: 17.6 infer/sec, latency 453741 usec
```

If we check here we can see that the preprocessing takes the majority of the processing time. 

Suggestions: 
- Lets create more instances of it.

Questions: 
- How do the backend work do one batch need to stay intact between all the processing steps in this case or can it be split up?
- Can we parallize the processing in some of the steps. 

As we can see we don't manage to increase the throughput even though we increas the concurrency. The bottleneck looks to be the preprocessing. In this we can try to increase the nbr of preprocessing instances, more tips on optimizations can be found [here](https://github.com/triton-inference-server/server/blob/main/docs/optimization.md)

After updating the `instance_group` part of the config should look like this: 

```ptxt
instance_group [
{
    kind: KIND_CPU
    count: 4
}
]
```


the following results where obtained: 

```
...
Inferences/Second vs. Client Average Batch Latency
Concurrency: 1, throughput: 13.8 infer/sec, latency 73402 usec
Concurrency: 2, throughput: 23.6 infer/sec, latency 84774 usec
Concurrency: 3, throughput: 24.6 infer/sec, latency 122203 usec
Concurrency: 4, throughput: 24.4 infer/sec, latency 164786 usec
Concurrency: 5, throughput: 24.2 infer/sec, latency 206315 usec
Concurrency: 6, throughput: 24.8 infer/sec, latency 240450 usec
Concurrency: 7, throughput: 24.4 infer/sec, latency 287432 usec
Concurrency: 8, throughput: 23.8 infer/sec, latency 336209 usec
```
as suspected increasing the concurrency increase the throughput when concurrency > 1. For single execution the throughput went down from 16.6 to 13.8 infer/sec, after rerunning one more time it increase slightly. It is believe that the difference is due to cold start of the serving.

Increasing to 4 gives: 

```
Inferences/Second vs. Client Average Batch Latency
Concurrency: 1, throughput: 15.2 infer/sec, latency 65986 usec
Concurrency: 2, throughput: 24.2 infer/sec, latency 82423 usec
Concurrency: 3, throughput: 31.4 infer/sec, latency 95638 usec
Concurrency: 4, throughput: 36.6 infer/sec, latency 109091 usec
Concurrency: 5, throughput: 37.2 infer/sec, latency 134947 usec
Concurrency: 6, throughput: 35.2 infer/sec, latency 168920 usec
Concurrency: 7, throughput: 37.2 infer/sec, latency 187942 usec
Concurrency: 8, throughput: 37.6 infer/sec, latency 214759 usec
```