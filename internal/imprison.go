package internal

import (
	"ImprisonSlowSQL/pkg/utils"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"strings"
)

//DB Db数据库连接池
var DB *sql.DB

//ImprisonSlowSQL 囚困Slow日志结构体
type ImprisonSlowSQL struct {
	Ch    chan struct{}
	Flags *Flags
}

// ThreadID SQL线程ID
type ThreadID string

// SlowSQLThreadInfo Thread表的字段结构体
type SlowSQLThreadInfo struct {
	processlistUser string
	processlistHost string
	processlistDB   string
	processlistInfo string
	resourceGroup   string
	processlistTime int
}

func (im *ImprisonSlowSQL) Imprison() {
	defer func() {
		im.Ch <- struct{}{}
	}()

	vcpu, err := utils.ExecCommand("cat /proc/cpuinfo  | grep processor | tail -n 1 | awk -F': ' '{print \\$NF}'")
	if err != nil {
		vcpu = "1"
	}
	log.Infof("The vm had cpus %s vcpu", vcpu)
	isConnect := im.ConnectMySQL()
	if !isConnect {

	}
	defer DB.Close()
	rows, sqlErr := DB.Query("SELECT * FROM information_schema.resource_groups WHERE RESOURCE_GROUP_NAME = 'rg_slowsql'")
	if sqlErr != nil {
		log.Errorf("检测rg_slowsql资源组错误, %v", sqlErr)
		return
	}
	if rows != nil && rows.Next() {
		log.Info("检测已存在MySQL的资源组 rg_slowsql")
		rows.Close()
	} else {
		log.Info("检测不存在资源组 rg_slowsql，现在创建`create resource group rg_slowsql`")
		DB.Query(fmt.Sprintf("CREATE RESOURCE GROUP rg_slowsql type=user vcpu=%s thread_priority=19 enable;", vcpu))
	}

	rows, sqlErr = DB.Query(fmt.Sprintf("SELECT THREAD_ID,PROCESSLIST_INFO,RESOURCE_GROUP,PROCESSLIST_TIME FROM performance_schema.threads WHERE PROCESSLIST_INFO REGEXP 'SELECT|INSERT|UPDATE|DELETE|ALTER' AND PROCESSLIST_TIME > %d", im.Flags.LongTime))
	if sqlErr != nil {
		log.Errorf("查询慢日志失败")
		return
	}
	if rows != nil && !rows.Next() {
		rows.Close()
		log.Info("未检测出当前执行中的卡顿慢SQL。")
		//TODO 输入到日志
		return
	}

	//这里将slow sql的语句设置到资源组那里
	for rows.Next() {
		var threadID ThreadID
		e := rows.Scan(&threadID)
		if e != nil {
			log.Errorf("查询有误")
			return
		}
		_, er := DB.Exec(fmt.Sprintf("SET resource group rg_slowsql for %s", threadID))
		if er == nil {
			log.Infof("ThreadId:%s 已设置到rg_slowsql的资源组", threadID)
		}
	}

	//显示目前有多少个slow sql并输入到日志
	rows, sqlErr = DB.Query(fmt.Sprintf("SELECT PROCESSLIST_USER,PROCESSLIST_HOST,PROCESSLIST_DB,PROCESSLIST_INFO,RESOURCE_GROUP,PROCESSLIST_TIME FROM performance_schema.threads WHERE PROCESSLIST_INFO REGEXP 'SELECT|INSERT|UPDATE|DELETE|REPLACE|ALTER' AND PROCESSLIST_TIME > %d", im.Flags.LongTime))
	if sqlErr != nil {
		log.Errorf("查询有误")
		return
	}
	for rows.Next() {
		log.Warnf("警告！出现卡顿慢SQL，请及时排查问题。")
		var info SlowSQLThreadInfo
		e := rows.Scan(&info.processlistUser, &info.processlistHost, &info.processlistDB, &info.processlistInfo, &info.resourceGroup, &info.processlistTime)
		if e != nil {
			log.Errorf("获取row有误")
			return
		}
		slowInfo := fmt.Sprintf("[用户名: %s, 来源IP: %s, 数据库名: %s, SQL语句: %s, 资源组：%s, 执行时间: %d]", info.processlistUser, info.processlistHost, info.processlistDB, info.processlistInfo, info.resourceGroup, info.processlistTime)
		log.Info(slowInfo)
	}

}

func (im *ImprisonSlowSQL) ConnectMySQL() bool {
	flags := im.Flags
	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"
	path := strings.Join([]string{flags.Username, ":", flags.Password, "@tcp(", flags.Host, ":", fmt.Sprintf("%d", flags.Port), ")/", flags.DbName, "?charset=utf8"}, "")

	//打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	DB, _ = sql.Open("mysql", path)
	//设置数据库最大连接数
	DB.SetConnMaxLifetime(10)
	//设置上数据库最大闲置连接数
	DB.SetMaxIdleConns(10)
	//验证连接
	if err := DB.Ping(); err != nil {
		log.Errorf("Open database Failure")
		return false
	}
	log.Info("Successfully connected to database")
	return true
}

func dropResourceGroup() error {
	_, err := DB.Exec(fmt.Sprintf("ALTER RESOURCE GROUP rg_slowsql DISABLE FORCE"))
	if err != nil {
		return errors.New(fmt.Sprintf("关闭rg_slowsql资源组有误, %v", err))
	}
	_, err = DB.Exec(fmt.Sprintf("DROP RESOURCE GROUP rg_slowsql"))
	if err != nil {
		return errors.New(fmt.Sprintf("删除rg_slowsql资源组有误, %v", err))
	}
	return nil
}
