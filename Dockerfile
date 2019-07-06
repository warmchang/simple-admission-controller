# build stage
FROM golang:1.10-stretch AS build-env
RUN mkdir -p /go/src/github.com/ChrisTheShark/simple-admission-controller
WORKDIR /go/src/github.com/ChrisTheShark/simple-admission-controller
COPY  . .
RUN useradd -u 10001 webhook
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o simplewebhook

FROM scratch
COPY --from=build-env /go/src/github.com/ChrisTheShark/simple-admission-controller .
COPY --from=build-env /etc/passwd /etc/passwd
USER webhook
ENTRYPOINT ["/simplewebhook"]