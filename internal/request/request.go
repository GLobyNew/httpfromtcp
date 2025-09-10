package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine  RequestLine
	RequestState RequestState
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.RequestState {
	case Initialized:
		rl, cb, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if cb == 0 {
			return 0, nil
		}
		r.RequestLine = *rl
		r.RequestState = Done
		return cb, nil
	case Done:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("error: unknown state")
	}
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type RequestState int

const (
	Initialized RequestState = iota
	Done
)

const SUPPORTED_HTTP_VERSION = "1.1"
const bufferSize = 8

func verifyRequestLine(requestLine *RequestLine) error {
	if requestLine.Method != strings.ToUpper(requestLine.Method) {
		return fmt.Errorf("method is not uppercase: %v", requestLine.Method)
	}

	if requestLine.HttpVersion != SUPPORTED_HTTP_VERSION {
		return fmt.Errorf("request has unsupported version: %v, but expected: %v", requestLine.HttpVersion, SUPPORTED_HTTP_VERSION)
	}
	return nil

}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return nil, 0, nil
	}
	requestLineText := string(data[:idx])
	requestLine, err := parseRequestLineString(requestLineText)
	if err != nil {
		return nil, 0, err
	}
	return requestLine, idx + 2, nil
}

func parseRequestLineString(str string) (rl *RequestLine, err error) {
	requestLineParts := strings.Split(str, " ")

	if len(requestLineParts) != 3 {
		return rl, fmt.Errorf("request has %v parts, instead of 3", len(requestLineParts))
	}

	httpVersionFull := strings.Split(requestLineParts[2], "/")
	if len(httpVersionFull) != 2 {
		return rl, fmt.Errorf("HTTP version has %v parts, instead of 2", len(httpVersionFull))
	}
	httpVersion := httpVersionFull[1]

	rl = &RequestLine{
		Method:        requestLineParts[0],
		RequestTarget: requestLineParts[1],
		HttpVersion:   httpVersion,
	}

	if err := verifyRequestLine(rl); err != nil {
		return rl, err
	}
	return rl, nil

}

func RequestFromReader(reader io.Reader) (*Request, error) {
	r := Request{
		RequestState: Initialized,
	}
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0
	for r.RequestState != Done {
		if readToIndex == cap(buf) {
			newBuf := make([]byte, cap(buf)*2, cap(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}
		n, err := reader.Read(buf[readToIndex:cap(buf)])
		if err != nil {
			if err == io.EOF {
				r.RequestState = Done
				break
			}
			return &Request{}, err
		}
		readToIndex += n
		pb, err := r.parse(buf[:readToIndex])
		if err != nil {
			return &Request{}, err
		}
		if pb > 0 {
			copy(buf, buf[pb:readToIndex])
			readToIndex -= pb
		}
	}
	return &r, nil

}
