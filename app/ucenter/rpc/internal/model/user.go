package model

const (
	GENERAL = iota
	REAL_NAME
	IDENTIFICATION
)

const (
	NORMAL_PARTER = "0"
	SUPER_PARTER  = "1"
	PSUPER_PARTER = "2"
)

const (
	NORMAL = iota
	ILLEGAL
)

type User struct {
	Id                         int64   `gorm:"id"`
	AliNo                      string  `gorm:"ali_no"`
	QrCodeUrl                  string  `gorm:"qr_code_url"`
	AppealSuccessTimes         int64   `gorm:"appeal_success_times"`
	AppealTimes                int64   `gorm:"appeal_times"`
	ApplicationTime            int64   `gorm:"application_time"`
	Avatar                     string  `gorm:"avatar"`
	Bank                       string  `gorm:"bank"`
	Branch                     string  `gorm:"branch"`
	CardNo                     string  `gorm:"card_no"`
	CertifiedBusinessApplyTime int64   `gorm:"certified_business_apply_time"`
	CertifiedBusinessCheckTime int64   `gorm:"certified_business_check_time"`
	CertifiedBusinessStatus    int64   `gorm:"certified_business_status"`
	ChannelId                  int64   `gorm:"channel_id"`
	Email                      string  `gorm:"email"`
	FirstLevel                 int64   `gorm:"first_level"`
	GoogleDate                 int64   `gorm:"google_date"`
	GoogleKey                  string  `gorm:"google_key"`
	GoogleState                int64   `gorm:"google_state"`
	IdNumber                   string  `gorm:"id_number"`
	InviterId                  int64   `gorm:"inviter_id"`
	IsChannel                  int64   `gorm:"is_channel"`
	JyPassword                 string  `gorm:"jy_password"`
	LastLoginTime              int64   `gorm:"last_login_time"`
	City                       string  `gorm:"city"`
	Country                    string  `gorm:"country"`
	District                   string  `gorm:"district"`
	Province                   string  `gorm:"province"`
	LoginCount                 int64   `gorm:"login_count"`
	LoginLock                  int64   `gorm:"login_lock"`
	Margin                     string  `gorm:"margin"`
	UserLevel                  int64   `gorm:"user_level"`
	MobilePhone                string  `gorm:"mobile_phone"`
	Password                   string  `gorm:"password"`
	PromotionCode              string  `gorm:"promotion_code"`
	PublishAdvertise           int64   `gorm:"publish_advertise"`
	RealName                   string  `gorm:"real_name"`
	RealNameStatus             int64   `gorm:"real_name_status"`
	RegistrationTime           int64   `gorm:"registration_time"`
	Salt                       string  `gorm:"salt"`
	SecondLevel                int64   `gorm:"second_level"`
	SignInAbility              int64   `gorm:"sign_in_ability"`
	Status                     int64   `gorm:"status"`
	ThirdLevel                 int64   `gorm:"third_level"`
	Token                      string  `gorm:"token"`
	TokenExpireTime            int64   `gorm:"token_expire_time"`
	TransactionStatus          int64   `gorm:"transaction_status"`
	TransactionTime            int64   `gorm:"transaction_time"`
	Transactions               int64   `gorm:"transactions"`
	Username                   string  `gorm:"username"`
	QrWeCodeUrl                string  `gorm:"qr_we_code_url"`
	Wechat                     string  `gorm:"wechat"`
	Local                      string  `gorm:"local"`
	Integration                int64   `gorm:"integration"`
	UserGradeId                int64   `gorm:"user_grade_id"`    // 等级id
	KycStatus                  int64   `gorm:"kyc_status"`       // kyc等级
	GeneralizeTotal            int64   `gorm:"generalize_total"` // 注册赠送积分
	InviterParentId            int64   `gorm:"inviter_parent_id"`
	SuperPartner               string  `gorm:"super_partner"`
	KickFee                    float64 `gorm:"kick_fee"`
	Power                      float64 `gorm:"power"`      // 个人矿机算力(每日维护)
	TeamLevel                  int64   `gorm:"team_level"` // 团队人数(每日维护)
	TeamPower                  float64 `gorm:"team_power"` // 团队矿机算力(每日维护)
	UserLevelId                int64   `gorm:"user_level_id"`
}

func (*User) TableName() string {
	return "user"
}

func NewUser() *User {
	return &User{}
}

func (m *User) FillSuperPartner(partner string) {
	if partner == "" {
		m.SuperPartner = NORMAL_PARTER
		m.Status = NORMAL
	} else {
		if partner != NORMAL_PARTER {
			m.SuperPartner = partner
			m.Status = ILLEGAL
		}
	}
}

func (m *User) UserLevelStr() string {
	if m.UserLevel == GENERAL {
		return "普通会员"
	}
	if m.UserLevel == REAL_NAME {
		return "实名"
	}
	if m.UserLevel == IDENTIFICATION {
		return "认证商家"
	}
	return ""
}

func (m *User) UserRate() int32 {
	if m.SuperPartner == NORMAL_PARTER {
		return 0
	}
	if m.SuperPartner == SUPER_PARTER {
		return 1
	}
	if m.SuperPartner == PSUPER_PARTER {
		return 2
	}
	return 0
}
