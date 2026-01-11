# Security Policy

## Supported Versions
| Version | Supported |
| --- | --- |
| Latest | Yes |

## Reporting a Vulnerability
Do not report security vulnerabilities via public GitHub issues.
Report via the Security tab -> Report a vulnerability.
You may also include a PR with a proposed fix in your private report.

### Information to Include
- Detailed description of the vulnerability type, such as issues with encryption (AES-256-GCM), key derivation, or memory-safe handling of secrets.
- Assessment of the potential impact, specifically whether it allows for secret exposure, unauthorized token generation, or encryption bypass.
- Clear, step-by-step instructions to reproduce the issue, including any specific CLI commands or input sequences required.
- Full technical environment details including the gauth version, Go compiler version, and the host operating system.

## Commitment
Acknowledgement within 48 hours, followed by technical verification, a code fix, and a public advisory once resolved.