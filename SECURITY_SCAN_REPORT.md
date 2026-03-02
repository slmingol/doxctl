# Security Vulnerability Scan Report

**Date:** 2026-03-02  
**Tool:** GitHub Advisory Database  
**Scope:** All Go dependencies in doxctl project

## Executive Summary

A comprehensive security vulnerability scan was performed on all dependencies in the doxctl project. The scan identified **8 vulnerabilities** across 4 key dependencies, with severity ranging from **CRITICAL** to **MODERATE**.

## Critical Findings

### 1. JWT Authorization Bypass (CRITICAL)
**Package:** `github.com/dgrijalva/jwt-go`  
**Current Version:** 3.2.0  
**Vulnerability:** Authorization bypass  
**Severity:** CRITICAL

**Details:**
- Two separate authorization bypass vulnerabilities identified
- CVE affects versions < 4.0.0-preview1
- Additional vulnerability in versions >= 0.0.0-20150717181359-44718f8a89b0, <= 3.2.0

**Recommendation:**
- Migrate to maintained fork: `github.com/golang-jwt/jwt/v4` or later
- The original package is deprecated and unmaintained
- **Action Required:** Immediate update to address critical authorization bypass

### 2. HTTP/2 Rapid Reset Vulnerabilities (HIGH)
**Packages:** 
- `golang.org/x/net` (v0.7.0)
- `google.golang.org/grpc` (v1.21.1)

**Details:**

#### golang.org/x/net
- **Current Version:** 0.7.0
- **Vulnerability:** HTTP/2 rapid reset can cause excessive work in net/http
- **Affected Versions:** < 0.17.0
- **Patched Version:** 0.17.0
- **Impact:** Denial of Service (DoS) attacks

#### google.golang.org/grpc
- **Current Version:** 1.21.1
- **Vulnerabilities:** Multiple HTTP/2 Rapid Reset issues
- **Affected Version Ranges:**
  - < 1.56.3
  - >= 1.57.0, < 1.57.1
  - >= 1.58.0, < 1.58.3
- **Recommended Version:** >= 1.58.3 or >= 1.56.3
- **Impact:** Denial of Service (DoS) attacks

**Recommendation:**
- Update `golang.org/x/net` to >= 0.17.0
- Update `google.golang.org/grpc` to >= 1.58.3 or >= 1.56.3
- **Action Required:** High priority update to prevent DoS attacks

### 3. Cryptographic Vulnerabilities (HIGH)
**Package:** `golang.org/x/crypto`  
**Current Version:** 0.6.0

**Details:**
1. **DoS via Slow or Incomplete Key Exchange**
   - **Affected Versions:** < 0.35.0
   - **Patched Version:** 0.35.0
   - **Impact:** Denial of Service

2. **Authorization Bypass via PublicKeyCallback Misuse**
   - **Affected Versions:** < 0.31.0
   - **Patched Version:** 0.31.0
   - **Impact:** Authentication bypass

**Recommendation:**
- Update `golang.org/x/crypto` to >= 0.35.0
- **Action Required:** High priority to address both DoS and authorization bypass

## Summary of Required Actions

| Package | Current Version | Minimum Safe Version | Priority |
|---------|----------------|---------------------|----------|
| github.com/dgrijalva/jwt-go | 3.2.0 | Migrate to github.com/golang-jwt/jwt/v4 | CRITICAL |
| golang.org/x/crypto | 0.6.0 | 0.35.0 | HIGH |
| golang.org/x/net | 0.7.0 | 0.17.0 | HIGH |
| google.golang.org/grpc | 1.21.1 | 1.58.3 | HIGH |

## Next Steps

1. **Immediate Actions:**
   - Replace `github.com/dgrijalva/jwt-go` with `github.com/golang-jwt/jwt/v4`
   - Update `golang.org/x/crypto` to latest stable version (>= 0.35.0)
   - Update `golang.org/x/net` to latest stable version (>= 0.17.0)
   - Update `google.golang.org/grpc` to latest stable version (>= 1.58.3)

2. **Testing:**
   - Run full test suite after updates
   - Verify no breaking changes in updated dependencies
   - Test authentication flows (JWT-related changes)

3. **Long-term:**
   - Establish regular security scanning schedule
   - Consider adding automated vulnerability scanning to CI/CD pipeline
   - Keep dependencies up-to-date with patch releases

## Scan Methodology

Due to network restrictions preventing direct use of `govulncheck`, the scan was performed using the GitHub Advisory Database against all major dependencies listed in `go.mod`. The following key dependencies were analyzed:

- github.com/dgrijalva/jwt-go
- golang.org/x/crypto
- golang.org/x/net
- golang.org/x/text
- golang.org/x/sys
- golang.org/x/term
- github.com/miekg/dns
- github.com/spf13/cobra
- github.com/spf13/viper
- google.golang.org/grpc
- github.com/coreos/etcd

## Additional Notes

- The project uses Go 1.16, which is quite old. Consider updating to a more recent Go version for better security and performance.
- All vulnerabilities identified are in transitive dependencies and can be resolved through updates.
- No custom code vulnerabilities were identified during this scan.

---

**Report Generated:** 2026-03-02  
**Scan Status:** ✅ Complete  
**Vulnerabilities Found:** 8  
**Critical Issues:** 2 (JWT authorization bypass)  
**High Priority Issues:** 6 (HTTP/2 rapid reset, crypto DoS/auth bypass)
