package common

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"strings"
)

func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

func ChangeToString(arg interface{}) string {
	switch arg.(type) {
	case string:
		return arg.(string)
	case int:
		return strconv.Itoa(arg.(int))
	default:
		data, _ := json.Marshal(arg)
		return string(data)
	}
}

// //指针类型判空
// func CheckTypeByReflectNil(arg interface{}) {
// 	if reflect.ValueOf(arg).IsNil() { //利用反射直接判空，指针用isNil
// 		// 函数解释：isNil() bool	判断值是否为 nil
// 		// 如果值类型不是通道（channel）、函数、接口、map、指针或 切片时发生 panic，类似于语言层的v== nil操作
// 		fmt.Printf("反射判断：数据类型为%s,数据值为：%v,nil：%v \n",
// 			reflect.TypeOf(arg).Kind(), reflect.ValueOf(arg), reflect.ValueOf(arg).IsValid())
// 	}
// }

// //基础类型判空
// func CheckTypeByReflectZero(arg interface{}) {
// 	if reflect.ValueOf(arg).IsZero() { //利用反射直接判空，基础数据类型用isZero
// 		fmt.Printf("反射判断：数据类型为%s,数据值为：%v,nil：%v \n",
// 			reflect.TypeOf(arg).Kind(), reflect.ValueOf(arg), reflect.ValueOf(arg).IsValid())
// 	}
// }
