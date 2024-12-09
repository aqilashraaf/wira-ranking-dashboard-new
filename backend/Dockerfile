FROM golang:1.21-alpine

WORKDIR /app

# Copy go mod and sum files from backend directory
COPY backend/go.mod backend/go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY backend/ .

# Copy the wait-for-postgres script
COPY wait-for-postgres.sh /wait-for-postgres.sh
RUN chmod +x /wait-for-postgres.sh

# Build the application
RUN go build -o main .

# Command to run
CMD ["./main"]
