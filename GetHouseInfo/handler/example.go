package handler

import (
	"context"

	"github.com/micro/go-log"

	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	_ "github.com/gomodule/redigo/redis"
	example "go-1/GetHouseInfo/proto/example"
	"go-1/homeweb/models"
	"go-1/homeweb/utils"
	"reflect"
	"strconv"
	"time"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetHouseInfo(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("获取房源详细信息 GetHouseInfo  api/v1.0/houses/:id ")

	//创建返回空间
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	/*从session中获取我们的user_id的字段 得到当前用户id*/
	/*通过session 获取我们当前登陆用户的user_id*/
	//构建连接缓存的数据
	redis_config_map := map[string]string{
		"key":      utils.G_server_name,
		"conn":     utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum":    utils.G_redis_dbnum,
		"password": "sher",
	}
	beego.Info(redis_config_map)
	redis_config, _ := json.Marshal(redis_config_map)

	//连接redis数据库 创建句柄
	bm, err := cache.NewCache("redis", string(redis_config))
	if err != nil {
		beego.Info("缓存创建失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	sessioniduserid := req.Sessionid + "user_id"

	value_id := bm.Get(sessioniduserid)
	beego.Info(value_id, reflect.TypeOf(value_id))
	id := int(value_id.([]uint8)[0])
	beego.Info(id, reflect.TypeOf(id))

	/*从请求中的url获取房源id*/

	houseid, _ := strconv.Atoi(req.Id)

	/*从缓存数据库中获取到当前房屋的数据*/

	house_info_key := fmt.Sprintf("house_info_%s", houseid)
	house_info_value := bm.Get(house_info_key)
	if house_info_value != nil {
		rsp.Userid = int64(id)
		rsp.Housedata = house_info_value.([]byte)

	}

	/*查询当前数据库得到当前的house详细信息*/
	//创建数据对象
	house := models.House{Id: houseid}
	//创建数据库句柄
	o := orm.NewOrm()
	o.Read(&house)
	/*关联查询 area user images fac等表*/
	o.LoadRelated(&house, "Area")
	o.LoadRelated(&house, "User")
	o.LoadRelated(&house, "Images")
	o.LoadRelated(&house, "Facilities")
	//o.LoadRelated(&house,"Orders")

	/*将查询到的结果存储到缓存当中*/
	housemix, err := json.Marshal(house)
	bm.Put(house_info_key, housemix, time.Second*3600)

	/*返回正确数据给前端*/

	rsp.Userid = int64(id)
	rsp.Housedata = housemix

	return nil

}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Example) Stream(ctx context.Context, req *example.StreamingRequest, stream example.Example_StreamStream) error {
	log.Logf("Received Example.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&example.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Example) PingPong(ctx context.Context, stream example.Example_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&example.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
