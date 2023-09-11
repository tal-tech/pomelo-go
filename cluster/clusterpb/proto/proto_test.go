package proto

import (
	"encoding/json"
	"testing"
)

func Test_RequestRequest(t *testing.T) {

	var (
		s = `{
        "namespace": "sys",
        "serverType": "custom",
        "service": "msgRemote",
        "method": "forwardMessage",
        "args": [
            {
                "id": 6,
                "type": 0,
                "compressRoute": 0,
                "route": "custom.recoverHandler.msgRecoverTemporary",
                "body": {
                    "liveId": "live-1",
                    "lecturerId": "l-1",
                    "tutorId": "t-1",
                    "stuId": "s-1-3"
                },
                "compressGzip": 0
            },
            {
                "id": 3,
                "frontendId": "cluster-server-connector-0",
                "uid": "10000*bench_0_1693555363_0",
                "settings": {
                    "uniqId": "5672155583468750658",
                    "rid": "bench_0_1693555363_0",
                    "rtype": 2,
                    "role": 1,
                    "ulevel": 1,
                    "uname": 10000,
                    "classid": "bench_0_1693555363_0",
                    "clientVer": "0.0.1",
                    "userVer": "1.1",
                    "liveType": 1
                }
            }
        ]
    }`
	)

	var in RequestRequest

	err := json.Unmarshal([]byte(s), &in)
	if err != nil {
		t.Fatal(err)
	}

	var (
		session = Session{}
		msg     = Message{}
	)

	arg := []interface{}{
		&msg,
		&session,
	}

	err = json.Unmarshal(in.Args, &arg)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(session, msg)
}
