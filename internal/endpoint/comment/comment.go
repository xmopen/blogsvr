package comment

import (
	"context"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/xmopen/commonlib/pkg/protocol/xmeventprotocol"

	"github.com/xmopen/blogsvr/internal/config"
	"github.com/xmopen/golib/pkg/xgoroutine"

	"github.com/gin-gonic/gin"
	"github.com/xmopen/blogsvr/internal/manager/commentmanager"
	"github.com/xmopen/blogsvr/internal/models/commentmod"
	"github.com/xmopen/blogsvr/internal/server/authsvr"
	ipnetutil "github.com/xmopen/blogsvr/internal/util/ipnet"
	"github.com/xmopen/commonlib/pkg/apphelper/ginhelper"
	"github.com/xmopen/commonlib/pkg/database/xmcomment"
	"github.com/xmopen/commonlib/pkg/errcode"
	"github.com/xmopen/commonlib/pkg/server/authserver"
)

var (
	commentAPIInstance *API
	commentOnce        sync.Once
)

// API comment api
type API struct {
}

type addCommentRequest struct {
	ArticleID int    `json:"article_id"`
	XMToken   string `json:"xm_token"`
	Content   string `json:"content"`
}

// New return a comment api instance
func New() *API {
	commentOnce.Do(func() {
		commentAPIInstance = &API{}
	})
	return commentAPIInstance
}

// GetArticleCommentList get article comment list by articleID
func (a *API) GetArticleCommentList(c *gin.Context) {
	articleIDStr := c.Query("article_id")
	if articleIDStr == "" {
		c.JSON(http.StatusOK, errcode.ErrorParam)
		return
	}
	xlog := ginhelper.Log(c)
	articleID, err := strconv.Atoi(articleIDStr)
	if err != nil {
		xlog.Errorf("article id err:[%+v] source:[%+v]", err, articleIDStr)
		c.JSON(http.StatusOK, errcode.ErrorParam)
		return
	}
	commentList, err := commentmanager.Manager().GetArticleComment(articleID)
	if err != nil {
		xlog.Errorf("get article comment list err:[%+v] article id:[%+v]", err, articleID)
		c.JSON(http.StatusOK, errcode.ErrorParam)
		return
	}
	c.JSON(http.StatusOK, errcode.Success(commentList))
}

// Comment insert comment to article
func (a *API) Comment(c *gin.Context) {
	request := &addCommentRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		c.JSON(http.StatusOK, errcode.ErrorParseJSONError)
		return
	}
	authResponse := &authserver.AuthSvrResponse{}
	err := authsvr.Server().GetUserInfoByToken(context.TODO(), &authserver.AuthSvrRequest{XMToken: request.XMToken},
		authResponse)
	xlog := ginhelper.Log(c)
	if err != nil {
		xlog.Errorf("get openxm user info from authsvr err:[%+v] xmtoken:[%+v]", err, request.XMToken)
		c.JSON(http.StatusOK, errcode.ErrorGORPCError)
		return
	}
	if authResponse == nil || authResponse.XMUserInfo == nil {
		c.JSON(http.StatusOK, errcode.ErrorXMUserTokenExpired)
		return
	}
	xmComment := &xmcomment.XMComment{
		ArticleID:  request.ArticleID,
		Account:    authResponse.XMUserInfo.UserAccount,
		Content:    request.Content,
		Status:     int(commentmod.ArticleCommentStatusOfUp),
		City:       ipnetutil.ParseCityFromIP(c.ClientIP()),
		CreateTime: time.Now(),
	}
	xlog.Infof("add comment from articleid:[%+v] value:[%+v]", request.ArticleID, xmComment)
	if err = commentmod.CreateXMComment(xmComment); err != nil {
		xlog.Errorf("creational xmcomment err:[%+v]", err)
		// TODO: errcode name update
		c.JSON(http.StatusOK, errcode.ErrorCreateArticleError)
		return
	}
	xgoroutine.SafeGoroutine(func() {
		config.BlogsRedis().Publish(context.TODO(), string(xmeventprotocol.XMEventKeyOfArticleCommentUpdate), "")
	})
	c.JSON(http.StatusOK, errcode.Success(nil))
}
