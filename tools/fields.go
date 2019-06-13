package tools

import (
	"encoding/json"
	"errors"
	"reflect"
)

/*
 获取struct字段名 json标签中定义的字段
*/
func GetFieldNameList(targetStruct interface{}) ([]string, error) {
	fields := make([]string, 0)

	fb, err := json.Marshal(targetStruct)
	if err != nil {
		return nil, err
	}

	mp := make(map[string]interface{})
	if err := json.Unmarshal(fb, &mp); err != nil {
		return nil, err
	}

	for field := range mp {
		fields = append(fields, field)
	}

	return fields, nil
}

// struct中的字符名
func GetFieldName(structName interface{}) ([]string, error) {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, errors.New("get struct fields name failed")
	}

	n := t.NumField()
	fieldList := make([]string, 0, n)
	for i := 0; i < n; i++ {
		fieldList = append(fieldList, t.Field(i).Name)
	}

	return fieldList, nil
}
