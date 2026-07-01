# Tunnel-Manager

This is a tool which allows us to create a new connection to Teleport using one of the following ways:

1. ALPN - Makes a TCP TLS connection to the proxy server an sends over a ALPN of our choosing
2. WebSocket - Makes a HTTPS connection to the proxy server and sets up a websocket for funneling over our data

Once the connection is setup, we bind a local TCP server or Unix socket to allow other tools to pass data over.

This gives us SOCAT like functionality, where we can chain multiple layers.

## Building

```
make build
```