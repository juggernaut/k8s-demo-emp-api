FROM golang:1.17 as build

WORKDIR /go/src/app
COPY go.sum go.mod .
RUN go mod download
COPY *.go .
RUN CGO_ENABLED=1 GOOS=linux go build -o k8s-demo-emp-api .

FROM gcr.io/distroless/base-debian11
WORKDIR /root/
COPY --from=build /go/src/app/k8s-demo-emp-api ./

EXPOSE 9090

CMD ["./k8s-demo-emp-api"]
