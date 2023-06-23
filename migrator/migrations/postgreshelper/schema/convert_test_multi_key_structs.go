// Code generated by pg-bindings generator. DO NOT EDIT.
package schema

import (
	"github.com/lib/pq"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres/pgutils"
)

// ConvertTestOneKeyStructFromProto converts a `*storage.TestOneKeyStruct` to Gorm model
func ConvertTestOneKeyStructFromProto(obj *storage.TestOneKeyStruct) (*TestOneKeyStructs, error) {
	serialized, err := obj.Marshal()
	if err != nil {
		return nil, err
	}
	model := &TestOneKeyStructs{
		Key1:              obj.GetKey1(),
		Key2:              obj.GetKey2(),
		StringSlice:       pq.Array(obj.GetStringSlice()).(*pq.StringArray),
		Bool:              obj.GetBool(),
		Uint64:            obj.GetUint64(),
		Int64:             obj.GetInt64(),
		Float:             obj.GetFloat(),
		Labels:            obj.GetLabels(),
		Timestamp:         pgutils.NilOrTime(obj.GetTimestamp()),
		Enum:              obj.GetEnum(),
		Enums:             pq.Array(pgutils.ConvertEnumSliceToIntArray(obj.GetEnums())).(*pq.Int32Array),
		String:            obj.GetString_(),
		Int32Slice:        pq.Array(obj.GetInt32Slice()).(*pq.Int32Array),
		OneofnestedNested: obj.GetOneofnested().GetNested(),
		Serialized:        serialized,
	}
	return model, nil
}

// ConvertTestOneKeyStruct_NestedFromProto converts a `*storage.TestOneKeyStruct_Nested` to Gorm model
func ConvertTestOneKeyStruct_NestedFromProto(obj *storage.TestOneKeyStruct_Nested, idx int, TestOneKeyStructKey1 string, TestOneKeyStructKey2 string) (*TestOneKeyStructsNesteds, error) {
	model := &TestOneKeyStructsNesteds{
		TestOneKeyStructsKey1: TestOneKeyStructKey1,
		TestOneKeyStructsKey2: TestOneKeyStructKey2,
		Idx:                     idx,
		Nested:                  obj.GetNested(),
		IsNested:                obj.GetIsNested(),
		Int64:                   obj.GetInt64(),
		Nested2Nested2:          obj.GetNested2().GetNested2(),
		Nested2IsNested:         obj.GetNested2().GetIsNested(),
		Nested2Int64:            obj.GetNested2().GetInt64(),
	}
	return model, nil
}

// ConvertTestOneKeyStructToProto converts Gorm model `TestOneKeyStructs` to its protobuf type object
func ConvertTestOneKeyStructToProto(m *TestOneKeyStructs) (*storage.TestOneKeyStruct, error) {
	var msg storage.TestOneKeyStruct
	if err := msg.Unmarshal(m.Serialized); err != nil {
		return nil, err
	}
	return &msg, nil
}
