package omihttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/stormi-li/omiv1/omihttp/serialization"
	"google.golang.org/protobuf/proto"
)

var JsonUnMarshal = func(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

var JsonMarshal = func(v any) ([]byte, error) {
	return json.Marshal(v)
}

var ProtobufMarshal = func(v any) ([]byte, error) {
	value, ok := v.(proto.Message)
	if ok {
		data, err := proto.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("序列化失败: %w", err)
		}
		return data, nil
	}
	return nil, fmt.Errorf("编码断言错误: 类型 %T 不是 proto.Message", v)
}

var ProtobufUnMarshal = func(data []byte, v any) error {
	value, ok := v.(proto.Message)
	if ok {
		err := proto.Unmarshal(data, value)
		if err != nil {
			return fmt.Errorf("反序列化失败: %w", err)
		}
		return nil
	}
	return fmt.Errorf("解码断言错误: 类型 %T 不是 proto.Message", v)
}

type Handler struct {
	ServeHTTP func(w http.ResponseWriter, r *http.Request, rw *ReadWriter)
}

type ReadWriter struct {
	w http.ResponseWriter
	r *http.Request
}

func NewReadWriter(w http.ResponseWriter, r *http.Request) *ReadWriter {
	return &ReadWriter{w: w, r: r}
}

func Write(w http.ResponseWriter, v any, sType serialization.Type) error {
	var data []byte
	var err error
	if sType == serialization.Protobuf {
		data, err = ProtobufMarshal(v)
	} else {
		data, err = JsonMarshal(v)
	}

	if err != nil {
		return err
	}
	// 写入响应
	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}

	return nil
}

func (rw *ReadWriter) Write(v any) error {
	_, ok := v.(proto.Message)
	if ok {
		return Write(rw.w, v, serialization.Protobuf)
	}
	return Write(rw.w, v, serialization.Json)
}

func Read(r *http.Request, v any, sType serialization.Type) error {
	data, err := io.ReadAll(r.Body)
	if len(data) == 0 {
		return fmt.Errorf("response body is empty")
	}
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}
	if sType == serialization.Protobuf {
		return ProtobufUnMarshal(data, v)
	} else {
		return JsonUnMarshal(data, v)
	}
}

func (rw *ReadWriter) Read(v any) error {
	_, ok := v.(proto.Message)
	if ok {
		return Read(rw.r, v, serialization.Protobuf)
	}
	return Read(rw.r, v, serialization.Json)
}

type Response struct {
	*http.Response
	SType serialization.Type
}

// OmiRead 读取响应的 Body 并解码到 v
func (response *Response) Read(v any) error {
	if response.Body == nil {
		return fmt.Errorf("response body is nil")
	}

	defer response.Body.Close()

	// 读取 Body 内容
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if len(data) == 0 {
		return fmt.Errorf("response body is empty")
	}
	if response.SType == serialization.Protobuf {
		return ProtobufUnMarshal(data, v)
	} else {
		return JsonUnMarshal(data, v)
	}
}

func Post(client *http.Client, url string, v any, sType serialization.Type) (*Response, error) {
	var data []byte
	var err error
	if sType == serialization.Protobuf {
		data, err = ProtobufMarshal(v)
	} else {
		data, err = JsonMarshal(v)
	}

	if err != nil {
		return nil, err
	}

	// 发起 POST 请求
	resp, err := client.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return &Response{Response: resp, SType: sType}, nil
}
