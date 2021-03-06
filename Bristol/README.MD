## Bristol ---- 模拟MySQL从库连接,监听MySQL binlog变更并解析

#### DDL 支持说明

当前只支持字段在表结构末尾追加新字段，如果配置的二进制位点是在DDL 之前的位点，会出现数据和字段对应不上
数据类型修改之后，如果配置的二进制位点是在DDL 之前的位点，可能会出现数据出错

#### 类型转换

从Binlog解析出来的数据,在转成 map[string]interface{} 的数据格式的时候
MySQL里的存储类型对应Golang里的数据类型不一样
括号里的代表是Golang里的数据类型


- [x] - TINYINT ( int8 | unit8 )
- [x] - SMALLINT ( int16 | unit16 )
- [x] - MEDIUMINT ( int32 | uint32 )
- [x] - INT ( int32 | uint32 )
- [x] - BIGINT ( int64 | uint64 )
- [x] - FLOAT ( float32 )
- [x] - DOUBLE ( float64 )
- [x] - REAL ( float64 )
- [x] - DECIMAL , NUMERIC( 返回string类型 )
- [x] - DATE , TIME , YEAR , DATETIME , TIMESTAMP   ( string类型 )
- [x] - CHAR , VARCHAR ( string )
- [x] - TEXT , TINYTEXT , MEDIUMINTTEXT , - LONGTEXT ( string )
- [x] - BLOB, TINYBLOB , MEDIUMINTBLOB , LONGBLOB ( string )
- [x] - ENUM ( string类型 )
- [x] - SET ( []string )
- [x] - BIT ( int64类型 )

#### 使用案例
`````go
import (
	"log"

	"github.com/Bristol/mysql"
	"time"
)

func callback(data *mysql.EventReslut) {
	log.Println(data)
}

func main() {
	filename := "mysql-bin.000022"
	var position uint32 = 13333
	reslut := make(chan error, 1)
	m := make(map[string]uint8, 0)
	m["testdbcreate"] = 1
	BinlogDump := &mysql.BinlogDump{
		DataSource:    "root:root@tcp(127.0.0.1:3306)/test",
		CallbackFun:   callback,
		ReplicateDoDb: m,
		OnlyEvent:     []mysql.EventType{mysql.QUERY_EVENT, mysql.WRITE_ROWS_EVENTv1, mysql.UPDATE_ROWS_EVENTv1, mysql.DELETE_ROWS_EVENTv1},
	}
	go BinlogDump.StartDumpBinlog(filename, position, 100,reslut,"",0)
	go func() {
		v := <-reslut
		log.Printf("monitor reslut:%s \r\n", v)
	}()
	for {
		time.Sleep(10 * time.Second)
	}
}

`````

#### 数据检验

请参考 [一键数据检验测试工具](https://github.com/brokercap/Bifrost/tree/master/Bristol/test)