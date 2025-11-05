package http

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Response struct {
	StatusCode int
	Headers    map[string]interface{}
	Body       interface{}
	RawBody    []byte
}

func parseStatusCode(statusCode interface{}) int {
	if statusCode == nil {
		return 0
	}

	str := fmt.Sprintf("%v", statusCode)
	if sc, err := strconv.Atoi(str); err == nil {
		return sc
	}
	if f, err := strconv.ParseFloat(str, 64); err == nil {
		return int(f)
	}

	return 0
}

func NewResponse(openapiResponse map[string]interface{}) *Response {
	resp := &Response{
		Headers: make(map[string]interface{}),
	}

	// Extract status code - handle interface{} type with various possible types
	if statusCode, ok := openapiResponse["statusCode"]; ok && statusCode != nil {
		if sc := parseStatusCode(statusCode); sc > 0 {
			resp.StatusCode = sc
		}
	}

	if headers, ok := openapiResponse["headers"]; ok {
		if h, ok := headers.(map[string]interface{}); ok {
			resp.Headers = h
		}
	}

	if body, ok := openapiResponse["body"]; ok {
		resp.Body = body
	}

	if rawBody, ok := openapiResponse["body"]; ok {
		switch v := rawBody.(type) {
		case []byte:
			resp.RawBody = v
		case string:
			resp.RawBody = []byte(v)
		}
	}

	return resp
}

func (r *Response) GetBodyString() string {
	if r.Body == nil {
		return ""
	}

	switch v := r.Body.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case map[string]interface{}, []interface{}:
		jsonData, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("%v", v)
		}
		return string(jsonData)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (r *Response) GetBodyJSON() ([]byte, error) {
	if r.Body == nil {
		return nil, nil
	}

	switch v := r.Body.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	case map[string]interface{}, []interface{}:
		return json.Marshal(v)
	default:
		return []byte(fmt.Sprintf("%v", v)), nil
	}
}

func (r *Response) GetHeader(key string) (string, bool) {
	if r.Headers == nil {
		return "", false
	}

	value, ok := r.Headers[key]
	if !ok {
		return "", false
	}

	if str, ok := value.(string); ok {
		return str, true
	}

	return fmt.Sprintf("%v", value), true
}

func (r *Response) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

func (r *Response) GetStatusCode() int {
	return r.StatusCode
}
