package albums

import (
	"fmt"
	"github.com/1975210542/OPMS/controllers"
	"github.com/1975210542/OPMS/models/albums"
	"github.com/1975210542/OPMS/models/messages"
	"github.com/1975210542/OPMS/utils"
)

type AddCommentController struct {
	controllers.BaseController
}

func (this *AddCommentController) Post() {
	albumid, _ := this.GetInt64("albumid")
	if albumid <= 0 {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "参数出错"}
		this.ServeJSON()
		return
	}
	content := this.GetString("comment")
	if "" == content {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请填写评论内容"}
		this.ServeJSON()
		return
	}

	var err error
	var comment albums.AlbumsComment
	comment.Id = utils.SnowFlakeId()
	comment.Userid = this.BaseController.UserUserId
	comment.Albumid = albumid
	comment.Content = content

	err = albums.AddAlbumComment(comment)

	if err == nil {
		//消息通知
		album, _ := albums.GetAlbum(albumid)
		var msg messages.Messages
		msg.Id = utils.SnowFlakeId()
		msg.Userid = this.BaseController.UserUserId
		msg.Touserid = album.Userid
		msg.Type = 1
		msg.Subtype = 12
		msg.Title = album.Title
		msg.Url = "/album/" + fmt.Sprintf("%d", albumid)
		messages.AddMessages(msg)
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "评价添加成功"}
	} else {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "添加失败"}
	}
	this.ServeJSON()
}

type AjaxLaudController struct {
	controllers.BaseController
}

func (this *AjaxLaudController) Post() {
	albumid, _ := this.GetInt64("albumid")
	if albumid <= 0 {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "参数出错"}
		this.ServeJSON()
		return
	}

	laudexit, _ := albums.GetAlbumLaud(albumid)
	if laudexit.Userid == this.BaseController.UserUserId {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "已经点过赞"}
		this.ServeJSON()
	}

	var err error
	var laud albums.AlbumsLaud
	laud.Id = utils.SnowFlakeId()
	laud.Userid = this.BaseController.UserUserId
	laud.Albumid = albumid

	err = albums.AddAlbumLaud(laud)

	if err == nil {
		//消息通知
		album, _ := albums.GetAlbum(albumid)
		var msg messages.Messages
		msg.Id = utils.SnowFlakeId()
		msg.Userid = this.BaseController.UserUserId
		msg.Touserid = album.Userid
		msg.Type = 2
		msg.Subtype = 22
		msg.Title = album.Title
		msg.Url = "/album/" + fmt.Sprintf("%d", albumid)
		messages.AddMessages(msg)
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "点赞成功"}
	} else {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "点赞失败"}
	}
	this.ServeJSON()
}
