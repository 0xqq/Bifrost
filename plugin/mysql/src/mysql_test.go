package src_test

import (
	"testing"
	"log"
	pluginDriver "github.com/brokercap/Bifrost/plugin/driver"
	"github.com/brokercap/Bifrost/sdk/pluginTestData"
	MyPlugin "github.com/brokercap/Bifrost/plugin/mysql/src"
	dbDriver "database/sql/driver"
	"github.com/brokercap/Bifrost/Bristol/mysql"
	"strings"
	"github.com/brokercap/Bifrost/util/dataType"
	"github.com/brokercap/Bifrost/server/history"
	"reflect"
	"fmt"
	"math/rand"
	"time"
	"strconv"
)

var url string = "root:root@tcp(10.40.2.41:3306)/bifrost_test"

var SchemaName string = "bifrost_test"
var TableName string = "binlog_field_test"
var mysqlConn mysql.MysqlConnection
/*
ddl

CREATE TABLE binlog_field_test(id UInt32,testtinyint Int8,testsmallint Int16,testmediumint Int32,testint Int32,testbigint Int64,testvarchar String,testchar String,testenum String,testset String,testtime String,testdate Date,testyear Int16,testtimestamp DateTime,testdatetime DateTime,testfloat Float64,testdouble Float64,testdecimal Float64,testtext String,testblob String,testbit Int64,testbool Int8,testmediumblob String,testlongblob String,testtinyblob String,test_unsinged_tinyint UInt8,test_unsinged_smallint UInt16,test_unsinged_mediumint UInt32,test_unsinged_int UInt32,test_unsinged_bigint UInt64) ENGINE = MergeTree() ORDER BY (id);
 */

func TestChechUri(t *testing.T){
	myConn := MyPlugin.MyConn{}
	if err := myConn.CheckUri(url);err!= nil{
		log.Println("TestChechUri err:",err)
	}else{
		log.Println("TestChechUri success")
	}
}

func TestGetSchemaList(t *testing.T)  {
	c := MyPlugin.NewMysqlDBConn(url)
	defer c.Close()
	list := c.GetSchemaList()
	if len(list) > 0{
		t.Log(list)
		t.Log("TestGetSchemaList success")
	}else{
		t.Error("TestGetSchemaList failed")
	}
}


func TestGetSchemaTableList(t *testing.T)  {
	c := MyPlugin.NewMysqlDBConn(url)
	defer c.Close()
	list := c.GetSchemaTableList("mysql")
	if len(list) > 0{
		t.Log(list)
		t.Log("TestGetSchemaTableList success")
	}else{
		t.Error("TestGetSchemaTableList failed")
	}
}

func TestGetTableFields(t *testing.T)  {
	c := MyPlugin.NewMysqlDBConn(url)
	c.Close()
	list := c.GetTableFields(SchemaName,TableName)
	if len(list) > 0{
		t.Log(list)
		t.Log("TestGetTableFields success")
	}else{
		t.Error("TestGetTableFields failed")
	}
}

func beforeTest()  {
	
}

func initDBTable(delTable bool) {
	c := mysql.NewConnect(url)
	sql1:= "CREATE DATABASE IF NOT EXISTS  `"+SchemaName+"`"
	_,err := c.Exec(sql1,[]dbDriver.Value{})
	if err != nil{
		log.Fatal(err)
	}
	sql2:="CREATE TABLE  IF NOT EXISTS `"+SchemaName+"`.`"+TableName+"`( `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT, `testtinyint` TINYINT(4) NOT NULL DEFAULT '-1', `testsmallint` SMALLINT(6) NOT NULL DEFAULT '-2', `testmediumint` MEDIUMINT(8) NOT NULL DEFAULT '-3', `testint` INT(11) NOT NULL DEFAULT '-4', `testbigint` BIGINT(20) NOT NULL DEFAULT '-5', `testvarchar` VARCHAR(400) NOT NULL, `testchar` CHAR(2) NOT NULL, `testenum` ENUM('en1', 'en2', 'en3') NOT NULL DEFAULT 'en1', `testset` SET('set1', 'set2', 'set3') NOT NULL DEFAULT 'set1', `testtime` TIME NOT NULL DEFAULT '00:00:00', `testdate` DATE NOT NULL DEFAULT '0000-00-00', `testyear` YEAR(4) NOT NULL DEFAULT '1989', `testtimestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, `testdatetime` DATETIME NOT NULL DEFAULT '0000-00-00 00:00:00', `testfloat` FLOAT(9, 2) NOT NULL DEFAULT '0.00', `testdouble` DOUBLE(9, 2) NOT NULL DEFAULT '0.00', `testdecimal` DECIMAL(9, 2) NOT NULL DEFAULT '0.00', `testtext` TEXT NOT NULL, `testblob` BLOB NOT NULL, `testbit` BIT(64) NOT NULL DEFAULT b'0', `testbool` TINYINT(1) NOT NULL DEFAULT '0', `testmediumblob` MEDIUMBLOB NOT NULL, `testlongblob` LONGBLOB NOT NULL, `testtinyblob` TINYBLOB NOT NULL, `test_unsinged_tinyint` TINYINT(4) UNSIGNED NOT NULL DEFAULT '1', `test_unsinged_smallint` SMALLINT(6) UNSIGNED NOT NULL DEFAULT '2', `test_unsinged_mediumint` MEDIUMINT(8) UNSIGNED NOT NULL DEFAULT '3', `test_unsinged_int` INT(11) UNSIGNED NOT NULL DEFAULT '4', `test_unsinged_bigint` BIGINT(20) UNSIGNED NOT NULL DEFAULT '5', PRIMARY KEY (`id`) ) ENGINE = MYISAM AUTO_INCREMENT = 0 CHARSET = utf8"
	if delTable == false{
		_,err = c.Exec(sql2,[]dbDriver.Value{})
		if err != nil{
			log.Fatal(err)
		}
	}else{
		sql3 := "DROP TABLE IF EXISTS `"+SchemaName+"`.`"+TableName+"`"
		_,err = c.Exec(sql3,[]dbDriver.Value{})
		if err != nil{
			log.Fatal(err)
		}
		_,err = c.Exec(sql2,[]dbDriver.Value{})
		if err != nil{
			log.Fatal(err)
		}
	}
	c.Close()
}


func getParam()  map[string]interface{}{
	type fieldStruct struct {
		ToField 		string
		FromMysqlField 	string
	}

	param := make(map[string]interface{},0)
	Field := make([]fieldStruct,0)
	Field = append(Field,fieldStruct{"id","id"})
	Field = append(Field,fieldStruct{"test_unsinged_bigint","test_unsinged_bigint"})
	Field = append(Field,fieldStruct{"test_unsinged_int","test_unsinged_int"})
	Field = append(Field,fieldStruct{"test_unsinged_mediumint","test_unsinged_mediumint"})
	Field = append(Field,fieldStruct{"test_unsinged_smallint","test_unsinged_smallint"})
	Field = append(Field,fieldStruct{"test_unsinged_tinyint","test_unsinged_tinyint"})
	Field = append(Field,fieldStruct{"testtinyint","testtinyint"})
	Field = append(Field,fieldStruct{"testsmallint","testsmallint"})
	Field = append(Field,fieldStruct{"testmediumint","testmediumint"})
	Field = append(Field,fieldStruct{"testint","testint"})
	Field = append(Field,fieldStruct{"testbigint","testbigint"})
	Field = append(Field,fieldStruct{"testbit","testbit"})
	Field = append(Field,fieldStruct{"testbool","testbool"})
	Field = append(Field,fieldStruct{"testvarchar","testvarchar"})
	Field = append(Field,fieldStruct{"testchar","testchar"})
	Field = append(Field,fieldStruct{"testtime","testtime"})
	Field = append(Field,fieldStruct{"testdate","testdate"})
	Field = append(Field,fieldStruct{"testyear","testyear"})

	Field = append(Field,fieldStruct{"testtimestamp","testtimestamp"})

	Field = append(Field,fieldStruct{"testdatetime","testdatetime"})
	Field = append(Field,fieldStruct{"testfloat","testfloat"})
	Field = append(Field,fieldStruct{"testdouble","testdouble"})
	Field = append(Field,fieldStruct{"testdecimal","testdecimal"})
	Field = append(Field,fieldStruct{"testtext","testtext"})
	Field = append(Field,fieldStruct{"testblob","testblob"})
	Field = append(Field,fieldStruct{"testmediumblob","testmediumblob"})
	Field = append(Field,fieldStruct{"testlongblob","testlongblob"})
	Field = append(Field,fieldStruct{"testtinyblob","testtinyblob"})
	Field = append(Field,fieldStruct{"testenum","testenum"})
	Field = append(Field,fieldStruct{"testset","testset"})

	sql := ""
	for _,f := range Field{
		sql += f.ToField+","
	}
	log.Println(sql)
	param["Field"] = Field

	PriKey := make([]fieldStruct,1)
	PriKey[0] = fieldStruct{"id","id"}
	param["PriKey"] = PriKey
	param["Schema"] = SchemaName
	param["Table"] = TableName

	return param
}

func getPluginConn() pluginDriver.ConnFun {
	myConn := MyPlugin.MyConn{}
	conn := myConn.Open(url)

	p,err := conn.SetParam(getParam())
	if err != nil{
		log.Println("set param fatal err")
		log.Fatal(err)
	}

	log.Println("p:",p)
	return conn
}


func TestCommit(t *testing.T){

	beforeTest()
	conn := getPluginConn()
	initDBTable(false)

	e := pluginTestData.NewEvent()

	conn.Insert(e.GetTestInsertData())
	conn.Del(e.GetTestDeleteData())
	conn.Update(e.GetTestUpdateData())
	conn.Insert(e.GetTestInsertData())
	conn.Insert(e.GetTestInsertData())
	conn.Insert(e.GetTestInsertData())

	_,err2 := conn.Commit()
	if err2 != nil{
		log.Fatal(err2)
	}
}


func TestInsertAndChekcData(t *testing.T){
	beforeTest()
	initDBTable(false)
	conn := getPluginConn()
	e := pluginTestData.NewEvent()
	insertdata := e.GetTestInsertData()
	conn.Insert(insertdata)
	_,err2 := conn.Commit()
	if err2 != nil{
		t.Fatal(err2)
	}

	checkResult,err := checkDataRight(insertdata.Rows[len(insertdata.Rows)-1])
	if err != nil{
		t.Fatal(err)
	}

	for _,v := range checkResult["ok"]{
		t.Log(v)
	}

	for _,v := range checkResult["error"]{
		t.Error(v)
	}
}


func TestInsertNullAndChekcData(t *testing.T){
	beforeTest()
	initDBTable(false)
	conn := getPluginConn()
	e := pluginTestData.NewEvent()
	e.SetIsNull(true)
	insertdata := e.GetTestInsertData()
	conn.Insert(insertdata)
	_,err2 := conn.Commit()
	if err2 != nil{
		t.Fatal(err2)
	}

	checkResult,err := checkDataRight(insertdata.Rows[len(insertdata.Rows)-1])
	if err != nil{
		t.Fatal(err)
	}

	for _,v := range checkResult["ok"]{
		t.Log(v)
	}

	for _,v := range checkResult["error"]{
		t.Error(v)
	}
}


func TestUpdateAndChekcData(t *testing.T){
	beforeTest()
	initDBTable(false)
	conn := getPluginConn()
	e := pluginTestData.NewEvent()
	insertdata := e.GetTestInsertData()
	conn.Insert(insertdata)

	updateData := e.GetTestUpdateData()
	conn.Update(updateData)
	_,err2 := conn.Commit()
	if err2 != nil{
		t.Fatal(err2)
	}

	checkResult,err := checkDataRight(updateData.Rows[len(updateData.Rows)-1])
	if err != nil{
		t.Fatal(err)
	}

	for _,v := range checkResult["ok"]{
		t.Log(v)
	}

	for _,v := range checkResult["error"]{
		t.Error(v)
	}
}


func TestDelAndChekcData(t *testing.T){
	beforeTest()
	initDBTable(false)
	conn := getPluginConn()
	e := pluginTestData.NewEvent()
	insertdata := e.GetTestInsertData()
	conn.Insert(insertdata)

	updateData := e.GetTestUpdateData()
	conn.Update(updateData)

	deleteData := e.GetTestDeleteData()
	conn.Del(deleteData)
	_,err2 := conn.Commit()
	if err2 != nil{
		t.Fatal(err2)
	}
	m,err:=getMysqlData(fmt.Sprint(deleteData.Rows[len(deleteData.Rows)-1]["id"]))
	if err != nil{
		t.Fatal(err)
	}

	if len(m) == 0{
		t.Log("test delete success")
	}else{
		t.Error("test delete error,delete failed")
	}
}

func getMysqlConn() mysql.MysqlConnection  {
	if mysqlConn == nil{
		mysqlConn = mysql.NewConnect(url)
	}
	return mysqlConn
}

func getMysqlData(id string)  (map[string]interface{},error){
	schema := SchemaName
	table := TableName
	conn := getMysqlConn()
	Fields := history.GetSchemaTableFieldList(conn,schema,table)
	sql := ""
	for index,Field := range Fields{
		if index == 0 {
			sql = Field.COLUMN_NAME
		}else{
			sql += ","+Field.COLUMN_NAME
		}
	}
	sql = "select "+sql +" from `"+schema+"`.`"+table +"` where id = "+id

	//sql := "select id,test_unsinged_bigint,test_unsinged_int,test_unsinged_mediumint,test_unsinged_tinyint,testtinyint,testsmallint,testmediumint,testint,testbigint,testbit,testbool,testvarchar,testchar,testtime,testdate,testyear,testtimestamp,testdatetime,testfloat,testdouble,testdecimal,testtext,testblob,testmediumblob,testlongblob,testtinyblob,testenum,testset from bifrost_test.binlog_field_test where id = 1"

	stmt,err := conn.Prepare(sql)
	if err != nil{
		return nil,err
	}
	defer stmt.Close()
	rows,err := stmt.Query([]dbDriver.Value{})
	if err != nil{
		return  nil,err
	}
	defer rows.Close()
	n := len(Fields)
	m := make(map[string]interface{}, n)
	for {
		dest := make([]dbDriver.Value, n, n)
		err := rows.Next(dest)
		if err != nil {
			break
		}
		for i, v := range Fields {
			if dest[i] == nil {
				m[v.COLUMN_NAME] = nil
				continue
			}
			switch v.DATA_TYPE {
			case "set":
				s := string(dest[i].([]byte))
				m[v.COLUMN_NAME] = strings.Split(s, ",")
				break
			default:
				m[v.COLUMN_NAME], _ = dataType.TransferDataType(dest[i].([]byte), v.ToDataType)
				break
			}
		}
		break
	}

	return m, nil
}

func getTableCount() (uint64,error){
	conn := getMysqlConn()
	sql := "select count(*) from `"+SchemaName+"`.`"+TableName +"`"
	stmt,err := conn.Prepare(sql)
	if err != nil{
		return 0,err
	}
	defer stmt.Close()
	rows,err := stmt.Query([]dbDriver.Value{})
	if err != nil{
		return  0,err
	}
	defer rows.Close()
	dest := make([]dbDriver.Value, 1, 1)
	rows.Next(dest)
	uint64, err := strconv.ParseUint(string(dest[0].([]byte)), 10, 64)
	return uint64,err

}

func checkDataRight(eventDataMap map[string]interface{}) (map[string][]string,error) {

	m,err :=getMysqlData(fmt.Sprint(eventDataMap["id"]))
	if err != nil{
		return nil,err
	}

	result := make(map[string][]string,0)
	result["ok"] = make([]string,0)
	result["error"] = make([]string,0)

	for key,val := range eventDataMap{
		if _,ok := m[key];!ok{
			s := fmt.Sprint("field:",key," not esxit")
			result["error"] = append(result["error"],s)
		}
		if reflect.TypeOf(val) == reflect.TypeOf(m[key]) && fmt.Sprint(val) == fmt.Sprint(m[key]){
			s := fmt.Sprint(key," == ",val," ( ",reflect.TypeOf(val)," ) ")
			result["ok"] = append(result["ok"],s)
		}else{
			s := fmt.Sprint(key," src: ",val," ( ",reflect.TypeOf(val)," ) "," != ",m[key]," ( ",reflect.TypeOf(m[key])," )")
			result["error"] = append(result["error"],s)
		}
	}

	return result,nil
}

func TestRandDataAndCheck(t *testing.T){

	var n int = 1000

	e := pluginTestData.NewEvent()

	beforeTest()
	initDBTable(true)

	conn := getPluginConn()

	for i:=0;i<n;i++{
		var eventData *pluginDriver.PluginDataType
		rand.Seed(time.Now().UnixNano()+int64(i))
		switch rand.Intn(3){
		case 0:
			eventData = e.GetTestInsertData()
			conn.Insert(eventData)
			break
		case 1:
			eventData = e.GetTestUpdateData()
			conn.Update(eventData)
			break
		case 2:
			eventData = e.GetTestDeleteData()
			conn.Del(eventData)
			break
		case 3:
			eventData = e.GetTestQueryData()
			conn.Query(eventData)
			break
		}
	}
	conn.Commit()

	count,err := getTableCount()
	if err != nil{
		t.Fatal(err)
	}

	if count != uint64(len(e.GetDataMap())){
		for k,v := range e.GetDataMap(){
			t.Log(k ," ",v)
		}
		t.Fatal("mysql Table Count:",count, " != srcDataCount:",len(e.GetDataMap()))
	}

	for _,data := range e.GetDataMap(){
		checkResult,err := checkDataRight(data)
		if err != nil{
			t.Error("data:",data,"err:",err)
			continue
		}
		if len(checkResult["error"]) > 0{
			t.Error("id:",data["id"]," failed")
			for _,v := range checkResult["error"]{
				t.Error(v)
			}
		}else{
			t.Log("id:",data["id"],data)
			t.Log("id:",data["id"]," success")
		}
	}

	t.Log("mysql Table Count:",count," srcDataCount:",len(e.GetDataMap()))

	t.Log("test over")
}