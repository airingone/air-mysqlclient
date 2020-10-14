# mysql client组件
## 1.组件描述
mysql client组件封装了mysql客户端，同时引进sql build，可以安全与便捷的生成sql语句。
## 2.如何使用
```
import (
    "github.com/airingone/config"
    "github.com/airingone/log"
    mysqlclient "github.com/airingone/air-mysqlclient"
)

func main() {
    config.InitConfig()                     //进程启动时调用一次初始化配置文件，配置文件名为config.yml，目录路径为../conf/或./
    log.InitLog(config.GetLogConfig("log")) //进程启动时调用一次初始化日志
    mysqlclient.InitMysqlClient("mysql_test1") //初始化mysql client,进程初始化的时候调用，有多个mysql则加对过参数
    defer mysqlclient.CloseMysqlClient()       //关闭mysql client,在进程退出的时候调用就行

    cli, err := mysqlclient.GetMysqlClient("mysql_test1")
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

    //or
    //get
    where7 := make(map[string]interface{})
    fields7 := []string{"id", "c_userid", "c_user_name", "c_state"}
    var users7 []DbUser
    err = mysqlclient.MysqlQuery("mysql_test1", "t_user", where7, fields7, 0, 2, &users7)
    if err != nil {
        log.Error("mysql Query err, err: %+v", err)
    } else {
        log.Error("mysql Query succ, user: %+v", users7)
    }

    //or
    mysqlConfig := config.GetMysqlConfig("mysql_test1")
    cli2, err := mysqlclient.NewMysqlClient(mysqlConfig.Addr, mysqlConfig.MaxIdleConns, mysqlConfig.MaxOpenConns)
    if err != nil {
        log.Error("new mysql client err, err: %+v", err)
        return
    }
    where8 := make(map[string]interface{})
    fields8 := []string{"id", "c_userid", "c_user_name", "c_state"}
    var users8 []DbUser
    err = cli2.Query("t_user", where8, fields8, 0, 2, &users8)
    if err != nil {
        log.Error("mysql Query err, err: %+v", err)
    } else {
        log.Error("mysql Query succ, user: %+v", users8)
    }
    cli2.Close()    

    //Update,Query,QueryCount,Delete等操作详见实现代码或测试列子。
}
```
更多使用请参考[测试例子](https://github.com/airingone/air-mysqlclient/blob/master/mysql_test.go)
