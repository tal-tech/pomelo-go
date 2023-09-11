package types

type (
	MsgCustom1Request struct {
		A string
		B string
	}

	MsgCustom1Response struct {
		Str   string
		Slice []string
		Map   map[string]interface{}
	}
)
type (
	MsgCustom2Request  struct{}
	MsgCustom2Response struct{}
)
type (
	MsgCustom3Request  struct{}
	MsgCustom3Response struct{}
)
type (
	MsgCustom4Request  struct{}
	MsgCustom4Response struct{}
)
