package global

import (
	"github.com/zm50/gte/core"
	"github.com/zm50/gte/trait"
)

var msgPool trait.ObjPool[trait.Message]

func SetMsgPool(size int, messageProvider func() trait.Message) {
	msgPool = core.NewObjPool[trait.Message](size, messageProvider)
}

func MsgPool() trait.ObjPool[trait.Message] {
	return msgPool
}
