package auler

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/HeapSoil/auler/internal/auler/controller/v1/user"
	"github.com/HeapSoil/auler/internal/auler/store"
	mw "github.com/HeapSoil/auler/internal/pkg/middleware"
	"github.com/HeapSoil/auler/internal/pkg/utils"
	pb "github.com/HeapSoil/auler/pkg/proto/auler/v1"
	"github.com/HeapSoil/auler/pkg/token"

	"github.com/HeapSoil/auler/internal/pkg/log"
)

var cfgFile string

// 创建*cobra.Command，之后可以使用Command的Execute方法在cmd启动该程序
// 业务的具体实现
func NewAulerCommand() *cobra.Command {
	cmd := &cobra.Command{
		// 命令的名字，长短描述
		Use:   "auler",
		Short: "A demo platform project",
		Long:  "A demo platform project with minimal backend",
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

	if err := initStore(); err != nil {
		return err
	}

	token.Init(viper.GetString("jwt-secret"), utils.XUsernameKey) // 初始化Gin：运行模式debug，创建引擎
	gin.SetMode(viper.GetString("runmode"))
	g := gin.New()

	// 添加中间件
	// gin.Recovery(): 用来捕获任何 panic，并恢复
	// mw.RequestID(): 为每个请求的Header复用或创建X-Request-ID
	mws := []gin.HandlerFunc{gin.Recovery(), mw.NoCache, mw.Cors, mw.Secure, mw.RequestID()}

	g.Use(mws...)

	if err := installRouters(g); err != nil {
		return err
	}

	// 创建 HTTP Server 实例、HTTPS Server gPRC Server 实例
	httpsrv := startServerHTTP(g)
	httpssrv := startServerHTTPS(g)
	grpcsrv := startServerGRPC()

	// 初步运行HTTP服务器
	// log.Infow("Start to listening the incoming requests on http address", "addr", viper.GetString("addr"))
	// 非阻塞的goroutine启动阻塞服务

	// channel接受系统信号
	// quit := make(chan os.Signal, 1)
	quit := make(chan os.Signal, 1)
	// 捕获ctrl+c和kill <pid>信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// 阻塞主程序，接受信号后执行
	<-quit
	log.Infow("Shutting down server...")

	// ctx创建，通知服务器 goroutine 处理完目前请求后关闭服务，限时10s
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 10 秒内关闭服务（将未处理完的请求处理完再关闭服务），超过 10 秒就超时退出
	if err := httpsrv.Shutdown(ctx); err != nil {
		log.Errorw("Insecure Server forced to showdown", "err", err)
		return err
	}
	if err := httpssrv.Shutdown(ctx); err != nil {
		log.Errorw("Secure Server forced to showdown", "err", err)
		return err
	}

	grpcsrv.GracefulStop()

	log.Infow("Server exiting")

	return nil
}

// startServerHTTP 创建 HTTP 服务器，并运行
func startServerHTTP(g *gin.Engine) *http.Server {
	// 创建 HTTP 服务器实例
	httpsrv := &http.Server{
		Addr:    viper.GetString("addr"),
		Handler: g,
	}

	log.Infow("Start to listening the incoming requests on http address", "addr", viper.GetString("addr"))

	// 在 goroutine 中启动 HTTP 服务器
	go func() {
		if err := httpsrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalw(err.Error())
		}
	}()

	return httpsrv
}

// startServerHTTPS 创建 HTTPS 服务器，并运行
func startServerHTTPS(g *gin.Engine) *http.Server {
	// 创建 HTTPS 服务器实例
	httpssrv := &http.Server{
		Addr:    viper.GetString("tls.addr"),
		Handler: g,
	}

	log.Infow("Start to listening the incoming requests on https address", "addr", viper.GetString("tls.addr"))

	// 获取 cert 和 key 存储位置
	certLoc, keyLoc := viper.GetString("tls.cert"), viper.GetString("tls.key")
	if certLoc != "" && keyLoc != "" {
		go func() {
			if err := httpssrv.ListenAndServeTLS(certLoc, keyLoc); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalw(err.Error())
			}
		}()
	}

	return httpssrv

}

// startServerGRPC 创建 gRPC 服务器，并运行
func startServerGRPC() *grpc.Server {
	lis, err := net.Listen("tcp", viper.GetString("grpc.addr"))
	if err != nil {
		log.Fatalw("Failed to listen", "err", err)
	}

	// 创建 GRPC 服务器
	grpcsrv := grpc.NewServer()
	pb.RegisterAulerServer(grpcsrv, user.New(store.S, nil))

	// 在 goroutine 中启动 HTTP 服务器
	log.Infow("Start to listening the incoming requests on grpc address", "addr", viper.GetString("grpc.addr"))

	go func() {
		if err := grpcsrv.Serve(lis); err != nil {
			log.Fatalw(err.Error())
		}
	}()

	return grpcsrv

}
