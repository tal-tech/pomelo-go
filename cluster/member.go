package cluster

import (
	"time"
)

type Member struct {
	//isMaster bool
	//memberInfo      *clusterpb.MemberInfo
	lastHeartbeatAt time.Time // cluster member report heartbeat time to the master
}

//func (m *Member) MemberInfo() *clusterpb.MemberInfo {
//	return m.memberInfo
//}

func (m *Member) String() string {
	//return fmt.Sprintf("Master: %t MemberInfo: %s LastHeartbeatAt: %s", m.isMaster, m.memberInfo.String(), m.lastHeartbeatAt.String())
	return ""
}
