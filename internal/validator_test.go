// @Title  请填写文件名称（需要改）
// @Description  请填写文件描述（需要改）
// @Author  请填写自己的真是姓名（需要改）  2021/10/29 下午4:35
// @Update  请填写自己的真是姓名（需要改）  2021/10/29 下午4:35
package internal

import "testing"

func TestVerifyParam(t *testing.T) {
	type args struct {
		flags *Flags
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := VerifyParam(tt.args.flags); (err != nil) != tt.wantErr {
				t.Errorf("VerifyParam() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
