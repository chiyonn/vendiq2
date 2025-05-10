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
		if !fv.IsValid() {
			continue
		}

		if fv.Kind() == reflect.Ptr {
			if fv.IsNil() {
				continue
			}
			fv = fv.Elem()
		}

		switch fv.Kind() {
		case reflect.Slice:
			for j := 0; j < fv.Len(); j++ {
				values.Add(tag, fmt.Sprint(fv.Index(j).Interface()))
			}
		case reflect.String:
			values.Set(tag, fv.String())
		case reflect.Bool:
			values.Set(tag, strconv.FormatBool(fv.Bool()))
		case reflect.Struct:
			if tm, ok := fv.Interface().(time.Time); ok {
				values.Set(tag, tm.Format(time.RFC3339))
			}
		default:
			values.Set(tag, fmt.Sprint(fv.Interface()))
		}
	}

	return values
}
