FROM golang:1.24.0-alpine AS builder
WORKDIR /build
ARG ENV=dev
COPY . .
RUN cat etc/cfg/conf.${ENV}.json > etc/cfg/conf.json && go build -o main src/cmd/main.go

FROM alpine:latest AS runner
WORKDIR /app
COPY --from=builder /build/etc etc/
COPY --from=builder /build/docs docs/
COPY --from=builder /build/main main
CMD ["./main", "-mode", "rest"]
