# tiny_server
最小構成の http プロトコルの勉強用サーバー

リクエストとレスポンスの形さえ合っていれば http サーバーとして機能することを確認できる

## Request
```go
type Request struct {
	Method  string
	Path    string
	Version string
	Headers map[string]string
	Body    string
}
```

## Response
```go
type Response struct {
	Version       string
	StatusCode    int
	StatusMessage string
	Headers       map[string]string
	Body          string
}
```
