package handler

import (
	"context"

	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/garyburd/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
	example "gomicro/DeleteSession/proto/example"
	"gomicro/IhomeWeb/utils"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) DeleteSession(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("DeleteSession  退出登陆 /api/v1.0/session")
	//返回值初始化
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	/*连接redis*/
	//配置缓存参数
	redis_conf := map[string]string{
		"key": utils.G_server_name,
		//127.0.0.1:6379
		"conn":     utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum":    utils.G_redis_dbnum,
		"password": "sher",
	}
	beego.Info(redis_conf)

	//将map进行转化成为json
	redis_conf_js, _ := json.Marshal(redis_conf)

	//创建redis句柄
	bm, err := cache.NewCache("redis", string(redis_conf_js))
	if err != nil {
		beego.Info("redis连接失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
	}

	//获取sessionid
	sessionid := req.Sessionid

	/*拼接key  删除session相关字段*/
	//user_id
	sessionuser_id := sessionid + "user_id"
	bm.Delete(sessionuser_id)
	//name
	sessionname := sessionid + "name"
	bm.Delete(sessionname)

	//mobile
	sessionmobile := sessionid + "mobile"
	bm.Delete(sessionmobile)

	return nil
}
