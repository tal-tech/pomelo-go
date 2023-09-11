package backend

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"pomelo-go/cluster/clusterpb/proto"
	"pomelo-go/component"
	"pomelo-go/internal/codec"
)

const (
	namespace = "sys"
	service   = "msgRemote"
)

var _ = component.Component(&Component{})

type ForwardMessageHandler func(context.Context, proto.Session, proto.Message) (interface{}, error)

type Component struct {
	serverType string
	router     map[string]ForwardMessageHandler
}

func NewComponent(serverType string) *Component {

	return &Component{
		serverType: serverType,
		router:     make(map[string]ForwardMessageHandler, 0),
	}

}

func (c *Component) Routes() (router map[string]component.Handler) {

	route := fmt.Sprintf("%s.%s.%s.%s", namespace, c.serverType, service, "forwardMessage")
	router = map[string]component.Handler{
		route: c.forwardMessageHandler,
	}

	return router
}

func (c *Component) AddRoutes(rs []Route) {

	for _, r := range rs {
		if _, ok := c.router[r.Method]; ok {
			panic(fmt.Errorf("handler: route already defined: %s", r.Method))
		}

		c.router[r.Method] = r.Handler
	}
}

func (c *Component) forwardMessageHandler(ctx context.Context, in json.RawMessage) (out json.RawMessage) {

	session := proto.Session{}
	msg := proto.Message{}

	if err := codec.Decode(in, []interface{}{&msg, &session}); err != nil {
		return result(err, nil)
	}

	if handler, ok := c.router[msg.Route]; !ok {
		return result(errors.New("invalid msg.route"), nil)

	} else {

		res, err := handler(ctx, session, msg)

		return result(err, res)
	}
}

func result(err error, msg interface{}) json.RawMessage {
	if err != nil {
		logx.Error("forwardMessage Component result failed,err:", err)

		return codec.Encode(err.Error(), nil)
	} else {
		return codec.Encode(nil, msg)
	}
}

type Route struct {
	Method  string
	Handler ForwardMessageHandler
}
