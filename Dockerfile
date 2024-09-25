# Use the official Golang image as a build stage
FROM golang:1.22-alpine AS builder
# First Stage as build file
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o main .

# Second Stage with lightweight docker container image
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
