#!/bin/sh

protoc -I=protobuf --go_out=backend/pb protobuf/superstellar.proto 
pbjs protobuf/superstellar.proto > webroot/js/superstellar_proto.json
