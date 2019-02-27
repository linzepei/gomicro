package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/weilaihui/fdfs_client"
)

/* 将url加上 http://IP:PROT/  前缀 */
//http:// + 127.0.0.1 + ：+ 8080 + 请求

func AddDomain2Url(url string) (domain_url string) {
	domain_url = "http://" + G_fastdfs_addr + ":" + G_fastdfs_port + "/" + url

	return domain_url
}

func Md5String(s string) string {
	//创建1个md5对象
	h := md5.New()
	h.Write([]byte(s))

	return hex.EncodeToString(h.Sum(nil))
}

//通过文件名的方式进行上传
func UploadByFilename(filename string) (GroupName, RemoteFileId string, err error) {
	//通过配置文件创建fdfs操作句柄
	fdfsClient, thiserr := fdfs_client.NewFdfsClient("/root/mygo/src/gomicro/IhomeWeb/conf/client.conf")
	if thiserr != nil {
		//说一下那里出问题了
		beego.Info("UploadByFilename( ) fdfs_client.NewFdfsClient  err", err)
		GroupName = ""
		RemoteFileId = ""
		err = thiserr
		return
	}

	//unc (this *FdfsClient) UploadByFilename(filename string) (*UploadFileResponse, error)
	//通过句柄上传文件（被上传的文件）
	uploadResponse, thiserr := fdfsClient.UploadByFilename(filename)
	if thiserr != nil {
		beego.Info("UploadByFilename( ) fdfsClient.UploadByFilename(filename)err", err)
		GroupName = ""
		RemoteFileId = ""
		err = thiserr
		return
	}

	beego.Info(uploadResponse.GroupName)
	beego.Info(uploadResponse.RemoteFileId)
	//回传
	return uploadResponse.GroupName, uploadResponse.RemoteFileId, nil

}

//上传二进制文件到fdfs中的操作
func UploadByBuffer(filebuffer []byte, fileExt string) (fileid string, err error) {
	fd_cilent, err := fdfs_client.NewFdfsClient("/home/itcast/go/src/sss/IhomeWeb/conf/client.conf")
	if err != nil {
		fmt.Println("创建句柄失败", err)
		fileid = ""
		return
	}

	fd_rsq, err := fd_cilent.UploadByBuffer(filebuffer, fileExt)
	if err != nil {
		fmt.Println("上传失败", err)
		fileid = ""
		return
	}

	fmt.Println(fd_rsq.GroupName)
	fmt.Println(fd_rsq.RemoteFileId)

	fileid = fd_rsq.RemoteFileId

	return fileid, nil
}
