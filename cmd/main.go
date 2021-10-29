package main

import (
	im "ImprisonSlowSQL/imprison"
	"ImprisonSlowSQL/pkg/utils"
	v "ImprisonSlowSQL/pkg/version"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

const programName string = "ImprisonSlowSQL"

var versionInfo v.Info

var (
	host       string
	username   string
	password   string
	dbName     string
	port       uint
	longTime   uint
	versionStr string
)

func main() {
	versionInfo = v.Get()
	versionStr = fmt.Sprintf("version:%s ,buildDate:%s ,gitTag:%s ,gitCommit:%s", versionInfo.Version, versionInfo.BuildDate, versionInfo.GitTag, versionInfo.GitCommit)

	initLog()
	flags, err := initCmd()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	if utils.IsEmpty(host) || utils.IsEmpty(password) {
		return
	}

	if err = im.VerifyParam(flags); err != nil {
		log.Error(err)
		os.Exit(1)
	}

	log.Infof("The parameter with host=[%s] , username=[%s] , password=[%s] , port=[%d] , dbName=[%s] , longTime=[%d]", host, username, password, port, dbName, longTime)
	if err := im.VerifyMySQLCAP(); err != nil {
		log.Errorf("Verify MySQL Failure. %v", err)
		os.Exit(1)
	}

	slowSQL := &im.ImprisonSlowSQL{
		Ch: make(chan struct{}),
	}

	go func() {
		slowSQL.Imprison(flags)
	}()

	select {
	case <-slowSQL.Ch:
		os.Exit(0)
	}

}

func newCommand(programName string) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   programName,
		Short: fmt.Sprintf("The %s Cli", programName),
		Long:  fmt.Sprintf("The %s and bind slow SQL statements to the last CPU core.\n%s", programName, versionStr),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return errors.New(fmt.Sprintf("Use \"%s [command] --help\" for more information about a command.: %v", programName, args))
			}
			return nil
		},
		Version: versionInfo.Version,
	}

	cmd.Flags().StringVarP(&host, "host", "i", "", "Target host IP")
	cmd.Flags().StringVarP(&username, "username", "u", "root", "Target MySQL login username")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Target MySQL login password")
	cmd.Flags().UintVarP(&port, "port", "P", 3306, "Target MySQL login port")
	cmd.Flags().StringVarP(&dbName, "dbname", "n", "", "Target MySQL specific login db name")
	cmd.Flags().UintVarP(&longTime, "longTime", "l", 10, "Target MySQL specific slow sql time seconds config")

	cmd.MarkFlagRequired("host")
	cmd.MarkFlagRequired("password")

	return cmd
}

func initCmd() (*im.Flags, error) {

	cmd := newCommand(programName)
	if err := cmd.Execute(); err != nil {
		return nil, errors.New(fmt.Sprintf("Use \"%s --help\" for more information about a command.\n", programName))
	}
	flags := &im.Flags{
		Host:     host,
		Username: username,
		Password: password,
		DbName:   dbName,
		Port:     port,
		LongTime: longTime,
	}

	return flags, nil
}

func initLog() {
	log.SetReportCaller(true)
	log.SetFormatter(&utils.LogFormatter{})

	switch im.LogLevel {
	case im.LogLevelDebug:
		log.SetLevel(log.DebugLevel)
	case im.LogLevelInfo:
		log.SetLevel(log.InfoLevel)
	case im.LogLevelWarn:
		log.SetLevel(log.WarnLevel)
	case im.LogLevelError:
		log.SetLevel(log.ErrorLevel)
	default:
		//初始化日志错误,退出程序
		panic(fmt.Sprintf("cannot set log level:%s, there have four types can set: debug, info, warn, error", im.LogLevel))
	}
	log.SetOutput(os.Stdout)
}
