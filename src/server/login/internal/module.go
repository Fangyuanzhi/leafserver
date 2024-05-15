package internal

import (
	"github.com/gomodule/redigo/redis"
	"github.com/name5566/leaf/module"
	"server/base"
)

var (
	skeleton  = base.NewSkeleton()
	ChanRPC   = skeleton.ChanRPCServer
	RedisConn *redis.Pool
)

type Module struct {
	*module.Skeleton
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton
}

func (m *Module) OnDestroy() {

}
