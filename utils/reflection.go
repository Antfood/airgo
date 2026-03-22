package utils

import (
	"fmt"
	"reflect"
	"strings"
)

/*
IsStruct checks if the given interface is a struct.

Example:

	type person struct {
	    Name string
	    Age  int
	}

	p := person{"John", 20}
	isPStruct := IsStruct(p)

	// Output: isPStruct = true
*/
func IsStruct(s any) bool {
	return reflect.TypeOf(s).Kind() == reflect.Struct
}

/*
StructToMap converts a struct to a map using the struct's field names as keys.

Example:

	type person struct {
	    Name string
	    Age  int
	}

	p := person{"John", 20}
	m, _ := StructToMap(p)

	// Output: m = {"Name": "John", "Age": 20}
*/
func StructToMap(s any) (map[string]any, error) {
	r := make(map[string]any)

	if !IsStruct(s) {
		return r, fmt.Errorf("StructToMap expects a struct, got '%v' -> %v",
			GetTypeName(s),
			reflect.ValueOf(s).Kind())
	}

	val := reflect.ValueOf(s)
	t := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := t.Field(i)
		fieldName := field.Name
		fieldVal := val.Field(i).Interface()
		r[fieldName] = fieldVal
	}

	return r, nil
}

/* TODO: Docs & Test */


type StructJsonOptions struct {
	IgnoreUpdateTag bool
}

func WithIgnore() StructJsonOptions {
	return StructJsonOptions{
		IgnoreUpdateTag: true,
	}
}

func WithoutIgnore() StructJsonOptions {
	return StructJsonOptions{
		IgnoreUpdateTag: false,
	}
}

func StructJsonToMap(s any, opts StructJsonOptions) (map[string]any, error) {
	result := make(map[string]any)

	if !IsStruct(s) {
		return result, fmt.Errorf("StructToMap expects a struct, got '%v' -> %v",
			GetTypeName(s),
			reflect.ValueOf(s).Kind())
	}

	val := reflect.ValueOf(s)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)

	   updateTag := field.Tag.Get("update")

		if updateTag == "ignore" && opts.IgnoreUpdateTag {
			continue
		}

		jsonTag := field.Tag.Get("json")

		if jsonTag == "" || jsonTag == "-" {
			jsonTag = field.Name
		} else {
			if commaIdx := strings.Index(jsonTag, ","); commaIdx > 0 {
				jsonTag = jsonTag[:commaIdx]
			}
		}
		fieldValue := val.Field(i).Interface()
		result[jsonTag] = fieldValue
	}

	return result, nil
}

/* TODO: Docs & Test */

func GetStructFieldNames[T any](s T) ([]string, error) {

	v := reflect.ValueOf(s)

	if !v.IsValid() {
		return nil, fmt.Errorf("GetStructFieldNames received zero Value")
	}

	if v.Kind() == reflect.Ptr {

		if v.IsNil() {
			return nil, fmt.Errorf("GetStructFieldNames received nil pointer")
		}

		v = v.Elem()
	}

	if !IsStruct(v.Interface()) {
		return nil, fmt.Errorf("GetStructFieldNames expects a struct, got '%v' -> %v",
			GetTypeName(v.Interface()), v.Kind())
	}

	t := v.Type()

	fieldNames := make([]string, 0, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldNames = append(fieldNames, field.Name)
	}

	return fieldNames, nil
}

/*
GetStructFieldJsonNames extracts field names from a struct that are annotated with JSON tags.
This function intentionally includes fields that have JSON annotations and ignores those without them.

Example:

	type person struct {
	    Name string `json:"name"`
	    Age  int    `json:"age"`
	}

	p := person{"John", 20}
	fields, _ := GetStructFieldJsonNames(p)

	// Output: fields = ["name", "age"]
*/
func GetStructFieldJsonNames[T any](s T) ([]string, error) {
	var fieldNames []string

	v := reflect.ValueOf(s)

	if !v.IsValid() {
		return fieldNames, fmt.Errorf("GetStructFieldJsonNames received zero Value")
	}

	if v.Kind() == reflect.Ptr {

		if v.IsNil() {
			return fieldNames, fmt.Errorf("GetStructFieldJsonNames received nil pointer")
		}

		v = v.Elem()
	}

	if !IsStruct(v.Interface()) {
		return fieldNames, fmt.Errorf("GetStructFieldJsonNames expects a struct, got '%v' -> %v",
			GetTypeName(v.Interface()),
			v.Kind())
	}

	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if name := field.Tag.Get("json"); name != "" {
			fieldNames = append(fieldNames, name)
		}
	}

	return fieldNames, nil
}

/*
GetTypeName returns the name of the type of the given object.

Example:

	type person struct {
	    Name string
	    Age  int
	}

	p := person{"John", 20}
	typeName := GetTypeName(p)

	// Output: typeName = "person"
*/

func GetTypeName[T any](obj T) string {
	t := reflect.TypeOf(obj)
	return t.Name()
}

/*
GetStructFieldValueByName returns the value of a struct field by its name.

Example:

   type person struct {
      Name string
      Age  int
   }

   p := person{"John", 20}

   name := GetStructFieldValueByName(p, "Name")
   age := GetStructFieldValueByName(p, "Age")
*/

func GetStructFieldValueByName(s any, fieldName string) any {

	if s == nil {
		return nil
	}

	v := reflect.ValueOf(s)

	if !v.IsValid() {
		return nil
	}

	if v.Kind() == reflect.Ptr {

		if v.IsNil() {
			return fmt.Errorf("GetStructFieldValueByName received nil pointer")
		}

		v = v.Elem()
	}

	if !IsStruct(v) {
		return nil
	}

	field := v.FieldByName(fieldName)

	if !field.IsValid() {
		return nil
	}

	return field.Interface()
}

/*
GetSliceStructValueByJsonName is a special function that will only accept a struct field that is a slice of structs.
It will then return a new instance of the inner struct type.

Example:

   type person struct {
      Name string `json:"MyName"`
      Age  int    `json:"MyAge"`
   }

   type people struct {
      People []person `json:"MyPeople"`
   }

   p := people{[]person{{"John", 20}, {"Jane", 30}}}
   s := GetSliceStructValueByJsonName(p, "MyPeople")

   // Output: s = person{"", 0}
*/

func GetSliceStructValueByJsonName(s any, fieldName string) any {

	v := reflect.ValueOf(s)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("json")

		if tag == fieldName {
			fieldVal := v.Field(i)
			if fieldVal.Kind() != reflect.Slice {
				return nil
			}

			elemType := field.Type.Elem()

			if elemType.Kind() == reflect.Ptr {
				elemType = elemType.Elem()
			}

			if elemType.Kind() != reflect.Struct {
				return nil
			}

			return reflect.New(elemType).Elem().Interface()

		}
	}
	return nil
}

/*
GetStructFieldValueByJsonName returns the value of a struct field by its json name.

Example:

   type person struct {
      Name string `json:"MyName"`
      Age  int    `json:"MyAge"`
   }

   p := person{"John", 20}

   name := GetStructFieldValueByJsonName(p, "MyName")
   age := GetStructFieldValueByJsonName(p, "MyAge")
*/

func GetStructFieldValueByJsonName(s any, jsonName string) interface{} {
	if s == nil {
		return nil
	}

	v := reflect.ValueOf(s)
	t := v.Type()

	if !v.IsValid() {
		return nil
	}

	if v.Kind() == reflect.Ptr {

		if v.IsNil() {
			return fmt.Errorf("GetStructFieldValueByJsonName received nil pointer")
		}

		v = v.Elem()
		t = v.Type()
	}

	if !IsStruct(v) {
		return nil
	}

	for i := 0; i < v.NumField(); i++ {
		// Check if this field's json tag matches the desired jsonName
		tag := t.Field(i).Tag.Get("json")
		if tag == jsonName || strings.Split(tag, ",")[0] == jsonName { // Handle cases with options, e.g., `json:"name,omitempty"`
			return v.Field(i).Interface()
		}
	}

	return nil
}

/*
HasStructEmptyFields checks if any field in a struct is empty (zero value for its type) and returns the name of the first empty field found.

Example:

	type person struct {
	    Name string
	    Age  int
	}

	p := person{"", 20}
	hasEmpty, fieldName := HasStructEmptyFields(p)

	// Output: hasEmpty = true, fieldName = "Name"
*/
func HasStructEmptyFields(s any) (bool, string) {

	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)

	if v.Kind() != reflect.Struct {
		return false, ""
	}

	for i := 0; i < v.NumField(); i++ {

		empty := reflect.Zero(v.Field(i).Type()).Interface()

		if reflect.DeepEqual(v.Field(i).Interface(), empty) {
			return true, t.Field(i).Name
		}
	}

	return false, ""
}
