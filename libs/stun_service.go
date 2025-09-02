package libs

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/pion/stun"
)

// STUNServer 表示 STUN 服务器实例
type STUNServer struct {
	conn      net.PacketConn
	isRunning bool
	mu        sync.RWMutex
	port      string
}

// NewSTUNServer 创建新的 STUN 服务器实例
func NewSTUNServer(port string) *STUNServer {
	return &STUNServer{
		port: port,
	}
}

// Start 启动 STUN 服务器
func (s *STUNServer) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return fmt.Errorf("STUN server is already running")
	}

	// 创建 UDP 监听器
	conn, err := net.ListenPacket("udp4", ":"+s.port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %v", s.port, err)
	}

	s.conn = conn
	s.isRunning = true

	// 启动 STUN 服务器 goroutine
	go s.run()

	log.Printf("STUN server started on port %s", s.port)
	return nil
}

// GetSTUNAddress 获取 STUN 服务器地址
func (s *STUNServer) GetSTUNAddress() string {
	return fmt.Sprintf("stun:%s", s.port)
}

// run 运行 STUN 服务器的主循环
func (s *STUNServer) run() {
	buf := make([]byte, 1024)

	for {
		// 读取数据包
		n, addr, err := s.conn.ReadFrom(buf)
		if err != nil {
			log.Printf("Error reading from connection: %v", err)
			continue
		}

		// 处理 STUN 消息
		s.handleSTUNMessage(addr, buf[:n])
	}
}

// handleSTUNMessage 处理 STUN 消息
func (s *STUNServer) handleSTUNMessage(addr net.Addr, buf []byte) {
	// 解析 STUN 消息
	msg := &stun.Message{Raw: buf}
	if err := msg.Decode(); err != nil {
		// 不是有效的 STUN 消息，可能是其他协议
		return
	}

	// 只处理绑定请求
	if msg.Type != stun.BindingRequest {
		return
	}

	// 创建响应消息
	response, err := stun.Build(stun.BindingSuccess, stun.TransactionID)
	if err != nil {
		log.Printf("Error building STUN response: %v", err)
		return
	}

	// 添加 XOR-MAPPED-ADDRESS 属性 - 使用正确的方法
	udpAddr := addr.(*net.UDPAddr)
	xorAddr := stun.XORMappedAddress{
		IP:   udpAddr.IP,
		Port: udpAddr.Port,
	}

	// 使用 AddTo 方法而不是 Add 方法
	if err := xorAddr.AddTo(response); err != nil {
		log.Printf("Error adding XOR-MAPPED-ADDRESS: %v", err)
		return
	}

	// 发送响应
	if _, err := s.conn.WriteTo(response.Raw, addr); err != nil {
		log.Printf("Error sending response: %v", err)
	}

	log.Printf("Handled STUN request from %s", addr.String())
}

// Close 关闭 STUN 服务器
func (s *STUNServer) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return nil
	}

	if err := s.conn.Close(); err != nil {
		return fmt.Errorf("error closing STUN server: %v", err)
	}

	s.isRunning = false
	log.Println("STUN server stopped")
	return nil
}
