# 🐮MS

## Dev

To run the server locally while developing with automatic server restarts when updating files:

```sh
reflex -r '\.(go|html)' -s go run cmd/cms/main.go start
```

To create the local database:

```sh
go run cmd/cms/main.go db-init
```
