package omihttp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func HandleFunc(pattern string, handleFunc func(w http.ResponseWriter, r *http.Request, rw ReadWriter)) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		handleFunc(w, r, rw)
	})
}

var rw = ReadWriter{}

var UnMarshalFunc = func(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

var MarshalFunc = func(v any) ([]byte, error) {
	return json.Marshal(v)
}

type ReadWriter struct{}

func (ReadWriter) Write(w http.ResponseWriter, v any) error {
	// 序列化为 MsgPack
	data, err := MarshalFunc(v)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	// 写入响应
	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}

	return nil
}

func (ReadWriter) Read(r *http.Request, v any) error {
	// 确保读取 Body 的内容
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}
	// 解码到目标对象
	if err := UnMarshalFunc(body, v); err != nil {
		return fmt.Errorf("failed to unmarshal body: %w", err)
	}

	return nil
}

type Response struct {
	*http.Response
}

// OmiRead 读取响应的 Body 并解码到 v
func (response *Response) UnMarshal(v any) error {
	if response.Body == nil {
		return fmt.Errorf("response body is nil")
	}

	defer response.Body.Close()

	// 读取 Body 内容
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := UnMarshalFunc(body, v); err != nil {
		return fmt.Errorf("failed to decode response body using msgpack: %w", err)
	}

	return nil
}
