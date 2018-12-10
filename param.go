/*
Package param used to parse a url.Values to a given struct
*/
package param

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// Unmarshal unmarchal url values to struct
func Unmarshal(u url.Values, o interface{}) error {
	val := reflect.ValueOf(o)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return reflectValueFromValues(u, val)
}

func reflectValueFromValues(u url.Values, v reflect.Value) error {
	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		vField := v.Field(i)
		if !vField.CanSet() {
			continue
		}

		tField := typ.Field(i)

		if tField.Type.Kind() == reflect.Struct {
			t := reflect.New(tField.Type).Elem()
			if err := reflectValueFromValues(u, t); err != nil {
				return err
			}
			vField.Set(t)
			continue
		}

		vl, exists := getVal(u, tField)
		if !exists {
			continue
		}

		switch tField.Type.Kind() {
		case reflect.String:
			vField.SetString(vl)
		case reflect.Bool:
			b, err := strconv.ParseBool(vl)
			if err != nil {
				return err
			}
			vField.SetBool(b)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			n, err := strconv.ParseUint(vl, 10, 64)
			if err != nil {
				return err
			}
			if vField.OverflowUint(n) {
				return fmt.Errorf("cast %s to uint overflow", tField.Name)
			}
			vField.SetUint(n)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, err := strconv.ParseInt(vl, 10, 64)
			if err != nil {
				return err
			}
			if vField.OverflowInt(n) {
				return fmt.Errorf("cast %s to int overflow", tField.Name)
			}
			vField.SetInt(n)
		case reflect.Float32, reflect.Float64:
			n, err := strconv.ParseFloat(vl, vField.Type().Bits())
			if err != nil {
				return err
			}
			if vField.OverflowFloat(n) {
				return fmt.Errorf("cast %s to float overflow", tField.Name)
			}
			vField.SetFloat(n)
		}
	}
	return nil
}

func getVal(u url.Values, f reflect.StructField) (string, bool) {
	tag, ok := f.Tag.Lookup("param")
	if !ok {
		if v, ok := u[f.Name]; ok {
			return v[0], true
		}
		if v, ok := u[strings.ToLower(f.Name[0:1])+f.Name[1:]]; ok {
			return v[0], true
		}
		return "", false
	}
	// exists tag
	slice := strings.Split(tag, ",")
	name := slice[0]
	opt := slice[1:]
	if name == "-" { // break
		return "", false
	}
	if name == "" {
		name = strings.ToLower(f.Name[0:1]) + f.Name[1:]
	}
	v := u.Get(name)
	if v == "" && len(opt) == 1 {
		v = opt[0]
	}
	return v, true
}
