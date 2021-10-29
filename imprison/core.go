package imprison

import (
	"ImprisonSlowSQL/pkg/utils"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"strings"
)

//Db数据库连接池
var DB *sql.DB

type ImprisonSlowSQL struct {
	Ch chan struct{}
}

func (im *ImprisonSlowSQL) Imprison(flags *Flags) {
	cores, err := utils.ExecCommand("cat /proc/cpuinfo  | grep processor | tail -n 1 | awk -F': ' '{print \\$NF}'")
	if err != nil {
		cores = "1"
	}
	log.Infof("The Vm Cpus %s", cores)
	connectMySQL(flags)

	Ch <- struct{}{}
}

func connectMySQL(flags *Flags) bool {
	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"
	path := strings.Join([]string{flags.Username, ":", flags.Password, "@tcp(", flags.Host, ":", string(flags.Port), ")/", flags.DbName, "?charset=utf8"}, "")

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
