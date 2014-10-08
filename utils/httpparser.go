package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type HttpParser struct {
	method  string
	host    string
	port    string
	path    string
	version string
}

const CLRF = string("\r\n")

var (
	// Parse `GET /url?a=1 HTTP/1.1`
	RegMethod, _ = regexp.Compile(`([A-Z]+)\s`)
	// Parse `CONNECT host.com:8080/url?a=1 HTTP/1.1`
	RegConn, _ = regexp.Compile(`([A-Z]+)\s([^\:\s]+)\:(\d+)\sHTTP\/(\d\.\d)`)
	// Parse `GET /url?a=1 HTTP/1.1`
	RegNormal, _ = regexp.Compile(`([A-Z]+)\s([^\s]+)\sHTTP\/(\d\.\d)`)
	// Parse `Host: host.com:8080`
	RegHost, _ = regexp.Compile(`Host\:\s+([^\n\s\r]+)`)
)

var (
	HEADERROR = errors.New("Unknown HTTP header format")
)

func Method(header string) (string, error) {
	elements := strings.Split(header, CLRF)
	if len(elements) <= 0 {
		return "", HEADERROR
	}

	re := RegMethod.FindAllStringSubmatch(elements[0], -1)
	if len(re) < 1 {
		return "", HEADERROR
	}
	return re[0][1], nil
}

func HeaderMap(header string) map[string]string {
	frames := strings.Split(header, CLRF)
	header_map := make(map[string]string)
	for _, frame := range frames {
		kv := strings.SplitN(frame, ":", 2)
		if len(kv) != 2 {
			continue
		}
		header_map[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}
	return header_map
}

func Parser(header []byte) (*HttpParser, error) {
	var (
		method  string
		host    string
		port    string
		path    string
		version string
		err     error
	)

	msg := string(header)

	method, err = Method(msg)
	if method == "" || err != nil {
		return nil, HEADERROR
	}

	header_map := HeaderMap(msg)

	if method == "CONNECT" {
		re := RegConn.FindAllStringSubmatch(msg, -1)
		if len(re) < 5 {
			return nil, HEADERROR
		}
		line := re[0]
		host = line[2]
		port = line[3]
		version = line[4]

		if host != "" && port != "" && version != "" {
			return &HttpParser{
				method:  method,
				host:    host,
				port:    port,
				path:    "",
				version: version}, nil
		}
	} else {
		re := RegNormal.FindAllStringSubmatch(msg, -1)
		if len(re) != 1 {
			fmt.Printf("unknow http header %q\n", re)
			return nil, HEADERROR
		}
		line := re[0]

		if len(line) != 4 {
			return nil, HEADERROR
		}

		path = line[2]
		version = line[3]

		value, exists := header_map["Host"]
		if exists != true {
			return nil, HEADERROR
		}

		host_port := strings.Split(value, ":")
		lhp := len(host_port)
		if lhp == 1 {
			host = host_port[0]
			port = "80"
		} else if lhp == 2 {
			host = host_port[0]
			port = host_port[1]
		} else {
			return nil, HEADERROR
		}

		if path != "" && host != "" && version != "" {
			return &HttpParser{
				method:  method,
				host:    host,
				port:    port,
				path:    path,
				version: version}, nil
		}
	}

	return nil, HEADERROR
}
