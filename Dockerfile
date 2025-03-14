FROM golang:1.24 AS builder

WORKDIR /app

COPY src/go.mod .
COPY src/go.sum .

RUN go mod download

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY src/ .

RUN CGO_ENABLED=0 go build -ldflags '-w -s' -o /app/main .

FROM alpine:3.12

WORKDIR /app

RUN mkdir /app/config

COPY src/config/config.yaml /app/config/config.yaml

COPY --from=builder /app/main /app/main

RUN chmod +x /app/main

CMD [ "/app/main" ]