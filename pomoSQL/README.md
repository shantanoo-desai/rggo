# Persistent Pomodoro App

Database used is `sqlite3`. Version: 3.44.0

## SQLITE3 commanda

### Create database

```bash
sqlite3 pomo.db
```

### List existing tables

```sql
sqlite> .tables
```

### Create Table

```sql
CREATE TABLE "interval" (
    "id" INTEGER,
    "start_time" DATETIME NOT NULL,
    "planned_duration" INTEGER DEFAULT 0,
    "actual_duration" INTEGER DEFAULT 0,
    "category" TEXT NOT NULL,
    "state" INTEGER DEFAULT 1,
    PRIMARY KEY ("id")
);
```

### Insert data

```sql

INSERT INTO interval VALUES(NULL, date('now'),25,25,'Pomodoro',3);
INSERT INTO interval VALUES(NULL, date('now'),5,5,'ShortBreak',3);
INSERT INTO interval VALUES(NULL, date('now'),15,15,'LongBreak',3);
```

### Query from table

```sql
SELECT * FROM interval
```

```sql
SELECT * FROM interval WHERE category='Pomodoro';
```

```sql
DELETE FROM interval
```

## CGO

to use `github.com/mattn/go-sqlite3` compile it using CGO.

To check if CGO is enabled:

```bash
go env CGO_ENABLED
```

Should return value `1`. 

to enable permanently:

```bash
go env -w CGO_ENABLED=1
```

to enable temporarily:

```bash
export CGO_ENABLED=1
```

Install `go-sqlite3`:

```bash
go get github.com/mattn/go-sqlite3
go install github.com/mattn/go-sqlite3
```

## Building and Testing

### inmemory testing

```bash
go test -v ./... -tags=inmemory && go build -tags=inmemory
```

### sqlite3

```bash
go test -v ./... && go build
```


