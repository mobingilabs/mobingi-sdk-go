package pretty

import (
	"fmt"
	"io"
	"reflect"
)

var Pad int = 2

func Indent(count int) string {
	pad := ""
	for i := 0; i < count; i++ {
		pad += " "
	}

	return pad
}

// Print prints the `field: value` of the input interface recursively. Recursion
// level `lvl` and `indent` are provided for indention in printing.
func Print(w io.Writer, s interface{}, lvl, indent int) {
	pad := Indent(lvl * indent)
	rt := reflect.TypeOf(s).Elem()
	rv := reflect.ValueOf(s).Elem()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i).Name
		value := rv.Field(i).Interface()
		switch rv.Field(i).Kind() {
		case reflect.Interface:
			fmt.Fprintf(w, "%s%s: %v\n", pad, field, value)
		case reflect.String:
			fmt.Fprintf(w, "%s%s: %s\n", pad, field, value)
		case reflect.Int32:
			fmt.Fprintf(w, "%s%s: %i\n", pad, field, value)
		case reflect.Struct:
			fmt.Fprintf(w, "%s[%s]\n", pad, field)
			v := rv.Field(i).Addr()
			Print(w, v.Interface(), lvl+1, indent)
		case reflect.Slice:
			fmt.Fprintf(w, "%s[%s]\n", pad, field)
			v := reflect.ValueOf(s)
			for i := 0; i < v.Len(); i++ {
				Print(w, v.Index(i).Interface(), lvl+1, indent)
			}
		case reflect.Ptr:
			v := reflect.ValueOf(s).Elem()
			if v.IsValid() {
				Print(w, v.Interface(), lvl+1, indent)
			} else {
				fmt.Fprintf(w, "%snil\n", pad)
			}
		default:
			fmt.Fprintf(w, "%s%#v\n", pad, s)
		}
	}
}
