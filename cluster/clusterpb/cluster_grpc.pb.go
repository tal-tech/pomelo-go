package clusterpb

import (
	"context"
	"pomelo-go/cluster/clusterpb/proto"
)

// MasterClient 与master的双向通信（请求响应式+推送式）
type MasterClient interface {
	// Register 向master注册服务信息
	Register(ctx context.Context, in *proto.RegisterRequest) (*proto.RegisterResponse, error)
	// Subscribe 订阅master中集群信息
	Subscribe(ctx context.Context, in *proto.SubscribeRequest) (*proto.SubscribeResponse, error)
	// Record 通知master启动完毕
	Record(ctx context.Context, in *proto.RecordRequest) (*proto.RecordResponse, error)
	// MonitorHandler 监听master中的集群变化
	MonitorHandler(ctx context.Context, in *proto.MonitorHandlerRequest) (*proto.MonitorHandlerResponse, error)
}

type MasterServer interface {
}

// MemberClient 与服务rpc的双向通信（请求响应式）
type MemberClient interface {
	// Request 发送Request rpc请求
	Request(ctx context.Context, in proto.RequestRequest) (proto.RequestResponse, error)
}

// MemberServer 服务rpc的双向通信（请求响应式）
type MemberServer interface {
	// RequestHandler 处理Request rpc请求
	RequestHandler(ctx context.Context, in proto.RequestRequest) (proto.RequestResponse, error)
}

type MasterClientAgent interface {
	MasterClient

	Connect() error
	Close() error
}

type MemberClientAgent interface {
	MemberClient

	Connect() error
	Close() error
}

//type MemberServerAgent interface {
//	MemberServer
//
//	Listen(advertiseAddr string) error
//}
