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
	example "gomicro/PostHouses/proto/example"
	"strconv"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PostHouses(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("PostHouses  发布房源  /api/v1.0/houses")
	/*初始化返回值*/
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	/*获取sessionid*/
	sessionid := req.Sessionid

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

	/*拼接key*/
	session_user_id := sessionid + "user_id"
	/*查询user_id*/
	user_id := bm.Get(session_user_id)
	/*转换user_id类型*/
	user_id_str, _ := redis.String(user_id, nil)
	id, _ := strconv.Atoi(user_id_str)

	/*解析对段发送过来的body*/
	var request = make(map[string]interface{})
	json.Unmarshal(req.Body, &request)

	/*准备插入数据库的对象*/
	house := models.House{}
	//"title":"上奥世纪中心",
	house.Title = request["title"].(string)
	//	"price":"666",
	price, _ := strconv.Atoi(request["price"].(string))
	house.Price = price * 100
	//	"address":"西三旗桥东建材城1号",
	house.Address = request["address"].(string)
	//	"room_count":"2",
	house.Room_count, _ = strconv.Atoi(request["room_count"].(string))
	//	"acreage":"60",
	house.Acreage, _ = strconv.Atoi(request["acreage"].(string))
	//	"unit":"2室1厅",
	house.Unit = request["unit"].(string)
	//	"capacity":"3",
	house.Capacity, _ = strconv.Atoi(request["capacity"].(string))
	//	"beds":"双人床2张",
	house.Beds = request["beds"].(string)
	//	"deposit":"200",
	deposit, _ := strconv.Atoi(request["deposit"].(string))
	house.Deposit = deposit * 100
	//	"min_days":"3",
	house.Min_days, _ = strconv.Atoi(request["min_days"].(string))
	//	"max_days":"0",
	house.Max_days, _ = strconv.Atoi(request["max_days"].(string))

	//"area_id":"5",
	area_id, _ := strconv.Atoi(request["area_id"].(string))
	area := models.Area{Id: area_id}
	house.Area = &area

	//"facility":["1","2","3","7","12","14","16","17","18","21","22"]
	facility := []*models.Facility{}

	for _, value := range request["facility"].([]interface{}) {
		//将设施编号转换成为应对的类型
		fid, _ := strconv.Atoi(value.(string))
		//创建临时变量 ， 使用设施编号创建的设施表对象的指针
		ftmp := &models.Facility{Id: fid}
		facility = append(facility, ftmp)
	}

	/*数据库插入操作*/
	user := models.User{Id: id}
	house.User = &user

	//创建orm句柄
	o := orm.NewOrm()
	house_id, err := o.Insert(&house)
	if err != nil {
		beego.Info("数据插入失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	/*多对多的插入*/
	m2m := o.QueryM2M(&house, "Facilities")

	_, err = m2m.Add(facility)
	if err != nil {
		beego.Info("房屋设施多对多插入失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	/*返回houses_id*/
	rsp.HousesId = strconv.Itoa(int(house_id))

	return nil
}
