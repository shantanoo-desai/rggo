# todoServer

adding local modules (`todo`)

```bash
go mod init github.com/shantanoo-desai/rggo/todoServer
```

```bash
go mod edit -require=github.com/shantanoo-desai/rggo/todo@v0.0.0
```

```bash
go mod edit -replace=github.com/shantanoo-desai/rggo/todo=../todo
```
