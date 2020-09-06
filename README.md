# ouo-scheduler
Custom kubernetes scheduler demo using kube-scheduler framework plugin.

Implement *QueueSort Plugin* and *PreFilter Plugin*.

# Usage

1. Build Docker image
```
make image
```

2. Change `spec.image` inside `deploy/deployment.yaml`
3. `kubectl apply -f deploy/.`
