# Build stage
FROM golang:1.24-alpine AS buildstage

# Set working directory matching your project name
WORKDIR /Sole-Spot

# First copy module files to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy all source files
COPY . .

# Install build dependencies
RUN apk --no-cache add ca-certificates

# Build the application (verify ./cmd1/main.go is correct path)
RUN go build -o /Sole-Spot/sole-spot-app ./cmd1/main.go

# Final stage
FROM alpine:3.18

# Install runtime dependencies
RUN apk --no-cache add ca-certificates

# Set working directory matching project name
WORKDIR /Sole-Spot

# Copy binary from build stage
COPY --from=buildstage /Sole-Spot/sole-spot-app .
COPY --from=buildstage /Sole-Spot/.env .


# Copy SSL certificates
COPY --from=buildstage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy templates
COPY --from=buildstage /Sole-Spot/templates ./templates/

# Run the application
CMD ["./sole-spot-app"]
