package handler

import (
	"context"

	"github.com/micro/go-log"

	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/gomodule/redigo/redis"
	example "go-1/PostHouses/proto/example"
	"go-1/homeweb/models"
	"go-1/homeweb/utils"
	"strconv"

	"github.com/astaxie/beego/orm"
	"reflect"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PostHouses(ctx context.Context, req *example.Request, rsp *example.Response) error {

	//打印被调用的函数
	beego.Info("PostHouses 发布房源信息 /api/v1.0/houses ")
	//创建返回空间
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	var Requestmap = make(map[string]interface{})
	json.Unmarshal(req.Max, &Requestmap)
	for key, value := range Requestmap {
		beego.Info(key, value)
	}

	house := models.House{}

	/*插入房源信息*/
	//"title":"上奥世纪中心",
	house.Title = Requestmap["title"].(string)
	//	"price":"666",
	price, _ := strconv.Atoi(Requestmap["price"].(string))
	house.Price = price * 100
	//	"address":"西三旗桥东建材城1号",
	house.Address = Requestmap["address"].(string)
	//	"room_count":"2",
	house.Room_count, _ = strconv.Atoi(Requestmap["room_count"].(string))
	//	"acreage":"60",
	house.Acreage, _ = strconv.Atoi(Requestmap["acreage"].(string))
	//	"unit":"2室1厅",
	house.Unit = Requestmap["unit"].(string)
	//	"capacity":"3",
	house.Capacity, _ = strconv.Atoi(Requestmap["capacity"].(string))
	//	"beds":"双人床2张",
	house.Beds = Requestmap["beds"].(string)
	//	"deposit":"200",
	deposit, _ := strconv.Atoi(Requestmap["deposit"].(string))
	house.Deposit = deposit * 100
	//	"min_days":"3",
	house.Min_days, _ = strconv.Atoi(Requestmap["min_days"].(string))
	//	"max_days":"0",
	house.Max_days, _ = strconv.Atoi(Requestmap["max_days"].(string))

	//设施
	//	"facility":["1","2","3","7","12","14","16","17","18","21","22"]
	facility := []*models.Facility{}

	for _, f_id := range Requestmap["facility"].([]interface{}) {
		fid, _ := strconv.Atoi(f_id.(string))
		fac := &models.Facility{Id: fid}
		facility = append(facility, fac)
	}

	//	"area_id":"5"，地区
	area_id, _ := strconv.Atoi(Requestmap["area_id"].(string))
	area := models.Area{Id: area_id}
	house.Area = &area

	//获得userid
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

	//添加user信息

	user := models.User{Id: userid}
	house.User = &user
	//创建数据库句柄
	o := orm.NewOrm()
	houseid, err := o.Insert(&house)
	if err != nil {

		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	beego.Info(houseid, reflect.TypeOf(houseid), house.Id)

	/*插入到房源与设施信息的多对多表中*/
	m2m := o.QueryM2M(&house, "Facilities")
	num, err := m2m.Add(facility)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	if num == 0 {

		rsp.Errno = utils.RECODE_NODATA
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	rsp.House_Id = int64(house.Id)

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
