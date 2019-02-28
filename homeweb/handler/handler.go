package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/afocus/captcha"
	"github.com/astaxie/beego"
	"github.com/julienschmidt/httprouter"
	example "github.com/micro/examples/template/srv/proto/example"
	"github.com/micro/go-grpc"
	"github.com/micro/go-micro/client"
	DELETESESSION "go-1/DeleteSession/proto/example"
	GETAREA "go-1/GetArea/proto/example"
	GETHOUSEINFO "go-1/GetHouseInfo/proto/example"
	GETHOUSES "go-1/GetHouses/proto/example"
	GETIMAGECD "go-1/GetImageCd/proto/example"
	GETINDEX "go-1/GetIndex/proto/example"
	GETSESSION "go-1/GetSession/proto/example"
	GETSMSCD "go-1/GetSmscd/proto/example"
	GETUSERHOUSES "go-1/GetUserHouses/proto/example"
	GETUSERINFO "go-1/GetUserInfo/proto/example"
	GETUSERORDER "go-1/GetUserOrder/proto/example"
	POSTAVATAR "go-1/PostAvatar/proto/example"
	POSTHOUSES "go-1/PostHouses/proto/example"
	POSTHOUSESIMAGE "go-1/PostHousesImage/proto/example"
	POSTLOGIN "go-1/PostLogin/proto/example"
	POSTORDERS "go-1/PostOrders/proto/example"
	POSTRET "go-1/PostRet/proto/example"
	POSTUSERAUTH "go-1/PostUserAuth/proto/example"
	PUTCOMMENT "go-1/PutComment/proto/example"
	PUTORDERS "go-1/PutOrders/proto/example"
	PUTUSERINFO "go-1/PutUserInfo/proto/example"
	"go-1/homeweb/models"
	"go-1/homeweb/utils"
	"image"
	"image/png"

	"reflect"

	"fmt"

	"io/ioutil"
)

//默认模板
func ExampleCall(w http.ResponseWriter, r *http.Request) {
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// call the backend service
	exampleClient := example.NewExampleService("go.micro.srv.template", client.DefaultClient)
	rsp, err := exampleClient.Call(context.TODO(), &example.Request{
		Name: request["name"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"msg": rsp.Msg,
		"ref": time.Now().UnixNano(),
	}

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

//获取地区
func GetArea(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("获取地区请求客户端 url：api/v1.0/areas")

	//创建新的grpc返回句柄
	server := grpc.NewService()
	//服务出初始化
	server.Init()

	//创建获取地区的服务并且返回句柄
	exampleClient := GETAREA.NewExampleService("go.micro.srv.GetArea", server.Client())

	//调用函数并且获得返回数据
	rsp, err := exampleClient.GetArea(context.TODO(), &GETAREA.Request{})
	if err != nil {
		http.Error(w, err.Error(), 502)
		return
	}
	//创建返回类型的切片
	area_list := []models.Area{}
	//循环读取服务返回的数据
	for _, value := range rsp.Data {
		tmp := models.Area{Id: int(value.Aid), Name: value.Aname, Houses: nil}
		area_list = append(area_list, tmp)
	}
	//创建返回数据map
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   area_list,
	}
	w.Header().Set("Content-Type", "application/json")

	// 将返回数据map发送给前端
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
}

//获取验证码图片
func GetImageCd(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	beego.Info("获取图片验证码 url：api/v1.0/imagecode/:uuid")
	//创建服务
	server := grpc.NewService()
	//服务初始化
	server.Init()

	//连接服务
	exampleClient := GETIMAGECD.NewExampleService("go.micro.srv.GetImageCd", server.Client())

	//获取前端发送过来的唯一uuid
	beego.Info(ps.ByName("uuid"))
	//通过句柄调用我们proto协议中准备好的函数
	//第一个参数为默认,第二个参数 proto协议中准备好的请求包
	rsp, err := exampleClient.GetImageCd(context.TODO(), &GETIMAGECD.Request{
		Uuid: ps.ByName("uuid"),
	})
	//判断函数调用是否成功
	if err != nil {
		beego.Info(err)
		http.Error(w, err.Error(), 502)
		return
	}

	//处理前端发送过来的图片信息
	var img image.RGBA

	img.Stride = int(rsp.Stride)

	img.Rect.Min.X = int(rsp.Min.X)
	img.Rect.Min.Y = int(rsp.Min.Y)
	img.Rect.Max.X = int(rsp.Max.X)
	img.Rect.Max.Y = int(rsp.Max.Y)

	img.Pix = []uint8(rsp.Pix)

	var image captcha.Image

	image.RGBA = &img

	//将图片发送给前端
	png.Encode(w, image)
}

//获取短信验证
func Getsmscd(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	beego.Info(" 获取短信验证   api/v1.0/smscode/:id ")

	//创建服务
	service := grpc.NewService()
	service.Init()

	//获取 前端发送过来的手机号
	mobile := ps.ByName("mobile")
	beego.Info(mobile)
	//后端进行正则匹配
	/*	//创建正则句柄
		myreg := regexp.MustCompile(`0?(13|14|15|17|18|19)[0-9]{9}`)
		//进行正则匹配
		bo := myreg.MatchString(mobile)

		//如果手机号错误则
		if bo == false {
			// we want to augment the response
			resp := map[string]interface{}{
				"errno":  utils.RECODE_NODATA,
				"errmsg": "手机号错误",
			}
			//设置返回数据格式
			w.Header().Set("Content-Type", "application/json")

			//将错误发送给前端
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				http.Error(w, err.Error(), 503)
				beego.Info(err)
				return
			}
			beego.Info("手机号错误返回")
			return
		}
	*/

	//获取url携带的验证码 和key（uuid）
	beego.Info(r.URL.Query())
	//获取url携带的参数
	text := r.URL.Query()["text"][0] //text=248484

	id := r.URL.Query()["id"][0] //id=9cd8faa9-5653-4f7c-b653-0a58a8a98c81

	//调用服务
	exampleClient := GETSMSCD.NewExampleService("go.micro.srv.GetSmscd", service.Client())
	rsp, err := exampleClient.GetSmscd(context.TODO(), &GETSMSCD.Request{
		Mobile: mobile,
		Id:     id,
		Text:   text,
	})

	if err != nil {
		http.Error(w, err.Error(), 502)
		beego.Info(err)
		//beego.Debug(err)
		return
	}
	//创建返回map
	resp := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}
	//设置返回格式
	w.Header().Set("Content-Type", "application/json")

	//将数据回发给前端
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), 503)
		beego.Info(err)
		return
	}
}

//注册请求
func Postreg(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	beego.Info(" 注册请求   /api/v1.0/users ")
	/*获取前端发送过来的json数据*/
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	for key, value := range request {
		beego.Info(key, value, reflect.TypeOf(value))
	}

	//由于前端每作所以后端进行下操作
	if request["mobile"] == "" || request["password"] == "" || request["sms_code"] == "" {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_NODATA,
			"errmsg": "信息有误请从新输入",
		}

		//如果不存在直接给前端返回
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		beego.Info("有数据为空")
		return
	}

	//创建服务
	service := grpc.NewService()
	service.Init()

	// 连接服务将数据发送给注册服务进行注册
	exampleClient := POSTRET.NewExampleService("go.micro.srv.PostRet", service.Client())
	rsp, err := exampleClient.PostRet(context.TODO(), &POSTRET.Request{
		Mobile:   request["mobile"].(string),
		Password: request["password"].(string),
		SmsCode:  request["sms_code"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 502)

		beego.Info(err)
		//beego.Debug(err)
		return
	}

	resp := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}

	//读取cookie
	cookie, err := r.Cookie("userlogin")

	//如果读取失败或者cookie的value中不存在则创建cookie
	if err != nil || "" == cookie.Value {
		cookie := http.Cookie{Name: "userlogin", Value: rsp.SessionID, Path: "/", MaxAge: 600}
		http.SetCookie(w, &cookie)
	}

	//设置回发数据格式
	w.Header().Set("Content-Type", "application/json")

	//将数据回发给前端
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), 503)
		beego.Info(err)
		return
	}

	return
}

//获取session
func GetSession(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("获取Session url：api/v1.0/session")

	//创建服务
	service := grpc.NewService()
	service.Init()

	//创建句柄
	exampleClient := GETSESSION.NewExampleService("go.micro.srv.GetSession", service.Client())

	//获取cookie
	userlogin, err := r.Cookie("userlogin")

	//如果不存在就返回
	if err != nil {

		//创建返回数据map
		response := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		return
	}

	//存在就发送数据给服务
	rsp, err := exampleClient.GetSession(context.TODO(), &GETSESSION.Request{
		Sessionid: userlogin.Value,
	})

	if err != nil {
		http.Error(w, err.Error(), 502)

		beego.Info(err)
		//beego.Debug(err)
		return
	}

	// we want to augment the response
	//将获取到的用户名返回给前端
	data := make(map[string]string)
	data["name"] = rsp.Data
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}

	w.Header().Set("Content-Type", "application/json")

	// 将返回数据map发送给前端
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
}

//登陆
func PostLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("登陆 api/v1.0/sessions")
	//获取前端post请求发送的内容
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	for key, value := range request {
		beego.Info(key, value, reflect.TypeOf(value))
	}
	//判断账号密码是否为空
	if request["mobile"] == "" || request["password"] == "" {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_NODATA,
			"errmsg": "信息有误请从新输入",
		}
		w.Header().Set("Content-Type", "application/json")

		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		beego.Info("有数据为空")
		return
	}

	//创建连接

	service := grpc.NewService()
	service.Init()
	exampleClient := POSTLOGIN.NewExampleService("go.micro.srv.PostLogin", service.Client())

	rsp, err := exampleClient.PostLogin(context.TODO(), &POSTLOGIN.Request{
		Password: request["password"].(string),
		Mobile:   request["mobile"].(string),
	})

	if err != nil {
		http.Error(w, err.Error(), 502)

		beego.Info(err)
		//beego.Debug(err)
		return
	}

	cookie, err := r.Cookie("userlogin")
	if err != nil || "" == cookie.Value {
		cookie := http.Cookie{Name: "userlogin", Value: rsp.SessionID, Path: "/", MaxAge: 600}
		http.SetCookie(w, &cookie)
	}
	beego.Info(rsp.SessionID)
	resp := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), 503)
		beego.Info(err)
		return
	}
}

//退出
func DeleteSession(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("---------------- DELETE  /api/v1.0/session Deletesession() ------------------")
	//创建返回空间

	server := grpc.NewService()
	server.Init()
	exampleClient := DELETESESSION.NewExampleService("go.micro.srv.DeleteSession", server.Client())

	//获取session
	userlogin, err := r.Cookie("userlogin")
	//如果没有数据说明没有的登陆直接返回错误
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		return
	}

	rsp, err := exampleClient.DeleteSession(context.TODO(), &DELETESESSION.Request{
		Sessionid: userlogin.Value,
	})

	if err != nil {
		http.Error(w, err.Error(), 502)

		beego.Info(err)
		//beego.Debug(err)
		return
	}
	//再次读取数据
	cookie, err := r.Cookie("userlogin")

	//数据不为空则将数据设置副的
	if err != nil || "" == cookie.Value {
		return
	} else {
		cookie := http.Cookie{Name: "userlogin", Path: "/", MaxAge: -1}
		http.SetCookie(w, &cookie)
	}

	//返回数据
	resp := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}
	//设置格式
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), 503)
		beego.Info(err)
		return
	}

	return
}

//获取用户信息
func GetUserInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("GetUserInfo  获取用户信息   /api/v1.0/user")
	//初始化服务
	service := grpc.NewService()
	service.Init()

	//创建句柄
	exampleClient := GETUSERINFO.NewExampleService("go.micro.srv.GetUserInfo", service.Client())

	//获取用户的登陆信息
	userlogin, err := r.Cookie("userlogin")

	//判断是否成功不成功就直接返回
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		return
	}

	//成功就将信息发送给前端
	rsp, err := exampleClient.GetUserInfo(context.TODO(), &GETUSERINFO.Request{
		Sessionid: userlogin.Value,
	})

	if err != nil {
		http.Error(w, err.Error(), 502)

		beego.Info(err)
		//beego.Debug(err)
		return
	}
	//

	// 准备1个数据的map
	data := make(map[string]interface{})
	//将信息发送给前端
	data["user_id"] = int(rsp.UserId)
	data["name"] = rsp.Name
	data["mobile"] = rsp.Mobile
	data["real_name"] = rsp.RealName
	data["id_card"] = rsp.IdCard
	data["avatar_url"] = utils.AddDomain2Url(rsp.AvatarUrl)

	resp := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}
	//设置格式
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), 503)
		beego.Info(err)
		return
	}

	return
}

//上传用户头像 PostAvatar
func PostAvatar(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("上传用户头像 PostAvatar /api/v1.0/user/avatar")

	//创建服务
	service := grpc.NewService()
	service.Init()

	//创建句柄
	exampleClient := POSTAVATAR.NewExampleService("go.micro.srv.PostAvatar", service.Client())

	//查看登陆信息
	userlogin, err := r.Cookie("userlogin")

	//如果没有登陆就返回错误
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		return
	}

	//接收前端发送过来的文集
	file, hander, err := r.FormFile("avatar")

	//判断是否接受成功
	if err != nil {
		beego.Info("Postupavatar   c.GetFile(avatar) err", err)

		resp := map[string]interface{}{
			"errno":  utils.RECODE_IOERR,
			"errmsg": utils.RecodeText(utils.RECODE_IOERR),
		}
		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		return
	}
	//打印基本信息
	//beego.Info(file ,hander)
	beego.Info("文件大小", hander.Size)
	beego.Info("文件名", hander.Filename)

	//二进制的空间用来存储文件
	filebuffer := make([]byte, hander.Size)

	//将文件读取到filebuffer里
	_, err = file.Read(filebuffer)
	if err != nil {
		beego.Info("Postupavatar   file.Read(filebuffer) err", err)
		resp := map[string]interface{}{
			"errno":  utils.RECODE_IOERR,
			"errmsg": utils.RecodeText(utils.RECODE_IOERR),
		}
		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		return
	}
	//调用函数传入数据
	rsp, err := exampleClient.PostAvatar(context.TODO(), &POSTAVATAR.Request{
		Sessionid: userlogin.Value,
		Filename:  hander.Filename,
		Filesize:  hander.Size,
		Avatar:    filebuffer,
	})
	if err != nil {
		http.Error(w, err.Error(), 502)

		beego.Info(err)
		//beego.Debug(err)
		return
	}
	//

	//准备回传数据空间
	data := make(map[string]interface{})
	//url拼接然回回传数据
	data["avatar_url"] = utils.AddDomain2Url(rsp.AvatarUrl)

	resp := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}

	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), 503)
		beego.Info(err)
		return
	}

	return
}

//更新用户名//PutUserInfo
func PutUserInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info(" 更新用户名 Putuserinfo /api/v1.0/user/name")
	//创建服务
	service := grpc.NewService()
	service.Init()
	// 接收前端发送内容
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 调用服务
	exampleClient := PUTUSERINFO.NewExampleService("go.micro.srv.PutUserInfo", service.Client())

	//获取用户登陆信息
	userlogin, err := r.Cookie("userlogin")
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		return
	}

	rsp, err := exampleClient.PutUserInfo(context.TODO(), &PUTUSERINFO.Request{
		Sessionid: userlogin.Value,
		Username:  request["name"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//接收回发数据
	data := make(map[string]interface{})
	data["name"] = rsp.Username

	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}
	w.Header().Set("Content-Type", "application/json")

	// 返回前端
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

//检查实名认证

func GetUserAuth(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("GetUserInfo  获取用户信息   /api/v1.0/user")
	//初始化服务
	service := grpc.NewService()
	service.Init()

	//创建句柄
	exampleClient := GETUSERINFO.NewExampleService("go.micro.srv.GetUserInfo", service.Client())

	//获取用户的登陆信息
	userlogin, err := r.Cookie("userlogin")

	//判断是否成功不成功就直接返回
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		return
	}

	//成功就将信息发送给前端
	rsp, err := exampleClient.GetUserInfo(context.TODO(), &GETUSERINFO.Request{
		Sessionid: userlogin.Value,
	})

	if err != nil {
		http.Error(w, err.Error(), 502)

		beego.Info(err)
		//beego.Debug(err)
		return
	}
	//

	// 准备1个数据的map
	data := make(map[string]interface{})
	//将信息发送给前端
	data["user_id"] = int(rsp.UserId)
	data["name"] = rsp.Name
	data["mobile"] = rsp.Mobile
	data["real_name"] = rsp.RealName
	data["id_card"] = rsp.IdCard
	data["avatar_url"] = utils.AddDomain2Url(rsp.AvatarUrl)

	resp := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}
	//设置格式
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), 503)
		beego.Info(err)
		return
	}

	return
}

//实名认证  PostUserAuth
func PostUserAuth(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info(" 实名认证 Postuserauth  api/v1.0/user/auth ")

	service := grpc.NewService()
	service.Init()

	//获取前端发送的数据
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// call the backend service
	exampleClient := POSTUSERAUTH.NewExampleService("go.micro.srv.PostUserAuth", service.Client())

	//获取cookie
	userlogin, err := r.Cookie("userlogin")
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		return
	}

	rsp, err := exampleClient.PostUserAuth(context.TODO(), &POSTUSERAUTH.Request{
		Sessionid: userlogin.Value,
		RealName:  request["real_name"].(string),
		IdCard:    request["id_card"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

// 获取当前用户所发布的房源 GetUserHouses
func GetUserHouses(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	beego.Info("获取当前用户所发布的房源 GetUserHouses /api/v1.0/user/houses")

	server := grpc.NewService()
	server.Init()

	// call the backend service
	exampleClient := GETUSERHOUSES.NewExampleService("go.micro.srv.GetUserHouses", server.Client())

	//获取cookie
	userlogin, err := r.Cookie("userlogin")
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		return
	}

	rsp, err := exampleClient.GetUserHouses(context.TODO(), &GETUSERHOUSES.Request{
		Sessionid: userlogin.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	house_list := []models.House{}
	json.Unmarshal(rsp.Mix, &house_list)

	var houses []interface{}
	for _, houseinfo := range house_list {
		fmt.Printf("house.user = %+v\n", houseinfo.Id)
		fmt.Printf("house.area = %+v\n", houseinfo.Area)
		houses = append(houses, houseinfo.To_house_info())
	}

	data_map := make(map[string]interface{})
	data_map["houses"] = houses

	// we want to augment the response
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data_map,
	}
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

//发布房源信息
func PostHouses(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("PostHouses 发布房源信息 /api/v1.0/houses ")
	//获取前端post请求发送的内容
	body, _ := ioutil.ReadAll(r.Body)

	//获取cookie
	userlogin, err := r.Cookie("userlogin")
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}
		//设置回传格式
		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		return
	}

	//创建连接

	service := grpc.NewService()
	service.Init()
	exampleClient := POSTHOUSES.NewExampleService("go.micro.srv.PostHouses", service.Client())

	rsp, err := exampleClient.PostHouses(context.TODO(), &POSTHOUSES.Request{
		Sessionid: userlogin.Value,
		Max:       body,
	})
	if err != nil {
		http.Error(w, err.Error(), 502)

		beego.Info(err)
		//beego.Debug(err)
		return
	}

	/*得到插入房源信息表的 id*/
	houseid_map := make(map[string]interface{})
	houseid_map["house_id"] = int(rsp.House_Id)

	resp := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   houseid_map,
	}
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), 503)
		beego.Info(err)
		return
	}
}

//发送房屋图片
func PostHousesImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	beego.Info("发送房屋图片PostHousesImage  /api/v1.0/houses/:id/images")

	//创建服务
	server := grpc.NewService()
	server.Init()

	// call the backend service
	exampleClient := POSTHOUSESIMAGE.NewExampleService("go.micro.srv.PostHousesImage", server.Client())
	//获取houserid
	houseid := ps.ByName("id")
	//获取sessionid
	userlogin, err := r.Cookie("userlogin")
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		return
	}

	file, hander, err := r.FormFile("house_image")
	if err != nil {
		beego.Info("Postupavatar   c.GetFile(avatar) err", err)

		resp := map[string]interface{}{
			"errno":  utils.RECODE_IOERR,
			"errmsg": utils.RecodeText(utils.RECODE_IOERR),
		}
		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		return
	}

	//beego.Info(file ,hander)
	beego.Info("文件大小", hander.Size)
	beego.Info("文件名", hander.Filename)
	//二进制的空间用来存储文件
	filebuffer := make([]byte, hander.Size)
	//将文件读取到filebuffer里
	_, err = file.Read(filebuffer)
	if err != nil {
		beego.Info("Postupavatar   file.Read(filebuffer) err", err)
		resp := map[string]interface{}{
			"errno":  utils.RECODE_IOERR,
			"errmsg": utils.RecodeText(utils.RECODE_IOERR),
		}
		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		return
	}

	rsp, err := exampleClient.PostHousesImage(context.TODO(), &POSTHOUSESIMAGE.Request{
		Sessionid: userlogin.Value,
		Id:        houseid,
		Image:     filebuffer,
		Filesize:  hander.Size,
		Filename:  hander.Filename,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	//准备返回值
	data := make(map[string]interface{})
	data["url"] = utils.AddDomain2Url(rsp.Url)
	// 返回数据map
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}
	w.Header().Set("Content-Type", "application/json")

	// 回发数据
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

//获取房源详细信息
func GetHouseInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	beego.Info("获取房源详细信息 GetHouseInfo  api/v1.0/houses/:id ")

	//创建服务
	server := grpc.NewService()
	server.Init()

	// call the backend service
	exampleClient := GETHOUSEINFO.NewExampleService("go.micro.srv.GetHouseInfo", server.Client())

	id := ps.ByName("id")

	//获取sessionid
	userlogin, err := r.Cookie("userlogin")
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		return
	}

	rsp, err := exampleClient.GetHouseInfo(context.TODO(), &GETHOUSEINFO.Request{
		Sessionid: userlogin.Value,
		Id:        id,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	house := models.House{}
	json.Unmarshal(rsp.Housedata, &house)

	data_map := make(map[string]interface{})
	data_map["user_id"] = int(rsp.Userid)
	data_map["house"] = house.To_one_house_desc()

	// we want to augment the response
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data_map,
	}
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
	return
}

//获取首页轮播
func GetIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("获取首页轮播 url：api/v1.0/houses/index")
	server := grpc.NewService()
	server.Init()

	exampleClient := GETINDEX.NewExampleService("go.micro.srv.GetIndex", server.Client())

	rsp, err := exampleClient.GetIndex(context.TODO(), &GETINDEX.Request{})
	if err != nil {
		beego.Info(err)
		http.Error(w, err.Error(), 502)
		return
	}
	data := []interface{}{}
	json.Unmarshal(rsp.Max, &data)

	//创建返回数据map
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}
	w.Header().Set("Content-Type", "application/json")

	// 将返回数据map发送给前端
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
}

//搜索房屋
func GetHouses(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	server := grpc.NewService()
	server.Init()

	// call the backend service
	exampleClient := GETHOUSES.NewExampleService("go.micro.srv.GetHouses", server.Client())

	//aid=5&sd=2017-11-12&ed=2017-11-30&sk=new&p=1
	aid := r.URL.Query()["aid"][0] //aid=5   地区编号
	sd := r.URL.Query()["sd"][0]   //sd=2017-11-1   开始世界
	ed := r.URL.Query()["ed"][0]   //ed=2017-11-3   结束世界
	sk := r.URL.Query()["sk"][0]   //sk=new    第三栏条件
	p := r.URL.Query()["p"][0]     //tp=1   页数

	rsp, err := exampleClient.GetHouses(context.TODO(), &GETHOUSES.Request{
		Aid: aid,
		Sd:  sd,
		Ed:  ed,
		Sk:  sk,
		P:   p,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	houses_l := []interface{}{}
	json.Unmarshal(rsp.Houses, &houses_l)

	data := map[string]interface{}{}
	data["current_page"] = rsp.CurrentPage
	data["houses"] = houses_l
	data["total_page"] = rsp.TotalPage

	// we want to augment the response
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

//发布订单
func PostOrders(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("PostOrders  发布订单 /api/v1.0/orders")

	//将post代过来的数据转化以下
	body, _ := ioutil.ReadAll(r.Body)

	userlogin, err := r.Cookie("userlogin")
	if err != nil || userlogin.Value == "" {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		return
	}

	service := grpc.NewService()
	service.Init()

	//调用服务
	exampleClient := POSTORDERS.NewExampleService("go.micro.srv.PostOrders", service.Client())
	rsp, err := exampleClient.PostOrders(context.TODO(), &POSTORDERS.Request{
		Sessionid: userlogin.Value,
		Body:      body,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	/*得到插入房源信息表的 id*/
	houseid_map := make(map[string]interface{})
	houseid_map["order_id"] = int(rsp.OrderId)

	// we want to augment the response
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   houseid_map,
	}
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

//获取订单

func GetUserOrder(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	beego.Info("/api/v1.0/user/orders   GetUserOrder 获取订单 ")
	server := grpc.NewService()
	server.Init()
	// call the backend service
	exampleClient := GETUSERORDER.NewExampleService("go.micro.srv.GetUserOrder", server.Client())

	//获取cookie
	userlogin, err := r.Cookie("userlogin")
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		return
	}
	//获取role
	role := r.URL.Query()["role"][0] //role

	rsp, err := exampleClient.GetUserOrder(context.TODO(), &GETUSERORDER.Request{
		Sessionid: userlogin.Value,
		Role:      role,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	order_list := []interface{}{}
	json.Unmarshal(rsp.Orders, &order_list)

	data := map[string]interface{}{}
	data["orders"] = order_list

	// we want to augment the response
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}

	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

//房东同意/拒绝订单
func PutOrders(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// decode the incoming request as json
	//接收请求携带的数据
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
	//获取cookie
	userlogin, err := r.Cookie("userlogin")
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 502)
			beego.Info(err)
			return
		}
		return
	}
	server := grpc.NewService()
	server.Init()

	// call the backend service
	exampleClient := PUTORDERS.NewExampleService("go.micro.srv.PutOrders", server.Client())

	rsp, err := exampleClient.PutOrders(context.TODO(), &PUTORDERS.Request{
		Sessionid: userlogin.Value,
		Action:    request["action"].(string),
		Orderid:   ps.ByName("id"),
	})
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 504)
		return
	}
}

//用户评价订单
func PutComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	beego.Info("PutComment  用户评价 /api/v1.0/orders/:id/comment")
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	service := grpc.NewService()
	service.Init()
	// call the backend service
	exampleClient := PUTCOMMENT.NewExampleService("go.micro.srv.PutComment", service.Client())

	//获取cookie
	userlogin, err := r.Cookie("userlogin")
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		return
	}

	rsp, err := exampleClient.PutComment(context.TODO(), &PUTCOMMENT.Request{

		Sessionid: userlogin.Value,
		Comment:   request["comment"].(string),
		OrderId:   ps.ByName("id"),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}
