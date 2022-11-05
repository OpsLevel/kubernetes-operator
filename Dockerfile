FROM golang:1.18-alpine as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download && go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.9.2

COPY ./ ./
RUN controller-gen object paths=./... &&\
    controller-gen crd paths=./... output:crd:artifacts:config=crds
RUN go build -o /app/opslevel && chmod +x /app/opslevel

FROM alpine:latest
RUN apk --no-cache add ca-certificates
CMD ["/usr/local/bin/opslevel"]

COPY --from=builder /app/crds /crds
COPY --from=builder /app/opslevel /usr/local/bin
