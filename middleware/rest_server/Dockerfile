# Builder stage
FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o rest_server .

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/rest_server .
RUN apk --no-cache add ca-certificates tzdata
CMD ["./rest_server"]