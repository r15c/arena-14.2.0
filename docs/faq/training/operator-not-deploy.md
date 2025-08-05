# FAQ: No matches for kind "TFJob" version "kubeflow.org/v1alpha2"

#### Problem

The error is described as below:

``` 
error: unable to recognize "/tmp/tf-dist-git.yaml392889812": no matches for kind "TFJob" in version "kubeflow.org/v1alpha2"
```

## Solution

```
git clone https://github.com/kubeflow/arena.git
kubectl delete -f kubernetes-artifacts/tf-operator/tf-operator.yaml
kubectl create -f kubernetes-artifacts/tf-operator/tf-operator.yaml
``` 