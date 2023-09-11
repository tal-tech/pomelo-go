package component

import (
	"context"
	"encoding/json"
)

type Handler func(ctx context.Context, in json.RawMessage) (out json.RawMessage)

type Component interface {
	Routes() (router map[string]Handler)
}
