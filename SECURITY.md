# Security notes

- Never commit real `.env` files or production secrets.
- Use a secrets manager in production.
- Replace development database credentials before deployment.
- Keep `CORS_ALLOWED_ORIGINS` restrictive in production.
- Run `govulncheck ./...` in CI if your organization allows the additional tool.
- Run migrations as a controlled deployment step.
- Review logs to ensure they do not include secrets or sensitive request bodies.
