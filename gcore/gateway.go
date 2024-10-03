package gcore

import (
	"fmt"
	"net"

	"github.com/go75/gte/gconf"
	"github.com/go75/gte/trait"
)

// Gateway 网关模块，处理客户端连接建立与注册
type Gateway struct {
	address net.TCPAddr
	version string

	listener *net.TCPListener

	connMgr trait.ConnMgr
}

var _ trait.Gateway = (*Gateway)(nil)

// NewGateway 创建网关实例
func NewGateway(connMgr trait.ConnMgr) trait.Gateway {	
	address := net.TCPAddr{
		IP:   net.ParseIP(gconf.Config.ListenIP()),
		Port: gconf.Config.ListenPort(),
	}

	return &Gateway{
		address: address,
		version: gconf.Config.NetworkVersion(),
		connMgr: connMgr,
	}
}

// ListenAndServe 监听并接收客户端连接，连接建立后注册到连接管理器
func (g *Gateway) ListenAndServe() error {
	var err error
	g.listener, err = net.ListenTCP(g.version, &g.address)
	if err != nil {
		return err
	}

	for {
		conn, err := g.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}

		g.connMgr.Add(conn)
	}
}

// Accept 接收客户端连接
func (g *Gateway) Accept() (trait.Connection, error) {
	conn, err := g.listener.AcceptTCP()
	if err != nil {
		fmt.Println("Accept error:", err)
		return nil, err
	}

	file, err := conn.File()
	if err != nil {
		fmt.Println("File error:", err)
		return nil, err
	}

	connection := NewConnection(int32(file.Fd()), conn)

	return connection, nil
}

// Stop 停止网关
func (g *Gateway) Stop() error {
	return g.listener.Close()
}
