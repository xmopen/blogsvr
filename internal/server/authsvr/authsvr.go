package authsvr

import (
	"context"
	"sync"

	"github.com/xmopen/blogsvr/internal/config"

	"github.com/xmopen/commonlib/pkg/server/authserver"
	"github.com/xmopen/gorpc/pkg/client"
)

var (
	authSvrInstance *AuthSvr
	authSvrOnce     sync.Once

	authSvrRPCServerAddr = config.Config().GetString("server.blogsvr.rpc.authsvrdns")
)

// AuthSvr auth server
type AuthSvr struct {
	rpcAuthServer *client.Client
}

// Server return an auth server.
func Server() *AuthSvr {
	authSvrOnce.Do(func() {
		authSvrInstance = &AuthSvr{}
		cli, _ := client.NewClient("tcpprotocol", authSvrRPCServerAddr, nil)
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
