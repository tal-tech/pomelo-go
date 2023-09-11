package handler

import (
	"context"
	"errors"
	"pomelo-go/app/custom/internal/handler/custom"
	"pomelo-go/app/custom/internal/svc"
	"pomelo-go/cluster/clusterpb/proto"
	"pomelo-go/component/remote/backend"
)

func RegisterHandlers(server *backend.Component, serverCtx *svc.ServiceContext) {

	server.AddRoutes([]backend.Route{
		{"custom.customHandler.msgCustom1", custom.MsgCustom1Handler(serverCtx)}, // 自定义Handler
		{"custom.customHandler.msgCustom2", custom.MsgCustom2Handler(serverCtx)}, // 自定义Handler
		{"custom.customHandler.msgCustom3", NilHandler(serverCtx)},               // 自定义Handler
		{"custom.customHandler.msgCustom4", NilHandler(serverCtx)},               // 自定义Handler
	})
}

func NilHandler(ctx *svc.ServiceContext) backend.ForwardMessageHandler {
	return func(ctx context.Context, session proto.Session, message proto.Message) (interface{}, error) {
		return nil, errors.New("nil function")
	}
}
