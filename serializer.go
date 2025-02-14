package grm

import (
	"encoding/json"
	"fmt"

	"github.com/vmihailenco/msgpack/v5"
	"google.golang.org/protobuf/proto"
)

// 内置序列化器实例
var (
	JSONSerializer        Serializer = &jsonSerializer{}
	MessagePackSerializer Serializer = &msgpackSerializer{}
	ProtobufSerializer    Serializer = &protobufSerializer{}
)

// Serializer 定义序列化接口
type Serializer interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

// 默认使用 JSON 序列化
type jsonSerializer struct{}

func (s *jsonSerializer) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (s *jsonSerializer) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// MessagePack 序列化
type msgpackSerializer struct{}

func (s *msgpackSerializer) Marshal(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}

func (s *msgpackSerializer) Unmarshal(data []byte, v interface{}) error {
	return msgpack.Unmarshal(data, v)
}

// Protobuf 序列化
type protobufSerializer struct{}

func (s *protobufSerializer) Marshal(v interface{}) ([]byte, error) {
	msg, ok := v.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("value %T is not a proto.Message", v)
	}
	return proto.Marshal(msg)
}

func (s *protobufSerializer) Unmarshal(data []byte, v interface{}) error {
	msg, ok := v.(proto.Message)
	if !ok {
		return fmt.Errorf("value %T is not a proto.Message", v)
	}
	return proto.Unmarshal(data, msg)
}
