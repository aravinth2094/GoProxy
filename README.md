# GoProxy Proxy Server

GoProxy is a lightweight proxy server written in Go that acts as a firewall and provides VPN-like functionality. It allows you to route your network traffic through a secure tunnel for enhanced privacy and security. This README file provides instructions on how to set up and use the GoProxy proxy server, along with information about the authentication and access rules.

## Installation

To install and use GoProxy, you need to have Go installed on your system. If you don't have Go installed, you can download it from the official website: https://golang.org/

1. Clone the GoProxy repository:
   ```
   git clone https://github.com/aravinth2094/GoProxy.git
   ```

2. Change to the GoProxy directory:
   ```
   cd goproxy
   ```

3. Build the GoProxy binary:
   ```
   make
   ```

## Usage

GoProxy provides several modes of operation: `proxy`, `tunnel server`, and `tunnel client`. Follow the instructions below to start the proxy server, set up the TLS tunnel, and configure authentication and access rules.

### Start the Proxy Server

To start the GoProxy proxy server, use the following command:
```
./goproxy proxy --listen :1080
```
This will start the proxy server on port 1080, and it will be ready to accept incoming connections.

### Set Up TLS Tunnel Server

To set up a TLS tunnel server that accepts proxy connections, use the following command:
```
./goproxy tunnel server --listen :1081 --target :1080
```
This will start the TLS tunnel server on port 1081. It will listen for incoming connections from clients and relay them to the proxy server running on port 1080.

### Set Up TLS Tunnel Client

To set up a TLS tunnel client that accepts your proxy traffic and relays it to the TLS tunnel server, use the following command:
```
./goproxy tunnel client --listen :1082 --target :1081
```
This will start the TLS tunnel client on port 1082. It will accept your proxy traffic and forward it to the TLS tunnel server running on port 1081.

### Proxy Traffic Flow

The traffic flow with the configured setup is as follows:

```
System -> TLS Tunnel client -> TLS Tunnel server -> proxy -> out
```

Your system's network traffic will be directed to the TLS tunnel client, which will relay it to the TLS tunnel server. The TLS tunnel server will then forward the traffic to the GoProxy proxy server, which will process the requests and send the outbound traffic.

### Authentication and Access Rules

To configure authentication and access rules for the proxy server, you can modify the `access.json` and `networks.json` files.

**access.json**
```json
[
    {
        "username": "john",
        "password": "$2y$10$abc123def456",
        "access": [
            {
                "network": "home",
                "ports": ["80", "443"]
            },
            {
                "network": "work",
                "ports": ["80", "8080", "443"]
            }
        ]
    },
    {
        ...
    }
]
```

The `access.json` file contains an array of objects representing user access configurations. Each object has a `username`, `password`, and an `access` array, which specifies the allowed networks and ports for that user.

**networks.json**
```json
[
    {
        "name": "home",
        "nodes": [
            {
                "name": "local",
                "ip": "192.168.0.100",
                "ports": ["80", "443", "1080"]
            },
            {
                "name": "router",
                "ip": "192.168.0.1",
                "ports": ["80", "443"]
            }
        ]
    },
    {
        "name": "work",
        "nodes": [
            {
                "name": "office",
                "ip": "10.0.0.100",
                "ports": ["80", "8080", "443"]
            },
            {
                "name": "gateway",
                "ip": "10.0.0.1",
                "ports": ["80", "443"]
            }
        ]
    }
]
```

The `networks.json` file defines the networks and their associated nodes. Each network has a `name` and an array of `nodes`. Each node object has a `name`, `ip` address, and an array of allowed `ports`.

When a user authenticates with the proxy server, their credentials are matched against the `access.json` file. If a match is found, the server checks if the user's allowed networks and ports match the requested destination. If the criteria match, the proxy server will forward the traffic accordingly.

## Conclusion

Congratulations! You have successfully set up and configured the GoProxy proxy server. You can now start using it to route your network traffic securely, while applying authentication and access restrictions based on the rules defined in the `access.json` and `networks.json` files.