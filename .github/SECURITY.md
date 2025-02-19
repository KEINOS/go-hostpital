# Security Policy

## Supported Versions and Statuses

| Version/Section | Status | Note |
| :------ | :----- | :--- |
| Go v1.24 to latest | [![go1.24+](https://github.com/KEINOS/go-hostpital/actions/workflows/unit-tests.yml/badge.svg)](https://github.com/KEINOS/go-hostpital/actions/workflows/unit-tests.yml "Unit tests on various Go versions") | |
| Golangci-lint v1.50.1 or later | [![golangci-lint](https://github.com/KEINOS/go-hostpital/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/KEINOS/go-hostpital/actions/workflows/golangci-lint.yml) | |
| Security advisories | [Enabled](https://github.com/KEINOS/go-hostpital/security/advisories) | |
| Dependabot alerts | [Enabled](https://github.com/KEINOS/go-hostpital/security/dependabot) | (Viewable only for admins) |
| Code scanning alerts | [Enabled](https://github.com/KEINOS/go-hostpital/security/code-scanning)<br>[![CodeQL-Analysis](https://github.com/KEINOS/go-hostpital/actions/workflows/codeQL-analysis.yml/badge.svg)](https://github.com/KEINOS/go-hostpital/actions/workflows/codeQL-analysis.yml) ||

## Update

- We [check the latest version of `go.mod` every week](https://github.com/KEINOS/go-hostpital/blob/main/.github/workflows/weekly-update.yml) and update it when it has passed all tests.

### Fail Fast Policy

We follow the "Fail Fast" policy. We constantly test the package with the latest Go versions and related modules. If the tests fail, we will update/fix the package ASAP.

If the latest version breaks the backward compatibility, we will bump up the minimum required version of Go.

## Reporting a Vulnerability, Bugs and etc

- [Issues](https://github.com/KEINOS/go-hostpital/issues)
  - [![Opened Issues](https://img.shields.io/github/issues/KEINOS/go-hostpital?color=lightblue&logo=github)](https://github.com/KEINOS/go-hostpital/issues "opened issues")
  - Plase attach a simple test that replicates the issue.
