# =========================
# Build stage
# =========================
FROM golang:1.25.4-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o tts-app .

# =========================
# Runtime stage
# =========================
FROM alpine:3.19

RUN apk add --no-cache espeak ca-certificates

WORKDIR /app

COPY --from=builder /app/tts-app .

EXPOSE 8080

USER nobody

CMD ["./tts-app"]
