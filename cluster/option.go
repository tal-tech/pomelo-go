package cluster

import "pomelo-go/component"

type RunOption func(*Server)

func WithComponent(component component.Component) RunOption {
	return func(node *Server) {
		node.components.Register(component, nil)
	}
}
