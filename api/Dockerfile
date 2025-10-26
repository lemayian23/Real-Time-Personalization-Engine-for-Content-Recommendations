FROM golang:1.21-alpine

WORKDIR /app

# Copy go mod files first
COPY go.mod ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o server ./cmd/server

EXPOSE 8080
CMD ["./server"]