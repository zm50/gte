package constant

const (
	ConnStartSignal uint8 = iota
	ConnStopSignal
	ConnNotActiveSignal
)

const (
	// 接收对应连接的心跳包时连接设置为活跃状态
	ConnActiveState uint32 = iota
	// 连接健康状态巡检时如果为活跃状态则设置为检测状态
	ConnInspectState
	// 连接关闭状态
	ConnCloseState
	// 连接健康状态巡检时如果为检测状态则设置为不活跃状态
	ConnNotActiveState
)
