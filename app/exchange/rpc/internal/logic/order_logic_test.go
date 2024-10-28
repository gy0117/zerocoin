package logic

import (
	"exchange-rpc/internal/config"
	"exchange-rpc/internal/svc"
	"flag"
	"fmt"
	"github.com/dtm-labs/client/dtmgrpc"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/exchange/types/order"
	"grpc-common/ucenter/types/wallet"
	"testing"
)

func TestOrder(t *testing.T) {

	var configFile = flag.String("f", "../../etc/conf.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	logx.MustSetup(c.LogConfig)

	ctx := svc.NewServiceContext(c)

	//host := ctx.Config.Etcd.Hosts[0]
	//key := ctx.Config.Etcd.Key
	//target := fmt.Sprintf("etcd://%s/%s", host, key)
	//logx.Info("target: ", target)

	ac(ctx)

}

func ac(ctx *svc.ServiceContext) {
	// DTM 改造
	//dtmServer, err := ctx.Config.DtmConf.BuildTarget()
	dtmServer := "etcd://localhost:2379/dtmservice"
	//if err != nil {
	//	logx.Error("dtmServer-0, err: ", err)
	//	return
	//}
	logx.Info("dtmServer -> ", dtmServer)
	gid := dtmgrpc.MustGenGid(dtmServer)
	logx.Info("gid -> ", gid)
	dtmSaga := dtmgrpc.NewSagaGrpc(dtmServer, gid)

	orderTarget, err := buildTarget(ctx)
	if err != nil {
		logx.Error("dtmServer-1, err: ", err)
		return
	}
	accountTarget, err := ctx.Config.UCenter.BuildTarget()
	if err != nil {
		logx.Error("dtmServer-2, err: ", err)
		return
	}

	var createOrderAddr = orderTarget + "/order.OrderService/CreateOrder"
	var createOrderRevertAddr = orderTarget + "/order.OrderService/CreateOrderRevert"
	var freezeUserWalletAddr = accountTarget + "/wallet.Wallet/FreezeUserAsset"
	var freezeUserWalletRevertAddr = accountTarget + "/wallet.Wallet/UnFreezeUserAsset"

	logx.Info("createOrderAddr: ", createOrderAddr)
	logx.Info("createOrderRevertAddr: ", createOrderRevertAddr)
	logx.Info("freezeUserWalletAddr: ", freezeUserWalletAddr)
	logx.Info("freezeUserWalletRevertAddr: ", freezeUserWalletRevertAddr)

	freezeReq := &wallet.FreezeUserAssetReq{
		Uid: 1,
		//Money: req.Amount,
		Symbol: "btc",
	}

	createOrderReq := &order.OrderReq{}

	saga := dtmSaga.Add(createOrderAddr, createOrderRevertAddr, freezeReq).Add(freezeUserWalletAddr, freezeUserWalletRevertAddr, createOrderReq)
	err = saga.Submit()
	if err != nil {
		logx.Error("saga, err: ", err)
		return
	}
}

func buildTarget(ctx *svc.ServiceContext) (string, error) {
	etcd := ctx.Config.Etcd
	if len(etcd.Hosts) == 0 {
		return "", errors.New("build target failed")
	}
	host := etcd.Hosts[0]
	key := etcd.Key
	target := fmt.Sprintf("etcd://%s/%s", host, key)
	return target, nil
}
