# Otel Keda Example Application

This repository contains an example application demonstrating the use of OpenTelemetry with KEDA for Kubernetes-based autoscaling.

## Dependencies

- [Helm](https://helm.sh/docs/intro/install/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [kind](https://kind.sigs.k8s.io/docs/user/quick-start/)
- [grpcurl](https://github.com/fullstorydev/grpcurl)
- [Nix](https://nixos.org/download.html) (Optional)
- [direnv](https://direnv.net/docs/installation.html) (Optional)

## Setup

### Environment Setup with Nix and direnv (Optional)
This repository makes use of [Nix](https://nixos.org/) and [direnv](https://direnv.net/) to manage the development environment and dependencies. Follow these steps to set them up:

1. Install Nix by following the instructions [here](https://nixos.org/download.html).
2. Install direnv by following the instructions [here](https://direnv.net/docs/installation.html).
3. Navigate to the project directory and allow direnv to load the environment:

  ```sh
  direnv allow
  ```

### Add Helm Repositories

```sh
helm repo add grafana https://grafana.github.io/helm-charts
helm repo add kedacore https://kedacore.github.io/charts
helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
helm repo update
```

### Create Kubernetes Cluster

```sh
kind create cluster --config=kind.yaml --name otel-keda-example
```

### Install supporting Helm Charts

```sh
helm upgrade grafana grafana/grafana --namespace grafana --create-namespace --install --values grafana-values.yaml
helm upgrade mimir grafana/mimir-distributed --namespace mimir --create-namespace --install --values mimir-values.yaml
helm upgrade keda kedacore/keda --namespace keda --create-namespace --install --values keda-values.yaml
helm upgrade opentelemetry-collector-daemonset open-telemetry/opentelemetry-collector --namespace opentelemetry-collector --create-namespace --install --values opentelemetry-collector-daemonset-values.yaml
helm upgrade opentelemetry-collector-deployment open-telemetry/opentelemetry-collector --namespace opentelemetry-collector --create-namespace --install --values opentelemetry-collector-deployment-values.yaml
```

### Install sample-app Helm Chart

From local repository
```sh
helm upgrade sample-app sample-app/helm --namespace sample-app --create-namespace --install --values sample-app-values.yaml
```

From GHCR
```sh
helm upgrade sample-app oci://ghcr.io/joaoestrela/otel-keda-example/helm-charts/sample-app --version 0.1.1 --namespace sample-app --create-namespace --install --values sample-app-values.yaml
```

### Create Keda Scaled Object
```sh
kubectl apply -f sampleAppScaledObject.yaml
```

### Retrieve Grafana Admin Password

```sh
kubectl get secret --namespace grafana grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo
```

### Port Forward Grafana

```sh
kubectl --namespace grafana port-forward svc/grafana 3000:9090
```

### Port Forward Sample Application Pods

```sh
POD_N=0
kubectl --namespace sample-app get pods -o jsonpath="{.items[$POD_N].metadata.name}" | xargs -I {} kubectl --namespace sample-app port-forward {} $((8080 + POD_N)):8080
```

### Interact with the Sample Application

```sh
POD_N=0
grpcurl -plaintext -d '{"value": 5}' localhost:$((8080 + POD_N)) counter.CounterService/IncreaseCounter
grpcurl -plaintext -d '{"value": 3}' localhost:$((8080 + POD_N)) counter.CounterService/DecreaseCounter
```

### Count Sample Application Pods

```sh
kubectl get pods --namespace sample-app --no-headers | wc -l
```

## Project Structure

- [`grafana-values.yaml`](grafana-values.yaml): Configuration for Grafana.
- [`keda-values.yaml`](keda-values.yaml): Configuration for KEDA.
- [`mimir-values.yaml`](mimir-values.yaml): Configuration for Mimir.
- [`opentelemetry-collector-daemonset-values.yaml`](opentelemetry-collector-daemonset-values.yaml): Configuration for OpenTelemetry Collector DaemonSet.
- [`opentelemetry-collector-deployment-values.yaml`](opentelemetry-collector-deployment-values.yaml): Configuration for OpenTelemetry Collector Deployment.
- [`sample-app/`](sample-app/): Sample application.
- [`sample-app/helm/`](sample-app/helm/): Helm chart for deploying the sample application.
- [`sample-app-values.yaml`](sample-app-values.yaml): Configuration for the sample application.

