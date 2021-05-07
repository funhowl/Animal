package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type Connection struct {
	wsConnect *websocket.Conn
	inChan    chan []byte
	outChan   chan []byte
	closeChan chan byte

	connmutex sync.Mutex // 对相同ip的登录关闭上锁

	mutex    sync.Mutex // 对closeChan关闭上锁
	isClosed bool       // 防止closeChan被关闭多次

	//首字母大写的才能被调用
	Accid      string
	LoginState bool //默认false
}

type msgResult struct {
	Key    float64
	Result float64
	Info   string
}

func InitConnection(wsConn *websocket.Conn) (conn *Connection, err error) {
	conn = &Connection{
		wsConnect:  wsConn,
		inChan:     make(chan []byte, 1000),
		outChan:    make(chan []byte, 1000),
		closeChan:  make(chan byte, 1),
		LoginState: false,
	}
	// 启动读协程
	go conn.readLoop()
	// 启动写协程
	go conn.writeLoop()
	return
}

func (conn *Connection) ReadMessage() (data []byte, err error) {

	select {
	case data = <-conn.inChan:
	case <-conn.closeChan:
		err = errors.New("connection is closeed")
	}
	return
}

func (conn *Connection) WriteMessage(data []byte) (err error) {

	select {
	case conn.outChan <- data:
	case <-conn.closeChan:
		err = errors.New("connection is closeed")
	}
	return
}

func (conn *Connection) Send(it interface{}) (err error) {
	str, _ := json.Marshal(it)
	data := []byte(str)
	return conn.WriteMessage(data)
}

func (conn *Connection) ResultMsg(result float64, info string) (err error) {
	var ans msgResult
	ans.Key = 2
	ans.Result = result
	ans.Info = info

	str, _ := json.Marshal(ans)
	data := []byte(str)
	return conn.WriteMessage(data)
}

func (conn *Connection) Close() {
	// 线程安全，可多次调用
	conn.LoginState = false
	conn.wsConnect.Close()
	// 利用标记，让closeChan只关闭一次
	conn.mutex.Lock()

	if !conn.isClosed {
		fmt.Println("Connection close")
		close(conn.closeChan)
		conn.isClosed = true
	}
	conn.mutex.Unlock()
}

// 内部实现
func (conn *Connection) readLoop() {
	var (
		data []byte
		err  error
	)
	for {
		if _, data, err = conn.wsConnect.ReadMessage(); err != nil {
			goto ERR
		}
		//阻塞在这里，等待inChan有空闲位置
		select {
		case conn.inChan <- data:
		case <-conn.closeChan: // closeChan 感知 conn断开
			goto ERR
		}

	}

ERR:
	conn.Close()
}

func (conn *Connection) writeLoop() {
	var (
		data []byte
		err  error
	)

	for {
		select {
		case data = <-conn.outChan:
		case <-conn.closeChan:
			goto ERR
		}
		if err = conn.wsConnect.WriteMessage(websocket.TextMessage, data); err != nil {
			goto ERR
		}
	}

ERR:
	conn.Close()

}
