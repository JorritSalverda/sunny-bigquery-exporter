## Installation

To install this application using Helm run the following commands: 

```bash
helm repo add jorritsalverda https://helm.jorritsalverda.com
kubectl create namespace sunny-bigquery-exporter

helm upgrade \
  sunny-bigquery-exporter \
  jorritsalverda/sunny-bigquery-exporter \
  --install \
  --namespace sunny-bigquery-exporter \
  --set secret.gcpServiceAccountKeyfile='{abc: blabla}' \
  --wait
```
