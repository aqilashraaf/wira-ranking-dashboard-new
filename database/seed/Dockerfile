FROM golang:1.21-alpine

# Install PostgreSQL client
RUN apk add --no-cache postgresql-client

WORKDIR /app

# Copy go mod and sum files from seeder directory
COPY database/seed/go.mod database/seed/go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY database/seed/ .

# Copy the wait-for-postgres script
COPY wait-for-postgres.sh ./wait-for-postgres.sh
RUN chmod +x ./wait-for-postgres.sh

# Build the application
RUN go build -o seed .

# Command to run
CMD ["./wait-for-postgres.sh", "db", "./seed"]
