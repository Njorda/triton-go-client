# Performance tuning


In this blog post we will see how we can se the performance analyzer tool for an ensemble model. The goal is to understand the models performance in find which part of the ensemble is the bottle neck for the throughput and see if we can find any potential solutions to move the bottleneck. Since performance improvements is a iterative process and we are more interested in showing the workflow in this blog post we will stop the performance hunt after the first iteration and instead leave it to the reader to continue the hunt. 

This will required that it is run on a linux machine since `--net=host` [is not supported using mac](https://docs.docker.com/network/host/)

> "The host networking driver only works on Linux hosts, and is not supported on Docker for Mac, Docker for Windows, or Docker EE for Windows Server."

```
docker run --runtime=nvidia -it --shm-size=1gb --rm --net=host -v$(pwd):/workspace/ -v/$(pwd)/model_repository:/models nvcr.io/nvidia/tritonserver:22.06-py3 bash
```

start the container with the performance analyzer aviliable. 

```
docker run -it --rm --net=host  nvcr.io/nvidia/tritonserver:22.05-py3-sdk /bin/bash
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
perf_analyzer -m ensemble_python_resnet50 -u localhost:8001 -i gRPC --concurrency-range 1:8 --shape INPUT:1
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
 > np.savetxt('array.txt',image_data, delimiter=', ', fmt='%.f')
```
Then we can exit the python shell copy the matrix to the input data file instead. Example matrix. The sample input is available [here](./sample_input.json)





