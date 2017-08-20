package pretty

import (
	"fmt"
	"io"
	"log"
	"reflect"
	"strings"
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
			/*
				v := reflect.ValueOf(s).Elem()
				if v.IsValid() {
					Print(w, v.Interface(), lvl+1, indent)
				} else {
					fmt.Fprintf(w, "%snil\n", pad)
				}
			*/
		default:
			fmt.Fprintf(w, "%s%#v\n", pad, s)
		}
	}
}

func JSON(v interface{}, lvl, indent int) string {
	var add string
	pad := Indent(lvl * indent)
	log.Println(len(pad))
	value := reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.Struct:
		if lvl == 1 {
			add = ""
		}
		str := add + fullName(value.Type()) + "[" + fmt.Sprintf("%d", lvl) + "]{\n"
		for i := 0; i < value.NumField(); i++ {
			l := string(value.Type().Field(i).Name[0])
			if strings.ToUpper(l) == l {
				str += pad + value.Type().Field(i).Name + "[" + fmt.Sprintf("%d", lvl) + "]: "
				str += JSON(value.Field(i).Interface(), lvl+1, indent)
				str += add + ",\n"
			}
		}
		str += pad + "}"
		return str
	case reflect.Map:
		if len(pad) > 0 {
			add = ""
		}
		str := add + "map[" + fullName(value.Type().Key()) + "]" + fullName(value.Type().Elem()) + "{\n"
		for _, k := range value.MapKeys() {
			str += pad + "\"" + k.String() + "\": "
			str += JSON(value.MapIndex(k).Interface(), lvl+1, indent)
			str += add + ",\n"
		}
		str += pad + "}"
		return str
	case reflect.Ptr:
		if e := value.Elem(); e.IsValid() {
			return "&" + JSON(e.Interface(), lvl+1, indent)
		}
		return "nil"
	case reflect.Slice:
		str := "[]" + fullName(value.Type().Elem()) + "{\n"
		for i := 0; i < value.Len(); i++ {
			str += JSON(value.Index(i).Interface(), lvl+1, indent)
			str += ",\n"
		}
		str += "}"
		return str
	default:
		return fmt.Sprintf("[def]%#v[def]", v)
	}
}

func pkgName(t reflect.Type) string {
	pkg := t.PkgPath()
	c := strings.Split(pkg, "/")
	return c[len(c)-1]
}

func fullName(t reflect.Type) string {
	if pkg := pkgName(t); pkg != "" {
		return pkg + "." + t.Name()
	}
	return t.Name()
}
