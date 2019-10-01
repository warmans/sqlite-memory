<h1 align="center">Joe Bot - Sqlite Storage</h1>
<p align="center">Sqlite memory adapater. https://github.com/go-joe/joe</p>
<p align="center">
	<a href="https://circleci.com/gh/warmans/sqlite-memory/tree/master"><img src="https://circleci.com/gh/warmans/sqlite-memory/tree/master.svg?style=shield"></a>
	<a href="https://goreportcard.com/report/github.com/warmans/sqlite-memory"><img src="https://goreportcard.com/badge/github.com/warmans/sqlite-memory"></a>
	<a href="https://godoc.org/github.com/warmans/sqlite-memory"><img src="https://img.shields.io/badge/godoc-reference-blue.svg?color=blue"></a>
	<a href="https://github.com/warmans/sqlite-memory/blob/master/LICENSE"><img src="https://img.shields.io/badge/license-BSD--3--Clause-blue.svg"></a>
</p>

Sqlite memory backend for: https://github.com/go-joe/joe


### Examples 

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
