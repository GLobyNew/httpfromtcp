package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const (
	SUPPORTED_HTTP_VERSION = "1.1"
)

func verifyRequestLine(requestLine *RequestLine) error {
	if requestLine.Method != strings.ToUpper(requestLine.Method) {
		return fmt.Errorf("method is not uppercase: %v", requestLine.Method)
	}

	if requestLine.HttpVersion != SUPPORTED_HTTP_VERSION {
		return fmt.Errorf("request has unsupported version: %v, but expected: %v", requestLine.HttpVersion, SUPPORTED_HTTP_VERSION)
	}
	return nil

}

func parseRequestLine(request string) (*RequestLine, error) {
	parts := strings.Split(request, "\r\n")
	requestLineStr := parts[0]
	requestLineParts := strings.Split(requestLineStr, " ")

	if len(requestLineParts) != 3 {
		return &RequestLine{}, fmt.Errorf("request has %v parts, instead of 3", len(requestLineParts))
	}

	httpVersionFull := strings.Split(requestLineParts[2], "/")
	if len(httpVersionFull) != 2 {
		return &RequestLine{}, fmt.Errorf("HTTP version has %v parts, instead of 2", len(httpVersionFull))
	}
	httpVersion := httpVersionFull[1]

	requestLine := RequestLine{
		Method:        requestLineParts[0],
		RequestTarget: requestLineParts[1],
		HttpVersion:   httpVersion,
	}

	if err := verifyRequestLine(&requestLine); err != nil {
		return &RequestLine{}, err
	}
	return &requestLine, nil

}

func RequestFromReader(reader io.Reader) (*Request, error) {
	r, err := io.ReadAll(reader)
	if err != nil {
		return &Request{}, err
	}
	strReq := string(r)

	reqLine, err := parseRequestLine(strReq)
	if err != nil {
		return &Request{}, err
	}

	request := Request{
		RequestLine: *reqLine,
	}
	return &request, nil

}
