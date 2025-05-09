package queryutil

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

func StructToQuery(i any) url.Values {
	values := url.Values{}

	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	for idx, field := range reflect.VisibleFields(t) {
		tag := field.Tag.Get("query")
		if tag == "" {
			continue
		}

		fv := v.Field(idx)
		if !fv.IsValid() || fv.IsNil() {
			continue
		}
		val := fv.Elem()

		switch val.Kind() {
		case reflect.Slice:
			for j := 0; j < val.Len(); j++ {
				values.Add(tag, fmt.Sprint(val.Index(j).Interface()))
			}
		case reflect.String:
			values.Set(tag, val.String())
		case reflect.Bool:
			values.Set(tag, strconv.FormatBool(val.Bool()))
		case reflect.Struct:
			if tm, ok := val.Interface().(time.Time); ok {
				values.Set(tag, tm.Format(time.RFC3339))
			}
		default:
			values.Set(tag, fmt.Sprint(val.Interface()))
		}
	}

	return values
}
