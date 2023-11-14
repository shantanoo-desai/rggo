# todoClient

## errata

the `ListDeleteTask` in the `cmd/integration_test.go` in the book checks for no errors, whereas as
if the task is deleted, the `listAction` API should throw and error which should be checked.

## Tests

### Integration

```bash
go test -v ./cmd -tags integration
```

No-Cache

```bash
go test -v ./cmd -tags integration -count=1
```

### Unit Tests

```bash
go test -v ./cmd
```
