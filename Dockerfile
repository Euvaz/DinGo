# syntax=docker/dockerfile:1

FROM --platform=linux/amd64 alpine:3.16

# Add compatibility for Go binary
RUN apk add --no-cache libc6-compat

# Create working directory
RUN mkdir -p /usr/src/dingo
WORKDIR /usr/src/dingo

# Copies project files
COPY . .

# Set permissions and execute
RUN chmod +x ./bin/DinGo
CMD ["./bin/DinGo"]
