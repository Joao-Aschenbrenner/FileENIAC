# SPDX-License-Identifier: MIT
# FileENIAC Governance

**Version**: 1.0
**Last updated**: 2026-06-28

## Project Structure

FileENIAC is an open-source project maintained by Joao Aschenbrenner.
The project follows a lightweight governance model appropriate for a
small, focused application.

## Decision-Making Process

All major decisions are made through:

1. **Issue Discussion**: Any community member can open an issue to propose
   changes, report bugs, or request features.
2. **Research**: For significant changes, the maintainer may research the
   topic and document the decision in an ADR (Architecture Decision Record).
3. **Consensus**: The maintainer considers community feedback. If consensus
   is reached, implementation proceeds.
4. **Implementation**: Changes are implemented in a feature branch,
   reviewed, and merged via pull request.
5. **Release**: Changes are released according to the release cadence.

For small changes (typos, documentation, minor fixes), the process is
streamlined: PR review and merge.

## Architecture Decision Records (ADRs)

Significant architectural and design decisions are documented as ADRs in
`docs/ADR-*.md`. Each ADR includes:

- **Status**: Proposed, Accepted, Deprecated, or Superseded
- **Context**: The problem or situation being addressed
- **Decision**: What was decided
- **Consequences**: Benefits and drawbacks of the decision

ADRs are the primary tool for documenting the "why" behind design choices.

## Roles

### Maintainer

- **Joao Aschenbrenner**: Project lead, primary developer
  - Makes final decisions on architecture and design
  - Reviews and merges pull requests
  - Manages releases and versioning
  - Sets strategic direction

### Contributors

Anyone who submits improvements via pull requests or documentation
fixes. Contributors retain copyright on their contributions per the
MIT License.

### Users

Anyone who uses FileENIAC. Users can:
- Report bugs and request features via GitHub Issues
- Participate in discussions
- Submit pull requests
- Vote on issues (using reactions)

## Roadmap

The roadmap is maintained in `docs/NEXT_SPRINTS.md`. Major items include:

| Version | Focus |
|---------|-------|
| v0.2.x | SFTP, WebDAV, S3 support |
| v0.3.x | Parallel uploads, resume |
| v0.4.x | Scheduler, file watcher |
| v0.5.x | Plugin system, bidirectional sync |

Roadmap is subject to change based on community feedback.

## Versioning

FileENIAC uses semantic versioning (SemVer):

- **MAJOR** version: Incompatible API changes
- **MINOR** version: New functionality in a backward-compatible manner
- **PATCH** version: Backward-compatible bug fixes

The current stable version is v0.1.x. See `SUPPORTED_VERSIONS.md` for
support windows.

## Release Process

Releases follow a structured process:

1. All CI/CD gates must pass (build, test, race detector)
2. Release notes are updated with changes since last version
3. Version numbers are updated in relevant files
4. Git tag is created (e.g., `v0.1.1`)
5. GitHub Release is published with artifacts
6. Checksums are generated and published

See `docs/RELEASE_CHECKLIST.md` for the full checklist.

## Community Standards

All participants are expected to follow the Code of Conduct
(`docs/legal/CODE_OF_CONDUCT.md`). All interactions should be:

- Respectful and constructive
- Focused on the technical topic
- Inclusive and welcoming

## Communication Channels

- **GitHub Issues**: Bug reports, feature requests, discussions
- **GitHub Discussions**: Q&A, ideas, announcements
- **Pull Requests**: Code contributions and reviews

## Changes to This Document

This governance document may be updated as the project evolves. Changes
will be announced in the repository. Significant changes will be discussed
via a GitHub issue before being merged.