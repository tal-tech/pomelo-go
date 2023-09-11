package clusterpb

import (
	"log"
	"time"
)

var (
	memberClient MemberClientAgent
)

func InitMqttMemberClient() {
	var (
		advertiseAddr = "127.0.0.1:10061"
	)

	c := NewMqttMemberClient(advertiseAddr)

	for {
		err := c.Connect()
		if err == nil {
			break
		}

		time.Sleep(5 * time.Second)
		log.Println("try connect again")
	}

	memberClient = c
}
