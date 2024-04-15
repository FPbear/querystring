package querystring

import (
	"testing"
)

func TestIsEmpty(t *testing.T) {
	var n Name
	if !n.IsEmpty() {
		t.Errorf("Expected empty name to be empty")
	}

	n = "non-empty"
	if n.IsEmpty() {
		t.Errorf("Expected non-empty name to not be empty")
	}
}

func TestConvert(t *testing.T) {
	tests := []struct {
		n     Name
		cases []struct {
			input  string
			expect string
		}
	}{
		{
			n: CamelCase,
			cases: []struct {
				input  string
				expect string
			}{
				{
					"helloWorld", "helloWorld",
				},
				{
					"hello2World", "hello2World",
				},
				{
					"helloworld", "helloworld",
				},
				{
					"hello_world", "helloWorld",
				},
				{
					"hello-world", "helloWorld",
				},
			},
		},
		{
			n: PascalCase,
			cases: []struct {
				input  string
				expect string
			}{
				{"helloWorld", "HelloWorld"},
				{"hello2World", "Hello2World"},
				{"helloworld", "Helloworld"},
				{"hello_world", "HelloWorld"},
				{"hello-world", "HelloWorld"},
				{"hello.world", "HelloWorld"},
				{"hello2_world", "Hello2World"},
			},
		},
		{
			n: SnakeCase,
			cases: []struct {
				input  string
				expect string
			}{
				{"helloWorld", "hello_world"},
				{"hello2World", "hello_2_world"},
				{"helloworld", "helloworld"},
				{"hello_world", "hello_world"},
				{"hello-world", "hello_world"},
				{"hello.world", "hello_world"},
				{"hello2_world", "hello_2_world"},
				{"hello2-world", "hello_2_world"},
			},
		},
	}
	for _, test := range tests {
		for _, c := range test.cases {
			actual := test.n.Convert(c.input)
			if actual != c.expect {
				t.Errorf("Failed to convert to %s case(input:%s, actual: %s, expect: %s)", test.n, c.input, actual, c.expect)
			}
		}

	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{"helloWorld", "helloWorld"},
		{"hello2World", "hello2World"},
		{"helloworld", "helloworld"},
		{"hello_world", "helloWorld"},
		{"hello-world", "helloWorld"},
		{"hello.world", "helloWorld"},
		{"hello2_world", "hello2World"},
		{"hello2-world", "hello2World"},
		{"hello2.world", "hello2World"},
		{"hello2world", "hello2world"},
		{"hello_2_world", "hello2World"},
		{"hello_world_hi", "helloWorldHi"},
		{"hello-world-hi", "helloWorldHi"},
		{"hello.world.hi", "helloWorldHi"},
		{"hello_world-hi.hi", "helloWorldHiHi"},
		{"foo-bar-", ""},
		{"foo-bar-1", "fooBar1"},
		{"1foo-bar", ""},
	}
	for _, test := range tests {
		actual := ToCamelCase(test.input)
		if actual != test.expect {
			t.Errorf("Failed to convert to camel case(input:%s, actual: %s, expect: %s)", test.input, actual, test.expect)
		}
	}
}

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{"helloWorld", "HelloWorld"},
		{"hello2World", "Hello2World"},
		{"helloworld", "Helloworld"},
		{"hello_world", "HelloWorld"},
		{"hello-world", "HelloWorld"},
		{"hello.world", "HelloWorld"},
		{"hello2_world", "Hello2World"},
		{"hello2-world", "Hello2World"},
		{"hello2.world", "Hello2World"},
		{"hello2world", "Hello2world"},
		{"hello_2_world", "Hello2World"},
		{"hello_world_hi", "HelloWorldHi"},
		{"hello-world-hi", "HelloWorldHi"},
		{"hello.world.hi", "HelloWorldHi"},
		{"hello_world-hi.hi", "HelloWorldHiHi"},
		{"foo-bar-", ""},
		{"foo-bar-1", "FooBar1"},
		{"1foo-bar", ""},
	}
	for _, test := range tests {
		actual := ToPascalCase(test.input)
		if actual != test.expect {
			t.Errorf("Failed to convert to pascal case(input:%s, actual: %s, expect: %s)", test.input, actual, test.expect)
		}
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{"helloWorld", "hello_world"},
		{"hello2World", "hello_2_world"},
		{"helloworld", "helloworld"},
		{"hello_world", "hello_world"},
		{"hello-world", "hello_world"},
		{"hello.world", "hello_world"},
		{"hello2_world", "hello_2_world"},
		{"hello2-world", "hello_2_world"},
		{"hello2.world", "hello_2_world"},
		{"hello2world", "hello_2_world"},
		{"hello_2_world", "hello_2_world"},
		{"hello_world_hi", "hello_world_hi"},
		{"hello-world-hi", "hello_world_hi"},
		{"hello.world.hi", "hello_world_hi"},
		{"hello_world-hi.hi", "hello_world_hi_hi"},
		{"foo-bar-", ""},
		{"foo-bar-1", "foo_bar_1"},
		{"1foo-bar", ""},
	}
	for _, test := range tests {
		actual := ToSnakeCase(test.input)
		if actual != test.expect {
			t.Errorf("Failed to convert to snake case(input:%s, actual: %s, expect: %s)", test.input, actual, test.expect)
		}
	}
}
