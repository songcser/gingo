package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/songcser/gingo/config"
	"github.com/songcser/gingo/pkg/model"
	"github.com/songcser/gingo/pkg/response"
	"github.com/songcser/gingo/utils"
	"go.uber.org/zap"
	"net/http"
)

type User interface {
	Login(c *gin.Context) error
	Register(c *gin.Context) error
	Auth(c *gin.Context) error
}

type BaseUser struct {
	model.BaseModel
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (b BaseUser) Login(c *gin.Context) error {
	var req LoginReq
	err := c.ShouldBind(&req)
	utils.CheckError(err)
	m := BaseUser{Name: req.Username}
	mapper := model.NewMapper[BaseUser](m, nil)
	user, err := mapper.SelectOne()
	if err != nil {
		return errors.Errorf("用户不存在: %s", req.Username)
	}
	if ok := utils.BcryptCheck(req.Password, user.Password); !ok {
		return errors.New("用户名或密码错误")
	}
	return b.tokenNext(c, user)
}

func (b BaseUser) Register(c *gin.Context) error {
	var req RegisterReq
	err := c.ShouldBind(&req)
	utils.CheckError(err)
	m := BaseUser{Name: req.Username}
	mapper := model.NewMapper[BaseUser](m, nil)
	_, err = mapper.SelectOne()
	if err == nil {
		return errors.Errorf("用户已存在: %s", req.Username)
	}
	user := BaseUser{
		Name:     req.Username,
		Password: utils.BcryptHash(req.Password),
		Email:    req.Email,
	}
	err = mapper.Insert(&user)
	return err
}

func (b BaseUser) Auth(c *gin.Context) error {
	token, err := c.Cookie("token")
	if err != nil {
		return err
	}
	j := utils.NewJWT()
	claims, err := j.ParseToken(token)
	if err != nil {
		return err
	}
	c.Set("user", claims.Username)
	return nil
}

func (b BaseUser) tokenNext(c *gin.Context, user BaseUser) error {
	j := utils.NewJWT()
	claims := j.CreateClaims(utils.BaseClaims{
		ID:       user.ID,
		Username: user.Name,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		config.GVA_LOG.Error("获取token失败!", zap.Error(err))
		response.FailWithMessage("获取token失败", c)
		return err
	}
	c.SetCookie("token", token, 3600, "/", config.GVA_CONFIG.Domain, false, true)
	c.Redirect(http.StatusFound, "/admin")
	return nil
}
