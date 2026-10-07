# 0003. Tests: next to the code, acceptance tests per ticket

Status: Accepted

## Convention
- **Module tests** follow each language's convention and are permanent:
  - Go: `foo_test.go` next to `foo.go`, run with `go test ./...`.
  - Frontend: `Foo.test.tsx` next to `Foo.tsx`, run with `npm test` (Vitest) in `ui/`.
  - Never call real services: use `httptest.Server`, fakes, or a stubbed `fetch`.
- **Acceptance tests** live in `test/acceptance/`, one file per ticket, with the build tag
  `//go:build acceptance`. They're written first, describe the ticket's behaviour end to end, and exist
  only on that ticket's branch: before merging, delete the file, moving anything worth keeping into module
  tests.

## Why
Module tests guard the code permanently. Acceptance tests pin down "done" for one ticket without piling
up slow end-to-end tests on `main`.

## How to apply
CI runs `go vet`, `go test -race ./...`, and `npm run lint`, `npm test` and `npm run build` in `ui/`.
`go test -tags acceptance ./test/acceptance/...` runs acceptance tests.
