package logic

import (
	"context"
	"fmt"

	"time"

	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/xh-polaris/sts-rpc/errorx"
	"github.com/xh-polaris/sts-rpc/internal/svc"
	"github.com/xh-polaris/sts-rpc/pb"
)

type GetUserCosStsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserCosStsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserCosStsLogic {
	return &GetUserCosStsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserCosStsLogic) GetUserCosSts(in *pb.GetUserCosStsReq) (*pb.GetUserCosStsResp, error) {
	if len(in.UserId) != 24 {
		return nil, errorx.ErrInvalidUserId
	}

	cosConfig := l.svcCtx.Config.CosConfig
	stsOption := &sts.CredentialOptions{
		// 临时密钥有效时长，单位是秒
		DurationSeconds: int64(time.Minute.Seconds()),
		Region:          cosConfig.Region,
		Policy: &sts.CredentialPolicy{
			Statement: []sts.CredentialPolicyStatement{
				{
					// 密钥的权限列表。简单上传和分片需要以下的权限，其他权限列表请看 https://caloud.tencent.com/document/product/436/31923
					Action: []string{
						// 简单上传
						"name/cos:PostObject",
						"name/cos:PutObject",
						// 分片上传
						"name/cos:InitiateMultipartUpload",
						"name/cos:ListMultipartUploads",
						"name/cos:ListParts",
						"name/cos:UploadPart",
						"name/cos:CompleteMultipartUpload",
					},
					Effect: "allow",
					// 密钥可控制的资源列表。此处开放名字为用户ID的文件夹及其子文件夹
					Resource: []string{
						fmt.Sprintf("qcs::cos:%s:uid/%s:%s/%s/%s",
							cosConfig.Region, cosConfig.AppId, cosConfig.BucketName, in.UserId, in.Path),
					},
				},
			},
		},
	}

	res, err := l.svcCtx.StsClient.GetCredential(stsOption)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserCosStsResp{
		SecretId:     res.Credentials.TmpSecretID,
		SecretKey:    res.Credentials.TmpSecretKey,
		SessionToken: res.Credentials.SessionToken,
	}, nil
}
