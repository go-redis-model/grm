package grm

import (
	"reflect"
	"testing"

	pb "github.com/go-redis-model/grm/example/protobuf/pb"
	"github.com/stretchr/testify/assert"
)

func TestSerializers(t *testing.T) {
	tests := []struct {
		name       string
		serializer Serializer
	}{
		{"JSON", JSONSerializer},
		{"MessagePack", MessagePackSerializer},
		{"Protobuf", ProtobufSerializer},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setupTestRedis()
			defer s.Close()

			db, _ := Open(
				&Options{Addr: s.Addr()},
				WithSerializer(tt.serializer),
			)

			// 测试模型（Protobuf 需特殊处理）
			var model interface{}
			switch tt.name {
			case "Protobuf":
				model = &pb.User{ID: 1, Name: "Alice"}
			default:
				model = &TestUser{ID: 1, Name: "Alice"}
			}

			err := db.Set(model)
			assert.NoError(t, err)

			fetchedRef := reflect.New(reflect.TypeOf(model).Elem())
			fetchedRef.Elem().FieldByName("ID").Set(reflect.ValueOf(uint32(1)))
			fetched := fetchedRef.Interface()
			err = db.Get(fetched)
			assert.NoError(t, err)
			assert.Equal(t, "Alice", fetchedRef.Elem().FieldByName("Name").String())
		})
	}
}
