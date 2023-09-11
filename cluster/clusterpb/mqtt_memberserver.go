package clusterpb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/timex"
	"pomelo-go/cluster/clusterpb/proto"
	"pomelo-go/tool"
)

type MqttMemberServer struct {
	memberServer MemberServer

	server *mqtt.Server
}

func (m *MqttMemberServer) Listen(advertiseAddr string) error {
	server := mqtt.New(nil)
	m.server = server

	ws := listeners.NewTCP("t1", advertiseAddr, nil)
	_ = ws

	err := server.AddHook(&hook{
		mqttMasterServer: m,
	}, nil)
	if err != nil {
		return err
	}

	err = server.AddListener(ws)
	if err != nil {
		return err
	}

	err = server.Serve()
	if err != nil {
		return err
	}
	return nil
}

func (m *MqttMemberServer) publishHandler(message rpcRequest) (*packets.Packet, error) {

	in := proto.RequestRequest(message.Msg)
	requestResponse, err := m.memberServer.RequestHandler(context.Background(), in)
	if err != nil {
		logx.Errorf("m.memberServer.RequestHandler failed, err:", err)
		return nil, err
	}

	response := rpcMessageResponse{
		Id:   *message.Id,
		Resp: json.RawMessage(requestResponse),
	}

	data, err := json.Marshal(response)
	if err != nil {
		logx.Errorf("server.Publish err response Marshal failed ,err:%s", err)
		return nil, err
	}

	res := &packets.Packet{
		FixedHeader: packets.FixedHeader{
			Type: packets.Publish,
		},
		TopicName: topic_RPC,
		Payload:   data,
	}

	return res, nil
}

func NewMqttMasterServer(memberServer MemberServer) *MqttMemberServer {

	return &MqttMemberServer{
		memberServer: memberServer,
		server:       nil,
	}
}

type hook struct {
	mqttMasterServer *MqttMemberServer

	mqtt.HookBase
}

func (h *hook) ID() string {
	return "mqtt_masterServer"
}

func (h *hook) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnConnect,
		mqtt.OnConnectAuthenticate,
		mqtt.OnDisconnect,

		mqtt.OnSubscribe,
		mqtt.OnSubscribed,
		mqtt.OnUnsubscribed,

		mqtt.OnPacketProcessed,

		mqtt.OnPublished,
		mqtt.OnPublish,
	}, []byte{b})
}

func (h *hook) OnConnect(cl *mqtt.Client, pk packets.Packet) error {
	h.Log.Info().Str("client", cl.ID).Msgf("client connected")
	return nil
}

func (h *hook) OnConnectAuthenticate(cl *mqtt.Client, pk packets.Packet) bool {
	// 全部放行
	logx.Infof("OnConnectAuthenticate , id:%s", cl.ID)
	return true
}

func (h *hook) OnDisconnect(cl *mqtt.Client, err error, expire bool) {
	h.Log.Info().Str("client", cl.ID).Bool("expire", expire).Err(err).Msg("client disconnected")
}

func (h *hook) OnSubscribe(cl *mqtt.Client, pk packets.Packet) packets.Packet {
	logx.Infof("OnSubscribe , id:%s", cl.ID)
	return pk
}

func (h *hook) OnSubscribed(cl *mqtt.Client, pk packets.Packet, reasonCodes []byte) {
	h.Log.Info().Str("client", cl.ID).Interface("filters", pk.Filters).Msgf("subscribed qos=%v", reasonCodes)

}

func (h *hook) OnUnsubscribed(cl *mqtt.Client, pk packets.Packet) {
	h.Log.Info().Str("client", cl.ID).Interface("filters", pk.Filters).Msg("unsubscribed")
}

func (h *hook) OnPacketProcessed(cl *mqtt.Client, pk packets.Packet, err error) {
	if err != nil {
		logx.Errorf("OnPacketProcessed failed , id:%s , err:%s", cl.ID, err)
		return
	}

	switch pk.FixedHeader.Type {

	case packets.Reserved:
		logx.Infof("OnPacketProcessed Reserved , id:%s, payload:%s", cl.ID, pk.Payload)

	case packets.Connect:
		logx.Infof("OnPacketProcessed Connect , id:%s, payload:%s", cl.ID, pk.Payload)

	case packets.Connack:
		logx.Infof("OnPacketProcessed Connack , id:%s, payload:%s", cl.ID, pk.Payload)

	case packets.Publish:

		m := rpcRequest{}

		err := json.Unmarshal(pk.Payload, &m)
		if err != nil {
			logx.Errorf("Publish json.Unmarshal failed ,err:%s", err)
			return
		}

		startTime := timex.Now()

		if m.Id == nil {

			logx.Errorf("[Notify] -- %s.%s.%s.%s - request:%s", m.Msg.Namespace, m.Msg.ServerType, m.Msg.Service, m.Msg.Method, tool.SimpleJson(m.Msg.Args))

		} else {

			id := *m.Id

			logx.Debugf("[Request] %d -- %s.%s.%s.%s - request:%s", id, m.Msg.Namespace, m.Msg.ServerType, m.Msg.Service, m.Msg.Method, tool.SimpleJson(m.Msg.Args))

			pkg, err := h.mqttMasterServer.publishHandler(m)

			var ( //  TODO 日志相关 优化
				duration = timex.ReprOfDuration(timex.Since(startTime))
				router   = fmt.Sprintf("%s.%s.%s.%s", m.Msg.Namespace, m.Msg.ServerType, m.Msg.Service, m.Msg.Method)
			)

			if err != nil {

				logx.Errorf("[Request ERR] %d %s -- %s - handler err:%s", id, duration, router, err.Error())
				return
			}
			err = cl.WritePacket(*pkg)
			if err != nil {
				logx.Errorf("[Request ERR] %d %s -- %s - writePacket err:%s", id, duration, router, err.Error())
				return
			}

			logx.Debugf("[Request OK] %d %s -- %s - response:%s", id, duration, router, string(pkg.Payload))
		}

	case packets.Puback:
		logx.Infof("OnPacketProcessed Puback , id:%s, payload:%s", cl.ID, pk.Payload)

	case packets.Pubrec:
		logx.Infof("OnPacketProcessed Pubrec , id:%s, payload:%s", cl.ID, pk.Payload)

	case packets.Pubrel:
		logx.Infof("OnPacketProcessed Pubrel , id:%s, payload:%s", cl.ID, pk.Payload)

	case packets.Pubcomp:
		logx.Infof("OnPacketProcessed Pubcomp , id:%s, payload:%s", cl.ID, pk.Payload)

	case packets.Subscribe:
		logx.Infof("OnPacketProcessed Subscribe , id:%s, payload:%s", cl.ID, pk.Payload)

	case packets.Suback:
		logx.Infof("OnPacketProcessed Suback , id:%s, payload:%s", cl.ID, pk.Payload)

	case packets.Unsubscribe:
		logx.Infof("OnPacketProcessed Unsubscribe , id:%s, payload:%s", cl.ID, pk.Payload)

	case packets.Unsuback:
		logx.Infof("OnPacketProcessed Unsuback , id:%s, payload:%s", cl.ID, pk.Payload)

	case packets.Pingreq:
		logx.Debugf("OnPacketProcessed Pingreq , id:%s, payload:%s", cl.ID, pk.Payload)

	case packets.Pingresp:
		logx.Infof("OnPacketProcessed Pingresp , id:%s, payload:%s", cl.ID, pk.Payload)

	case packets.Disconnect:
		logx.Infof("OnPacketProcessed Disconnect , id:%s, payload:%s", cl.ID, pk.Payload)

	case packets.Auth:
		logx.Infof("OnPacketProcessed Auth , id:%s, payload:%s", cl.ID, pk.Payload)

	case packets.WillProperties:
		logx.Infof("OnPacketProcessed WillProperties , id:%s, payload:%s", cl.ID, pk.Payload)

	default:
		logx.Infof("OnPacketProcessed default , id:%s, payload:%s", cl.ID, pk.Payload)

	}

}

func (h *hook) OnPublish(cl *mqtt.Client, pk packets.Packet) (packets.Packet, error) {
	h.Log.Info().Str("client", cl.ID).Str("payload", string(pk.Payload)).Msg("received from client")

	pkx := pk
	if string(pk.Payload) == "hello" {
		pkx.Payload = []byte("hello world")
		h.Log.Info().Str("client", cl.ID).Str("payload", string(pkx.Payload)).Msg("received modified packet from client")
	}

	return pkx, nil
}

func (h *hook) OnPublished(cl *mqtt.Client, pk packets.Packet) {
	h.Log.Info().Str("client", cl.ID).Str("payload", string(pk.Payload)).Msg("published to client")
}

type rpcRequest struct {
	Id  *int                 `json:"id,omitempty"`
	Msg proto.RequestRequest `json:"msg"`
}
