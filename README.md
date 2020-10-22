# GoProxy
## Socks5 Proxy Server Implemented in GoLang

Core Mechanism Credits: https://github.com/oov/socks5

### <u>Usage Menu</u>
```
-h    Print this help
-host string
    Proxy Server Bind Host (default "0.0.0.0")
-kafka
    Enable/Disable kafka streaming (default true)
-kafka-server string
    Kafka Bootstrap Server IP:PORT (default "localhost:9092")
-kafka-topic string
    Kafka Streaming Topic (default "proxyMonitor")
-port string
    Proxy Server Bind Port (default "1080")
```

* This application blocks all domains listed from [StevenBlack Blacklist](https://raw.githubusercontent.com/StevenBlack/hosts/master/alternates/fakenews-gambling-porn/hosts)
* The blacklist cache is auto updated every 24 hours
* Configure SOCKS5 proxy in browser or whole device
* You can also use [Monitoring Application](https://github.com/aravinth2094/GoProxyMonitor) to obtain the PacFile via ```http://localhost/api/pac```

### See [Monitoring Application](https://github.com/aravinth2094/GoProxyMonitor) that consumes traffic from Apache Kafka and flushes to InfluxDB