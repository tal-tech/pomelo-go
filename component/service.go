package component

//type (
//	Service struct {
//		Name     string             // name of service
//		Type     reflect.Type       // type of the receiver
//		Receiver reflect.Value      // receiver of methods for the service
//		Handlers map[string]Handler // registered methods
//		Options  options            // options
//	}
//)
//
//func NewService(comp Component, opts []Option) *Service {
//	s := &Service{
//		Type:     reflect.TypeOf(comp),
//		Receiver: reflect.ValueOf(comp),
//	}
//
//	// apply options
//	for i := range opts {
//		opt := opts[i]
//		opt(&s.Options)
//	}
//	if name := s.Options.name; name != "" {
//		s.Name = name
//	} else {
//		s.Name = reflect.Indirect(s.Receiver).Type().Name()
//	}
//
//	s.Handlers = comp.Handlers()
//
//	return s
//}
