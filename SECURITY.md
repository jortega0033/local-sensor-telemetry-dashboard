# Security policy

## Supported versions

Only the default branch (`main`) is supported. There are no releases or published packages.

## Reporting a vulnerability

Report it privately through GitHub: open the repository's **Security** tab and choose **Report a vulnerability** (private vulnerability reporting). Do not open a public issue for a vulnerability.

Include what you did, what you expected and what happened, and the commit you tested. The server is designed to listen on `127.0.0.1` only, to reject cross-origin mutations and to never build a filesystem path from an HTTP value; reports that break those rules are the most useful. Nothing is implemented yet, so reports are about the plan in `docs/TECHNICAL_SPEC.md` until the first tickets land.
