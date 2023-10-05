package authsvr

import (
	"context"
	"sync"

	"github.com/xmopen/commonlib/pkg/server/authserver"
	"github.com/xmopen/gorpc/pkg/client"
)

var (
	authSvrInstance *AuthSvr
	authSvrOnce     sync.Once
)

// AuthSvr auth server
type AuthSvr struct {
	rcpAuthServer *client.Client
}

// Server return a auth server.
func Server() *AuthSvr {
	authSvrOnce.Do(func() {
		authSvrInstance = &AuthSvr{}
		// TODO: 待优化.
		cli, _ := client.NewClient("tcp", ":18849", nil)
		cli.Trace = true
		authSvrInstance.rcpAuthServer = cli
	})
	return authSvrInstance
}

// GetUserInfoByAccount get user info by account from gorpc
func (a *AuthSvr) GetUserInfoByAccount(ctx context.Context, request *authserver.AuthSvrRequest, response *authserver.AuthSvrResponse) error {
	return a.rcpAuthServer.Call(ctx, authserver.AuthSvrName,
		string(authserver.AuthSvrMethodTypeOfGetUserInfoByXMAccount), request, response)
}

// GetUserInfoByToken get user info by token from gorpc.
func (a *AuthSvr) GetUserInfoByToken(ctx context.Context, request *authserver.AuthSvrRequest, response *authserver.AuthSvrResponse) error {
	return a.rcpAuthServer.Call(ctx, authserver.AuthSvrName,
		string(authserver.AuthSvrMethodTypeOfGetUserInfoByXMToken), request, response)
}
