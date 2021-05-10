package user

import (
	"Animal/common"

	mapset "github.com/deckarep/golang-set"
)

type UserMgr struct {
	Accid string
}

var UserList map[string]*Connection
var UserRowData map[string]*common.UserRow
var UserAcciSet mapset.Set

func init() {
	UserList = make(map[string]*Connection)
	UserRowData = make(map[string]*common.UserRow)
	UserAcciSet = mapset.NewSet()

	var accidRow UserMgr
	sqlStr := "select Accid from users" //占位符
	rowObj, err := common.MysqlDb.Query(sqlStr)
	if err != nil {
		panic(err.Error())
	}

	for rowObj.Next() {
		//调用scan函数拿到结果,映射到结构体中
		rowObj.Scan(&accidRow.Accid)
		// fmt.Println("Accid:", accidRow.Accid)
		UserAcciSet.Add(accidRow.Accid)
	}
}

func GetRowById(Accid string) (row *common.UserRow, err error) {
	if row, ok := UserRowData[Accid]; ok {
		return row, err
	}
	if UserAcciSet.Contains(Accid) {
		if row, err = common.SelUserRow(Accid); err != nil {
			panic(err.Error())
		}
	} else {
		if row, err = common.InitUserRow(Accid); err != nil {
			panic(err.Error())
		}
		UserAcciSet.Add(Accid)
	}
	UserRowData[Accid] = row
	return row, err
}

func GetUser(Accid string) (conn *Connection, err error) {

	if conn, ok := UserList[Accid]; ok {
		return conn, err
	}
	conn = &Connection{
		Accid:      Accid,
		LoginState: false,
	}
	UserList[Accid] = conn
	return conn, err
}

func Send(Accid string, it interface{}) (conn *Connection, err error) {
	if conn, err = GetUser(Accid); err != nil {
		panic(err.Error())
	}
	if err = conn.Send(it); err != nil {
		panic(err.Error())
	}
	return
}
