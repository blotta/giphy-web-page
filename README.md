# Giphy Web Page

A webserver in Go that returns a page with a greeter and a random gif from the Giphy API

## Requirements

This server uses the `gorilla/mux` package.
```bash
go get -u github.com/gorilla/mux
```

## Configuration
Enter a Giphy API key in the `config.json` file.

## Build

```bash
go build main.go giphyapi.go
```

## Run
```bash
./main
```