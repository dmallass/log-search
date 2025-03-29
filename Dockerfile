# --- BUILD STAGE ---
FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git build-base ripgrep curl
WORKDIR /app
COPY . .
RUN go build -o log_search_app .

# --- FINAL STAGE ---
FROM alpine:latest
RUN apk add --no-cache ripgrep curl
WORKDIR /app
COPY --from=builder /app/log_search_app .
RUN chmod +x ./log_search_app
EXPOSE 8080
# Run the Go application binary 
CMD ["./log_search_app"]