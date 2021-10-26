// @Title  请填写文件名称（需要改）
// @Description  请填写文件描述（需要改）
// @Author  请填写自己的真是姓名（需要改）  2021/10/25 上午9:54
// @Update  请填写自己的真是姓名（需要改）  2021/10/25 上午9:54
package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

const programName string = "ImprisonSlowSQL"
const version string = "1.0.0"

var (
	host     string
	username string
	password string
	dbName   string
	port     uint16
	longTime uint16
)

func main() {
	cmd := newCommand(programName)
	if err := cmd.Execute(); err != nil {
		fmt.Printf("Use \"%s --help\" for more information about a command.\n", programName)
		os.Exit(1)
	}

	fmt.Printf("The parameter with host=[%s] , username=[%s] , password=[%s] , port=[%d] , dbName=[%s] , longTime=[%d]\n\n", host, username, password, port, dbName, longTime)

}

func newCommand(programName string) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   programName,
		Short: fmt.Sprintf("The %s Cli", programName),
		Long:  fmt.Sprintf("The %s and bind slow SQL statements to the last CPU core.\nThe Version: %s ", programName, version),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("Use \"%s [command] --help\" for more information about a command.: %v", programName, args)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&host, "host", "i", "", "Target host IP")
	cmd.Flags().StringVarP(&username, "username", "u", "root", "Target MySQL login username")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Target MySQL login password")
	cmd.Flags().Uint16VarP(&port, "port", "P", 3306, "Target MySQL login port")
	cmd.Flags().StringVarP(&dbName, "dbname", "n", "", "Target MySQL specific login db name")
	cmd.Flags().Uint16VarP(&longTime, "longTime", "l", 10, "Target MySQL specific slow sql time seconds config")

	cmd.MarkFlagRequired("host")
	cmd.MarkFlagRequired("password")

	return cmd
}
