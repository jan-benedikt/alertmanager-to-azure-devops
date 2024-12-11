FROM golang:1.23-alpine AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -ldflags="-w -s" -o /usr/local/bin/awp ./cmd


FROM alpine:3.18

COPY --from=builder /usr/local/bin/awp /usr/local/bin/awp

RUN addgroup -g 10001 vigilantes && adduser -D -G vigilantes -u 10001 batman
RUN chown -R batman:vigilantes /usr/local/bin/awp
RUN chmod +x /usr/local/bin/awp

USER batman:vigilantes

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/awp"]
