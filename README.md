# RSSParser

### Setup

Clone the repo.

```
git git@github.com:Sunchasing/RSSParser.git
```

### Tests:
From the root directory:
```
go test ./...
```

### Build
```
go build
```
### Use
After import, run the `Parse` function, providing an array of URLs to it. 
The function spawn a goroutine for each URL that has been provided, 
and returns an array of `RSSItem` and an `error`.

### Example
To see an example of the parser working (internet connection needed), run main.go.
```
go run cmd/main.go
```
