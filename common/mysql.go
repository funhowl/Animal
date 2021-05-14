package common

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var MysqlDb *sql.DB
var MysqlDbErr error

// 初始化链接
func init() {

	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", USER_NAME, PASS_WORD, HOST, PORT, DATABASE, CHARSET)

	// 打开连接失败
	MysqlDb, MysqlDbErr = sql.Open("mysql", dbDSN)
	//defer MysqlDb.Close();
	if MysqlDbErr != nil {
		log.Println("dbDSN: " + dbDSN)
		panic("数据源配置不正确: " + MysqlDbErr.Error())
	}

	// 最大连接数
	MysqlDb.SetMaxOpenConns(100)
	// 闲置连接数
	MysqlDb.SetMaxIdleConns(20)
	// 最大连接周期
	MysqlDb.SetConnMaxLifetime(100 * time.Second)

	if MysqlDbErr = MysqlDb.Ping(); nil != MysqlDbErr {
		panic("数据库链接失败: " + MysqlDbErr.Error())
	}

}

func SelUserRow(Accid string) (*UserRow, error) {

	var (
		urow     UserRow
		rowObj   *sql.Rows
		Userinfo sql.RawBytes
		Tasks    sql.RawBytes
		err      error
	)
	sqlStr := "select Accid, Gold, Userinfo, Tasks from users where Accid = ?" //占位符
	//  rowObj, err := MysqlDb.Query(sqlStr, Accid)
	if rowObj, err = MysqlDb.Query(sqlStr, Accid); err != nil {
		panic(err.Error())
	}

	for rowObj.Next() {
		rowObj.Scan(&urow.Accid, &urow.Gold, &Userinfo, &Tasks)
	}
	json.Unmarshal(Tasks, &urow.Tasks)
	json.Unmarshal(Userinfo, &urow.Userinfo)

	return &urow, err
}

func InitUserRow(Accid string) (*UserRow, error) {

	var (
		taskrow *TaskRow
		inforow *Userinfo
		row     *UserRow
		err     error
		// rowObj  sql.Result
	)

	if taskrow, err = initTaskRow(); err != nil {
		return row, err
	}

	if inforow, err = initInfoRow(); err != nil {
		return row, err
	}

	row = &UserRow{
		Accid:    Accid,
		Gold:     500,
		Tasks:    taskrow,
		Userinfo: inforow,
	}

	//1. 写sql语句
	sqlStr := "insert into users (Accid, Gold, Userinfo, Tasks) values(?, ?, ?,?);" //占位符

	//2. JSON格式处理
	userinfo, _ := json.Marshal(row.Userinfo)
	tasks, _ := json.Marshal(row.Tasks)

	//3.调用函数
	if _, err = MysqlDb.Exec(sqlStr, row.Accid, row.Gold, userinfo, tasks); err != nil {
		return row, err
	}

	return row, err
}

func initTaskRow() (ans *TaskRow, err error) {
	ans = &TaskRow{
		Record: "",
		Point:  1000,
	}
	return
}

func initInfoRow() (ans *Userinfo, err error) {
	ans = &Userinfo{
		Icon: "2",
		Name: "大肚皮",
	}
	return
}

func UpdateUsersDB(urow *UserRow, data UserRow) (*UserRow, error) { // 更新urow至数据库

	var (
		err     error
		sqlbody string
		sqlarr  []string
		uplist  *UserRow
	)
	uplist = &data
	sqlarr = make([]string, 0)

	if urow.Gold != uplist.Gold {
		sqlarr = append(sqlarr, " Gold = "+ChangeToString(uplist.Gold))
	}

	if urow.Tasks != uplist.Tasks {
		_type := reflect.TypeOf(*uplist.Tasks)
		uplistValue := reflect.ValueOf(*uplist.Tasks)
		urowValue := reflect.ValueOf(*urow.Tasks)
		sqlarr = getFieldsSql("Tasks", _type, uplistValue, urowValue, sqlarr)
	}
	if urow.Userinfo != uplist.Userinfo {
		_type := reflect.TypeOf(*uplist.Userinfo)
		uplistValue := reflect.ValueOf(*uplist.Userinfo)
		urowValue := reflect.ValueOf(*urow.Userinfo)
		sqlarr = getFieldsSql("Userinfo", _type, uplistValue, urowValue, sqlarr)
	}
	sqlbody = strings.Join(sqlarr, ",")

	sqlStr := "update users set " + sqlbody + " where Accid = ?"
	fmt.Println("sqlStr:", sqlStr)

	// MysqlDb.Exec(sqlStr, sqlbody, urow.Accid)
	if _, err = MysqlDb.Exec(sqlStr, urow.Accid); err != nil {
		panic(err.Error())
	}
	*urow = data

	return urow, err
}

func getFieldsSql(FieldsName string, _type reflect.Type, uplistValue reflect.Value, urowValue reflect.Value, sqlarr []string) []string {

	sqlpaths := make([]string, 0)
	del_sqlpaths := make([]string, 0)
	for k := 0; k < _type.NumField(); k++ {
		newField := uplistValue.Field(k).Interface()
		oldField := urowValue.Field(k).Interface()
		if newField != oldField {
			if reflect.ValueOf(newField).IsZero() {
				del_sqlpaths = append(del_sqlpaths, "'$.\""+_type.Field(k).Name+"\"'")
			} else if reflect.TypeOf(newField).Kind() == reflect.Struct {
				sqlpaths = append(sqlpaths, "'$.\""+_type.Field(k).Name+"\"'", "CAST('"+ChangeToString(newField)+"' as JSON)")
			} else {
				sqlpaths = append(sqlpaths, "'$.\""+_type.Field(k).Name+"\"'", "\""+ChangeToString(newField)+"\"")
			}
		}
	}

	if len(del_sqlpaths) > 0 {
		sqlarr = append(sqlarr, FieldsName+"=JSON_REMOVE("+FieldsName+","+strings.Join(del_sqlpaths, ",")+")")
	}
	if len(sqlpaths) > 0 {
		sqlarr = append(sqlarr, FieldsName+"=JSON_SET("+FieldsName+","+strings.Join(sqlpaths, ",")+")")
	}
	return sqlarr
}
