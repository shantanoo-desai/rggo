# colStats

## Testing

```bash
go test -v
```

## Benchmarking

```bash
go test -bench .
```

```bash
go test -bench . -benchtime 10x
```

## CPU Profiling

Generating file:

```bash
go test -bench . -benchtime 10x -cpuprofile cpu00.pprof
```

```bash
go tool pprof cpu00.pprof
```

### pprof commands

- `top`
- `top -cum`
- `list <function>`
- `web`
- `quit`

## Memory Profiling

Generating file:

```bash
go test -bench . -benchtime 10x -memprofile mem00.pprof
```
