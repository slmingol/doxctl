# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- CONTRIBUTING.md with development guidelines and PR process
- CHANGELOG.md following Keep a Changelog format
- SECURITY.md with security policy and vulnerability reporting process
- ARCHITECTURE.md with system design documentation

## [0.0.59-alpha] - 2023-02-28

### Changed
- Updated to Go 1.25.0
- Improved goreleaser configuration

## [0.0.58.1-alpha] - Prior to 2023-02-28

### Added
- VPN connectivity diagnostics
- DNS resolver and search path validation
- Network connectivity testing
- Service reachability checks
- Server reachability checks
- Configuration file support (.doxctl.yaml)
- Cross-platform support (macOS, Linux)
- Multi-architecture builds (amd64, arm64)
- Docker container support
- Homebrew tap for easy installation
- GitHub Actions CI/CD pipeline
- Code coverage reporting
- GoReleaser integration for automated releases

### Features
- VPN route validation
- DNS resolution testing
- Ping-based connectivity checks
- Customizable configuration via YAML
- Color-coded output for better readability
- Table-formatted results

## Project History

This project was created in 2021 to help diagnose connectivity problems stemming from VPN and DNS setups on laptops. It provides a comprehensive suite of diagnostic tools for:

- VPN connectivity across geo-locations
- DNS resolver and search path configuration
- Network routing and connectivity
- Well-known server reachability

---

## How to Use This Changelog

### For Users
- Check the [Unreleased] section to see what's coming in the next release
- Review version sections to see what changed in each release
- Look for security updates under the "Security" category

### For Contributors
- Add your changes to the [Unreleased] section
- Follow the Keep a Changelog format
- Use the appropriate category (Added, Changed, Deprecated, Removed, Fixed, Security)
- Include issue/PR references where applicable

### Categories
- **Added** for new features
- **Changed** for changes in existing functionality
- **Deprecated** for soon-to-be removed features
- **Removed** for now removed features
- **Fixed** for any bug fixes
- **Security** for vulnerability fixes

[Unreleased]: https://github.com/slmingol/doxctl/compare/0.0.59-alpha...HEAD
[0.0.59-alpha]: https://github.com/slmingol/doxctl/releases/tag/0.0.59-alpha
[0.0.58.1-alpha]: https://github.com/slmingol/doxctl/releases/tag/0.0.58.1-alpha
