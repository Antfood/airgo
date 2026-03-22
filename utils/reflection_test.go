package utils

import (
	"testing"

	. "github.com/Antfood/airgo/testutils/testutils"
)

type someStruct struct {
	Name string `json:"Name"`
	Age  int    `json:"Age"`
}

func TestUtilsReflection(t *testing.T) {
	t.Run("StructJsonToMap", testStructJsonToMap)
	t.Run("IsStruct", testIsStruct)
	t.Run("StructToMap", testStructToMap)
	t.Run("GetStructFieldJsonNames", testGetStructFieldJsonNames)
	t.Run("GetStructFieldNames", testGetStructFieldNames)
	t.Run("GetTypeName", testGetTypeName)
	t.Run("GetStructFieldValueByName", testGetStructFieldValueByName)
	t.Run("GetStructFieldValueByJsonName", testGetStructValueByJsonName)
	t.Run("HasStructEmptyFields", testHasStructEmptyFields)
}

func testStructJsonToMap(t *testing.T) {
	type testStruct struct {
		Name     string `json:"name"`
		Age      int    `json:"age"`
		Email    string `json:"email,omitempty"`
		Password string `json:"-"`
		NoTag    string
		Updated  string `json:"updated" update:"ignore"`
	}

	s := testStruct{
		Name:     "John",
		Age:      30,
		Email:    "john@example.com",
		Password: "secret",
		NoTag:    "value",
		Updated:  "2024-01-01",
	}

	// Test without IgnoreUpdateTag
	result, err := StructJsonToMap(s, StructJsonOptions{IgnoreUpdateTag: false})
	Ok(t, err)
	Equals(t, "John", result["name"])
	Equals(t, 30, result["age"])
	Equals(t, "john@example.com", result["email"])
	Equals(t, "value", result["NoTag"])
	Equals(t, "2024-01-01", result["updated"])
	_, hasPassword := result["-"]
	Assert(t, !hasPassword, "Password field with json:\"-\" should use field name")

	// Test with IgnoreUpdateTag
	result, err = StructJsonToMap(s, StructJsonOptions{IgnoreUpdateTag: true})
	Ok(t, err)
	_, hasUpdated := result["updated"]
	Assert(t, !hasUpdated, "Field with update:\"ignore\" should be skipped")

	// Test with non-struct
	_, err = StructJsonToMap("not a struct", StructJsonOptions{})
	Assert(t, err != nil, "Expected error for non-struct input")
}

func testIsStruct(t *testing.T) {
	s := someStruct{"John", 20}
	Assert(t, IsStruct(s), "Expected IsStruct to return true")

	slice := []someStruct{s}
	Assert(t, !IsStruct(slice), "Expected IsStruct to return false")
}

func testStructToMap(t *testing.T) {

	s := someStruct{"John", 20}

	result, err := StructToMap(s)
	Ok(t, err)

	Assert(t, result["Name"] == "John", "Expected Name to be John but found %s", result["Name"])
	Assert(t, result["Age"] == 20, "Expected Age to be 20 but found %d", result["Age"])
}

func testGetStructFieldJsonNames(t *testing.T) {

	s := someStruct{"John", 20}

	result, err := GetStructFieldJsonNames(s)
	Ok(t, err)

	Equals(t, []string{"Name", "Age"}, result)
}

func testGetStructFieldNames(t *testing.T) {

	s := someStruct{"John", 20}

	result, err := GetStructFieldNames(s)
	Ok(t, err)

	Equals(t, []string{"Name", "Age"}, result)
}

func testGetTypeName(t *testing.T) {
	s := someStruct{"John", 20}

	result := GetTypeName(s)
	Assert(t, result == "someStruct", "Expected type name to be someStruct but found %s", result)
}

func testHasStructEmptyFields(t *testing.T) {
	s := someStruct{"", 10}
	ok, empty := HasStructEmptyFields(s)
	Assert(t, ok, "Expected HasStructEmptyFields to return true")
	Assert(t, empty == "Name", "Expected empty field to be Name but found %s", empty)
}

func testGetStructFieldValueByName(t *testing.T) {

	type someStruct struct {
		Name   string
		Age    int
		Nested struct {
			Details string
		}
	}

	s := someStruct{"John", 20, struct{ Details string }{"Some Details"}}

	result := GetStructFieldValueByName(s, "Name")

	Equals(t, "John", result)

	result = GetStructFieldValueByName(s, "Nested")
	Equals(t, struct{ Details string }{Details: "Some Details"}, result)

	// must handle both values and pointers
	ptr := &someStruct{"John", 20, struct{ Details string }{"Some Details"}}

	result = GetStructFieldValueByName(ptr, "Nested")
	Equals(t, struct{ Details string }{Details: "Some Details"}, result)

	result = GetStructFieldValueByName(nil, "Name")
	Assert(t, result == nil, "Expected result to be nil but found %v", result)

	result = GetStructFieldValueByName(s, "NotAField")
	Assert(t, result == nil, "Expected result to be nil but found %v", result)
}

func testGetStructValueByJsonName(t *testing.T) {

	type nested struct {
		Details string `json:"myDetails"`
	}

	type someStruct struct {
		Name   string `json:"myName"`
		Age    int    `json:"myAge"`
		Nested nested `json:"myNested"`
	}

	s := someStruct{Name: "John",
		Age:    20,
		Nested: nested{Details: "Some Details"},
	}

	result := GetStructFieldValueByJsonName(s, "myName")
	Equals(t, "John", result)

	result = GetStructFieldValueByJsonName(s, "myNested")
	Equals(t, nested{Details: "Some Details"}, result)

}
