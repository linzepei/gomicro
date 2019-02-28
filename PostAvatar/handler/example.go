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
	example "go-1/PostAvatar/proto/example"
	"go-1/homeweb/models"
	"go-1/homeweb/utils"
	"path"
	"reflect"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PostAvatar(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("上传用户头像 PostAvatar /api/v1.0/user/avatar")
	//初始化返回正确的返回值
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	//检查下数据是否正常
	beego.Info(len(req.Avatar), req.Filesize)

	/*获取文件的后缀名*/ //dsnlkjfajadskfksda.sadsdasd.sdasd.jpg
	beego.Info("后缀名", path.Ext(req.Filename))

	/*存储文件到fastdfs当中并且获取 url*/
	//.jpg
	fileext := path.Ext(req.Filename)
	//group1 group1/M00/00/00/wKgLg1t08pmANXH1AAaInSze-cQ589.jpg
	//上传数据
	FileId, err := models.UploadByBuffer(req.Avatar, fileext[1:])
	if err != nil {
		beego.Info("Postupavatar  models.UploadByBuffer err", err)
		rsp.Errno = utils.RECODE_IOERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	//beego.Info(Group)

	/*通过session 获取我们当前现在用户的uesr_id*/
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
	//拼接key
	sessioniduserid := req.Sessionid + "user_id"

	//获得当前用户的userid
	value_id := bm.Get(sessioniduserid)
	beego.Info(value_id, reflect.TypeOf(value_id))

	id := int(value_id.([]uint8)[0])
	beego.Info(id, reflect.TypeOf(id))

	//创建表对象
	user := models.User{Id: id, Avatar_url: FileId}
	/*将当前fastdfs-url 存储到我们当前用户的表中*/
	o := orm.NewOrm()
	//将图片的地址存入表中
	_, err = o.Update(&user, "avatar_url")
	if err != nil {
		beego.Info("数据更新失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)

	}

	//回传图片地址
	rsp.AvatarUrl = FileId

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
