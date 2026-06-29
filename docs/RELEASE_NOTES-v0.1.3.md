## FileENIAC v0.1.3 - Legal Hardening Patch

This release corrects legal documentation to match the actual behavior of the software.

### Corrections

* **SPDX**: `SPDX-LICENSE-IDENTIFIER` -> `SPDX-License-Identifier` in all documents
* **Vault**: Removed references to non-existent `FILEENIAC_VAULT_PASSWORD`. The vault uses auto-generated AES-256-GCM keys per workspace. No unencrypted fallback.
* **OAuth -> PAT**: Corrected all references from "GitHub OAuth" to "GitHub personal access token" to match actual code behavior
* **SQLite advice**: Replaced "edit the SQLite database directly" with "use application settings"
* **CCPA removed**: Changed to "Privacy Rights (LGPD/GDPR)"
* **Termination/Amendments**: MIT rights preserved; changes apply to future releases only
* **EULA removed**: `INSTALLER_EULA.md` replaced with `INSTALLER_NOTICE.md` — no conflict with MIT license

### New Documents

* `docs/legal/THIRD_PARTY_SERVICES.md`
* `docs/legal/SECURITY_AND_CREDENTIALS.md`

### Installer

* `installer-license.txt` rewritten with corrected legal notice
* WebView2 bootstrapper included (no more DLL errors)

### Checksums

| File | SHA-256 |
|------|---------|
| FileENIAC_0.1.3_x64-setup.exe | 8CEC546608DEC4558131968EA8BE88B236D8C8173574DE3830CE21AD6CAFD4F2 |

### Links

* [MIT License](https://github.com/Joao-Aschenbrenner/FileENIAC/blob/main/LICENSE)
* [Terms of Use](https://github.com/Joao-Aschenbrenner/FileENIAC/blob/main/docs/legal/TERMS_OF_USE.md)
* [Privacy Policy](https://github.com/Joao-Aschenbrenner/FileENIAC/blob/main/docs/legal/PRIVACY_POLICY.md)
* [LGPD Compliance](https://github.com/Joao-Aschenbrenner/FileENIAC/blob/main/docs/legal/LGPD.md)
* [Third-Party Services](https://github.com/Joao-Aschenbrenner/FileENIAC/blob/main/docs/legal/THIRD_PARTY_SERVICES.md)
* [Installer Notice](https://github.com/Joao-Aschenbrenner/FileENIAC/blob/main/docs/legal/INSTALLER_NOTICE.md)
