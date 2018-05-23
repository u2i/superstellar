ARG VERSION=latest
FROM superstellar-backend-builder:$VERSION
WORKDIR /go/src/superstellar
RUN go get github.com/onsi/ginkgo github.com/onsi/gomega
CMD go test superstellar/...
