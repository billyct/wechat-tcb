# wechat-tcb

小程序云开发 HTTP API SDK

## 安装

```
go get -u github.com/yyiidev/wechat-tcb
```

## 使用

```go
t := tcb.New(&tcb.Config{
    AppID: "",
    AppSecret: "",
    EnvID: "",
    Cache: c
})

...

t.UploadFile(path)
t.UploadFileWithFile(path, file)
t.DatabaseCollectionAdd(collection)
t.DatabaseAdd(query)
```