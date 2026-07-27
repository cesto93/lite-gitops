FROM golang:1.26-alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o lite-gitops .

FROM alpine:3.21

RUN apk add --no-cache docker-cli docker-compose git openssh
COPY --from=builder /build/lite-gitops /usr/local/bin/lite-gitops

ENTRYPOINT ["lite-gitops"]
