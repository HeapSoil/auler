package auler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	mw "github.com/HeapSoil/auler/internal/pkg/middleware"

	"github.com/HeapSoil/auler/internal/pkg/log"
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

			// 初始化日志包的logger，在应用退出时将缓存中的日志写入磁盘防止丢失
			log.Init(logOptionsFromDefaultConfig())
			defer log.Sync()

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

	// 初始化Gin：运行模式debug，创建引擎
	gin.SetMode(viper.GetString("runmode"))
	g := gin.New()

	// 添加中间件
	// gin.Recovery(): 用来捕获任何 panic，并恢复
	// mw.RequestID(): 为每个请求的Header复用或创建X-Request-ID
	mws := []gin.HandlerFunc{gin.Recovery(), mw.NoCache, mw.Cors, mw.Secure, mw.RequestID()}

	g.Use(mws...)

	// Handler注册：404，/healthz
	g.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 10003, "message": "Page not found."})
	})

	g.GET("/healthz", func(c *gin.Context) {
		log.C(c).Infow("Healthz function called")
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 创建HTTP Server实例
	httpsrv := &http.Server{Addr: viper.GetString("addr"), Handler: g}

	// 初步运行HTTP服务器
	log.Infow("Start to listening the incoming requests on http address", "addr", viper.GetString("addr"))
	if err := httpsrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalw(err.Error())
	}

	return nil
}
