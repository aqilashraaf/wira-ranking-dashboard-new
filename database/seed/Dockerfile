FROM golang:1.21-alpine

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies and verify modules
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the application
RUN go mod tidy && go build -o seed

# Run the seeder
CMD ["./seed"]
