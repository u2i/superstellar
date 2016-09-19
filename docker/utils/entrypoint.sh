#/bin/sh
go get
go build
go install

/go/bin/superstellar_utils backend 100 50ms
