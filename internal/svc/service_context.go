package svc

import (
	"deep/internal/config"
	"deep/internal/infra/logger"
	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	Config config.Config
	Logger *logger.ZapLogger
}

func NewServiceContext(c config.Config) *ServiceContext {
	logInstance, err := logger.New()
	if err != nil {
		panic(err) // 现在你的 New() 是不会返回 err，直接写也无妨
	}

	logx.SetWriter(logInstance.Writer())

	return &ServiceContext{
		Config: c,
	}
}
