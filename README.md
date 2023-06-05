### httpclient

对 golang http 库进行封装，更优雅的调用网络接口

### 最新版本

```
go get github.com/topology-zero/httpclient@v1.0.4
```

### 示例使用

#### GET 请求

```go
query := map[string]string{
    "a": "1",
    "b": "2",
}
request := httpclient.NewHttpRequest(httpclient.GET, "https://foo.com", httpclient.WithQuery(query))
resp, err := request.DoHttpRequest()
if err != nil {
    panic(err)
}
log.Println(string(resp))
```

#### POST `JSON`请求

```go
type reqStruct struct {
    id int
    name string
}

req := reqStruct{
    id:   10086,
    name: "中国移不动",
}

request := httpclient.NewHttpRequest(httpclient.POST, "https://foo.com", httpclient.WithJson(req))
resp, err := request.DoHttpRequest()
if err != nil {
    panic(err)
}
log.Println(string(resp))
```

#### POST `FROM-DATA`请求 -- 附带文件

```go
req := httpclient.UploadFile{
    Field:    "file",
    FileName: "image.png",
    File:     bytes.NewReader([]byte{}),
}
request := httpclient.NewHttpRequest(httpclient.POST, "https://foo.com", httpclient.WithMultipartFrom(&req, nil))
resp, err := request.DoHttpRequest()
if err != nil {
    panic(err)
}
log.Println(string(resp))
```

#### POST `FROM-DATA`请求 -- 不附带文件

```go
req := map[string]string{
    "a": "1",
    "b": "2",
}
request := httpclient.NewHttpRequest(httpclient.POST, "https://foo.com", httpclient.WithMultipartFrom(nil, req))
resp, err := request.DoHttpRequest()
if err != nil {
    panic(err)
}
log.Println(string(resp))
```

#### POST `FROM-DATA`请求 -- 附带文件 + 其他字段

```go
file := httpclient.UploadFile{
    Field:    "file",
    FileName: "image.png",
    File:     bytes.NewReader([]byte{}),
}
data := map[string]string{
    "a": "1",
    "b": "2",
}
request := httpclient.NewHttpRequest(httpclient.POST, "https://foo.com", httpclient.WithMultipartFrom(&file, data))
resp, err := request.DoHttpRequest()
if err != nil {
    panic(err)
}
```

#### POST `x-www-form-urlencoded`请求

```go
req := map[string]any{
    "a": 1,
    "b": "2",
}
request := httpclient.NewHttpRequest(httpclient.POST, "https://foo.com", httpclient.WithFromData(req))
resp, err := request.DoHttpRequest()
if err != nil {
    panic(err)
}
log.Println(string(resp))
```

#### PUT 请求

```go
type reqStruct struct {
    name string
}

req := reqStruct{
    name: "中国移不动",
}
request := httpclient.NewHttpRequest(httpclient.PUT, "https://foo.com/10086", httpclient.WithJson(req))
resp, err := request.DoHttpRequest()
if err != nil {
    panic(err)
}
log.Println(string(resp))
```

#### DELETE 请求

```go
request := httpclient.NewHttpRequest(httpclient.DELETE, "https://foo.com/10086")
resp, err := request.DoHttpRequest()
if err != nil {
    panic(err)
}
log.Println(string(resp))
```

### 自定义错误日志处理

```go
request := httpclient.NewHttpRequest(httpclient.DELETE, "https://foo.com/10086", httpclient.WithLog(logrus.StandardLogger()))
resp, err := request.DoHttpRequest()
if err != nil {
    panic(err)
}
log.Println(string(resp))
```
