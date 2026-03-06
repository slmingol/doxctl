# Docker Configuration

This directory contains the Dockerfile for building doxctl container images.

## Quick Start

### Pull Pre-built Image

From Docker Hub:
```bash
docker pull slmingol/doxctl:latest
```

From GitHub Container Registry:
```bash
docker pull ghcr.io/slmingol/doxctl:latest
```

### Run Container

```bash
# Show help
docker run --rm slmingol/doxctl:latest --help

# Run DNS diagnostics
docker run --rm slmingol/doxctl:latest dns

# Run VPN diagnostics
docker run --rm slmingol/doxctl:latest vpn

# Run with custom config
docker run --rm \
  -v $(pwd)/.doxctl.yaml:/root/.doxctl.yaml:ro \
  slmingol/doxctl:latest dns

# Check version
docker run --rm slmingol/doxctl:latest version
```

### Build Locally

```bash
# Build from project root
docker build -f docker/Dockerfile -t doxctl:local .

# Run local build
docker run --rm doxctl:local --help
```

## Multi-Architecture Support

Images are built for multiple architectures:
- `linux/amd64` (x86_64)
- `linux/arm64` (ARM 64-bit)

Docker will automatically pull the correct image for your platform.

## Image Tags

- `latest` - Latest stable release
- `X.Y.Z` - Specific version (e.g., `1.0.0`)
- `X.Y.Z-amd64` - Architecture-specific tag
- `X.Y.Z-arm64` - Architecture-specific tag

## Image Registries

Images are published to:
- **Docker Hub**: `slmingol/doxctl`
- **GitHub Container Registry**: `ghcr.io/slmingol/doxctl`

## Volumes

Mount your config file to use custom settings:
```bash
-v /path/to/.doxctl.yaml:/root/.doxctl.yaml:ro
```

## Versioning

All Docker images are tagged with the same semantic version as the CLI binary. The version is injected during the build process via GoReleaser and is displayed when running:

```bash
docker run --rm slmingol/doxctl:latest version
```

## Build Process

Images are automatically built and published by GoReleaser when a new version tag is pushed:

1. Developer merges PR to main
2. Auto-version workflow creates new semantic version tag
3. Build-release workflow triggers on tag
4. GoReleaser builds binaries and Docker images
5. Images pushed to Docker Hub and GHCR with version tags
