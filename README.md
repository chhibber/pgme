# Prometheus GPU Metrics Exporter (PGME)

PGME is a GPU Metrics exporters that leverages the nvidai-smi binary.  The initial work and key metric gathering code is
derived from:

 * https://github.com/zhebrak/nvidia_smi_exporter

Nvidia-smi command used to gather metrics:
```
nvidia-smi --query-gpu=name,index,temperature.gpu,utilization.gpu,utilization.memory,memory.total,memory.free,memory.used --format=csv,noheader,nounits
```

I have added the  following in an attempt to make it a more robust service:
* configuration via environment variables
* Makefile for local build
* liveness HTTP request probe for Kubernetes(k8s)
* graceful shutdown of http server
* exporter details at http://[[ip of server]]:[[port]/
* Integration with AWS Codebuild and Publishing to DockerHub or AWS ECR via different buildspec files

Working On:
* Kubernetes service and helm configuration


## Building

Local MAC Build (Generates a binary that works on OSX based systems)
```
git clone https://github.com/chhibber/pgme.git
cd pgme
make build-mac
```

Local Linux Build (Genrates a binary that works on Linux systems)
```
https://github.com/chhibber/pgme.git
cd pgme
make build
```

Local Docker Build (Generates a docker image)
```
https://github.com/chhibber/pgme.git
cd pgme
make docker-build IMAGE_REPO_NAME=[[ repo_name/app_name ]] IMAGE_TAG=[[ version info ]]

# Example run
nvidia-docker run -p 9101:9101 chhibber/pgme
2018/01/05 21:32:31 Starting the service...
2018/01/05 21:32:31 - PORT set to 9101.  If  environment variable PORT is not set the default is 9101
2018/01/05 21:32:31 The service is listening on 9101
...
```

## Running the binary directly
* The default port is 9101

You can change the port by defining the environment variabl PORT in front of the binary.
```
> PORT=9101 ./pgme
```

#### Runnign via Docker (Needed to expose the GPU to the running container)
nvidia-docker run -p 9101:9101 chhibber/pgme:2017.01


### Available Metrics - http://localhost:9101/metrics
```
temperature_gpu{gpu="TITAN X (Pascal)[0]"} 41
utilization_gpu{gpu="TITAN X (Pascal)[0]"} 0
utilization_memory{gpu="TITAN X (Pascal)[0]"} 0
memory_total{gpu="TITAN X (Pascal)[0]"} 12189
memory_free{gpu="TITAN X (Pascal)[0]"} 12189
memory_used{gpu="TITAN X (Pascal)[0]"} 0
temperature_gpu{gpu="TITAN X (Pascal)[1]"} 78
utilization_gpu{gpu="TITAN X (Pascal)[1]"} 95
utilization_memory{gpu="TITAN X (Pascal)[1]"} 59
memory_total{gpu="TITAN X (Pascal)[1]"} 12189
memory_free{gpu="TITAN X (Pascal)[1]"} 1738
memory_used{gpu="TITAN X (Pascal)[1]"} 10451
temperature_gpu{gpu="TITAN X (Pascal)[2]"} 83
utilization_gpu{gpu="TITAN X (Pascal)[2]"} 99
utilization_memory{gpu="TITAN X (Pascal)[2]"} 82
memory_total{gpu="TITAN X (Pascal)[2]"} 12189
memory_free{gpu="TITAN X (Pascal)[2]"} 190
memory_used{gpu="TITAN X (Pascal)[2]"} 11999
temperature_gpu{gpu="TITAN X (Pascal)[3]"} 84
utilization_gpu{gpu="TITAN X (Pascal)[3]"} 97
utilization_memory{gpu="TITAN X (Pascal)[3]"} 76
memory_total{gpu="TITAN X (Pascal)[3]"} 12189
memory_free{gpu="TITAN X (Pascal)[3]"} 536
memory_used{gpu="TITAN X (Pascal)[3]"} 11653
```

### Prometheus example config

```
- job_name: "gpu_exporter"
  static_configs:
  - targets: ['localhost:9101']
```

