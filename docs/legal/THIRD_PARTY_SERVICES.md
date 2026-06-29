# SPDX-License-Identifier: MIT
# FileENIAC — Third-Party Services

**Version**: 1.0
**Last updated**: 2026-06-28

## Overview

FileENIAC interacts with third-party services only when explicitly configured
by the user. This document describes those interactions, the data involved,
and the responsibilities of each party.

## Services

### GitHub

**Interaction**: GitHub authentication via personal access token, repository
management (clone, push, pull), issue tracking, update checks.

**Data sent**: GitHub token, repository metadata, git objects.

**Data received**: Repository data, version metadata, API responses.

**Governing terms**: [GitHub Terms of Service](https://docs.github.com/en/site-policy/github-terms/github-terms-of-service)

**Developer responsibility**: The developer does not operate GitHub and has no
control over GitHub's data practices.

### Git Providers (GitLab, Bitbucket, self-hosted)

**Interaction**: Git operations (clone, push, pull) to any Git provider
configured by the user.

**Data sent**: Git objects, repository data, authentication credentials.

**Governing terms**: Terms of the respective Git provider or your organization.

**Developer responsibility**: The developer has no access to these transactions
and is not responsible for the security practices of third-party Git providers.

### FTPS Servers

**Interaction**: File transfer operations (upload, download, sync) to FTPS
servers configured by the user.

**Data sent**: File contents, paths, authentication credentials.

**Governing terms**: Terms of your hosting provider or internal policies.

**Developer responsibility**: The developer has no access to FTPS credentials
or transferred files. FTPS server security is the responsibility of the
server operator.

### GitHub Releases API

**Interaction**: Optional version check on application startup.

**Data sent**: None (no personal data, no credentials, no identifiers).

**Data received**: Release version metadata (tag names, release notes).

**Developer responsibility**: Standard HTTP request to GitHub's public API.
GitHub receives network metadata (IP address, user agent) per their own
privacy policy.

### Operating System

**Interaction**: File system operations, process management, network access.

**Data access**: The OS has full access to FileENIAC data. Standard OS
security controls (file permissions, disk encryption, antivirus) apply.

### Tauri / WebView2

**Interaction**: Application runtime, native window management, WebView
rendering.

**Data access**: The WebView runtime renders local content only. No data
is sent to external servers by the runtime itself.

## Disclaimer

The maintainers are not responsible for the privacy, security, or operational
practices of third-party services configured by the user.

You should review the privacy policies and terms of service of any third-party
service you configure FileENIAC to interact with.
