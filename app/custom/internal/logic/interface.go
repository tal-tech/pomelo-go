package logic

import "pomelo-go/app/custom/internal/types"

type Record interface {
	MsgCustom1(req types.MsgCustom1Request) (res types.MsgCustom1Response, err error)
	MsgCustom2(req types.MsgCustom2Request) (res types.MsgCustom2Response, err error)
	MsgCustom3(req types.MsgCustom3Request) (res types.MsgCustom3Response, err error)
	MsgCustom4(req types.MsgCustom4Request) (res types.MsgCustom4Response, err error)
}
