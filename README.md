# superstellar
Massive multiplayer galactic game written in Golang

## Installation & running
1. Clone this repository to your `$GOPATH/src` directory
1. Run `go get`
1. Run `go build && go install`
1. Run `$GOPATH/bin/superstellar`
1. Go to [localhost:8080](http://localhost:8080/)

## Compiling protobufs

### Golang

1. Go to superstellar src directory.
1. `brew install protobuf`
1. `go get -u github.com/golang/protobuf/{proto,protoc-gen-go}`
1. `protoc -I=protobuf --go_out=backend/proto protobuf/superstellar.proto`

### JavaScript

1. Clone https://github.com/dcodeIO/ProtoBuf.js/
1. `npm install`
1. Go to superstellar src directory.
1. `node [path to Protobuf.js repo]/bin/pbjs -s proto -t json protobuf/superstellar.proto webroot/js/superstellar_proto.json`
