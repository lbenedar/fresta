package actions

import (
	"fmt"
	"reflect"
	"strings"
)

func findByTag(obj any, tag string) (int, error) {
	structType := reflect.TypeOf(obj).Elem()
	for i := 0; i < structType.NumField(); i++ {
		fieldTags := structType.Field(i).Tag.Get("json")
		if tag == strings.Split(fieldTags, ",")[0] {
			return i, nil
		}
	}
	return -1, fmt.Errorf("No such field: %s in obj", tag)
}

func setField(obj any, name string, value any) error {
	structValue := reflect.ValueOf(obj).Elem()
	idx, err := findByTag(obj, name)
	if err != nil {
		return err
	}
	structFieldValue := structValue.Field(idx)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return ErrObjTypeNotMatch
	}

	if value == nil {
		structFieldValue.Set(reflect.Zero(structFieldType))
	} else {
		structFieldValue.Set(val)
	}
	return nil
}

func isPointer(obj any) bool {
	iv := reflect.ValueOf(obj)
	if !iv.IsValid() {
		return false
	}

	return iv.Kind() == reflect.Pointer
}

func FillStruct(m map[string]any, obj any) error {
	if !isPointer(obj) {
		return ErrObjNotPointer
	}

	for k, v := range m {
		if v == nil {
			continue
		}

		err := setField(obj, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func Dereference(input any) any {
	v := reflect.ValueOf(input)

	// Check if the value is a pointer
	if v.Kind() == reflect.Pointer {
		// v.Elem() returns the value the pointer points to
		return v.Elem().Interface()
	}

	return input
}
