package ctx

import (
	"fmt"
	"reflect"
	"strconv"
)

type ErrorCode int

type Error struct {
	Code   ErrorCode
	Errors Errors
}

func ErrorStatusCode(code int) int {
	i, _ := strconv.Atoi(fmt.Sprintf("%d", code)[:3])
	return i
}

func NewErrors(s ...interface{}) *Errors {
	r := &Errors{
		[]interface{}{},
	}

	for _, v := range s {
		rt := reflect.TypeOf(v)
		switch rt.Kind() {
		case reflect.Slice, reflect.Array:
			v1 := reflect.ValueOf(v)
			for i := 0; i < v1.Len(); i++ {
				r.s = append(r.s, v1.Index(i).Interface())
			}
		default:
			r.s = append(r.s, v)
		}
	}

	return r
}

type Errors struct {
	s []interface{}
}

func (d *Errors) Add(s interface{}) *Errors {
	d.s = append(d.s, s)
	return d
}

func (d *Errors) Error() []interface{} {
	return d.s
}

func (d *Errors) NotEmpty() bool {
	return len(d.s) > 0
}
