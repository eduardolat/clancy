# Use golang:1.25-trixie as the base image
FROM golang:1.25-trixie

# Declare ARG to make it available in the RUN commands
ARG TARGETPLATFORM

# Set environment variables
ENV DEBIAN_FRONTEND="noninteractive"

RUN set -e && \
    # Update and install system dependencies
    apt-get update -qq && \
    apt-get install -y -qq --no-install-recommends \
        wget curl zip unzip p7zip-full git tree && \
    # Install downloadable binaries in a temporary directory
    mkdir -p /tmp/bin-downloads && cd /tmp/bin-downloads && \
    if [ "${TARGETPLATFORM}" = "linux/arm64" ]; then \
        ARCH="arm64"; \
    else \
        ARCH="amd64"; \
    fi && \
    # Task
    wget --no-verbose https://github.com/go-task/task/releases/download/v3.46.4/task_linux_${ARCH}.tar.gz && \
    tar -xzf task_linux_${ARCH}.tar.gz && \
    mv ./task /usr/local/bin/task && \
    # Golangci-lint
    wget --no-verbose https://github.com/golangci/golangci-lint/releases/download/v2.8.0/golangci-lint-2.8.0-linux-${ARCH}.tar.gz && \
    tar -xzf golangci-lint-2.8.0-linux-${ARCH}.tar.gz && \
    mv ./golangci-lint-2.8.0-linux-${ARCH}/golangci-lint /usr/local/bin/golangci-lint && \
    # Final setup: chmod, git config, and cleanup
    chmod +x /usr/local/bin/* && \
    git config --global --add safe.directory '*' && \
    cd / && rm -rf /tmp/bin-downloads /var/lib/apt/lists/*

WORKDIR /app

# Add the startup script on every bash session
RUN rm -f /root/.bashrc
COPY .devcontainer/.bashrc /root/.bashrc

# Command just to keep the container running
CMD ["sleep", "infinity"]
