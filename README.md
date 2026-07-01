# condition-generic

[![Latest Release](https://img.shields.io/github/v/release/SemRels/condition-generic?label=version\&color=blue)](https://github.com/SemRels/condition-generic/releases/latest)

Runs one or more shell commands and passes only when every configured command exits successfully.

This plugin is distributed as the standalone Go binary `semrel-plugin-condition-generic`. Semrel executes the binary as a subprocess, provides plugin configuration through `SEMREL_PLUGIN_*` environment variables, provides release context through `SEMREL_*` environment variables, reads standard output, and treats exit code `0` as success and any non-zero exit code as failure. Install the binary in `~/.semrel/plugins/` or anywhere on your `$PATH`.

## Installation

### Binary

```bash
go install github.com/SemRels/condition-generic/cmd/plugin@latest
```

### Docker

Pre-built, multi-platform images (linux/amd64, linux/arm64) are published to the GitHub Container Registry on every release:

```bash
docker pull ghcr.io/semrels/condition-generic:latest
```

Images are signed with [cosign](https://github.com/sigstore/cosign) and include a full SBOM attestation. Verify the signature:

```bash
cosign verify ghcr.io/semrels/condition-generic:latest \
  --certificate-identity-regexp 'https://github.com/SemRels/condition-generic/.github/workflows/release.yml.*' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com
```


## Configuration

```yaml
plugins:
  - name: condition-generic
    path: ~/.semrel/plugins/semrel-plugin-condition-generic
    env:
      SEMREL_PLUGIN_COMMAND: |
        test -f .github/workflows/release.yml
        git diff --quiet HEAD
```

## `SEMREL_PLUGIN_*` variables

| Name | Required | Description | Default |
| --- | --- | --- | --- |
| `SEMREL_PLUGIN_COMMAND` | Primary, required if not using env-var mode | One or more newline-separated shell commands. Each non-empty line is executed independently via `sh -c`, and all commands must exit with status 0 for the condition to pass. | None |
| `SEMREL_PLUGIN_ENV_VAR` | Legacy, optional | Environment variable name to compare for backward-compatible equality checks. Used only when `SEMREL_PLUGIN_COMMAND` is not set. | None |
| `SEMREL_PLUGIN_ENV_VALUE` | Legacy, optional | Expected value for `SEMREL_PLUGIN_ENV_VAR`. Used only when `SEMREL_PLUGIN_COMMAND` is not set. | Empty string |

## `SEMREL_*` release context used

This plugin does not consume any `SEMREL_*` release context variables directly.

## Example behavior

Semrel executes each configured command before releasing. If every non-empty command line exits with 0 the pipeline continues; otherwise the release is blocked at the first failing command. For backward compatibility, if `SEMREL_PLUGIN_COMMAND` is not configured, the plugin falls back to comparing `SEMREL_PLUGIN_ENV_VAR` against `SEMREL_PLUGIN_ENV_VALUE`.

## License

Apache-2.0
