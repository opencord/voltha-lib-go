# Helm Chart configuration

## Prometheus based stats server

Below is the configuration for running exposing the statistics on port 8081.
A running prometheus on kubernetes will grab the statistics from this service looking at the annotations.

### deployment.yaml

``` yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test
spec:
  template:
    metadata:
      {{- with .Values.deployment.annotations }}
      annotations: {{- toYaml . | nindent 8 }}
      {{- end }}
    spec:
      containers:
        - name: testContainer
          ports:
          - name: promserver
            containerPort: 8081
```

### values.yaml

``` yaml
deployment:
  annotations:
    prometheus.io/path: "/metrics"
    prometheus.io/port: "8081"
    prometheus.io/scrape: "true"
```

### service.yaml

``` yaml
apiVersion: v1
kind: Service
metadata:
  name: test
spec:
  type: ClusterIP
  ports:
    - port: 8081
      targetPort: promserver
      protocol: TCP
      name: promserver
```
