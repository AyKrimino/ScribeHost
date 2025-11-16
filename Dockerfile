# --- Build ---
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/server server.go

# --- Runtime ---
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

RUN addgroup -S appgroup && adduser -S -G appgroup appuser

WORKDIR /app

COPY --from=builder /app/bin/server .
COPY --from=builder /app/migrations ./migrations

RUN chown -R appuser:appgroup ./

USER appuser

EXPOSE 4000

CMD [ "./server" ]
