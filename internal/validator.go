package internal

import (
	"ImprisonSlowSQL/pkg/utils"
	"errors"
	log "github.com/sirupsen/logrus"
	"unicode"
)

func VerifyParam(flags *Flags) error {
	var errMsg string
	if !utils.CheckIPv4(flags.Host) {
		errMsg = "The IP format of host is incorrect"
	} else if !unicode.IsDigit(rune(flags.Port)) {
		errMsg = "The Port is not number"
	} else if !(flags.Port < 2<<15 && flags.Port > 0) {
		errMsg = "The Port is less than 65536 or greater than 0 "
	} else if flags.LongTime < 0 {
		errMsg = "The Longtime is greater than 0"
	}
	if errMsg != "" {
		return errors.New(errMsg)
	}

	return nil
}

func VerifyMySQLCAP() error {
	result, err := utils.ExecCommand("/usr/sbin/getcap $(which mysqld)")
	errMsg := "检测到mysql没有开启CAP_SYS_NICE功能。请使用下面的命令进行开启设置：\"setcap cap_sys_nice+ep /usr/sbin/mysqld\". 并需要重启mysqld服务生效."
	if err != nil {
		return errors.New(errMsg)
	}
	if !utils.IsEmpty(result) {
		log.Info("mysql已经开启CAP_SYS_NICE功能")
	} else {
		return errors.New(errMsg)
	}

	return nil
}
