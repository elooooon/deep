package main

import (
	"deep/internal/middleware"
	"flag"
	"github.com/zeromicro/go-zero/core/logx"

	"deep/internal/config"
	"deep/internal/handler"
	"deep/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/deep.yaml", "the config file")

//goland:noinspection HttpUrlsUsage
func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf, rest.WithNotFoundHandler(middleware.NotFoundHandler()))
	defer server.Stop()

	server.Use(middleware.ZapLogger())

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	logx.Infof("Starting server at http://%s:%d", c.Host, c.Port)
	server.Start()
}
