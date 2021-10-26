// @Title  请填写文件名称（需要改）
// @Description  请填写文件描述（需要改）
// @Author  请填写自己的真是姓名（需要改）  2021/10/25 下午3:47
// @Update  请填写自己的真是姓名（需要改）  2021/10/25 下午3:47
package pkg

import "testing"

func TestExecCommand(t *testing.T) {
	if err := ExecCommand("/bin/ls -la"); err != nil {
		t.Errorf("ExecCommand() error = [%v]", err)
	}

	if err := ExecCommand("/usr/sbin/getcap $(which mysqld)"); err != nil {
		t.Errorf("ExecCommand() error = [%v]", err)
	}
}


