package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var UnMarshalFunc = func(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

var MarshalFunc = func(v any) ([]byte, error) {
	return json.Marshal(v)
}

type ReadWriter struct {
	UnMarshalFunc func(data []byte, v any) error
	MarshalFunc   func(v any) ([]byte, error)
}

func (rw *ReadWriter) Write(w http.ResponseWriter, v any) error {
	// 序列化为 MsgPack
	data, err := rw.MarshalFunc(v)
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

func (rw *ReadWriter) Read(r *http.Request, v any) error {
	// 确保读取 Body 的内容
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}
	// 解码到目标对象
	if err := rw.UnMarshalFunc(body, v); err != nil {
		return fmt.Errorf("failed to unmarshal body: %w", err)
	}

	return nil
}
