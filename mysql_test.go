package air_mysqlclient

import (
	"github.com/airingone/config"
	"github.com/airingone/log"
	"testing"
	"time"
)

/*
mysql -h 127.0.0.1  -P 3306 -uroot -ptest.2020  --default-character-set='utf8' -Dair_test
//用户数据表
CREATE TABLE `t_user` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '',
  `c_userid` varchar(64) NOT NULL DEFAULT '' COMMENT '',
  `c_user_name` varchar(64) NOT NULL DEFAULT '' COMMENT '',
  `c_user_tel` varchar(64) NOT NULL DEFAULT '' COMMENT '',
  `c_state` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '',
  `time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
ALTER TABLE `t_user` ADD unique(`c_userid`);
*/

type DbUser struct {
	Id       uint32    `ddb:"id"`
	UserId   string    `ddb:"c_userid"`
	UserName string    `ddb:"c_user_name"`
	UserTel  string    `ddb:"c_user_tel"`
	State    string    `ddb:"c_state"`
	Time     time.Time `ddb:"time"`
}

//mysql client 测试
func TestMysqlClient(t *testing.T) {
	config.InitConfig()                     //配置文件初始化
	log.InitLog(config.GetLogConfig("log")) //日志初始化

	InitMysqlClient("mysql_test1") //初始化mysql client,进程初始化的时候调用，有多个mysql则加对过参数
	defer CloseMysqlClient()       //关闭mysql client,在进程退出的时候调用就行

	cli, err := GetMysqlClient("mysql_test1")
	if err != nil {
		log.Error("get mysql client err, err: %+v", err)
		return
	}

	//insert
	values1 := make(map[string]interface{})
	values1["c_userid"] = "123456"
	values1["c_user_name"] = "test1"
	values1["c_user_tel"] = "2123456"
	values1["c_state"] = 1
	err = cli.Insert("t_user", values1)
	if err != nil {
		log.Error("mysql Insert err, err: %+v", err)
	}

	values2 := make(map[string]interface{})
	values2["c_userid"] = "12345678"
	values2["c_user_name"] = "test2"
	values2["c_user_tel"] = "212345678"
	values2["c_state"] = 1
	err = cli.Insert("t_user", values2)
	if err != nil {
		log.Error("mysql Insert err, err: %+v", err)
	}

	//count
	where1 := make(map[string]interface{})
	where1["c_state"] = 1
	count, err := cli.QueryCount("t_user", where1)
	if err != nil {
		log.Error("mysql QueryCount err, err: %+v", err)
	} else {
		log.Error("mysql QueryCount succ, count: %d", count)
	}

	where2 := make(map[string]interface{})
	where2["c_state"] = 0
	count2, err := cli.QueryCount("t_user", where2)
	if err != nil {
		log.Error("mysql QueryCount err, err: %+v", err)
	} else {
		log.Error("mysql QueryCount succ, count: %d", count2)
	}

	//get
	where3 := make(map[string]interface{})
	fields := []string{"id", "c_userid", "c_user_name", "c_state"}
	var users []DbUser
	err = cli.Query("t_user", where3, fields, 0, 2, &users)
	if err != nil {
		log.Error("mysql Query err, err: %+v", err)
	} else {
		log.Error("mysql Query succ, user: %+v", users)
	}

	//update
	where4 := make(map[string]interface{})
	where4["c_userid"] = "12345678"
	values3 := make(map[string]interface{})
	values3["c_state"] = 0
	values3["c_user_tel"] = "137"
	err = cli.Update("t_user", where4, values3)
	if err != nil {
		log.Error("mysql Update err, err: %+v", err)
	} else {
		log.Error("mysql Update succ")
	}

	//get
	where5 := make(map[string]interface{})
	where5["c_state"] = 0
	fields2 := []string{"c_userid", "c_user_name", "c_state"}
	var users2 []DbUser
	err = cli.Query("t_user", where5, fields2, 0, 2, &users2)
	if err != nil {
		log.Error("mysql Query err, err: %+v", err)
	} else {
		log.Error("mysql Query succ, user: %+v", users2)
	}

	//del
	where6 := make(map[string]interface{})
	where6["c_userid"] = "12345678"
	err = cli.Delete("t_user", where6)
	if err != nil {
		log.Error("mysql Delete err, err: %+v", err)
	} else {
		log.Error("mysql Delete succ")
	}

}
