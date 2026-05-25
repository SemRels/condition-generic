# condition-generic

Generic CI condition plugin for Semantic Release.

Validates generic local or self-hosted CI runtime conditions before a Semantic Release is executed.

## Documentation

- Docs (coming soon): <https://github.com/SemRels/semrel/tree/main/docs/plugins/condition-generic>
- Template source: <https://github.com/SemRels/plugin-template>

## Repository Layout

`	ext
cmd/plugin/              Plugin entry point
internal/plugin/         Business logic scaffold
internal/grpc/           gRPC transport scaffold
proto/v1                 Symlink to the SemRel protobuf contract
.github/workflows/       CI, release, and security automation
`

## Development

`ash
go build ./cmd/plugin
go test ./...
`

## Configuration Example

`yaml
plugins:
  - name: condition-generic
    type: condition
    config:
      require_branch: main
      require_clean_worktree: true
      environment_flag: CI
`

## Status

This repository is bootstrapped from SemRels/plugin-template and is ready for implementation.