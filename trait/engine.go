package trait

// Engine 服务器引擎模块抽象层
type Engine interface {
	RouterGroup

	// 运行服务器
	Run() error
	// 停止服务器
	Stop()
}
