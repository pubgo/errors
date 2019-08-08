package internal

import "reflect"

var errTags = struct {
	UnknownTypeCode string
}{
	"errors_unknown_type",
}

var _errTags = make(map[string]bool)

func ErrTagRegistry(tags ...interface{}) {
	for _, tag := range tags {
		if IsNone(tag) {
			continue
		}

		var _tags []string
		t := reflect.ValueOf(tag)
		switch t.Kind() {
		case reflect.String:
			_tags = append(_tags, tag.(string))
		case reflect.Ptr, reflect.Struct:
			for i := 0; i < t.NumField(); i++ {
				_tags = append(_tags, t.Field(i).String())
			}
		}

		for _, t := range _tags {
			if _, ok := _errTags[t]; ok {
				T(ok, "tag %s has existed", t)
			}
			_errTags[t] = true
		}
	}
}

func ErrTags() map[string]bool {
	return _errTags
}

func ErrTagsMatch(tag string) bool {
	_, ok := _errTags[tag]
	return ok
}

func init() {
	ErrTagRegistry(errTags.UnknownTypeCode)
}
