# kube-proxy-websocket-testcase

This is a reproduce of the issue in Kubernetes proxy related to query string is not being forwarded in case of WebSocket connections.

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

Start the proxy:

```
go run main.go
```

Do requests:

```
curl http://localhost:8080/example/ws
curl http://localhost:8080/example/ws\?q\=x
curl http://localhost:8080/example/ws\?q\=x -H "Connection: Upgrade"
curl http://localhost:8080/example/ws\?q\=x -H "Upgrade: websocket" -H "Sec-WebSocket-Version: 13" -H "Sec-WebSocket-Key: x3JJHMbDL1EzLkh9GBhXDw=="
```

You should see these logs:

```
query: map[]
Failed to set websocket upgrade: %+v websocket: the client is not using the websocket protocol: 'upgrade' token not found in 'Connection' header
[GIN] 2021/09/19 - 14:28:42 | 400 |     237.402µs |       127.0.0.1 | GET      "/ws"
query: map[q:[x]]
Failed to set websocket upgrade: %+v websocket: the client is not using the websocket protocol: 'upgrade' token not found in 'Connection' header
[GIN] 2021/09/19 - 14:28:46 | 400 |     147.181µs |       127.0.0.1 | GET      "/ws?q=x"
query: map[]
Failed to set websocket upgrade: %+v websocket: the client is not using the websocket protocol: 'websocket' token not found in 'Upgrade' header
[GIN] 2021/09/19 - 14:28:50 | 400 |     470.374µs |       127.0.0.1 | GET      "/ws"
query: map[q:[x]]
Failed to set websocket upgrade: %+v websocket: the client is not using the websocket protocol: 'upgrade' token not found in 'Connection' header
[GIN] 2021/09/19 - 14:28:56 | 400 |     297.282µs |       127.0.0.1 | GET      "/ws?q=x"
```

In conclusion, `Connection: Upgrade` header is causing the query string to disappear.
