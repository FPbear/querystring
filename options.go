package querystring

type option struct {
	useName   Name
	skipField string
	tag       string
}

type Option func(*option)

func WithUseName(useName Name) Option {
	return func(o *option) {
		o.useName = useName
	}
}

func WithSkipField(skipField string) Option {
	return func(o *option) {
		o.skipField = skipField
	}
}

func WithTag(tag string) Option {
	return func(o *option) {
		o.tag = tag
	}
}

func defaultOption() *option {
	opt := &option{
		useName:   SnakeCase,
		skipField: "-",
		tag:       "url",
	}
	return opt
}
