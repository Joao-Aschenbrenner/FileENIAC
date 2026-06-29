# External Blind Audit — Instructions

## Purpose

This document provides the instructions for an external AI to perform a blind
audit of FileENIAC. The auditor should receive **only** the repository URL and
this document — no prior history, no context of previous sprints, no explanation
of what was recently changed.

## Repository

https://github.com/Joao-Aschenbrenner/FileENIAC

## Audit Instructions for the AI

You are an independent security and quality auditor. You have been given
access to the FileENIAC repository. You do NOT know anything about its
development history, recent changes, or previous audits.

Please perform a thorough audit covering the following areas:

### 1. Architecture
- Is the project structure logical and well-organized?
- Are there architectural anti-patterns?

### 2. Security
- Hardcoded credentials, tokens, secrets
- SQL injection vectors
- Path traversal risks
- Authentication weaknesses
- Cryptographic practices
- Credential handling
- Log leakage

### 3. License
- Is the LICENSE file correct and complete?
- Does the project comply with its stated license (MIT)?
- Are third-party dependencies properly attributed?
- Are there any license conflicts?

### 4. Documentation (Legal)
- Review `docs/legal/` folder: are the documents consistent with each other?
- Do they match what the code actually does?
- Are there contradictions between legal docs and code behavior?
- Is the installer notice correct and MIT-compatible?
- Check for any misleading claims (e.g., claims about encryption, data collection)

### 5. Installer
- Does the installer reference appropriate legal notice?
- Are version numbers consistent?
- Are checksums provided?

### 6. Tests
- Do the backend Go tests pass?
- Do the frontend tests pass?

### 7. Docker
- Are Dockerfiles secure (non-root user, no secrets baked in)?
- Is docker-compose.yml functional?

### 8. Frontend
- Are there security issues in the React/TypeScript frontend?
- Are dependencies up to date?

### 9. Performance
- Are there obvious performance bottlenecks?
- Resource leaks (goroutines, connections)?

### 10. False Positives
- If you find an issue, confirm it's real
- Mark anything that looks suspicious but is actually correct as "false positive"

## Output Format

Generate `docs/audits/EXTERNAL_BLIND_AUDIT.md` with the following structure:

```markdown
# External Blind Audit — FileENIAC

## Summary

[One-paragraph summary of findings]

## Findings

### Critical (blocking release)

### High

### Medium

### Low

### False Positives (items that look wrong but are correct)

### Inconclusive (needs human review)

## Recommendations

- Keep release as-is
- Patch release recommended
- Block release (do not publish)

## Plan

[If corrections are needed, list them in priority order]
```
