package user

import "fmt"

// 实例化一个通过字符串映射函数切片的map
var eventByName = make(map[float64][]func(*Connection, map[string]interface{}))

// 注册事件，提供事件名和回调函数
func RegisterEvent(name float64, callback func(*Connection, map[string]interface{})) {

	// 通过名字查找事件列表
	list := eventByName[name]

	// 在列表切片中添加函数
	list = append(list, callback)

	// 将修改的事件列表切片保存回去
	eventByName[name] = list
}

// 调用事件
func CallEvent(name float64, conn *Connection, param map[string]interface{}) {

	// 通过名字找到事件列表
	list := eventByName[name]

	// 遍历这个事件的所有回调
	for _, callback := range list {

		// 传入参数调用回调
		fmt.Println("callback:", param)
		callback(conn, param)
	}
}
