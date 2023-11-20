# Notify

## listing files

```bash
go list -f '{{.GoFiles }}' ./...
```

```bash
go list -tags=inmemory -f '{{ .GoFiles }}' ./...
```

```bash
go list -tags=containers -f '{{ .GoFiles }}' ./...
```

## list distributions

```bash
go tool dist list
```

```bash
go en GOOS
```

```bash
go env GOARCH
```

## Static / Dynamic Linking

### Static Linking

```bash
CGO_ENABLED=0 go build
```

```bash
file pomo
```

### Dynamic Linking

```bash
CGO_ENABLED=1 go build
```

```bash
file pomo
```

## Cross-Compilation

```bash
GOOS=windows GOARCH=arm64 go build -tags=inmemory
```


## Strip Debug Symbols

```bash
go build -ldflags="-s -w" -tags=containers
```

## Containerized Application

```bash
docker build -t pomo/pomo:latest -f containers/Dockerfile .
```

```bash
docker run --rm -it localhost/pomo/pomo
```

```bash
docker build -t pomo/pomo:latest -f container/Dockerfile.builder .
```

### Scratch Container

```bash
docker build -t pomo/pomo:latest -f container/Dockerfile.scratch
```
```
