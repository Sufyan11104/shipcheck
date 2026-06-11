FROM golang:1.22-alpine AS build

WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
RUN go build -o /tmp/good-service ./cmd/good-service

FROM alpine:3.20

RUN addgroup -S app && adduser -S -G app app
WORKDIR /app
COPY --from=build /tmp/good-service ./good-service

ENV APP_PORT=8080
HEALTHCHECK --interval=30s --timeout=3s CMD wget -qO- http://127.0.0.1:8080/health || exit 1

USER app
EXPOSE 8080
CMD ["./good-service"]
