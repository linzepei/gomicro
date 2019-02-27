package handler

import (
	"context"

	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	_ "github.com/garyburd/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
	"gomicro/IhomeWeb/models"
	"gomicro/IhomeWeb/utils"
	example "gomicro/PostUserAuth/proto/example"
	"strconv"
	"time"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PostUserAuth(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("PostUserAuth 实名认证 /api/v1.0/user/auth")

	/*初始化返回值*/
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	/*获取sessionid*/
	sesionid := req.Sessionid

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
		return nil
	}

	/*通过sessionid拼接key 查询user_id*/
	sessionuser_id := sesionid + "user_id"

	user_id := bm.Get(sessionuser_id)
	beego.Info(user_id)
	//beego.Info(reflect.TypeOf(user_id),user_id)
	userid_str, _ := redis.String(user_id, nil)
	beego.Info(userid_str)
	//beego.Info(reflect.TypeOf(userid_str),userid_str)
	id, _ := strconv.Atoi(userid_str)
	beego.Info(id)

	/*通过user_id 更新表 将身份证号和姓名更新到表上*/
	//创建user表单对象
	user := models.User{Id: id, Id_card: req.IdCard, Real_name: req.RealName}

	//创建orm句柄
	o := orm.NewOrm()
	_, err = o.Update(&user, "real_name", "id_card")
	if err != nil {
		beego.Info("身份信息更新失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	/*刷新下session时间*/
	bm.Put(sessionuser_id, userid_str, time.Second*600)

	return nil
}
