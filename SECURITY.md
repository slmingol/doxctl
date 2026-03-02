# Security Policy

## Supported Versions

We take security seriously and are committed to addressing security vulnerabilities promptly. The following versions of doxctl are currently supported with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 0.0.59-alpha (latest) | :white_check_mark: |
| < 0.0.59-alpha | :x: |

**Note**: As this project is in alpha stage, we recommend always using the latest version to ensure you have the most recent security fixes.

## Reporting a Vulnerability

We appreciate your efforts to responsibly disclose security vulnerabilities. If you discover a security issue, please follow these steps:

### How to Report

**DO NOT** open a public GitHub issue for security vulnerabilities.

Instead, please report security vulnerabilities by:

1. **Email**: Send a detailed report to **github@lamolabs.org**
   - Include "SECURITY" in the subject line
   - Provide a clear description of the vulnerability
   - Include steps to reproduce the issue
   - Describe the potential impact
   - Suggest a fix if you have one

2. **GitHub Security Advisory** (Preferred): Use GitHub's private vulnerability reporting
   - Go to the [Security tab](https://github.com/slmingol/doxctl/security)
   - Click "Report a vulnerability"
   - Fill out the advisory form with detailed information

### What to Include

A good security report should include:

- **Description**: Clear description of the vulnerability
- **Impact**: What could an attacker accomplish?
- **Affected Versions**: Which versions are affected?
- **Reproduction Steps**: Step-by-step instructions to reproduce
- **Proof of Concept**: Code or commands that demonstrate the issue (if applicable)
- **Suggested Fix**: Your recommendation for fixing the issue (optional)
- **Disclosure Timeline**: Your expectations for public disclosure

### Example Report

```
Subject: SECURITY - Command Injection in DNS Resolver Check

Description:
The DNS resolver check command does not properly sanitize user input,
allowing command injection through the --resolver flag.

Affected Versions: 0.0.58-alpha and earlier

Steps to Reproduce:
1. Run: doxctl dns --resolver "8.8.8.8; malicious-command"
2. Observe that malicious-command is executed

Impact:
An attacker could execute arbitrary commands on the user's system
if they can control the --resolver parameter.

Suggested Fix:
Implement input validation and use parameterized command execution
instead of shell command interpolation.
```

## Response Timeline

We are committed to addressing security vulnerabilities in a timely manner:

| Stage | Timeline |
|-------|----------|
| **Initial Response** | Within 48 hours of report |
| **Vulnerability Confirmation** | Within 1 week |
| **Fix Development** | Depends on complexity (1-4 weeks) |
| **Security Patch Release** | As soon as fix is ready and tested |
| **Public Disclosure** | 90 days after patch release (or earlier with reporter agreement) |

### Process

1. **Acknowledgment** (48 hours)
   - We'll acknowledge receipt of your report
   - Assign a tracking identifier

2. **Assessment** (1 week)
   - We'll assess the vulnerability
   - Determine severity and affected versions
   - Provide initial feedback on validity

3. **Fix Development** (1-4 weeks)
   - Develop and test a fix
   - Keep you informed of progress
   - Request your validation if needed

4. **Release** (ASAP)
   - Release patched version
   - Update security advisories
   - Credit reporter (if desired)

5. **Disclosure** (90 days)
   - Coordinate public disclosure
   - Publish CVE if applicable
   - Update documentation

## Security Best Practices

When using doxctl, we recommend:

### For Users

- **Always use the latest version** to get security fixes
- **Validate your configuration file** before use
- **Review DNS and VPN settings** before running diagnostics
- **Run with least privilege** - don't use sudo unless necessary
- **Keep dependencies updated** if building from source
- **Use official releases** from GitHub or Homebrew

### For Contributors

- **Never commit secrets** or credentials to the repository
- **Validate all user input** before using in commands
- **Use parameterized commands** instead of shell interpolation
- **Run security scanners** on pull requests
- **Follow secure coding practices** per CONTRIBUTING.md
- **Review dependencies** for known vulnerabilities

## Security Features

doxctl implements the following security practices:

- **Input Validation**: User inputs are validated before use
- **Least Privilege**: Tool runs with user privileges (no sudo required for most operations)
- **No Credential Storage**: Does not store sensitive credentials
- **Safe Configuration**: Configuration files are read with appropriate permissions
- **Dependency Scanning**: Regular security scans of dependencies
- **Code Review**: All changes undergo code review before merge

## Known Limitations

As a diagnostic tool, doxctl:

- Requires network access to perform connectivity tests
- May execute network commands that could be logged by system administrators
- Reads system network configuration (DNS, routing tables)
- Uses ICMP ping which may be restricted on some networks

These are expected behaviors for a network diagnostic tool, not vulnerabilities.

## Security Contact

- **Primary Contact**: github@lamolabs.org
- **GitHub Security**: Use GitHub Security Advisories feature
- **Maintainer**: Sam Mingolelli

## Hall of Fame

We recognize and thank security researchers who have responsibly disclosed vulnerabilities:

<!-- This section will be updated as vulnerabilities are reported and fixed -->

*No vulnerabilities reported yet.*

---

Thank you for helping keep doxctl and its users safe! 🔒
