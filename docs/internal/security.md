# Security Notes

## HTTP Basic Auth credentials

Scenaria v1.0 stores HTTP Basic Auth credentials in the application settings
JSON so recording and playback can reuse them for a host. This is explicit
plaintext at-rest storage, not OS-protected secret storage.

Release policy:

- Passwords are accepted from explicit user input only.
- `LoadSettings` and HTTP auth read DTOs do not return the password value to the
  frontend; they return only the username and a `hasPassword` flag.
- URLs entered as `user:pass@host` are parsed, the credentials are stored for the
  host, and the URL used by the app is stripped before it is persisted or logged.
- Diagnostic backend logs and frontend journal/status text redact common secret
  forms such as passwords, bearer/basic auth headers, cookies, tokens, and URL
  userinfo.

Future protected storage migration should preserve existing settings safely and
move passwords to Windows Credential Manager/DPAPI, macOS Keychain, and Linux
Secret Service.
