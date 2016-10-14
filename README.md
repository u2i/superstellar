[![Build Status](http://jenkins.talkwit.tv/buildStatus/icon?job=u2i/superstellar/master)](http://jenkins.talkwit.tv/job/u2i/job/superstellar/job/master/)

# superstellar
Massive multiplayer galactic game written in Golang

## Installation & running
1. Clone this repository to your `$GOPATH/src` directory
1. `cd` to that directory
2. Run `go get`
3. Run `go build && go install`
4. Run `$GOPATH/bin/superstellar`
5. `cd webroot`
6. `npm install`
7. `npm run dev`
8. Go to [http://localhost:8090](http://localhost:8090)

## Running stress test util
You can run a stress test util that spawns any number of clients which connect to the server and send ramdomly correct user input messages.

1. `cd superstellar_utils`
1. `go build && go install`
1. Run `$GOPATH/bin/superstellar_utils 127.0.0.1 100 50ms` for spawning 100 clients, with 50 ms interval.

## Live profiling 
It's possible to dump various information from the running server, e.g. stacktraces of all goroutines which might be useful in case of a deadlock. 

1. Run server
1. Go to [http://localhost:8080/debug/pprof/](http://localhost:8080/debug/pprof/)

## Using JS `__DEBUG__` flag

If you run `DEBUG=true npm run dev` you will see additional debugging
informations. You can add your own debugging info in code. Just detect that
we're in the debug mode:

```javascript
if (__DEBUG__) {
   console.log("I'm in debug mode!");
}
```

## Compiling protobufs

### Golang

1. Go to superstellar src directory.
1. `brew install protobuf`
1. `go get -u github.com/golang/protobuf/{proto,protoc-gen-go}`
1. `protoc -I=protobuf --go_out=backend/pb protobuf/superstellar.proto` (you
   need to have $GOPATH/bin in your $PATH so `protoc-gen-go` can be found)

### JavaScript

1. `npm install -g protobufjs`
1. Go to superstellar src directory.
1. `pbjs protobuf/superstellar.proto > webroot/js/superstellar_proto.json`
