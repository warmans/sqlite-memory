# sqlite-memory

Sqlite memory backend for: https://github.com/go-joe/joe


Example: 

```go
b := &ExampleBot{
	Bot: joe.New("example", sqlite.Memory(":memory:")),
}
```

or for a persistant DB:

```go
b := &ExampleBot{
	Bot: joe.New("example", sqlite.Memory("file:example.db")),
}
```
