package authsvr

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/xmopen/blogsvr/internal/config"
	"github.com/xmopen/commonlib/pkg/server/authserver"
	"github.com/xmopen/gorpc/pkg/client"
	"github.com/xmopen/gorpc/pkg/netpoll"
)

var (
	authSvrInstance *AuthSvr
	authSvrOnce     sync.Once
	// authSvrRPCServerAddr 兜底的RPC配置信息,如果环境变量中不存在则使用该配置的值
	authSvrRPCServerAddr = config.Config().GetString("server.blogsvr.rpc.authsvrdns")
)

func init() {
	// TODO: 采用临时的解决方案：从环境变量中获取对应的ServiceIP:TargetPort来实现服务发现
	// 未来这里应该归纳到config中
	authSvrRPCHost := os.Getenv("AUTHSVR_SERVICE_HOST")
	authSvrRPCPort := os.Getenv("AUTHSVR_SERVICE_PORT_AUTHSVRRPCPORT")
	if authSvrRPCHost == "" || authSvrRPCPort == "" {
		return
	}
	authSvrRPCServerAddr = fmt.Sprintf("%s:%s", authSvrRPCHost, authSvrRPCPort)
}

// AuthSvr auth server
type AuthSvr struct {
	rpcAuthServer *client.Client
}

// Server return an auth server.
func Server() *AuthSvr {
	authSvrOnce.Do(func() {
		authSvrInstance = &AuthSvr{}
		cli, _ := client.NewClient(netpoll.TCP, authSvrRPCServerAddr, nil)
		cli.Trace = true
		authSvrInstance.rpcAuthServer = cli
	})
	return authSvrInstance
}

// GetUserInfoByAccount get user info by account from gorpc
func (a *AuthSvr) GetUserInfoByAccount(ctx context.Context, request *authserver.AuthSvrRequest, response *authserver.AuthSvrResponse) error {
	return a.rpcAuthServer.Call(ctx, authserver.AuthSvrName,
		string(authserver.AuthSvrMethodTypeOfGetUserInfoByXMAccount), request, response)
}

// GetUserInfoByToken get user info by token from gorpc.
func (a *AuthSvr) GetUserInfoByToken(ctx context.Context, request *authserver.AuthSvrRequest, response *authserver.AuthSvrResponse) error {
	return a.rpcAuthServer.Call(ctx, authserver.AuthSvrName,
		string(authserver.AuthSvrMethodTypeOfGetUserInfoByXMToken), request, response)
}
