// @Title  请填写文件名称（需要改）
// @Description  请填写文件描述（需要改）
// @Author  请填写自己的真是姓名（需要改）  2021/10/25 下午3:47
// @Update  请填写自己的真是姓名（需要改）  2021/10/25 下午3:47
package utils

import "testing"

func TestExecCommand(t *testing.T) {
	if msg, err := ExecCommand("/bin/ls -la"); err != nil {
		t.Errorf("ExecCommand() error = [%v]", err)
	} else {
		t.Fatalf("ExecCommand()  result = [%s]", msg)
	}

	if msg, err := ExecCommand("/usr/sbin/getcap $(which mysqld)"); err != nil {
		t.Errorf("ExecCommand() error = [%v]", err)
	} else {
		t.Fatalf("ExecCommand()  result = [%s]", msg)
	}
}
