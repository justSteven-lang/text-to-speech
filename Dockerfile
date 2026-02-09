# =========================
# Build stage
# =========================
FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod ./
RUN go mod download || true

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o tts-app .

# =========================
# Runtime stage
# =========================
FROM alpine:3.19

RUN apk add --no-cache espeak ca-certificates

WORKDIR /app

COPY --from=builder /app/tts-app .

USER nobody

CMD ["./tts-app"]
