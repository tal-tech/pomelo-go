package custom

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"pomelo-go/app/custom/internal/logic"
	"pomelo-go/app/custom/internal/svc"
	"pomelo-go/app/custom/internal/types"
)

type custom struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func (r *custom) MsgCustom1(req types.MsgCustom1Request) (res types.MsgCustom1Response, err error) {
	return types.MsgCustom1Response{
		Str:   "AAA",
		Slice: []string{"BBB", "CCC"},
		Map: map[string]interface{}{
			"DDD": "EEEE",
		},
	}, nil
}

func (r *custom) MsgCustom2(req types.MsgCustom2Request) (res types.MsgCustom2Response, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *custom) MsgCustom3(req types.MsgCustom3Request) (res types.MsgCustom3Response, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *custom) MsgCustom4(req types.MsgCustom4Request) (res types.MsgCustom4Response, err error) {
	//TODO implement me
	panic("implement me")
}

func NewRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) logic.Record {
	logger := logx.WithContext(ctx)
	return &custom{
		Logger: logger,
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
