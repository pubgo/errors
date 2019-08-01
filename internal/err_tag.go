package internal

var errTags = struct {
	UnknownTypeCode string
}{
	"errors_unknown_type",
}

var _errTags = make(map[string]bool)

func ErrTagRegistry(tags ...string) {
	for _, tag := range tags {
		if _, ok := _errTags[tag]; ok {
			T(ok, "tag %s has existed", tag)
		}
		_errTags[tag] = true
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
