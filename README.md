# Prometheus GPU Metrics Exporter (PGME)

PGME is a GPU Metrics exporters that leverages the nvidai-smi binary.  The primary work for metrics gathering is derived
from:
 * https://github.com/zhebrak/nvidia_smi_exporter


I have added the  following in an attempt to make it a more robust service:
* configuration via environment variables
* Makefile for local build
* liveness HTTP request probe for Kubernetes(k8s)
* graceful shutdown of http server
* exporter details at http://[[ip of server]]:[[port]/
* Integration with AWS Codebuild and Publishing to DockerHub or AWS ECR via different buildspec files
* Kubernetes service and helm configuration

## Build

Local Build
```
go build -v nvidia_smi_exporter
```
Multistage Docker Build
```
make docker-build IMAGE_REPO_NAME=[[ repo_name/app_name ]] IMAGE_TAG=[[ version info ]]
```

## Run
* Default port is 9101

#### Local
```
> PORT=9101 ./nvidia_smi_exporter
```

#### Docker
nvidia-docker run -p 9101:9101 chhibber/pgme:2017.01

### localhost:9101/metrics
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

### Exact command
```
nvidia-smi --query-gpu=name,index,temperature.gpu,utilization.gpu,utilization.memory,memory.total,memory.free,memory.used --format=csv,noheader,nounits
```

### Prometheus example config

```
- job_name: "gpu_exporter"
  static_configs:
  - targets: ['localhost:9101']
```

