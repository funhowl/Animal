package scene

import (
	"fmt"
	"time"

	"Animal/common"
	"Animal/user"
)

//{"Name":"Wang Wu","Key":1000,"Weight":150.07}
type login struct { // 回包格式
	Key    float64
	Sucess bool
	Now    time.Time
	Data   *common.UserRow
}

type Data struct {
	Key    float64
	Sucess bool
	Now    time.Time
}

func init() {
	// self := new(Login)
	user.RegisterEvent(common.C2Slogin, new(login).login) // 登录请求
}

func (*login) login(conn *user.Connection, param map[string]interface{}) {

	var (
		err     error
		ans     login
		oldconn *user.Connection
		urow    *common.UserRow
		uct     Data
		updata  *common.UserRow
		Accid   string
	)

	if acc, ok := param["code"]; ok { //现在都是第三方登录 应该写个登录服务器 这个accid 还有个token 是前端从登录服务器请求来的 还要验证一下 这里就随便写
		Accid = acc.(string)
	} else {
		conn.ResultMsg(common.MessageParamsErr, common.MessageParamsErrMsg)
		return
	}

	if oldconn, err = user.GetUser(Accid); err != nil {
		panic(err.Error())
	}

	if oldconn.LoginState {
		oldconn.ResultMsg(common.ClientConflict, common.ClientConflictMsg)
		oldconn.LoginState = false
	}
	conn.Accid = Accid
	conn.LoginState = true
	user.UserList[Accid] = conn
	// name, _ := param.Get("Name").String()

	ans.Key = common.S2Clogin
	ans.Sucess = true
	ans.Now = time.Now()

	if urow, err = user.GetRowById(Accid); err != nil {
		panic(err.Error())
	}
	ans.Data = urow

	//test UpdateUsersDB
	uplist := *urow

	uptask := *(uplist.Tasks)
	uptask.Record = "材不材间过此生"
	uptask.Point = uptask.Point + 1

	uct.Key = 110
	uct.Sucess = true
	uct.Now = time.Now()
	uptask.Uct = uct

	uplist.Gold = 10
	uplist.Tasks = &uptask

	if updata, err = common.UpdateUsersDB(urow, uplist); err != nil {
		panic(err.Error())
	}
	fmt.Println(updata)

	if err = conn.Send(ans); err != nil {
		panic(err.Error())
	}
}

// func merge(a, b Login) Login {

// 	jb, err := json.Marshal(b)
// 	if err != nil {
// 		fmt.Println("Marshal error b:", err)
// 	}
// 	err = json.Unmarshal(jb, &a)
// 	if err != nil {
// 		fmt.Println("Unmarshal error b-a:", err)
// 	}

// 	return a
// }
