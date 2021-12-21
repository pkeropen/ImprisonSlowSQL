// @Title  请填写文件名称（需要改）
// @Description  请填写文件描述（需要改）
// @Author  请填写自己的真是姓名（需要改）  2021/10/29 下午4:25
// @Update  请填写自己的真是姓名（需要改）  2021/10/29 下午4:25
package test

import (
	im "ImprisonSlowSQL/internal"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"sync"
	"testing"
)

//DB Db数据库连接池
var DB *sql.DB

func TestInsertData(t *testing.T) {
	flags := &im.Flags{
		Host:     "121.201.50.239",
		Username: "root",
		Port:     43306,
		Password: "58_v29rC",
		DbName:   "shop-fstv",
	}

	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"
	path := strings.Join([]string{flags.Username, ":", flags.Password, "@tcp(", flags.Host, ":", fmt.Sprintf("%d", flags.Port), ")/", flags.DbName, "?charset=utf8"}, "")

	//打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	DB, _ = sql.Open("mysql", path)
	//设置数据库最大连接数
	DB.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	DB.SetMaxIdleConns(100)
	//验证连接
	if err := DB.Ping(); err != nil {
		fmt.Println("Open database Failure")
	}
	fmt.Println("Successfully connected to database")

	wg := sync.WaitGroup{}

	for i := 1; i < 100; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			for i := 0; i < 10000; i++ {
				DB.Exec("insert into test(name,age) values(?,?)", "test", 1)
			}
		}()
	}

	wg.Wait()

	t.Fatal("test")
}
