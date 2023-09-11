package proto

import "encoding/json"

type MonitorAction string

const (
	Type_Monitor = "monitor"

	ServerType_Connector = "connector"
	ServerType_Chat      = "chat"
	ServerType_Recover   = "custom"

	MonitorAction_addServer     MonitorAction = "addServer"
	MonitorAction_removeServer  MonitorAction = "removeServer"
	MonitorAction_replaceServer MonitorAction = "replaceServer"
	MonitorAction_startOve      MonitorAction = "startOver"
)

const (
	BEFORE_FILTER        = "__befores__"
	AFTER_FILTER         = "__afters__"
	GLOBAL_BEFORE_FILTER = "__globalBefores__"
	GLOBAL_AFTER_FILTER  = "__globalAfters__"
	ROUTE                = "__routes__"
	BEFORE_STOP_HOOK     = "__beforeStopHook__"
	MODULE               = "__modules__"
	SERVER_MAP           = "__serverMap__"
	RPC_BEFORE_FILTER    = "__rpcBefores__"
	RPC_AFTER_FILTER     = "__rpcAfters__"
	MASTER_WATCHER       = "__masterwatcher__"
	MONITOR_WATCHER      = "__monitorwatcher__"
)

// ClusterServerInfo 集群服务信息
type ClusterServerInfo map[string]interface{}

// Register 向master注册服务信息
type (
	RegisterRequest struct {
		ServerInfo ClusterServerInfo

		Token string
	}

	RegisterResponse struct{}
)

// Subscribe 订阅master中集群信息
type (
	SubscribeRequest struct {
		Id string `json:"id"`
	}

	SubscribeResponse map[string]ClusterServerInfo // 集群内其他服务信息
)

// Record 通知master启动完毕
type (
	RecordRequest struct {
		Id string `json:"id"`
	}

	RecordResponse struct{}
)

// MonitorHandler 监听master中的集群变化
type (
	MonitorHandlerRequest struct {
		CallBackHandler func(action MonitorAction, serverInfos []ClusterServerInfo)
	}

	MonitorHandlerResponse struct{}
)

// Request 发送Request rpc请求
type (
	RequestRequest struct {
		Namespace  string          `json:"namespace"`
		ServerType string          `json:"serverType"`
		Service    string          `json:"service"`
		Method     string          `json:"method"`
		Args       json.RawMessage `json:"args"` // []interface{}{}
	}

	RequestResponse json.RawMessage // []interface{}{}
)

type Session struct {
	Id         int     `json:"id"`
	FrontendId string  `json:"frontendId"`
	Uid        string  `json:"uid"`
	Settings   Setting `json:"settings"`
}

type Setting struct {
	UniqId    string `json:"uniqId"`
	Rid       string `json:"rid"`
	Rtype     int    `json:"rtype"`
	Role      int    `json:"role"`
	Ulevel    int    `json:"ulevel"`
	Uname     int    `json:"uname"`
	Classid   string `json:"classid"`
	ClientVer string `json:"clientVer"`
	UserVer   string `json:"userVer"`
	LiveType  int    `json:"liveType"`
}

type Message struct {
	Id            int             `json:"id"`
	Type          int             `json:"type"`
	CompressRoute int             `json:"compressRoute"`
	Route         string          `json:"route"`
	CompressGzip  int             `json:"compressGzip"`
	Body          json.RawMessage `json:"body"`
}
