# condition-generic

Runs a custom shell command and passes only when the command exits successfully.

This plugin is distributed as the standalone Go binary `semrel-plugin-condition-generic`. Semrel executes the binary as a subprocess, provides plugin configuration through `SEMREL_PLUGIN_*` environment variables, provides release context through `SEMREL_*` environment variables, reads standard output, and treats exit code `0` as success and any non-zero exit code as failure. Install the binary in `~/.semrel/plugins/` or anywhere on your `$PATH`.

## Installation

```bash
go install github.com/SemRels/condition-generic/cmd/plugin@latest
```

## Configuration

```yaml
plugins:
  - name: condition-generic
    path: ~/.semrel/plugins/semrel-plugin-condition-generic
    env:
      SEMREL_PLUGIN_COMMAND: "test -f .github/workflows/release.yml"
```

## `SEMREL_PLUGIN_*` variables

| Name | Required | Description | Default |
| --- | --- | --- | --- |
| `SEMREL_PLUGIN_COMMAND` | Required | Shell command that must exit with status 0 for the condition to pass. | None |

## `SEMREL_*` release context used

This plugin does not consume any `SEMREL_*` release context variables directly.

## Example behavior

Semrel executes the configured command before releasing. If the command exits with 0 the pipeline continues; otherwise the release is blocked.

## License

Apache-2.0
