# SPDX-License-Identifier: MIT
# FileENIAC Contributing Guide

**Version**: 1.0
**Last updated**: 2026-06-28

Thank you for contributing to FileENIAC! This guide will help you set up
your development environment and understand the contribution process.

## Quick Links

- [Code of Conduct](CODE_OF_CONDUCT.md)
- [Security Policy](SECURITY.md)
- [Governance](GOVERNANCE.md)
- [Open Issues](https://github.com/Joao-Aschenbrenner/FileENIAC/issues)

## Development Environment

### Prerequisites

- **Go**: 1.26+
- **Rust**: 1.70+ (for desktop app)
- **Node.js**: 18+ (for frontend)
- **pnpm**: 8+ (for frontend dependencies)
- **Docker**: 20.10+ (for container builds)

### Backend Setup

```bash
# Clone the repository
git clone https://github.com/Joao-Aschenbrenner/FileENIAC.git
cd FileENIAC

# Build the backend
go build -o fileeniac .

# Run tests (sequential)
go test ./...

# Run tests with race detector
go test -race ./...

# Run vet
go vet ./...
```

### Desktop App Setup

```bash
# Install frontend dependencies
cd apps/desktop
pnpm install

# Run in development mode
pnpm run tauri dev

# Build production desktop app
pnpm run tauri -- build
```

### Docker Build

```bash
docker build . -t fileeniac:test
docker compose up -d
```

## Branch Strategy

- `main`: Stable release branch
- `sprint-*`: Feature branches for sprints
- `fix/*`: Bug fix branches
- `docs/*`: Documentation-only branches

## Commit Messages

FileENIAC uses conventional commits:

```
<type>(<scope>): <description>

[optional body]
```

### Types

| Type | Description |
|------|-------------|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation changes |
| `test` | Adding or updating tests |
| `refactor` | Code refactoring |
| `perf` | Performance improvement |
| `security` | Security-related changes |
| `chore` | Build, CI, tooling changes |

### Examples

```
feat(sessions): add clear workspace endpoint
fix(validate): prevent path traversal in SafeRelativePath
docs(adr): add ADR-015 for plugin architecture
test(api): add integration tests for session activate
```

## Pull Request Process

### Before Submitting

1. **Check existing issues**: Avoid duplicate work. For significant
   changes, open an issue first to discuss the approach.
2. **Run tests locally**: Ensure `go test -race ./...` passes.
3. **Run vet**: Ensure `go vet ./...` reports no issues.
4. **Build**: Ensure `go build ./...` completes without errors.
5. **ADR for architecture changes**: If your change affects the
   architecture, create an ADR first (`docs/ADR-*.md`).

### PR Description

Include in your PR description:

- **Summary**: What does this change do?
- **Motivation**: Why is this change needed?
- **Changes**: Detailed list of changes
- **Testing**: How was this tested?
- **Screenshots**: For UI changes
- **Breaking changes**: If any

### PR Checklist

- [ ] Code follows existing style (run `go fmt` and `go vet`)
- [ ] Tests pass: `go test -race ./...`
- [ ] Build succeeds: `go build ./...`
- [ ] New tests added for new functionality
- [ ] Documentation updated (if needed)
- [ ] ADR created (if architecture change)
- [ ] No hardcoded credentials or secrets

### Review Process

1. Maintainer reviews the PR
2. Feedback may be provided for improvements
3. Once approved, maintainer merges
4. PR is closed

## Architecture Decisions

For any significant architectural change, an ADR (Architecture Decision Record)
is required. See `docs/ADR-*.md` for existing ADRs.

To create a new ADR:

1. Create `docs/ADR-NNN-title.md` (use next sequential number)
2. Include sections: Status, Context, Decision, Consequences
3. Submit as part of your PR or as a separate discussion

## Testing Guidelines

### Unit Tests

```bash
go test ./internal/validate/... -v
```

### Integration Tests

Integration tests require a running environment:

```bash
go test ./integration/... -v
```

### Race Detector

Always run with `-race` before submitting:

```bash
go test -race ./...
```

## Coding Standards

### Go

- Use `go fmt` for formatting
- Use `go vet` for static analysis
- No hardcoded credentials — use environment variables
- All exported functions have doc comments
- Error wrapping with `fmt.Errorf("...: %w", err)`

### Frontend (TypeScript/React)

- Use TypeScript strict mode
- Component files: PascalCase (e.g., `SessionList.tsx`)
- Hook files: camelCase with `use` prefix (e.g., `useSession.ts`)
- All props and state typed

## Docker Guidelines

- Dockerfile uses `golang:1.26-alpine` base
- No hardcoded credentials in Dockerfiles
- Use multi-stage builds where appropriate
- Validate with `docker compose up -d`

## Documentation

Documentation is as important as code. Update relevant docs:

- `docs/ADR-*.md` for architecture changes
- `README.md` for user-facing changes
- `RELEASE_NOTES.md` for all releases
- `docs/legal/*` for legal changes

## License

By submitting a contribution, you agree that your contribution will be
licensed under the MIT License. You must have the right to license
your contribution.

## Questions?

- Open a GitHub Discussion for questions
- Open a GitHub Issue for bugs
- Read existing ADRs for architectural context