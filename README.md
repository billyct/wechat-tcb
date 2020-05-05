# wechat-tcb

小程序云开发 HTTP API SDK

## 安装

```
go get -u github.com/yyiidev/wechat-tcb
```

## 使用

设置环境变量
```
APP_ID=小程序的 app_id
APP_SECRET=小程序的 app_secret
```

```go
t := tcb.New()
...

t.UploadFile(key, file)
t.DatabaseCollectionAdd(collection)
t.DatabaseAdd(query)
```