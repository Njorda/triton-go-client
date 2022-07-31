# Triton model analyzer

In this blog post we will dig down in to the [model analyzer](https://github.com/triton-inference-server/model_analyzer) tool for optimizing models. Triton ships two tools for performance analysis: 

- [perf_analyzer](https://github.com/triton-inference-server/server/blob/main/docs/perf_analyzer.md) - cli for performance measurements
- [model_analyzer](https://github.com/triton-inference-server/model_analyzer) - model performance optimization and analysis(utilizing perf_analyzer internally sweeping over model config file parameters)

In this tutorial we will build upon the model developer in the [triton-ensemble](https://github.com/Njorda/triton-ensemble).

The model analysis will test different configurations and bring up the different models and run the tests on it, thus it will need to have `perf_analyzer` and `tritonserver` cli to run. 

```
docker run --runtime=nvidia -it --shm-size=1gb --rm -p8000:8000 -p8001:8001 -p8002:8002 -v$(pwd):/workspace/ -v/$(pwd)/model_repository:/models nvcr.io/nvidia/tritonserver:22.06-py3 bash
pip install numpy pillow torchvision triton-model-analyzer
```

```
model-analyzer profile -m /models --profile-models ensemble_python_resnet50
```

Sadly `model_analyzer` don't support optimization of ensembles, check [issue 400](https://github.com/triton-inference-server/model_analyzer/issues/400) and [issue 162](https://github.com/triton-inference-server/model_analyzer/issues/162) for updates on this topic. 

Instead we have to break the problem down and optimize each part of the ensemble individually in order to optimize the ensemble. We will create a config with the parameters for the optimization process, more info about available settings can be found [here](https://github.com/triton-inference-server/model_analyzer/blob/main/docs/config.md#example-1). 

```
model_repository: /path/to/model/repository/
profile_models:
  model_1:
    parameters:
      batch_sizes: 4
perf_analyzer_flags:
  shape:
    - INPUT0:1024,1024
    - INPUT1:1024,1024
```

```
model-analyzer profile -m /models --profile-models resnet50_trt
```


