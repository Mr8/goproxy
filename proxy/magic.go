package main

import (
	"bytes"
	"github.com/Mr8/goproxy/utils"
	"regexp"
	"strings"
)

var (
	RegProxy      = regexp.MustCompile(`Proxy\-Connection\:.+\r\n`)
	RegKeepAlived = regexp.MustCompile(`Keep\-Alive\:.+\r\n`)
	RegHttpHost   = regexp.MustCompile(`(?i:http)\:\/\/[^\/]+/`)
)

// Method which used to transfer
// HTTP Header with Proxy position to an normal HTTP request
func TransferHTTP(requests []byte, hp *utils.HttpParser) ([]byte, error) {
	lenReq := len(requests)
	posBody := hp.BodyLen
	if posBody < 0 || posBody > lenReq {
		posBody = len(requests) - 1
	}

	// Get HTTP Header string
	header := string(requests[0:posBody])
	lenHeaderOri := len(header)

	// Remove Proxy-Connection segment
	header = RegProxy.ReplaceAllString(header, "")
	// Remove Keep-Alived segment
	header = RegKeepAlived.ReplaceAllString(header, "")
	// Replace absolutly URL to relative URL
	url := RegHttpHost.ReplaceAllString(hp.Path, "")
	header = strings.Replace(header, hp.Path, url, 1)
	// Add Connection: close segment
	header = strings.Replace(header, "\r\n", "\r\nConnection: close\r\n", 1)

	lenHeaderTran := len(header)
	// retBuf := make([]byte, lenReq-lenHeaderOri+lenHeaderTran)

	retBuf := bytes.NewBuffer([]byte{})
	retBuf.WriteString(header)
	retBuf.Write(requests[posBody:])
	return retBuf.Next(lenReq - lenHeaderOri + lenHeaderTran), nil
}
