package logic

import (
	"context"
	"time"
	"ucenter-rpc/internal/domain"
	"ucenter-rpc/internal/verify"
	"zero-common/tools"
	"zero-common/zerr"

	"grpc-common/ucenter/types/register"
	"ucenter-rpc/internal/svc"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

const RegisterVerifyCode = "REGISTER::VERIFY::"

var ErrUserRegister = zerr.NewCodeErr(zerr.USER_REGISTER_ERROR)
var ErrUserHasRegistered = zerr.NewCodeErr(zerr.USER_HAS_REGISTERED_ERROR)
var ErrUserVerifyCode = zerr.NewCodeErr(zerr.USER_VERIFY_CODE_ERROR)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	captchaVerify *verify.MachineVerify
	userDomain    *domain.UserDomain
}

func NewRegisterByPhoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:           ctx,
		svcCtx:        svcCtx,
		Logger:        logx.WithContext(ctx),
		captchaVerify: verify.NewMachineVerify(),
		userDomain:    domain.NewUserDomain(svcCtx.DB),
	}
}

func (l *RegisterLogic) RegisterByPhone(in *register.RegisterReq) (*register.RegisterResp, error) {

	// 使用postman，就不走人机验证了
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
	//	if !isVerify {
	//		return nil, errors.New("人机验证不通过")
	//	}
	//	logx.Info("人机验证通过...")
	//}

	// 2. 验证码校验
	//ctx, cancel := context.WithTimeout(l.ctx, time.Second*3)
	//defer cancel()

	resultVal := ""
	key := RegisterVerifyCode + in.Phone

	err := l.svcCtx.Cache.GetCtx(l.ctx, key, &resultVal)
	if err != nil || resultVal == "" {
		return nil, errors.Wrapf(ErrUserRegister, "验证失败 phone: %s, err: %v", in.Phone, err)
	}
	if in.Code != resultVal {
		return nil, errors.Wrapf(ErrUserVerifyCode, "验证码校验失败 phone: %s", in.Phone)
	}
	// 3. 验证码通过，进行注册即可，手机号首先验证此手机号是否注册过
	user, err := l.userDomain.FindByPhone(context.Background(), in.Phone)
	if err != nil {
		return nil, errors.Wrapf(ErrUserRegister, "查询失败 phone: %s, err: %v", in.Phone, err)
	}
	if user != nil {
		return nil, errors.Wrapf(ErrUserHasRegistered, "用户已经注册了 phone: %s, err: %v", in.Phone, err)
	}

	// 4. 注册
	err = l.userDomain.Register(l.ctx, in.Username, in.Phone, in.Password, in.Country, in.Promotion, in.SuperPartner)
	if err != nil {
		return nil, errors.Wrapf(ErrUserRegister, "用户注册失败 phone: %s, err: %v", in.Phone, err)
	}
	logx.Info("RPC-REGISTER | register success!")

	return &register.RegisterResp{}, nil
}

// SendCode 验证码逻辑：
//
//	收到手机号和国家标识
//	生成验证码
//	根据对应的国家和手机号调用对应的短信平台发送验证码
//	将验证码存入redis，过期时间10分钟
//	返回成功
func (l *RegisterLogic) SendCode(in *register.CodeReq) (*register.CodeResp, error) {

	// 1. 生成验证码
	verifyCode := tools.GenerateVerifyCode()
	logx.Infof("verifyCode: %s\n", verifyCode)

	go func() {
		// 2. 发送验证码
		logx.Info("RPC-REGISTER | rpc sendCode 发送验证码")
	}()

	// 3. 将验证码存入redis，过期时间10分钟
	key := RegisterVerifyCode + in.Phone
	err := l.svcCtx.Cache.SetWithExpireCtx(l.ctx, key, verifyCode, 10*time.Minute)
	if err != nil {
		return nil, errors.Wrapf(ErrUserRegister, "发送验证码处理失败 phone: %s, err: %v", in.Phone, err)
	}
	return &register.CodeResp{
		SmsCode: verifyCode,
	}, nil
}
