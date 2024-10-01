package auler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// 定义放置 auler 服务配置的默认目录
	recommendedHomeDir = ".auler"
	// 指定了 auler 服务的默认配置文件名
	defaultConfigName = "auler.yaml"
)

// initConfig 设置需要读取的配置文件名、环境变量，并读取配置文件内容到 viper 中.
func initConfig() {
	// 全局cfgFile
	if cfgFile != "" {
		// 用户指定了配置文件路径，则从用户给定的路径读取
		viper.SetConfigFile(cfgFile)
	} else {
		// 获取用户主目录，添加搜索路径
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(filepath.Join(home, recommendedHomeDir))
		// 把当前目录加入到配置文件的搜索路径中，设置配置文件名称与格式
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(defaultConfigName)
	}

	// 读取环境变量，即前缀为AULER或auler的环境变量
	viper.AutomaticEnv()
	viper.SetEnvPrefix("AULER")

	// 替换.和-变成下划线_
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// 读取配置文件
	// ReadInConfig: 如果指定了配置文件名，则使用指定的配置文件，否则在注册的搜索路径中搜索
	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	// 打印
	fmt.Fprintln(os.Stdout, "Using config file:", viper.ConfigFileUsed())
}
