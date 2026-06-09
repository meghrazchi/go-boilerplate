# Contributing

Recommended workflow:

```bash
cp .env.example .env
go mod tidy
make install-tools
make precommit-install
make verify
```

Before opening a pull request, run:

```bash
make verify
make test-integration
```
