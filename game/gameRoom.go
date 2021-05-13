package game

import (
	"Animal/common"
	"Animal/user"
	"encoding/json"
	"fmt"
	"reflect"
)

type room struct {
	Roomid int
	// count int
	Left  string // room 里面最多两个玩家 就直接用两个字段
	Right string
}

var roomList map[int]*room
var maxroomid int

type roomessage struct {
	Key  float64
	Data interface{}
}

func init() {

	roomList = make(map[int]*room)
	roomid := 5000

	for i := 0; i < 6; i++ { // 先 生成6个空房间
		var roomNew room
		roomNew.Roomid = roomid + i + 1
		roomList[roomNew.Roomid] = &roomNew
		if roomNew.Roomid > maxroomid {
			maxroomid = roomNew.Roomid
		}
	}

	user.RegisterEvent(common.C2Sroomessage, new(roomessage).roomessage) // 房间信息请求
	user.RegisterEvent(common.C2Sjoinroom, new(roomessage).joinroom)     //  加入房间请求
}

func (*roomessage) roomessage(conn *user.Connection, param map[string]interface{}) {

	var (
		err error
		ans roomessage
	)

	ans.Key = common.S2Croomessage
	fmt.Println("roomList:", roomList[5001].Roomid)
	a, _ := json.Marshal(roomList)
	fmt.Println("string(a):", string(a))
	ans.Data = roomList

	if err = conn.Send(ans); err != nil {
		panic(err.Error())
	}
}

func (*roomessage) joinroom(conn *user.Connection, param map[string]interface{}) {
	var (
		err        error
		ans        roomessage
		urow       *common.UserRow
		istartgame bool
	)

	if urow, err = user.GetRowById(conn.Accid); err != nil {
		panic(err.Error())
	}

	if !reflect.ValueOf(urow.GameMsg).IsZero() && urow.GameMsg.Roomid > 5000 {
		conn.ResultMsg(common.RoomRepeatErr, common.RoomRepeatErrMsg)
		return
	}

	ans.Key = common.S2Croomessage

	roomid := int(param["Roomid"].(float64))
	roomS := roomList[roomid]
	if reflect.ValueOf(roomS.Left).IsZero() { // 没人
		roomS.Left = conn.Accid

		if roomid == maxroomid && roomid < 5009 { // 由于本人前端水平 就只开10个房间
			for i := 0; i < 2; i++ { // 生成2个空房间
				var roomNew room
				roomNew.Roomid = roomid + i + 1
				roomList[roomNew.Roomid] = &roomNew
				if roomNew.Roomid > maxroomid {
					maxroomid = roomNew.Roomid
				}
			}
		}
	} else if reflect.ValueOf(roomS.Right).IsZero() { // 有一个

		roomS.Right = conn.Accid
		// start game
		istartgame = true
		// left := roomS.left

	} else { // 人满
		conn.ResultMsg(common.RoomEnoughErr, common.RoomEnoughErrMsg)
		return
	}

	roomList[roomid] = roomS
	ans.Data = roomList

	// uplist := *urow
	// uplist.GameMsg.Roomid = roomid
	// *urow = uplist

	var GameMsg common.GameRow
	GameMsg.Roomid = roomid
	GameMsg.Left = roomS.Left
	GameMsg.Right = roomS.Right
	urow.GameMsg = &GameMsg

	user.UserRowData[conn.Accid] = urow

	if err = conn.Send(ans); err != nil {
		panic(err.Error())
	}
	if istartgame {
		urowL := user.UserRowData[roomS.Left]
		urowL.GameMsg = &GameMsg
		user.UserRowData[roomS.Left] = urowL
		startgame(roomS)
	}
}

//离开房间暂时就不写了
