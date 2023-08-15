package bencode

import (
	"reflect"
	"sort"
	"strings"
	"unicode"
)

func sortStrings(ss []string) {
	if len(ss) <= strSliceLen {
		// for i := 1; i < len(ss); i++ {
		// 	for j := i; j > 0 && ss[j] < ss[j-1]; j-- {
		// 		ss[j], ss[j-1] = ss[j-1], ss[j]
		// 	}
		// }
		// below is the code above, but (almost) without bound checks

		for i := 1; i < len(ss); i++ {
			for j := i; j > 0; j-- {
				if ss[j] >= ss[j-1] {
					break
				}
				ss[j], ss[j-1] = ss[j-1], ss[j]
			}
		}
	} else {
		sort.Strings(ss)
	}
}

func fieldTag(field reflect.StructField) (string, bool) {
	tag := field.Tag.Get("bencode")

	var opts string
	switch {
	case tag == "":
		return field.Name, true
	case tag == "-":
		return "", false
	default:
		if idx := strings.Index(tag, ","); idx != -1 {
			tag, opts = tag[:idx], tag[idx+1:]
		}
	}

	switch {
	case strings.Contains(opts, ",omitempty"):
		return "", false
	case !isValidTag(tag):
		return field.Name, true
	default:
		return tag, true
	}
}

func isValidTag(key string) bool {
	if key == "" {
		return false
	}

	for _, c := range key {
		if c != ' ' && c != '$' && c != '-' && c != '_' && c != '.' &&
			!unicode.IsLetter(c) && !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func isNil(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	default:
		return false
	}
}
