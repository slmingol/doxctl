FROM alpine:latest

# Install ca-certificates for HTTPS support
RUN apk --no-cache add ca-certificates

# Copy the binary
COPY doxctl-linux /usr/bin/doxctl

# Set the entrypoint
ENTRYPOINT ["/usr/bin/doxctl"]

# Default command shows help
CMD ["--help"]
