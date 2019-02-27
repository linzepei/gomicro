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
	example "gomicro/GetUserHouses/proto/example"
	"gomicro/IhomeWeb/models"
	"gomicro/IhomeWeb/utils"
	"strconv"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetUserHouses(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("获取用户已发部房源 GetUserHouses api/v1.0/user/houses")

	/*初始化 返回值*/
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	/*获取sessionid*/
	sessionid := req.Sessionid
	/*连接redis*/
	//配置缓存参数
	redis_conf := map[string]string{
		"key": utils.G_server_name,
		//127.0.0.1:6379
		"conn":  utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum": utils.G_redis_dbnum,
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

	/*拼接key*/
	sessionid_user_id := sessionid + "user_id"
	/*查询对应的user_id*/
	user_id := bm.Get(sessionid_user_id)
	//转换格式
	user_id_str, _ := redis.String(user_id, nil)
	id, _ := strconv.Atoi(user_id_str)

	/*查询数据库*/
	//创建orm句柄
	o := orm.NewOrm()
	qs := o.QueryTable("house")

	houses_list := []models.House{}
	/*获得当前用户房屋信息*/
	_, err = qs.Filter("user_id", id).All(&houses_list)
	if err != nil {
		beego.Info("查询房屋数据失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	/*json编码成为二进制返回*/
	house, _ := json.Marshal(houses_list)
	//返回二进制数据
	rsp.Mix = house

	return nil
}
