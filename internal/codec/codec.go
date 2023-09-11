package codec

import (
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
)

func Decode(in json.RawMessage, out []interface{}) error {

	err := json.Unmarshal(in, &out)
	if err != nil {
		return err
	}

	return nil
}

func Encode(in ...interface{}) (out json.RawMessage) {

	out, err := json.Marshal(in)
	if err != nil {
		logx.Error("result json.Marshal failed ,err:", err)
	}

	return out
}
