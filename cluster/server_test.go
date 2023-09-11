package cluster

import (
	"context"
	"encoding/json"
	"fmt"
	"pomelo-go/cluster/clusterpb/proto"
	"pomelo-go/component"
	"testing"
)

type MyComponent struct {
}

func (m MyComponent) Routes() (router map[string]component.Handler) {
	return nil
}

func (m MyComponent) Init() {
	fmt.Println("MyComponent.Init")
}

func (m MyComponent) AfterInit() {
	fmt.Println("MyComponent.AfterInit")
}

func (m MyComponent) BeforeShutdown() {
	fmt.Println("MyComponent.BeforeShutdown")
}

func (m MyComponent) Shutdown() {
	fmt.Println("MyComponent.Shutdown")
}

func TestNode_Startup(t *testing.T) {

	c := &component.Components{}
	c.Register(&MyComponent{}, nil)

	serverid := "cluster-server-custom-996"

	opt := Config{
		IsMaster:      false,
		ServerId:      serverid,
		AdvertiseAddr: "localhost:3005", // master 地址
		ServerInfo: proto.ClusterServerInfo{
			"serverType": proto.ServerType_Recover,
			"id":         serverid,
			"type":       proto.Type_Monitor,
			"pid":        99,
			"info": map[string]interface{}{ // 本地服务信息
				"serverType": proto.ServerType_Recover,
				"id":         serverid,
				"env":        "local",
				"host":       "127.0.0.1",
				"port":       8081,

				"channelType":   2,
				"cloudType":     1,
				"clusterCount":  1,
				"restart-force": "true",
			},
		},
		RetryInterval: 5,
		RetryTimes:    60,
		Token:         "agarxhqb98rpajloaxn34ga8xrunpagkjwlaw3ruxnpaagl29w4rxn",
		Listen:        "127.0.0.1:8081", // 本地服务地址
	}

	n := &Server{
		cnf:        opt,
		components: c,
	}

	err := n.Startup()
	if err != nil {
		t.Fatal(err)
	}

	args := []interface{}{
		"stu1*kick_testsss",
		"cluster-server-connector-0",
		"kick_testsss",
		true,
		2,
		1,
		0,
		"123",
		"abc",
		"2.9.8.7",
		map[string]interface{}{
			"uniqId":    "231FF2BB-BA09-598D-9EB6-3B0299D292E7ssss",
			"rid":       "kick_testsss",
			"rtype":     2,
			"role":      1,
			"ulevel":    0,
			"uname":     "123",
			"classid":   "abc",
			"clientVer": "2.9.8.7",
			"userVer":   "1.0",
			"liveType":  "COMBINE_SMALL_CLASS_MODE"},
		"0"}

	body, err := json.Marshal(args)
	if err != nil {
		t.Fatal(err)
	}

	res, err := n.RemoteProcess(context.Background(), proto.RequestRequest{
		Namespace:  "user",
		ServerType: "chat",
		Service:    "chatRemote",
		Method:     "add",
		Args:       body,
	})

	if err != nil {
		t.Fatal(err)
	}

	t.Log(res)

	select {}
}
