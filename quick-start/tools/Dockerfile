#Build Delve
FROM golang:1.16.4 AS build-dev
RUN git clone https://github.com/go-delve/delve.git	&& \
    cd ./delve/cmd/dlv && \
    go build -ldflags '-extldflags "-static"' && \
    go install

FROM registry.erda.cloud/erda/dice/erda:1.1-20210721-725137a
COPY --from=build-dev /go/bin/dlv /app