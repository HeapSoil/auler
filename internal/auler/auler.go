package auler

import (
	"fmt"

	"github.com/spf13/cobra"
)

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

	return cmd
}

// 实际的业务代码入口函数
func run() error {
	fmt.Println("Hello Auler!")
	return nil
}
