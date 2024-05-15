package verify

import (
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"zero-common/tools"
)

const (
	LoginScene    = 1
	RegisterScene = 2
)

type MachineReq struct {
	Id        string `json:"id"`
	SecretKey string `json:"secretkey"`
	Scene     int    `json:"scene"`
	Token     string `json:"token"`
	Ip        string `json:"ip"`
}

type MachineResp struct {
	Success int    `json:"success"`
	Score   int    `json:"score"`
	Msg     string `json:"msg"`
}

type MachineVerify struct {
}

func NewMachineVerify() *MachineVerify {
	return &MachineVerify{}
}

// Verify  人机验证
// vid: 验证单元的VID
// secretKey: 验证单元的Key
// token: 用户在前端验证成功后取得的token
// ip: 获取用户的remote address
// scene: 配置验证单元的场景ID
func (mv *MachineVerify) Verify(vid, secretKey, server, token, ip string, scene int) bool {

	req := &MachineReq{
		Id:        vid,
		SecretKey: secretKey,
		Scene:     scene,
		Token:     token,
		Ip:        ip,
	}
	fmt.Printf("Verify---machineReq: %+v", req)

	resp, err := tools.Post(server, req)
	if err != nil {
		logx.Error("failed to post, err: " + err.Error())
		return false
	}
	result := &MachineResp{}
	if err = json.Unmarshal(resp, &result); err != nil {
		logx.Error("failed to Unmarshal, err: " + err.Error())
		return false
	}
	fmt.Printf("Verify---result: %+v", result)
	return result.Success == 1
}
