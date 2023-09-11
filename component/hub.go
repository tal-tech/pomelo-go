package component

type CompWithOptions struct {
	Comp Component
	Opts []Option
}

type Components struct {
	comps []CompWithOptions
}

// Register registers a component to hub with options
func (cs *Components) Register(c Component, options ...Option) {
	cs.comps = append(cs.comps, CompWithOptions{c, options})
}

// List returns all components with it's options
func (cs *Components) List() []CompWithOptions {
	return cs.comps
}
