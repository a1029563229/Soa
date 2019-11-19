package sutils

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func HasKey(m map[string]interface{}, key string) bool {
	if _, ok := m[key]; !ok {
		return false
	}
	return true
}

func Includes(A []string, val string) bool {
	for _, v := range A {
		if string(v) == val {
			return true
		}
	}
	return false
}

func Goid() int {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic recover:panic info:%v", err)
		}
	}()

	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}

func ToBson(structure interface{}) bson.M {
	result := make(bson.M)
	t := reflect.TypeOf(structure)
	v := reflect.ValueOf(structure)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.IsZero() {
			continue
		}

		tag := t.Field(i).Tag
		name := t.Field(i).Name
		key := tag.Get("bson")

		fmt.Println("kind:", field.Kind())

		switch field.Kind() {
		case reflect.Int, reflect.Int64:
			v := field.Int()
			result[key] = v
			break
		case reflect.String:
			v := field.String()
			result[key] = v
			break
		case reflect.Struct:
			v := getField(structure, name)
			result[key] = v
			break
		}
	}

	return result
}

// get struct field value
func getField(v interface{}, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	fieldValue := f.Interface()

	switch v := fieldValue.(type) {
	case int64:
		return strconv.FormatInt(v, 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int:
		return strconv.FormatInt(int64(v), 10)
	case string:
		return v
	case bool:
		if v {
			return "true"
		}
		return "false"
	case time.Time:
		return v.String()
	default:
		return ""
	}
}
