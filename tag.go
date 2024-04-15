package querystring

import (
	"reflect"
	"strings"
)

// TagOptions represents the options of a tag.
// The options are separated by commas.
type TagOptions []string

// Contains checks whether the tagOptions contains the specified option.
func (o TagOptions) Contains(option string) bool {
	for _, s := range o {
		if s == option {
			return true
		}
	}
	return false
}

// Tag represents the tag interface.
// The tag interface provides the methods to get the tag value, parse the tag, and skip the field.
// The tag value is the value of the tag in the struct field.
// The parse tag method parses the tag value and returns the tag name and tag options.
// The skip method returns the value of the tag to skip the field.
type Tag interface {
	// Get returns the tag value of the struct field.
	Get(field reflect.StructField) (string, bool)
	// ParseTag parses the tag value and returns the tag name and tag options.
	ParseTag(tag string) (string, TagOptions)
	// Skip returns the value of the tag to skip the field.
	Skip() string
}

// defaultTag represents the default tag.
// The default tag contains the tag type, the skip field value, and the use name.
// The tag type is the type of the tag in the struct field.
// The skip field value is the value of the tag to skip the field.
// The useName is the name to use when the tag is empty.
type defaultTag struct {
	tagType string
	skip    string
	useName Name
}

func NewTag(opts ...Option) Tag {
	opt := defaultOption()
	for _, o := range opts {
		o(opt)
	}
	return &defaultTag{
		tagType: opt.tag,
		skip:    opt.skipField,
		useName: opt.useName,
	}
}

func (t *defaultTag) Get(field reflect.StructField) (string, bool) {
	tag := field.Tag.Get(t.tagType)
	if tag == t.skip {
		return "", false
	}
	if tag == "" && !t.useName.IsEmpty() {
		return t.useName.Convert(field.Name), true
	}
	return tag, true
}

func (t *defaultTag) ParseTag(tag string) (string, TagOptions) {
	s := strings.Split(tag, ",")
	if len(s) == 0 {
		return "", nil
	}
	return s[0], s[1:]
}

func (t *defaultTag) Skip() string {
	return t.skip
}
