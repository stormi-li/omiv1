package omirpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/protobuf/proto"
)

var DefaultUnMarshal = func(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

var DefaultMarshal = func(v any) ([]byte, error) {
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

func Marshal(v any) ([]byte, error) {
	data, err := ProtobufMarshal(v)
	if err != nil {
		data, err = DefaultMarshal(v)
	}
	return data, err
}

func Unmarshal(data []byte, v any) error {
	err := ProtobufUnMarshal(data, v)
	if err != nil {
		err = DefaultUnMarshal(data, v)
	}
	return err
}

func Write(w http.ResponseWriter, v any) error {
	// 序列化为 MsgPack
	data, err := Marshal(v)
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

func Read(r *http.Request, v any) error {
	// 确保读取 Body 的内容
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}
	return Unmarshal(body, v)
}

type Response struct {
	*http.Response
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

	return Unmarshal(data, v)
}

func Call(client *http.Client, url string, v any) (*Response, error) {
	// 将 v 序列化为 JSON 数据
	data, err := Marshal(v)
	if err != nil {
		return nil, err
	}

	// 发起 POST 请求
	resp, err := client.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return &Response{Response: resp}, nil
}
