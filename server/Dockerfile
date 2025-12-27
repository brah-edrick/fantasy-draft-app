# Start with the official Go image
FROM golang:1.24-alpine

# Install "Air" for hot reloading
RUN go install github.com/air-verse/air@v1.61.1

# Set working directory inside the container
WORKDIR /app

# Copy dependency files first (better caching)
COPY go.mod go.sum ./
# Download dependencies (this will warn if go.sum matches go.mod, but won't fail if go.sum is empty and no deps are in go.mod yet)
RUN go mod download

# Copy the rest of the app
COPY . .

# Run "air" with explicit config
CMD ["air", "-c", "air.toml"]
