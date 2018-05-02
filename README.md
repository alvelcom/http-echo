http-echo
=========

http-echo is a simple service that sends client and service IP back in
response. It's useful for testing kubernetes services.

```
$ curl http://127.0.0.1:8080/
{
    "local_ips": [
        "10.10.10.64",
        "192.168.255.70"
    ],
    "remote_addr": "127.0.0.1:57917",
    "server_hostname": "alvo-mbp"
}
```

Here is kubernetes manifest snippet for your convenience:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: http-echo
  labels:
    app: http-echo
spec:
  replicas: 1
  selector:
    matchLabels:
      app: http-echo
  template:
    metadata:
      labels:
        app: http-echo
    spec:
      containers:
      - name: http-echo
        image: alvelcom/http-echo
        ports:
        - containerPort: 8080
----
kind: Service
apiVersion: v1
metadata:
  name: http-echo
spec:
  type: ClusterIP
  selector:
    app: http-echo
  ports:
  - protocol: TCP
    port: 8080
    targetPort: 8080
```
