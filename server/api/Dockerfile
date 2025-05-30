ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /usr/src/app

RUN apt-get update && apt-get install -y ca-certificates 

COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .

COPY private.pem ./private.pem
COPY public.pem ./public.pem

RUN go build -v -o /run-app .


FROM debian:bookworm
RUN apt-get update && apt-get install -y ca-certificates && apt-get clean
WORKDIR /app
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /run-app /usr/local/bin/
COPY --from=builder /usr/src/app/private.pem ./private.pem
COPY --from=builder /usr/src/app/public.pem ./public.pem
CMD ["run-app"]
