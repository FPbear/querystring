# QueryString

QueryString is a Go library for converting a struct or a map into URL values, with support for custom encoders, time formatting, and handling of different data types.


## Installation

Use the package manager [go get](https://github.com/FPbear/querystring) to install QueryString.

```bash
go get github.com/FPbear/querystring
```
## Usage
### Example

```golang
import "github.com/FPbear/querystring"

type Input struct {
    Hello string     `url:"hello"`
    Foo   string     `url:"foo,omitempty"`
    Empty string     `url:"empty"`
}

func main() {
    input := Input{
        Hello: "world",
        Foo:   "",
        Empty: "",
    }

    values, err := querystring.Values(input)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(values.Encode())
    // Output: hello=world&empty=
}
```
You can also customize the encoder, for example:
```golang
import "github.com/FPbear/querystring"

type subEncode struct {
	Hello string `form:"hello"`
	World string `form:"world"`
}

func (s *subEncode) Encode() ([]string, error) {
	return []string{s.Hello + "-" + s.World}, nil
}

type Input struct {
    Hello string     `form:"hello"`
    Foo   string     `form:"foo,omitempty"`
    Empty string     `form:"empty"`
    Sub   *subEncode `form:"sub,omitempty"`
}

func main() {
    input := Input{
        Hello: "world",
        Foo:   "",
        Empty: "",
        Sub: &subEncode{
            Hello: "hello",
            World: "world",
        },
    }

    con := NewConverter(NewTag(WithTag("form")))
    values, err := con.Values(input)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(values.Encode())
    // Output: hello=world&empty=&sub=hello-world
}
```