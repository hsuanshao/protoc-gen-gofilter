package main

import (
	"testing"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// mockFieldDescriptor is a minimal mock of protoreflect.FieldDescriptor
type mockFieldDescriptor struct {
	protoreflect.FieldDescriptor
	kind        protoreflect.Kind
	isList      bool
	isMap       bool
	hasOptional bool
}

func (m mockFieldDescriptor) Kind() protoreflect.Kind {
	return m.kind
}

func (m mockFieldDescriptor) IsList() bool {
	return m.isList
}

func (m mockFieldDescriptor) IsMap() bool {
	return m.isMap
}

func (m mockFieldDescriptor) HasOptionalKeyword() bool {
	return m.hasOptional
}

func TestBuildZeroValueStmt(t *testing.T) {
	tests := []struct {
		name     string
		field    *protogen.Field
		expected string
	}{
		{
			name: "String field",
			field: &protogen.Field{
				GoName: "MyString",
				Desc: mockFieldDescriptor{
					kind: protoreflect.StringKind,
				},
			},
			expected: `x.MyString = ""`,
		},
		{
			name: "Int32 field",
			field: &protogen.Field{
				GoName: "MyInt",
				Desc: mockFieldDescriptor{
					kind: protoreflect.Int32Kind,
				},
			},
			expected: `x.MyInt = 0`,
		},
		{
			name: "Bool field",
			field: &protogen.Field{
				GoName: "MyBool",
				Desc: mockFieldDescriptor{
					kind: protoreflect.BoolKind,
				},
			},
			expected: `x.MyBool = false`,
		},
		{
			name: "List field",
			field: &protogen.Field{
				GoName: "MyList",
				Desc: mockFieldDescriptor{
					isList: true,
				},
			},
			expected: `x.MyList = nil`,
		},
		{
			name: "Map field",
			field: &protogen.Field{
				GoName: "MyMap",
				Desc: mockFieldDescriptor{
					isMap: true,
				},
			},
			expected: `x.MyMap = nil`,
		},
		{
			name: "Message field",
			field: &protogen.Field{
				GoName: "MyMessage",
				Desc: mockFieldDescriptor{
					kind: protoreflect.MessageKind,
				},
			},
			expected: `x.MyMessage = nil`,
		},
		{
			name: "Optional field",
			field: &protogen.Field{
				GoName: "MyOptional",
				Desc: mockFieldDescriptor{
					kind:        protoreflect.StringKind,
					hasOptional: true,
				},
			},
			expected: `x.MyOptional = nil`,
		},
		{
			name: "Oneof field",
			field: &protogen.Field{
				GoName: "MyOneof",
				Desc: mockFieldDescriptor{
					kind: protoreflect.StringKind,
				},
				Oneof: &protogen.Oneof{}, // Non-nil Oneof implies it's inside a oneof
			},
			expected: `x.MyOneof = nil`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildZeroValueStmt(tt.field)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
