# P2P messaging

Console chat with auto reconnection and TLS encryption written in go

Based on https://github.com/faroyam/golang-p2p-chat

### Installing

```
go get github.com/faroyam/golang-p2p-chat
```

## Getting Started

1. Generate SSL keys for TSL 
```
bash makeCerts.sh
```
2. Run P2P Messaging
```
go build p2pmsg.go && ./p2pmsg <username> <port to accept traffic> <remote IP:Port> 
```
