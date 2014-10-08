package utils

import (
	"strings"
	"testing"
)

func Test_Method(t *testing.T) {
	var (
		method string
		err    error
		txt    string
	)
	txt = `GET /asdasfd21af/asdasdasd?a=1 HTTP/1.1\r\n
    User-Agent: curl/7.22.0 (x86_64-pc-linux-gnu) libcurl/7.22.0\r\n
    Host: 127.0.0.1:1121\r\n
    Accept: */*\r\n`

	if method, err = Method(txt); method != "GET" || err != nil {
		t.Error("Method parser error, method", method)
	} else {
		t.Log("Parsed method", method)
	}

	txt = `POST /asdasfd21af/asdasdasd?a=1 HTTP/1.1\r\n`

	if method, err = Method(txt); method != "POST" || err != nil {
		t.Error("Method parser error, method", method)
	} else {
		t.Log("Parsed method", method)
	}

	// Bad HTTP header
	txt = "unknon; /daaksdj/!a=1 HTTP/1.1\r\n"
	if method, err = Method(txt); method != "" {
		t.Error("Method parser error, method", method)
	}
}

func Test_HeaderMap(t *testing.T) {
	txt := `GET /asdasfd21af/asdasdasd?a=1 HTTP/1.1\r\n
    User-Agent: curl/7.22.0 (x86_64-pc-linux-gnu) libcurl/7.22.0\r\n
    Host: 127.0.0.1:1121\r\n
    Accept: */*\r\n`
	// test HeaderMap
	hm := HeaderMap(txt)
	if _, exists := hm["Host"]; exists != true {
		t.Error("Parsed header map", hm)
	}
}

func Test_Parser(t *testing.T) {
	txt := `GET /asdasfd21af/asdasdasd?a=1 HTTP/1.1\r\n
    User-Agent: curl/7.22.0 (x86_64-pc-linux-gnu) libcurl/7.22.0\r\n
    Host: 127.0.0.1:1121\r\n
    Accept: */*\r\n`

	// test Parser
	header, err := Parser([]byte(txt))
	if header == nil || err != nil {
		t.Error("Parsed error")
	} else {
		t.Log(header.method, header.host, header.port, header.version)
	}

	if header.method != "GET" {
		t.Error("Parsed error, method error")
	}
	if strings.TrimSpace(header.host) != "127.0.0.1" {
		t.Error("Parsed error, host error")
	}
	if header.port != "1121" {
		t.Error("Parsed error, port error")
	}
	if header.version != "1.1" {
		t.Error("Parsed error, version error")
	}
}
