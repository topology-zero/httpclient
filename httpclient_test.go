package httpclient

import (
	"fmt"
	"testing"
)

func TestNewHttpRequest(t *testing.T) {
	request := NewHttpRequest(GET, "https://www.112233.com")
	resp, err := request.DoHttpRequest()
	fmt.Println(string(resp))
	fmt.Println(err)
}

func TestNewHttpRequest2(t *testing.T) {
	request := NewHttpRequest(GET, "http://localhost:8857/notfound")
	resp, err := request.DoHttpRequest()
	fmt.Println(string(resp))
	fmt.Println(err)
}

func TestNewHttpRequest3(t *testing.T) {
	request := NewHttpRequest(GET, "http://localhost:8857/notcontent")
	resp, err := request.DoHttpRequest()
	fmt.Println(string(resp))
	fmt.Println(err)
}
