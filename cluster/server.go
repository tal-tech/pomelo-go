package cluster

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"pomelo-go/cluster/clusterpb"
	"pomelo-go/cluster/clusterpb/proto"
	"pomelo-go/component"
	"time"
)

type Config struct {
	IsMaster   bool                    // 目前不支持master
	Listen     string                  // node服务 address (RPC)
	ServerId   string                  // node服务id名称
	ServerInfo proto.ClusterServerInfo // node 服务信息用于向master注册

	AdvertiseAddr string // node服务对应的master地址
	RetryInterval int    // master 重试间隔 default 3s
	RetryTimes    int    // master 重试间隔 default 10次

	Token string // master 通信token
}

type Server struct {
	cnf        Config
	components *component.Components
	die        chan bool // 关闭标记

	handler      *LocalHandler          // 处理本地或远程Handler调用
	masterClient clusterpb.MasterClient // 与master通信的客户端 对应pomelo的monitor
	rpcClient    *rpcClient
	server       *clusterpb.MqttMemberServer // 本地rpc服务端

	//sessions map[int64]*session.Session
}

func MustNewServer(c Config, opts ...RunOption) *Server {
	server, err := NewServer(c, opts...)
	if err != nil {
		logx.Must(err)
	}

	return server
}

func NewServer(c Config, opts ...RunOption) (*Server, error) {

	node := &Server{
		cnf:        c,
		components: &component.Components{},
		die:        make(chan bool),
	}

	for _, opt := range opts {
		opt(node)
	}

	if node.cnf.RetryInterval == 0 {
		return nil, errors.New("RetryInterval is 0")
	}

	if node.cnf.RetryTimes == 0 {
		return nil, errors.New("RetryTimes is 0")
	}

	return node, nil
}

func (s *Server) Start() {
	err := s.Startup()

	if err != nil {
		logx.Must(err)
	}

	select {
	case <-s.die:
		logx.Info("The app will shutdown in a few seconds")
	}

}

// Stop stops the Server.
func (s *Server) Stop() {

	//// reverse call `BeforeShutdown` hooks
	//components := s.components.List()
	//length := len(components)
	//for i := length - 1; i >= 0; i-- {
	//	components[i].Comp.BeforeShutdown()
	//}
	//// reverse call `Shutdown` hooks
	//for i := length - 1; i >= 0; i-- {
	//	components[i].Comp.Shutdown()
	//}

	//_, err = client.Unregister(context.Background(), request)

	close(s.die)
	_ = logx.Close()
}

func (s *Server) Startup() error {
	if s.cnf.Listen == "" {
		return errors.New("service address cannot be empty in master node")
	}

	s.rpcClient = newRPCClient()
	s.handler = NewHandler(s)

	components := s.components.List()
	for _, c := range components {
		err := s.handler.register(c.Comp, c.Opts)
		if err != nil {
			return err
		}
	}

	if err := s.initNode(); err != nil {
		return err
	}

	//// Initialize all components
	//for _, c := range components {
	//	c.Comp.Init()
	//}
	//for _, c := range components {
	//	c.Comp.AfterInit()
	//}

	return nil
}

func (s *Server) Handler() *LocalHandler {
	return s.handler
}

// RemoteProcess 远程调用
func (s *Server) RemoteProcess(ctx context.Context, in proto.RequestRequest) (proto.RequestResponse, error) {
	return s.handler.remoteProcess(ctx, in)
}

func (s *Server) RequestHandler(ctx context.Context, in proto.RequestRequest) (proto.RequestResponse, error) {

	return s.handler.localProcess(ctx, in)
}

func (s *Server) initNode() error {
	if !s.cnf.IsMaster && s.cnf.AdvertiseAddr == "" {
		return errors.New("invalid AdvertiseAddr")
	}

	s.server = clusterpb.NewMqttMasterServer(s)

	err := s.server.Listen(s.cnf.Listen)
	if err != nil {
		return err
	}

	mqttMasterClient := clusterpb.NewMqttMasterClient(s.cnf.AdvertiseAddr)

	retryTimes := s.cnf.RetryTimes

	for retryTimes > 0 {
		err := mqttMasterClient.Connect()
		if err == nil {
			break
		}

		time.Sleep(time.Duration(s.cnf.RetryInterval) * time.Second)
		logx.Info("try connect again, retryTimes :", retryTimes)

		retryTimes--
	}

	_, err = mqttMasterClient.MonitorHandler(context.Background(), &proto.MonitorHandlerRequest{
		CallBackHandler: func(action proto.MonitorAction, serverInfos []proto.ClusterServerInfo) { // 收到master推送的消息变更

			switch action {
			case proto.MonitorAction_addServer:
				for i := 0; i < len(serverInfos); i++ {

					remoteService, err := transformRemoteServiceInfo(serverInfos[i])
					if err != nil {
						logx.Error("transformRemoteServiceInfo failed,err:", err)
						continue
					}
					s.handler.addRemoteService(remoteService)
				}

			case proto.MonitorAction_removeServer:
				for i := 0; i < len(serverInfos); i++ {
					s.handler.delMember("")
				}
			case proto.MonitorAction_replaceServer:

			case proto.MonitorAction_startOve:
			}

		},
	})
	if err != nil {
		return err
	}

	_, err = mqttMasterClient.Register(context.Background(), &proto.RegisterRequest{
		ServerInfo: s.cnf.ServerInfo,
		Token:      s.cnf.Token,
	})
	if err != nil {
		return err
	}

	// 获取注册信息
	subscribeResponse, err := mqttMasterClient.Subscribe(context.Background(), &proto.SubscribeRequest{
		Id: s.cnf.ServerId,
	})

	// 初始化handler
	rs := make([]RemoteServiceInfo, 0, len(*subscribeResponse))
	for _, info := range *subscribeResponse {

		remoteService, err := transformRemoteServiceInfo(info)
		if err != nil {
			logx.Error("transformRemoteServiceInfo failed,err:", err)
			continue
		}

		rs = append(rs, remoteService)
	}

	s.handler.initRemoteService(rs)

	s.masterClient = mqttMasterClient
	return nil
}

// Enable current server accept connection
// 与pomelo端通信
func (s *Server) listenAndServe() {

}

func transformRemoteServiceInfo(info proto.ClusterServerInfo) (res RemoteServiceInfo, err error) {

	var (
		host       string
		port       int
		serverType string
	)

	if v, ok := info["host"]; !ok {
		return RemoteServiceInfo{}, errors.New("invalid host")
	} else {
		host = v.(string)
	}

	if v, ok := info["port"]; !ok {
		return RemoteServiceInfo{}, errors.New("invalid port")
	} else {
		port = int(v.(float64))
	}
	if v, ok := info["serverType"]; !ok {
		return RemoteServiceInfo{}, errors.New("invalid serverType")
	} else {
		serverType = v.(string)
	}

	return RemoteServiceInfo{
		ClusterServerInfo: info,
		ServerType:        serverType,
		ServiceAddr:       fmt.Sprintf("%s:%d", host, port),
	}, nil

}
