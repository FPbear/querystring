package querystring

import (
	"fmt"
	"net/url"
	"reflect"
	"time"
)

var timeType = reflect.TypeOf(time.Time{})

// Encoder is an interface implemented by any type that wishes to encode
var encoderType = reflect.TypeOf((*Encoder)(nil)).Elem()

// Encoder is an interface implemented by any type that wishes to encode
// itself into URL values in a non-standard way.
// The Encode method returns the encoded values.
type Encoder interface {
	Encode() ([]string, error)
}

type Converter struct {
	tag Tag
}

func NewConverter(tag Tag) *Converter {
	return &Converter{
		tag: tag,
	}
}

// Values converts a value to url.Values.
// The value can be a map, a struct, a pointer to a struct, or a pointer to a map.
// The value can also implement the Encoder interface.
// If the value is a map, the key must be a string and the value must be a string.
// If the value is a struct, the field must have a tag with the key "url" or a custom tag type.
// The tag value is the name of the field in the url.Values.
func (c *Converter) Values(v interface{}) (url.Values, error) {
	if val, ok := v.(url.Values); ok {
		return val, nil
	}
	values := make(url.Values)
	if v == nil {
		return values, nil
	}
	vf := reflect.ValueOf(v)
	if vf.Kind() == reflect.Map {
		if vf.Type().Key().Kind() != reflect.String {
			return nil, fmt.Errorf("map key must be a string")
		}
		iter := vf.MapRange()
		for iter.Next() {
			if iter.Value().Kind() != reflect.String {
				return nil, fmt.Errorf("map value must be a string")
			}
			values.Add(iter.Key().String(), iter.Value().String())
		}
		return values, nil
	}
	if vf.Kind() == reflect.Ptr {
		if vf.IsNil() {
			return values, nil
		}
		vf = vf.Elem()
	}
	if vf.Kind() != reflect.Struct {
		return nil, fmt.Errorf("unsupported type %T", vf.Kind())
	}
	if err := c.reflectValue(values, vf); err != nil {
		return nil, err
	}
	return values, nil
}

func (c *Converter) reflectValue(values url.Values, val reflect.Value) error {
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)
		if sf.PkgPath != "" && !sf.Anonymous { // unexported
			continue
		}
		sv := val.Field(i)
		tag, ok := c.tag.Get(sf)
		if !ok {
			continue
		}
		name, opts := c.tag.ParseTag(tag)
		if opts.Contains("omitempty") && isEmptyValue(sv) {
			continue
		}

		if sv.Type() == timeType {
			values.Add(name, valueToString(sv))
			continue
		}
		if sv.Type().Implements(encoderType) {
			enc := sv.Interface().(Encoder)
			encoded, err := enc.Encode()
			if err != nil {
				return err
			}
			for _, v := range encoded {
				values.Add(name, v)
			}
			continue
		}

		if sv.Kind() == reflect.Ptr {
			if sv.IsNil() {
				break
			}
			sv = sv.Elem()
		}
		switch sv.Kind() {
		case reflect.Slice, reflect.Array:
			for index := 0; index < sv.Len(); index++ {
				values.Add(name, valueToString(sv.Index(index)))
			}
		case reflect.String:
			values.Add(name, sv.String())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			values.Add(name, fmt.Sprintf("%d", sv.Int()))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			values.Add(name, fmt.Sprintf("%d", sv.Uint()))
		case reflect.Bool:
			values.Add(name, fmt.Sprintf("%t", sv.Bool()))
		default:
			continue

		}
	}
	return nil
}

// isEmptyValue checks if a value should be considered empty for the purposes
// of omitting fields with the "omitempty" option.
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	default:
	}

	type zeroAble interface {
		IsZero() bool
	}

	if z, ok := v.Interface().(zeroAble); ok {
		return z.IsZero()
	}

	return false
}

// valueToString converts a reflect.Value to a string.
// The value can be a string, an integer, a float, or a boolean.
// If the value is a time.Time, the function returns the time in time.RFC3339Nano format.
func valueToString(value reflect.Value) string {
	if value.Type() == timeType {
		return value.Interface().(time.Time).Format(time.RFC3339Nano)
	}
	switch value.Kind() {
	case reflect.String:
		return value.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", value.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", value.Uint())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", value.Float())
	case reflect.Bool:
		return fmt.Sprintf("%t", value.Bool())
	default:
		return ""
	}
}

// Values converts a value to url.Values.
// The value can be a map, a struct, a pointer to a struct, or a pointer to a map.
// The value can also implement the Encoder interface.
// If the value is a map, the key must be a string and the value must be a string.
// If the value is a struct, the field must have a tag with the key "url".
func Values(v interface{}) (url.Values, error) {
	converter := NewConverter(NewTag())
	return converter.Values(v)
}
