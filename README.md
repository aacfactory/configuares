# 概要

配置工具

## 特性

* 开放：可根据实际情况进行额外扩展。
* 多环境支持：指定激活的配置文件，并合并进根配置。
* 多配置语言：当前支持 JSON 与 YAML 。

## 安装

```go
go get github.com/aacfactory/configuares
```

## 使用

```go
path, err := filepath.Abs("./_example/json")
if err != nil {
    // handle error
    return
}

// 构建配置存储器
store := configuares.NewFileStore(path, "app", '.')

// 构建配置读取器
retriever, retrieverErr := configuares.NewRetriever(configuares.RetrieverOption{
    Active: "dev",
    Format: "JSON",
    Store:  store,
})

if retrieverErr != nil {
    // handle error
    return
}

// 获取配置内容
config, configErr := retriever.Get()
if configErr != nil {
    // handle error
    return
}


```