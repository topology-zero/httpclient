package httpclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type method string

const (
	GET    = method("GET")
	POST   = method("POST")
	PUT    = method("PUT")
	DELETE = method("DELETE")
)

var (
	Not200 = errors.New("status code != 200")
)

type Request struct {
	method string
	url    string
	body   io.Reader
	header map[string]string
	query  map[string]string
	retry  int

	client   *http.Client
	response *http.Response
	logger   Logger
}

type HTTPOption func(r *Request)

// NewHttpRequest 构建一个 Request 对象
func NewHttpRequest(m method, url string, opt ...HTTPOption) *Request {
	r := &Request{
		method: string(m),
		url:    url,
		client: &http.Client{},
		header: make(map[string]string),
		query:  make(map[string]string),
		logger: DefaultLog(),
	}
	for _, option := range opt {
		option(r)
	}
	return r
}

// With 添加配置项
func (r *Request) With(opt ...HTTPOption) *Request {
	for _, option := range opt {
		option(r)
	}
	return r
}

// DoHttpRequest 发送HTTP请求
func (r *Request) DoHttpRequest() ([]byte, error) {
	localRetry := 0
	do := func() ([]byte, error) {
		request, err := http.NewRequest(r.method, r.url, r.body)
		if err != nil {
			return nil, err
		}

		for k, v := range r.header {
			request.Header.Add(k, v)
		}

		if r.query != nil && len(r.query) > 0 {
			for k, v := range r.query {
				query := request.URL.Query()
				query.Add(k, v)
				request.URL.RawQuery = query.Encode()
			}
		}

		response, err := r.client.Do(request)
		if err != nil {
			return nil, err
		}
		r.response = response

		if response.StatusCode != http.StatusOK {
			return nil, Not200
		}

		resp, err := io.ReadAll(response.Body)
		defer response.Body.Close()
		if err != nil {
			return nil, err
		}

		return resp, nil
	}

	for {
		res, err := do()
		if err == nil {
			return res, nil
		}
		if localRetry == r.retry {
			return res, err
		}

		localRetry++
		r.logger.Error(fmt.Sprintf("请求[%s]失败,进行%d次重试...", r.url, localRetry))
		time.Sleep(time.Duration(200*localRetry) * time.Millisecond)
	}
}

// GetResponse 获取返回
func (r *Request) GetResponse() *http.Response {
	return r.response
}

// WithLog 设置日志输出
func WithLog(logger Logger) HTTPOption {
	return func(r *Request) {
		r.logger = logger
	}
}

// WithQuery 设置 GET 请求参数
func WithQuery(query map[string]string) HTTPOption {
	return func(r *Request) {
		r.query = query
	}
}

// WithTimeout 设置请求超时
func WithTimeout(d time.Duration) HTTPOption {
	return func(r *Request) {
		r.client.Timeout = d
	}
}

// WithRetry 设置重试次数 默认不重试
func WithRetry(retryTime int) HTTPOption {
	return func(r *Request) {
		r.retry = retryTime
	}
}

// WithJson 设置请求结构为 JSON
func WithJson(body any) HTTPOption {
	return func(r *Request) {
		marshal, _ := json.Marshal(body)
		r.body = bytes.NewBuffer(marshal)
		r.header["Content-Type"] = "application/json;charset=UTF-8"
	}
}

// WithFromData 设置请求结构为 from-urlencoded
func WithFromData(body map[string]any) HTTPOption {
	return func(r *Request) {
		r.body = strings.NewReader(Map2Str(body))
		r.header["Content-Type"] = "application/x-www-form-urlencoded;charset=UTF-8"
	}
}

// WithHeader 添加请求头
func WithHeader(key, val string) HTTPOption {
	return func(r *Request) {
		r.header[key] = val
	}
}

// WithMultipartFrom 设置请求结构为 form-data
func WithMultipartFrom(file *UploadFile, fromData map[string]string) HTTPOption {
	body := &bytes.Buffer{}
	newWriter := multipart.NewWriter(body)

	// 文件处理
	if file != nil {
		formFile, _ := newWriter.CreateFormFile(file.Field, file.FileName)
		_, _ = io.Copy(formFile, file.File)
	}

	// 普通字段处理
	if fromData != nil && len(fromData) > 0 {
		for k, v := range fromData {
			_ = newWriter.WriteField(k, v)
		}
	}

	_ = newWriter.Close()

	return func(r *Request) {
		r.body = body
		r.header["Content-Type"] = newWriter.FormDataContentType()
	}
}

func Map2Str(m map[string]any) string {
	var strArr []string
	for k, v := range m {
		strArr = append(strArr, k+"="+ConvertString(v))
	}
	return strings.Join(strArr, "&")
}

func ConvertString(value any) string {
	switch value := value.(type) {
	case string:
		return value
	case int:
		return strconv.Itoa(value)
	case int64:
		return strconv.Itoa(int(value))
	case json.Number:
		return value.String()
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	default:
		return ""
	}
}
