# convertfs
**convertfs** is a file converter (packed and unpacked eigenstrat <-> binary and plain plink).

## Usage
```
convertfs -p PATH_TO_EIGENSTRAT_WITHOUT_EXTENSION
```

## Build
```
git clone https://github.com/stasundr/go-convertfs.git convertfs
cd convertfs
go build
```

### Cross platform compile
```
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o convertfs.linux main.go
```