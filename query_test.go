package querystring

import (
	"net/url"
	"testing"
	"time"
)

func TestMapToValues(t *testing.T) {
	input := map[string]string{
		"hello": "world",
		"foo":   "bar",
		"empty": "",
	}

	values, err := Values(input)
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range input {
		if values.Get(k) != v {
			t.Errorf("expected %q, got %q", v, values.Get(k))
		}
	}
}

func TestValues(t *testing.T) {
	input := url.Values{
		"hello": []string{"world"},
		"foo":   []string{"bar"},
		"empty": []string{""},
	}

	values, err := Values(input)
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range input {
		if values.Get(k) != v[0] {
			t.Errorf("expected %q, got %q", v, values.Get(k))
		}
	}
}

type subEncode struct {
	Hello string `form:"hello"`
	World string `form:"world"`
}

func (s *subEncode) Encode() ([]string, error) {
	return []string{s.Hello + "-" + s.World}, nil
}

func TestStructToValues(t *testing.T) {
	type Input struct {
		Hello string     `form:"hello"`
		Foo   string     `form:"foo,omitempty"`
		Empty string     `form:"empty"`
		Om    int64      `form:"om,omitempty"`
		Array [3]int     `form:"array"`
		Slice []string   `form:"slice"`
		Sub   *subEncode `form:"sub"`
		small string     `form:"small"`
		Skip  int        `form:"-"`
		Time  time.Time  `form:"time,omitempty"`
	}
	expected := make(url.Values)
	now := time.Now()
	in := Input{
		Hello: "world",
		Foo:   "bar",
		Empty: "",
		Array: [3]int{1, 2, 3},
		Slice: []string{"a", "b", "c"},
		Sub: &subEncode{
			Hello: "subHello",
			World: "subWorld",
		},
		small: "small",
		Skip:  10,
		Time:  now,
	}
	expected["hello"] = []string{"world"}
	expected["foo"] = []string{"bar"}
	expected["empty"] = []string{""}
	expected["array"] = []string{"1", "2", "3"}
	expected["slice"] = []string{"a", "b", "c"}
	expected["sub"] = []string{"subHello-subWorld"}
	expected["time"] = []string{now.Format(time.RFC3339Nano)}
	con := NewConverter(NewTag(WithTag("form")))
	values, err := con.Values(in)
	if err != nil {
		t.Fatal(err)
	}
	for k, vs := range values {
		expectedVs, ok := expected[k]
		if !ok {
			t.Errorf("unexpected key %q", k)
			continue
		}
		if len(vs) != len(expectedVs) {
			t.Errorf("expected %q to have %d values, got %d", k, len(expectedVs), len(vs))
			continue
		}
		for i, v := range vs {
			if v != expectedVs[i] {
				t.Errorf("expected %q, got %q", expectedVs[i], v)
			}
		}
	}
}
