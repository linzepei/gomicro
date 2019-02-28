package handler

import (
	"context"

	"github.com/micro/go-log"

	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	_ "github.com/gomodule/redigo/redis"
	example "go-1/GetUserOrder/proto/example"
	"go-1/homeweb/models"
	"go-1/homeweb/utils"
	"reflect"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetUserOrder(ctx context.Context, req *example.Request, rsp *example.Response) error {

	beego.Info("==============/api/v1.0/user/orders  GetOrders post succ!!=============")
	//创建返回空间
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	//根据session得到当前用户的user_id
	//构建连接缓存的数据
	redis_config_map := map[string]string{
		"key": utils.G_server_name,
		//"conn":"127.0.0.1:6379",
		"conn":     utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum":    utils.G_redis_dbnum,
		"password": "sher",
	}
	beego.Info(redis_config_map)
	redis_config, _ := json.Marshal(redis_config_map)
	beego.Info(string(redis_config))

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
	userid := int(value_id.([]uint8)[0])
	beego.Info(userid, reflect.TypeOf(userid))

	//得到用户角色
	beego.Info(req.Role)

	o := orm.NewOrm()
	orders := []models.OrderHouse{}
	order_list := []interface{}{} //存放订单的切片

	if "landlord" == req.Role {
		//角色为房东
		//现找到自己目前已经发布了哪些房子
		landLordHouses := []models.House{}
		o.QueryTable("house").Filter("user__id", userid).All(&landLordHouses)

		housesIds := []int{}
		for _, house := range landLordHouses {
			housesIds = append(housesIds, house.Id)
		}
		//在从订单中找到房屋id为自己房源的id
		o.QueryTable("order_house").Filter("house__id__in", housesIds).OrderBy("ctime").All(&orders)
	} else {
		//角色为租客
		_, err := o.QueryTable("order_house").Filter("user__id", userid).OrderBy("ctime").All(&orders)
		if err != nil {
			beego.Info(err)
		}

	}
	//循环将数据放到切片中
	for _, order := range orders {
		o.LoadRelated(&order, "User")
		o.LoadRelated(&order, "House")
		order_list = append(order_list, order.To_order_info())
	}

	rsp.Orders, _ = json.Marshal(order_list)

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
