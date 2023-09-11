package component

type (
	options struct {
		//name string // component name
	}

	// Option used to customize handler
	Option func(options *options)
)

//// WithName used to rename component name
//func WithName(name string) Option {
//	return func(opt *options) {
//		opt.name = name
//	}
//}
