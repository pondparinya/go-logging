package ecslog

//---------------------------------------------------------

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

//---------------------------------------------------------

var CENSORS []string = []string{
	"address",
	// "bank account no.",
	"account",
	// "biometric data",
	// "birthdate",
	"birth",
	// "car licence no.",
	// "driver license no.",
	"license",
	// "cookie id",
	// "credit card no.",
	"card",
	"email",
	// "id card no.",
	// "passport no.",
	// "tax id",
	"citizen",
	"passport",
	"tax",
	// "ip adress",
	// "mac address",
	// "mobile no.",
	"mobile",
	// "username",
	// "password",
	"user",
	"pass",
	// "nationality",
	"nation",
	"race",
	// "firstname",
	// "surname",
	"name",
	// "height",
	// "weight",
}

//---------------------------------------------------------

func isInterfaceNil(i interface{}) bool {
	if i == nil {
		return true
	}

	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}

	return false
}

//---------------------------------------------------------

func containsAny(vs []string, n string) bool {
	if vs == nil {
		return false
	}

	n = strings.ToLower(n)
	for _, v := range vs {
		if strings.Contains(n, v) {
			return true
		}
	}

	return false
}

func Filter(vs []string) func(string) bool {
	return func(n string) bool {
		return containsAny(vs, n)
	}
}

func Censor(n string) bool {
	return containsAny(CENSORS, n)
}

//---------------------------------------------------------

func filterFields(filter func(string) bool, i interface{}) error {

	if isInterfaceNil(i) {
		return nil
	}

	r := reflect.ValueOf(i)
	if r.Kind() != reflect.Ptr {
		return fmt.Errorf("expect to be an address")
	}

	r = r.Elem()
	if !r.CanSet() {
		return fmt.Errorf("expect to be settable")
	}

	t := r.Type()
	k := t.Kind()
	if k == reflect.Struct {
		for i, N := 0, t.NumField(); i < N; i += 1 {
			f := r.Field(i)
			if !f.CanInterface() ||
				isInterfaceNil(f.Interface()) {
				continue
			}

			n := strings.ToLower(t.Field(i).Name)
			if filter(n) {
				f.Set(reflect.Zero(f.Type()))
				continue
			}

			if t.Kind() == reflect.Struct {
				if err := filterFields(filter, f.Interface()); err != nil {
					return fmt.Errorf("%s: %v", n, err.Error())
				}
			}
		}
	}

	return nil
}

func FilterFields(filter func(string) bool, i interface{}) error {
	return filterFields(filter, i)
}

func CensorFields(i interface{}) error {
	return FilterFields(Censor, i)
}

//---------------------------------------------------------

func FilterValue(filter func(string) bool, x interface{}) string {
	if isInterfaceNil(x) {
		return ""
	}

	if t, ok := x.(time.Time); ok {
		return t.Format("\"2006-01-02T150405\"")
	}

	if t, ok := x.(*time.Time); ok {
		return t.Format("\"2006-01-02T150405\"")
	}

	r := reflect.ValueOf(x)
	t := r.Type()
	k := t.Kind()

	if k == reflect.Interface ||
		k == reflect.Ptr ||
		k == reflect.UnsafePointer {
		return FilterValue(filter, r.Elem())
	}

	if k == reflect.String {
		return fmt.Sprintf("\"%v\"", r)
	}

	if k == reflect.Bool ||
		k == reflect.Chan ||
		k == reflect.Func ||
		k == reflect.Int ||
		k == reflect.Int8 ||
		k == reflect.Int16 ||
		k == reflect.Int32 ||
		k == reflect.Int64 ||
		k == reflect.Uint ||
		k == reflect.Uint8 ||
		k == reflect.Uint16 ||
		k == reflect.Uint32 ||
		k == reflect.Uint64 ||
		k == reflect.Uintptr ||
		k == reflect.Float32 ||
		k == reflect.Float64 ||
		k == reflect.Complex64 ||
		k == reflect.Complex128 {
		return fmt.Sprintf("%v", r)
	}

	if k == reflect.Array ||
		k == reflect.Slice {

		s := ""
		for i, n := 0, r.Len(); i < n; i += 1 {
			if s != "" {
				s += ","
			}
			s += FilterValue(filter, r.Index(i).Interface())
		}

		return fmt.Sprintf("[%s]", s)
	}

	if k == reflect.Map {

		s := ""
		for _, k := range r.MapKeys() {
			i := k.Interface()
			if v, ok := i.(string); ok && filter(v) {
				continue
			}

			if s != "" {
				s += ","
			}
			s += fmt.Sprintf("%s:", FilterValue(filter, i))
			s += FilterValue(filter, r.MapIndex(k).Interface())
		}

		return fmt.Sprintf("{%s}", s)
	}

	if k == reflect.Struct {

		s := ""
		for i, N := 0, r.NumField(); i < N; i += 1 {
			n := t.Field(i).Name
			if filter(n) {
				continue
			}

			f := r.Field(i)
			if !f.CanInterface() {
				continue
			}

			v := r.Field(i).Interface()
			if s != "" {
				s += ","
			}
			s += fmt.Sprintf("\"%s\":", n)
			s += FilterValue(filter, v)
		}

		return fmt.Sprintf("{%s}", s)
	}

	return ""
}

func CensorValue(x interface{}) string {
	return FilterValue(Censor, x)
}

//---------------------------------------------------------

func FilterValues(filter func(string) bool, args ...interface{}) []interface{} {

	buff := make([]interface{}, len(args))

	for i, arg := range args {
		if isInterfaceNil(arg) {
			buff[i] = arg
			continue
		}

		r := reflect.ValueOf(arg)
		k := r.Type().Kind()

		if k == reflect.Struct {
			arg = FilterValue(filter, arg)
		}

		buff[i] = arg
	}

	return buff
}

func CensorValues(args ...interface{}) []interface{} {
	return FilterValues(Censor, args...)
}

//---------------------------------------------------------
// End-of-file
//---------------------------------------------------------
