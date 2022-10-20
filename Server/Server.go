package main

import (
	"encoding/json"
	"log"
	"net"
	"os/exec"
)

const MessageLimitSize = 1024

//var socket []net.Conn

type socket struct {
	conn     net.Conn
	clientId int
}

type ReturnMsg struct {
	Code int                    `json:"code"`
	Msg  string                 `json:"msg,omitempty"`
	Data map[string]interface{} `json:"data,omitempty"`
}

var sockets []socket

func main() {

	addr := "localhost:8088"

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	// 清理当前输出屏幕
	_ = exec.Command("clear").Run()

	for {
		// 接受请求
		// 保存请求句柄
		// 收到客户端消息时，对正在活跃的客户端进行消息群发
		// 客户端退出时，要及时清理保存的句柄信息
		accept, err := listen.Accept()
		if err != nil {
			panic(err)
		}
		go resolveRequest(accept)
	}
}

// 能否在连接的时候发送客户端的一些信息，这样有助于标识客户端 todo
func resolveRequest(conn net.Conn) {

	// 分配连接ID
	// 存储 handler 并转化
	clientId := len(sockets)
	sockets = append(sockets, socket{
		conn:     conn,
		clientId: clientId,
	})
	data := make(map[string]interface{})
	data["clientId"] = sockets[clientId].clientId
	// 下发至客户端被分配到的id标识
	sockets[clientId].sendMsg(0, "issue Client Id", data)
	// todo 设置连接等待时长限制

	// 等待接受发送过来的消息
	for {
		buf := make([]byte, MessageLimitSize)
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("Read Err : " + err.Error())
			break
		}
		data := make(map[string]interface{})
		data["from"] = clientId
		// 使用[]byte直接进行json转码会发生前后转码不一致的情况
		data["msg"] = string(buf[:n])
		// 接收到消息后进行广播至其他客户端
		broadcast(ReturnMsg{
			Code: 1,
			Msg:  "BroadCast From Client",
			Data: data,
		}, clientId)
	}
	//defer conn.Close()
	//defer destructConnection(clientId)
}

// destructConnection 处理关闭的连接句柄，并将其从当前活跃客户端中清除
func destructConnection(clientId int) {
	sockets = append(sockets[:clientId], sockets[clientId+1:]...)
}

// 消息广播，下发至当前活跃的所有客户端连接  code = 1
func broadcast(msg ReturnMsg, hoster int) {
	for key, soc := range sockets {
		if key != hoster {
			soc.sendMsg(msg.Code, msg.Msg, msg.Data)
		}
	}
}

// sendMsg 向客户端发送信息
func (s *socket) sendMsg(code int, msg string, data map[string]interface{}) {
	toJson, err := json.Marshal(ReturnMsg{
		Code: code,
		Msg:  msg,
		Data: data,
	})
	if err != nil {
		panic(err)
	}
	s.conn.Write(toJson)
}
