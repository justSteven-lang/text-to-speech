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

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o tts-app .

# =========================
# Runtime stage
# =========================
FROM alpine:3.19

RUN apk add --no-cache espeak ca-certificates wget

WORKDIR /app

COPY --from=builder /app/tts-app .


EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
CMD wget --spider -q http://localhost:8080/health || exit 1

RUN adduser -D appuser
USER appuser

CMD ["./tts-app"]
