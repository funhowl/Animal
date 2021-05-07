package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"Animal/common"
	"Animal/user"

	_ "Animal/scene"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("hello"))
	var (
		wsConn  *websocket.Conn
		err     error
		conn    *user.Connection
		message []byte
		ip      string
	)
	ip = common.ClientIP(r)
	fmt.Println(ip)

	fmt.Println("create new conn")

	// 完成ws协议的握手操作
	// Upgrade:websocket
	if wsConn, err = upgrader.Upgrade(w, r, nil); err != nil {
		return
	}

	if conn, err = user.InitConnection(wsConn); err != nil {
		goto ERR
	}

	// 启动线程，不断发消息
	// go func() {

	// 	for {
	// 		user.RegisterEvent(common.GAME_RUN, new(common.GAME_RUN_EVENT).GAME_RUN) // 帧事件
	// 		time.Sleep(100 * time.Second)
	// 	}
	// }()

	for {
		if message, err = conn.ReadMessage(); err != nil {
			goto ERR
		}
		fmt.Println(string(message)) //{"Name":"Wang Wu","Key":1000,"Weight":150.07}  {"Name":"Wang Wu","Key":1002,"Weight":150.07}
		// fmt.Print(ws.Connection.Accid)

		res := make(map[string]interface{})
		_ = json.Unmarshal(message, &res)

		// res, err := simplejson.NewJson(message)
		// if err != nil {
		// 	fmt.Printf("%v\n", err)
		// 	goto ERR
		// }
		// fmt.Println(reflect.TypeOf(res))
		// key, _ := res.Get("Key").Int64()
		key := res["Key"].(float64)

		if !conn.LoginState && key != common.C2Slogin {
			conn.ResultMsg(common.LoginErr, common.LoginErrMsg)
		} else {
			user.CallEvent(key, conn, res)
		}
		// if err = conn.WriteMessage([]byte("get this:" + ip)); err != nil {
		// 	goto ERR
		// }

	}

ERR:
	conn.Close()
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	http.ListenAndServe("0.0.0.0:7777", nil)
}
