FROM golang:1.17

WORKDIR /go/src/app
COPY *.go go.sum go.mod .
RUN go install -v
EXPOSE 9090

CMD ["k8s-demo-emp-api"]
