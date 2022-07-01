# Configures

Configures for Golang

## Features

* Open
* Multi environment support。
* Json
* Yaml

## Install

```go
go get github.com/aacfactory/configures
```

## Usage

```go
path, err := filepath.Abs("./_example/json")
if err != nil {
    // handle error
    return
}

store := configures.NewFileStore(path, "app", '.')

retriever, retrieverErr := configures.NewRetriever(configures.RetrieverOption{
    Active: "dev",
    Format: "JSON",
    Store:  store,
})

if retrieverErr != nil {
    // handle error
    return
}

config, configErr := retriever.Get()
if configErr != nil {
    // handle error
    return
}

```