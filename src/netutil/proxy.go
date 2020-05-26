package netutil

import (
	"net"
	"log"
	"io"
	"flag"
	"os"
)

type TcpProxy struct {
	server_addr		string
	remote_addr		string
}

/**
 * 服务启动逻辑
 */
func (p *TcpProxy) Start() {
	flag.StringVar(&p.server_addr, "s", ":10000", "-s=0.0.0.0:10000 指定代理服务监听的地址")
	flag.StringVar(&p.remote_addr, "r", "", "-r=0.0.0.0:9999 指定远程服务器监听地址")
	flag.Parse()

	if p.remote_addr == "" {
		log.Println("远程服务端IP和端口不能空")
		os.Exit(1)
	}
	log.Println(p.server_addr)

	listener, err := net.Listen("tcp", p.server_addr)
	if err != nil {
		log.Println("listen error:" , err)
	}
	log.Println("listen success")

	for {
		clientConn, err := listener.Accept();

		if err != nil {
			log.Println("accept error:" , err)
			break
		}
		log.Println("client accept success,", clientConn.RemoteAddr())

		go p.handlerConn(clientConn);
	}
}

/**
 * 连接远程服务器，设置代理
 */
func (p *TcpProxy) handlerConn(clientConn net.Conn) {
	defer func() {
		clientConn.Close();
		log.Println("关闭客户端连接")
	}()

	serverConn, err := net.Dial("tcp", p.remote_addr)
	if err != nil {
		log.Println("server connect error:" , err)
	}
	log.Println("server connect success,ip", serverConn.RemoteAddr())

	defer func() {
		serverConn.Close();
		log.Println("关闭远程服务器连接")
	}()

	channel := make(chan bool, 1)

	go p.handlerServerConn(serverConn, clientConn, channel)
	go p.handlerServerConn(clientConn, serverConn, channel)

	<- channel
}

/**
 * 客户端和远程服务端数据转发
 */
func (proxy *TcpProxy) handlerServerConn(serverConn net.Conn, clientConn net.Conn, channel chan bool) {
	_, err := io.Copy(serverConn, clientConn)
	log.Println("转发信息失败:", err)

	channel <- true
}


/**
 * 服务初始化启动
 */
func Start() {
	tcpProxy := new(TcpProxy)
	tcpProxy.Start();
}
