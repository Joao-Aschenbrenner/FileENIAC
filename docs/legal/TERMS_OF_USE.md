# SPDX-License-Identifier: MIT
# FileENIAC Terms of Use

**Effective date**: 2026-06-28
**Version**: 1.1

## 1. Acceptance

By installing, copying, or using FileENIAC, you agree to be bound by these
Terms of Use. If you do not agree, do not install or use the software.

## 2. License

FileENIAC source code is licensed under the MIT License.
See the `LICENSE` file for the full license text.

These Terms apply to the official FileENIAC binary distributions, installer,
branding, documentation, support channels, and official project infrastructure.
They do not restrict rights granted under the MIT License for the source code.

## 3. Allowed Use

FileENIAC may be used for:

- Managing local workspace environments
- Operating Git repositories
- Authenticating with GitHub for repository management
- Deploying files via FTPS to servers you own or have permission to use
- Running on personal computers or servers you control
- Integration into commercial or non-commercial projects (under MIT terms)

## 4. Prohibited Use

You may NOT use FileENIAC to:

- Access GitHub accounts or FTPS servers you do not have authorization to use
- Deploy content to servers without proper authorization
- Attempt to extract, decrypt, or access credentials that do not belong to you
- Use the software in any way that violates applicable law
- Remove or obscure any copyright, trademark, or license notices

## 5. Your Data

You retain full ownership of all data processed by FileENIAC. The software
operates exclusively on data stored locally on your device. The developer
does not claim any ownership over your workspace data, credentials, or
configuration.

Your use of GitHub, GitLab, or any FTPS server must comply with the
respective Terms of Service of those platforms.

## 6. Security and Credential Management

You are responsible for:

- Keeping your device secure and protected from unauthorized access
- Ensuring the FileENIAC data directory is not accessible to other users
- Using strong, unique credentials for FTPS servers
- Reviewing and limiting access permissions for deployed content

The Vault encrypts credentials using AES-256-GCM with a key that is
auto-generated when a workspace is created. This key is stored locally
in your workspace configuration and is unique per installation.

**Never share your credentials or workspace configuration files.**
The developer will never ask for your credentials.

If you believe a credential has been compromised, rotate it immediately
in both FileENIAC and the affected service.

## 7. No Warranty

FILEENIAC IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR
A PARTICULAR PURPOSE AND NONINFRINGEMENT.

The developer does not warrant that:

- The software will meet your specific requirements
- The software will be uninterrupted, timely, secure, or error-free
- The results obtained from the use of the software will be accurate
  or reliable
- The quality of any deploys, syncs, or other operations will meet
  your expectations

## 8. Limitation of Liability

IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
DAMAGES, OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
OTHERWISE, ARISING FROM, OUT OF, OR IN CONNECTION WITH FILEENIAC OR THE USE
OR OTHER DEALINGS IN FILEENIAC.

This includes, without limitation:

- Data loss or corruption
- Failed deployments or sync operations
- Unauthorized access resulting from improper credential management
- Any business interruption, loss of revenue, or loss of data

## 9. Indemnification

To the maximum extent permitted by law, you are responsible for your use
of FileENIAC, including deployments, credentials, server access, and
compliance with third-party terms. You agree to indemnify, defend, and
hold harmless the developer from and against any claims, damages, losses,
costs, and expenses arising from your use of FileENIAC or your violation
of these Terms.

## 10. Use by Businesses

Businesses may use FileENIAC under the MIT License. There is no separate
commercial license required. When using FileENIAC in a business context:

- You are responsible for ensuring your use complies with applicable laws
- Your employees and contractors must agree to these Terms before using
  the software
- You are responsible for credential management across your organization

## 11. Third-Party Services

FileENIAC interacts with third-party services you configure:

- **GitHub**: Your use of GitHub is governed by GitHub's Terms of Service
- **FTPS Servers**: Your use of FTPS servers is governed by the terms
  set by the server operator

The developer is not affiliated with, endorsed by, or responsible for the
actions of any third-party service.

## 12. Export Controls

You are responsible for ensuring your use of FileENIAC complies with all
applicable export control laws and regulations.

## 13. Termination

These Terms are effective until terminated. Your rights under the MIT License
are perpetual, provided you comply with all terms. The maintainers may
discontinue official releases, support, update services, or project
infrastructure. This does not terminate rights already granted under the
MIT License for copies of the source code you have received.

Sections 7 (No Warranty), 8 (Limitation of Liability), and 9 (Indemnification)
survive any termination.

## 14. Governing Law

These Terms shall be governed by and construed in accordance with the laws
of Brazil, without regard to its conflict of law provisions.

## 15. Dispute Resolution

Any dispute arising from these Terms or your use of FileENIAC shall first be
attempted to be resolved through good-faith negotiation. If negotiation fails,
the dispute shall be submitted to the jurisdiction of the courts of Brazil.

## 16. Entire Agreement

These Terms, together with the MIT License and Privacy Policy, constitute
the entire agreement between you and the developer regarding your use of
FileENIAC.

## 17. Amendments

Changes to these Terms apply to future official releases or future use of
official project services. They do not retroactively change the MIT License
rights granted for previously received source code.

## 18. Contact

For questions regarding these Terms of Use:
- Open an issue: https://github.com/Joao-Aschenbrenner/FileENIAC/issues