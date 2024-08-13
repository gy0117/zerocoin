package logic

import (
	"context"
	"github.com/pkg/errors"

	"github.com/golang-jwt/jwt/v4"
	"grpc-common/ucenter/types/login"
	"time"
	"ucenter-rpc/internal/domain"
	"ucenter-rpc/internal/svc"
	"ucenter-rpc/internal/verify"
	"zero-common/tools"
	"zero-common/zerr"

	"github.com/zeromicro/go-zero/core/logx"
)

var ErrUserLogin = zerr.NewCodeErr(zerr.USER_LOGIN_ERROR)
var ErrPhoneNotExist = zerr.NewCodeErr(zerr.USER_PHONE_NOT_EXIST_ERROR)
var ErrPassword = zerr.NewCodeErr(zerr.USER_PASSWORD_ERROR)
var ErrGenerateToken = zerr.NewCodeErr(zerr.TOKEN_GENERATE_ERROR)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	captchaVerify *verify.MachineVerify
	userDomain    *domain.UserDomain
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:           ctx,
		svcCtx:        svcCtx,
		Logger:        logx.WithContext(ctx),
		captchaVerify: verify.NewMachineVerify(),
		userDomain:    domain.NewUserDomain(svcCtx.DB),
	}
}

// Login 登录逻辑
// 输入账号和密码
// 校验人机
// 根据账号查询用户的salt等值
// 密码进行匹配
// 匹配成功，使用jwt生成token
// 返回登录所用信息
func (l *LoginLogic) Login(in *login.LoginReq) (*login.LoginResp, error) {
	// 如果使用postman的话，就不走人机验证
	logx.Info("in.Env = " + in.Env)
	// TODO 暂时不走人机验证
	//if in.Env != "dev_postman" {
	//	// 1. 人机验证
	//	isVerify := l.captchaVerify.Verify(
	//		l.svcCtx.Config.CaptchaVerify.Vid,
	//		l.svcCtx.Config.CaptchaVerify.SecretKey,
	//		in.Captcha.Server,
	//		in.Captcha.Token,
	//		in.Ip,
	//		verify.RegisterScene,
	//	)
	//
	//	if !isVerify {
	//		return nil, errors.New("人机验证不通过")
	//	}
	//}

	// 2. 密码验证
	// in.Username是手机号
	user, err := l.userDomain.FindByPhone(l.ctx, in.Username)
	if err != nil {
		return nil, errors.Wrapf(ErrUserLogin, "查询手机号失败 phone: %s, err: %v", in.Username, err)
	}
	if user == nil {
		return nil, errors.Wrapf(ErrPhoneNotExist, "手机号不存在 phone: %s", in.Username)
	}
	if ok := tools.Verify(in.Password, user.Salt, user.Password, nil); !ok {
		return nil, errors.Wrapf(ErrPassword, "密码错误 phone: %s", in.Username)
	}

	// 3. 登录成功，将jwt token返回给前端
	//token, err := l.generateToken(user.Id, user.Username)
	token, err := l.getJwtToken(l.svcCtx.Config.Jwt.AccessSecret, time.Now().Unix(), l.svcCtx.Config.Jwt.AccessExpire, user.Id)
	if err != nil {
		return nil, errors.Wrapf(ErrGenerateToken, "token生成失败 phone: %s, err: %v", in.Username, err)
	}

	return &login.LoginResp{
		Token:         token,
		Id:            user.Id,
		Username:      user.Username,
		UserLevel:     user.UserLevelStr(),
		UserRate:      user.UserRate(),
		RealName:      user.RealName,
		Country:       user.Country,
		Avatar:        user.Avatar,
		PromotionCode: user.PromotionCode,
		SuperPartner:  user.SuperPartner,
	}, nil
}

type Claims struct {
	ID       int64  `json:"id"`
	UserName string `json:"username"`
	jwt.StandardClaims
}

func (l *LoginLogic) getJwtToken(secretKey string, iat, seconds, userId int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}

//func (l *LoginLogic) generateToken(id int64, username string) (token string, err error) {
//	// 不能加password，header和payload都是base64编码，没有加密的，
//	// 因此payload里面不能加敏感信息
//	claims := Claims{
//		ID:       id,
//		UserName: username,
//		StandardClaims: jwt.StandardClaims{
//			ExpiresAt: time.Now().Add(time.Duration(l.svcCtx.Config.Jwt.AccessExpire)).Unix(), // 过期时间
//			Issuer:    l.svcCtx.Config.Jwt.Issuer,                                             // 签发人
//		},
//	}
//
//	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//
//	token, err = tokenClaims.SignedString(l.svcCtx.Config.Jwt.AccessSecret)
//	return
//}

func (l *LoginLogic) parseToken(token string) (*Claims, error) {

	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return l.svcCtx.Config.Jwt.AccessSecret, nil
	})
	if err != nil {
		return nil, err
	}
	// 验证token
	if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
		return claims, nil
	}
	return nil, errors.New("valid token")
}
