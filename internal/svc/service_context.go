package svc

import (
	"github.com/xh-polaris/sts-rpc/internal/config"

	"github.com/tencentyun/qcloud-cos-sts-sdk/go"
)

type ServiceContext struct {
	Config    config.Config
	StsClient *sts.Client
	StsOption *sts.CredentialOptions
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		StsClient: sts.NewClient(
			c.CosConfig.SecretId,
			c.CosConfig.SecretKey,
			nil),
	}
}
