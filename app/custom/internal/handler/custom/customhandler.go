package custom

import (
	"context"
	"encoding/json"
	"pomelo-go/app/custom/internal/logic/custom"
	"pomelo-go/app/custom/internal/svc"
	"pomelo-go/app/custom/internal/types"
	"pomelo-go/cluster/clusterpb/proto"
	"pomelo-go/component/remote/backend"
)

func MsgCustom1Handler(serverCtx *svc.ServiceContext) backend.ForwardMessageHandler {
	return func(ctx context.Context, session proto.Session, message proto.Message) (interface{}, error) {
		var req types.MsgCustom1Request
		if err := json.Unmarshal(message.Body, &req); err != nil {
			return nil, err
		}

		l := custom.NewRecordLogic(ctx, serverCtx)
		resp, err := l.MsgCustom1(req)
		return resp, err
	}
}

func MsgCustom2Handler(serverCtx *svc.ServiceContext) backend.ForwardMessageHandler {
	return func(ctx context.Context, session proto.Session, message proto.Message) (interface{}, error) {
		var req types.MsgCustom2Request
		if err := json.Unmarshal(message.Body, &req); err != nil {
			return nil, err
		}

		l := custom.NewRecordLogic(ctx, serverCtx)
		resp, err := l.MsgCustom2(req)
		return resp, err
	}
}
