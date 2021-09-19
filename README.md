# kube-proxy-websocket-testcase

This is an attempt of reproducing the issue in Kubernetes proxy related to query string is not being forwarded in case of WebSocket connections.

Build and publish the Docker image

```
cd server/ && docker build . -t mertyildiran/websocket-testcase:latest && docker push mertyildiran/websocket-testcase:latest && cd ..
```

Apply the Kubernetes manifests:

```
kubectl apply -f server/kubernetes
```

Watch the logs:

```
kubectl logs --follow websocket-testcase-55c775dcc5-ff9rp -n websocket-testcase-ns
```

Do requests:

```
curl http://localhost:8080/example/ws
curl http://localhost:8080/example/ws\?q\=x
```

See if it's printing:

```
query: map[]
Failed to set websocket upgrade: %+v websocket: the client is not using the websocket protocol: 'upgrade' token not found in 'Connection' header
[GIN] 2021/09/19 - 14:03:25 | 400 |      91.111µs |       127.0.0.1 | GET      "/ws"
query: map[q:[x]]
Failed to set websocket upgrade: %+v websocket: the client is not using the websocket protocol: 'upgrade' token not found in 'Connection' header
[GIN] 2021/09/19 - 14:03:31 | 400 |       44.21µs |       127.0.0.1 | GET      "/ws?q=x"
```
