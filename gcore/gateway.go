package gcore

import (
	"fmt"
	"net"
	"net/http"
	"reflect"

	"github.com/go75/gte/gconf"
	"github.com/go75/gte/glog"
	"github.com/go75/gte/trait"
	"github.com/gorilla/websocket"
)

// TCPGateway 网关模块，处理客户端TCP连接建立与注册
type TCPGateway struct {
	address net.TCPAddr
	version string

	listener *net.TCPListener

	connMgr trait.ConnMgr
	taskMgr trait.TaskMgr
}

var _ trait.Gateway = (*TCPGateway)(nil)

// NewTCPGateway 创建网关实例
func NewTCPGateway(connMgr trait.ConnMgr, taskMgr trait.TaskMgr) trait.Gateway {	
	address := net.TCPAddr{
		IP:   net.ParseIP(gconf.Config.ListenIP()),
		Port: gconf.Config.ListenPort(),
	}

	return &TCPGateway{
		address: address,
		version: gconf.Config.NetworkVersion(),
		connMgr: connMgr,
		taskMgr: taskMgr,
	}
}

// ListenAndServe 监听TCP连接并接收客户端连接，连接建立后注册到连接管理器
func (g *TCPGateway) ListenAndServe() error {
	var err error
	g.listener, err = net.ListenTCP(g.version, &g.address)
	if err != nil {
		return err
	}

	for {
		conn, err := g.Accept()
		if err != nil {
			glog.Error("Accept error:", err)
			continue
		}

		g.connMgr.Add(conn)
	}
}

// Accept 接收客户端连接
func (g *TCPGateway) Accept() (trait.Connection, error) {
	conn, err := g.listener.AcceptTCP()
	if err != nil {
		glog.Error("AcceptTCP error:", err)
		return nil, err
	}

	file, err := conn.File()
	if err != nil {
		glog.Error("Failed to get file descriptor:", err)
		return nil, err
	}

	connection := NewTCPConnection(uint64(file.Fd()), conn, g.connMgr, g.taskMgr)

	return connection, nil
}

// WebsocketGateway 网关模块，处理客户端Websocket连接建立与注册
type WebsocketGateway struct {
	upgrader *websocket.Upgrader
	address string
	connCh chan *websocket.Conn

	connMgr trait.ConnMgr
	taskMgr trait.TaskMgr
}

var _ trait.Gateway = (*WebsocketGateway)(nil)

func NewWebsocketGateway(connMgr trait.ConnMgr, taskMgr trait.TaskMgr) trait.Gateway {
	return &WebsocketGateway{
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		address: fmt.Sprintf("%s:%d", gconf.Config.ListenIP(), gconf.Config.ListenPort()),
		connCh:  make(chan *websocket.Conn, 1024),
		connMgr: connMgr,
		taskMgr: taskMgr,
	}
}

func (g *WebsocketGateway) ListenAndServe() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := g.upgrader.Upgrade(w, r, nil)
		if err != nil {
			glog.Error("websocket upgrade error:", err)
			return
		}

		g.connCh <- conn
	})

	go func() {
		for {
			conn, err := g.Accept()
			if err != nil {
				glog.Error("Accept websocket error:", err)
				continue
			}

			g.connMgr.Add(conn)
		}
	}()

	err := http.ListenAndServe(g.address, nil)

	return err
}

func (g *WebsocketGateway) Accept() (trait.Connection, error) {
	conn := <-g.connCh

	fd := g.websocketFD(conn)

	connection := NewWebsocketConnection(uint64(fd), conn, g.connMgr, g.taskMgr)

	return connection, nil
}

func (g *WebsocketGateway) websocketFD(conn *websocket.Conn) int32 {
	tcpConn := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	return int32(pfdVal.FieldByName("Sysfd").Int())
}
