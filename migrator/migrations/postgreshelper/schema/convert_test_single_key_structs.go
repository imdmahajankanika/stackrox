// Code generated by pg-bindings generator. DO NOT EDIT.
package schema

import (
	"github.com/lib/pq"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres/pgutils"
	"github.com/stackrox/rox/pkg/postgres/schema"
)

// ConvertTestSingleKeyStructFromProto converts a `*storage.TestSingleKeyStruct` to Gorm model
func ConvertTestSingleKeyStructFromProto(obj *storage.TestSingleKeyStruct) (*schema.TestSingleKeyStructs, error) {
	serialized, err := obj.MarshalVT()
	if err != nil {
		return nil, err
	}
	model := &schema.TestSingleKeyStructs{
		Key:         obj.GetKey(),
		Name:        obj.GetName(),
		StringSlice: pq.Array(obj.GetStringSlice()).(*pq.StringArray),
		Bool:        obj.GetBool(),
		Uint64:      obj.GetUint64(),
		Int64:       obj.GetInt64(),
		Float:       obj.GetFloat(),
		Labels:      obj.GetLabels(),
		Timestamp:   pgutils.NilOrTime(obj.GetTimestamp()),
		Enum:        obj.GetEnum(),
		Enums:       pq.Array(pgutils.ConvertEnumSliceToIntArray(obj.GetEnums())).(*pq.Int32Array),
		Serialized:  serialized,
	}
	return model, nil
}

// ConvertTestSingleKeyStructToProto converts Gorm model `TestSingleKeyStructs` to its protobuf type object
func ConvertTestSingleKeyStructToProto(m *schema.TestSingleKeyStructs) (*storage.TestSingleKeyStruct, error) {
	var msg storage.TestSingleKeyStruct
	if err := msg.UnmarshalVT(m.Serialized); err != nil {
		return nil, err
	}
	return &msg, nil
}
