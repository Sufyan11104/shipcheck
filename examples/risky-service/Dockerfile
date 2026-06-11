FROM golang:1.22-alpine

WORKDIR /app
ARG DB_PASSWORD=example-password
ENV API_TOKEN=example-token

COPY . .
COPY .env .env
RUN go build -o risky-service ./cmd/risky-service

CMD ["./risky-service"]
