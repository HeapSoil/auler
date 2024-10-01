package auler

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// 创建*cobra.Command，之后可以使用Command的Execute方法在cmd启动该程序
// 业务的具体实现
func NewAulerCommand() *cobra.Command {
	cmd := &cobra.Command{
		// 命令的名字，长短描述
		Use:   "auler",
		Short: "A demo scheduler project",
		Long:  "A demo scheduler project, aiming to coordinate distributed tasks",
		// 命令出错时处理， 保持命令出错时一眼就能看到错误信息
		SilenceUsage: true,
		// Run函数
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
		// 命令运行时不，设置需要指定命令行参数
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
	}

	// 初始化时，使得 initConfig 函数在每个命令运行时都会被调用以读取配置
	cobra.OnInitialize(initConfig)
	// 打印Use config file, ....

	// 持久化标志
	// Cobra 支持持久性标志(PersistentFlag)，该标志可用于它所分配的命令以及该命令下的每个子命令
	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "The path to the auler configuration file. Empty string for no configuration file.")

	// Cobra 也支持本地标志，本地标志只能在其所绑定的命令上使用
	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	return cmd
}

// 实际的业务代码入口函数
func run() error {

	settings, _ := json.Marshal(viper.AllSettings())
	fmt.Println(string(settings))
	fmt.Println(viper.GetString("db.username"))

	fmt.Println("Hello Auler!")
	return nil
}
