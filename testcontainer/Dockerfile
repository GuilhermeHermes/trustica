# Build stage for test server
FROM golang:1.23-alpine AS server-builder

WORKDIR /build
COPY /testcontainer/testserver/ .
RUN go build -o /testserver main.go


# Build stage for trustica
FROM golang:1.23-alpine AS trustica-builder

WORKDIR /build
COPY go.mod go.sum* ./
RUN go mod download || true

COPY . .
RUN go build -o /trustica ./cmd/trustica


# Final test container
FROM ubuntu:22.04

# Avoid interactive prompts
ENV DEBIAN_FRONTEND=noninteractive

# Install all the tools we want to test
RUN apt-get update && apt-get install -y \
    # Python
    python3 \
    python3-pip \
    python3-venv \
    # Git
    git \
    # Node.js (via nodesource for newer version)
    curl \
    ca-certificates \
    gnupg \
    # OpenSSL and utilities
    openssl \
    && rm -rf /var/lib/apt/lists/*

# Install Node.js 20.x
RUN mkdir -p /etc/apt/keyrings \
    && curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg \
    && echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_20.x nodistro main" | tee /etc/apt/sources.list.d/nodesource.list \
    && apt-get update \
    && apt-get install -y nodejs \
    && rm -rf /var/lib/apt/lists/*

# Create directories
RUN mkdir -p /certs /app /tests

# Copy certificates
COPY testcontainer/certs/ /certs/

# Copy test server
COPY --from=server-builder /testcontainer/testserver /usr/local/bin/testserver

# Copy trustica binary
COPY --from=trustica-builder /testcontainer/trustica /usr/local/bin/trustica

# Copy test scripts
COPY testcontainer/scripts/ /tests/

# Add testserver.local to hosts
RUN echo "127.0.0.1 testserver.local" >> /etc/hosts

# Make scripts executable
RUN chmod +x /tests/*.sh

WORKDIR /app

# Default command runs the full test suite
CMD ["/tests/run_all_tests.sh"]
