package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	CosConfig struct {
		AppId      string
		BucketName string
		Region     string
		SecretId   string
		SecretKey  string
	}
}
