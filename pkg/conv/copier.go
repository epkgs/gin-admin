package conv

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
)

func NewOptions() *Options {
	return &Options{
		SkipFields:     []string{},
		SkipEqual:      false,
		KeepNil:        false,
		GetFieldByName: GetFieldByName,
	}
}

type Options struct {
	SkipFields     []string // 忽略的字段
	SkipEqual      bool     // 忽略相等字段
	KeepNil        bool     // nil 值也进行赋值
	GetFieldByName func(val reflect.Value, name string) (found reflect.Value, fieldName string)
}

func (opt *Options) ShouldSkipField(field string) bool {
	for _, fieldName := range opt.SkipFields {
		if fieldName == field {
			return true
		}
	}
	return false
}

func (opt *Options) ShouldSkipNil(val reflect.Value) bool {

	k := val.Kind()
	switch k {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		// 跳过 src nil 属性
		if !opt.KeepNil && val.IsNil() {
			return true
		}
	}

	return false
}

// GetFieldByName 尝试通过字段名获取结构体中的字段值。
// 如果字段名直接匹配或通过解析JSON标签匹配到相应的字段，则返回该字段的reflect.Value。
// 如果未找到匹配的字段，则返回一个无效的reflect.Value。
// 此函数主要用于反射操作，帮助动态访问结构体字段。
//
// 参数:
//
//	val reflect.Value: 代表结构体实例的反射值。
//	name string: 字段名，可以是结构体字段的直接名称或JSON标签名称。
//
// 返回值:
//
//	reflect.Value: 如果找到匹配的字段，则返回字段的反射值；否则返回一个无效的反射值。
func GetFieldByName(val reflect.Value, name string) (found reflect.Value, fieldName string) {
	// 首先尝试直接通过字段名获取字段值。
	if field := val.FieldByName(name); field.IsValid() {
		return field, name
	}

	// 如果直接获取失败，则尝试通过JSON标签名称获取字段值。
	return GetFieldByJsonName(val, name)
}

func GetFieldByJsonName(val reflect.Value, name string) (found reflect.Value, fieldName string) {

	return GetFieldByTag(val, func(field *reflect.StructField) bool {
		tag := field.Tag.Get("json")
		tags := strings.Split(tag, ",")
		if len(tags) == 0 {
			return false
		}

		return tags[0] == name
	})
}

func GetFieldByTag(
	val reflect.Value,
	matchFunc func(field *reflect.StructField) bool,
) (found reflect.Value, fieldName string) {
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)

		// 匿名字段
		if field.Anonymous {
			if embedField, name := GetFieldByTag(val.Field(i), matchFunc); embedField.IsValid() {
				return embedField, name
			}
		}

		if !matchFunc(&field) {
			continue
		}

		return val.Field(i), field.Name
	}

	return reflect.Value{}, ""
}

func getJsonName(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	tags := strings.Split(tag, ",")
	if len(tags) == 0 {
		return ""
	}

	if tags[0] == "-" {
		return ""
	}

	return tags[0]
}

type Metadata struct {
	SrcKeys []string
	DstKeys []string
}

// Copy 将一个结构体实例的数据传输到另一个结构体实例中。(浅拷贝)。
// 该函数主要用于在不同结构体之间复制数据，通常用于数据转换或对象映射。
//
// src: 数据传输对象，接受 struct | map
//
// dst: 目标对象的指针，接受 struct | map 的指针
//
// options: 复制的选项。
//
// 返回值 fields 是传输的字段名列表，err 是可能发生的错误。
func Copy[T any](src any, dst *T, options ...func(opt *Options)) (*Metadata, error) {

	option := NewOptions()

	for _, fn := range options {
		fn(option)
	}

	// 获取 src 的反射值
	srcVal := reflect.ValueOf(src)
	// 解包 src 直到获取结构体实例的反射值
	for srcVal.Kind() == reflect.Ptr {
		srcVal = srcVal.Elem()
	}

	// 获取 dst 的反射值
	dstVal := reflect.ValueOf(dst)
	// 确保 dst 是一个指针
	if dstVal.Kind() != reflect.Ptr {
		return nil, errors.New("dst must be a pointer")
	}

	// 解包 dst 直到获取结构体实例的反射值
	for dstVal.Kind() == reflect.Ptr {
		dstVal = dstVal.Elem()
	}

	switch srcVal.Kind() {
	case reflect.Struct:
		switch dstVal.Kind() {
		case reflect.Struct:
			return struct2struct(srcVal, dstVal, option)
		case reflect.Map:
			return struct2map(srcVal, dstVal, option)
		default:
			return nil, errors.New("dst must be a pointer to a struct or a map")
		}
	case reflect.Map:

		switch dstVal.Kind() {
		case reflect.Struct:
			return map2struct(srcVal, dstVal, option)
		case reflect.Map:
			return map2map(srcVal, dstVal, option)
		default:
			return nil, errors.New("dst must be a pointer to a struct or a map")
		}
	default:
		return nil, errors.New("src must be a struct or map")
	}
}

func Assign[T any](dst *T, src any, options ...func(opt *Options)) (*Metadata, error) {
	return Copy(src, dst, options...)
}

func map2map(srcVal, dstVal reflect.Value, opt *Options) (*Metadata, error) {

	meta := new(Metadata)

	if dstVal.IsNil() {
		dstVal.Set(reflect.MakeMap(dstVal.Type()))
	}
	iter := srcVal.MapRange()
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		// 跳过 src nil 属性
		if opt.ShouldSkipNil(value) {
			continue
		}

		keyStr := key.Interface().(string)

		if opt.ShouldSkipField(keyStr) {
			continue
		}

		if opt.SkipEqual && reflect.DeepEqual(value, dstVal.MapIndex(key)) {
			continue
		}

		dstVal.SetMapIndex(key, value)

		meta.SrcKeys = append(meta.SrcKeys, keyStr)
		meta.DstKeys = append(meta.DstKeys, keyStr)

	}

	return meta, nil
}

func map2struct(srcVal, dstVal reflect.Value, opt *Options) (*Metadata, error) {

	meta := new(Metadata)

	iter := srcVal.MapRange()
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		for value.Kind() == reflect.Interface {
			value = value.Elem()
		}

		// 跳过 src nil 属性
		if opt.ShouldSkipNil(value) {
			continue
		}

		srcFieldName := key.Interface().(string)

		dstFieldVal, dstFieldName := opt.GetFieldByName(dstVal, srcFieldName)
		if !dstFieldVal.IsValid() {
			continue
		}

		if opt.ShouldSkipField(srcFieldName) || opt.ShouldSkipField(dstFieldName) {
			continue
		}

		if opt.SkipEqual && reflect.DeepEqual(value, dstFieldVal.Interface()) {
			continue
		}

		for dstFieldVal.Kind() == reflect.Ptr {
			ele := dstFieldVal.Elem()
			if !ele.IsValid() {
				dstFieldVal.Set(reflect.New(dstFieldVal.Type().Elem()))
			} else {
				dstFieldVal = ele
			}
		}

		if srcFieldName == "authority" {
			dtype := dstFieldVal.Type()
			vtype := value.Type()
			fmt.Printf("dst type: %v; value type: %v;  \n", dtype, vtype)
		}

		dstFieldVal.Set(value)

		meta.SrcKeys = append(meta.SrcKeys, srcFieldName)
		meta.DstKeys = append(meta.DstKeys, dstFieldName)

	}

	return meta, nil
}

func struct2map(srcVal, dstVal reflect.Value, opt *Options) (*Metadata, error) {

	meta := new(Metadata)

	if dstVal.IsNil() {
		dstVal.Set(reflect.MakeMap(dstVal.Type()))
	}

	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Type().Field(i)
		srcFieldVal := srcVal.Field(i)

		// 跳过 src nil 属性
		if opt.ShouldSkipNil(srcFieldVal) {
			continue
		}

		// 匿名字段
		if srcField.Anonymous && srcField.Type.Kind() == reflect.Struct {

			_meta, err := struct2map(srcFieldVal, dstVal, opt)
			if err != nil {
				return nil, err
			}
			meta.SrcKeys = append(meta.SrcKeys, _meta.SrcKeys...)
			meta.DstKeys = append(meta.DstKeys, _meta.SrcKeys...)

		} else {

			dstFieldName := getJsonName(srcField)
			if dstFieldName == "" {
				dstFieldName = strcase.ToLowerCamel(srcField.Name)
			}

			if opt.ShouldSkipField(srcField.Name) || opt.ShouldSkipField(dstFieldName) {
				continue
			}

			for srcFieldVal.Kind() == reflect.Ptr {
				srcFieldVal = srcFieldVal.Elem()
			}

			if opt.SkipEqual && reflect.DeepEqual(srcFieldVal.Interface(), dstVal.MapIndex(reflect.ValueOf(dstFieldName))) {
				continue
			}

			dstVal.SetMapIndex(reflect.ValueOf(dstFieldName), srcFieldVal)

			meta.SrcKeys = append(meta.SrcKeys, srcField.Name)
			meta.DstKeys = append(meta.DstKeys, dstFieldName)
		}
	}

	return meta, nil
}

func struct2struct(srcVal, dstVal reflect.Value, opt *Options) (*Metadata, error) {

	meta := new(Metadata)

	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Type().Field(i)
		srcFieldVal := srcVal.Field(i)

		// 跳过 src nil 属性
		if opt.ShouldSkipNil(srcFieldVal) {
			continue
		}

		dstFieldVal, dstFieldName := opt.GetFieldByName(dstVal, srcField.Name)
		if !dstFieldVal.IsValid() {
			continue
		}

		// 匿名字段
		if srcField.Anonymous && srcField.Type.Kind() == reflect.Struct {

			_meta, err := struct2struct(srcFieldVal, dstFieldVal, opt)
			if err != nil {
				return nil, err
			}
			meta.SrcKeys = append(meta.SrcKeys, _meta.SrcKeys...)
			meta.DstKeys = append(meta.DstKeys, _meta.SrcKeys...)

		} else {

			if opt.ShouldSkipField(dstFieldName) {
				continue
			}

			for srcFieldVal.Kind() == reflect.Ptr {
				srcFieldVal = srcFieldVal.Elem()
			}

			for dstFieldVal.Kind() == reflect.Ptr {
				dstFieldVal = dstFieldVal.Elem()
			}

			if opt.SkipEqual && reflect.DeepEqual(srcFieldVal.Interface(), dstFieldVal.Interface()) {
				continue
			}

			dstFieldVal.Set(srcFieldVal)

			meta.SrcKeys = append(meta.SrcKeys, srcField.Name)
			meta.DstKeys = append(meta.DstKeys, dstFieldName)
		}
	}

	return meta, nil
}
