package projects

import (
	"fmt"
	"github.com/1975210542/OPMS/controllers"
	//. "opms/models/projects"
	"github.com/1975210542/OPMS/models/projects"
	"github.com/1975210542/OPMS/utils"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/utils/pagination"
)

type MyProjectController struct {
	controllers.BaseController
}

func (this *MyProjectController) Get() {
	userid := this.BaseController.UserUserId
	_, _, projects := projects.ListMyProject(userid, 1, 100)
	this.Data["projects"] = projects
	this.Data["countProject"] = len(projects)

	this.TplName = "projects/myproject.tpl"
}

type ManageProjectController struct {
	controllers.BaseController
}

func (this *ManageProjectController) Get() {
	//权限检测
	if !strings.Contains(this.GetSession("userPermission").(string), "project-manage") {
		this.Redirect("/my/task", 302)
		return
		//this.Abort("401")
	}
	page, err := this.GetInt("p")
	status := this.GetString("status")
	keywords := this.GetString("keywords")
	if err != nil {
		page = 1
	}

	offset, err1 := beego.AppConfig.Int("pageoffset")
	if err1 != nil {
		offset = 15
	}

	condArr := make(map[string]string)
	condArr["status"] = status
	condArr["keywords"] = keywords

	countProject := projects.CountProject(condArr)
	paginator := pagination.SetPaginator(this.Ctx, offset, countProject)
	_, _, projects := projects.ListProject(condArr, page, offset)

	this.Data["paginator"] = paginator
	this.Data["condArr"] = condArr
	this.Data["projects"] = projects
	this.Data["countProject"] = countProject

	this.TplName = "projects/project.tpl"
}

type AjaxStatusProjectController struct {
	controllers.BaseController
}

func (this *AjaxStatusProjectController) Post() {
	//权限检测
	if !strings.Contains(this.GetSession("userPermission").(string), "project-edit") {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "无权设置"}
		this.ServeJSON()
		return
	}
	id, _ := this.GetInt64("id")
	if id <= 0 {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请选择项目"}
		this.ServeJSON()
		return
	}
	status, _ := this.GetInt("status")
	if status <= 0 || status >= 5 {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请选择操作状态"}
		this.ServeJSON()
		return
	}

	err := projects.ChangeProjectStatus(id, status)

	if err == nil {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "项目状态更改成功"}
	} else {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "项目状态更改失败"}
	}
	this.ServeJSON()
}

type AddProjectController struct {
	controllers.BaseController
}

func (this *AddProjectController) Get() {
	//权限检测
	if !strings.Contains(this.GetSession("userPermission").(string), "project-add") {
		this.Abort("401")
	}
	var project projects.Projects
	project.Started = time.Now().Unix()
	project.Ended = time.Now().Unix()
	this.Data["project"] = project
	this.TplName = "projects/project-form.tpl"
}

func (this *AddProjectController) Post() {
	//权限检测
	if !strings.Contains(this.GetSession("userPermission").(string), "project-add") {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "无权设置"}
		this.ServeJSON()
		return
	}
	name := this.GetString("name")
	if "" == name {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请填写项目名称"}
		this.ServeJSON()
		return
	}
	aliasname := this.GetString("aliasname")
	if "" == aliasname {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请填写项目别名"}
		this.ServeJSON()
		return
	}
	startedstr := this.GetString("started")
	if "" == startedstr {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请填写开始时间"}
		this.ServeJSON()
		return
	}
	startedtime := utils.GetDateParse(startedstr)

	endedstr := this.GetString("ended")
	if "" == endedstr {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请填写结束时间"}
		this.ServeJSON()
		return
	}
	endedtime := utils.GetDateParse(endedstr)

	desc := this.GetString("desc")
	if "" == desc {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请填写项目描述"}
		this.ServeJSON()
		return
	}

	userid := this.BaseController.UserUserId

	var err error
	//雪花算法ID生成
	id := utils.SnowFlakeId()

	var pro projects.Projects
	pro.Id = id
	pro.Userid = userid
	pro.Name = name
	pro.Aliasname = aliasname
	pro.Started = startedtime
	pro.Ended = endedtime
	pro.Desc = desc

	err = projects.AddProject(pro)

	if err == nil {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "项目信息添加成功", "id": fmt.Sprintf("%d", id)}
	} else {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "项目信息添加失败"}
	}
	this.ServeJSON()
}

type EditProjectController struct {
	controllers.BaseController
}

func (this *EditProjectController) Get() {
	//权限检测
	if !strings.Contains(this.GetSession("userPermission").(string), "project-edit") {
		this.Abort("401")
	}
	idstr := this.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idstr)
	project, err := projects.GetProject(int64(id))
	if err != nil {
		this.Redirect("/404.html", 302)
	}
	_, _, teams := projects.ListProjectTeam(project.Id, 1, 100)
	this.Data["teams"] = teams
	this.Data["project"] = project
	this.TplName = "projects/project-form.tpl"
}

func (this *EditProjectController) Post() {
	//权限检测
	if !strings.Contains(this.GetSession("userPermission").(string), "project-edit") {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "无权设置"}
		this.ServeJSON()
		return
	}
	id, _ := this.GetInt64("id")
	if id <= 0 {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "参数出错"}
		this.ServeJSON()
		return
	}
	_, err := projects.GetProject(id)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "项目不存在"}
		this.ServeJSON()
		return
	}
	name := this.GetString("name")
	if "" == name {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请填写项目名称"}
		this.ServeJSON()
		return
	}
	aliasname := this.GetString("aliasname")
	if "" == aliasname {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请填写项目别名"}
		this.ServeJSON()
		return
	}
	startedstr := this.GetString("started")
	if "" == startedstr {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请填写开始时间"}
		this.ServeJSON()
		return
	}
	startedtime := utils.GetDateParse(startedstr)

	endedstr := this.GetString("ended")
	if "" == endedstr {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请填写结束时间"}
		this.ServeJSON()
		return
	}
	endedtime := utils.GetDateParse(endedstr)

	desc := this.GetString("desc")
	if "" == desc {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请填写项目描述"}
		this.ServeJSON()
		return
	}
	projuserid, _ := this.GetInt64("projuserid")
	produserid, _ := this.GetInt64("produserid")
	testuserid, _ := this.GetInt64("testuserid")
	publuserid, _ := this.GetInt64("publuserid")

	var pro projects.Projects
	pro.Name = name
	pro.Aliasname = aliasname
	pro.Started = startedtime
	pro.Ended = endedtime
	pro.Desc = desc
	pro.Projuserid = projuserid
	pro.Produserid = produserid
	pro.Testuserid = testuserid
	pro.Publuserid = publuserid

	err = projects.UpdateProject(id, pro)

	if err == nil {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "项目修改成功"}
	} else {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "项目修改失败"}
	}
	this.ServeJSON()
}

//项目详情
type ShowProjectController struct {
	controllers.BaseController
}

func (this *ShowProjectController) Get() {
	idstr := this.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idstr)
	project, err := projects.GetProject(int64(id))
	if err != nil {
		this.Abort("404")
	}
	if !strings.Contains(this.GetSession("userPermission").(string), "project-manage") {
		this.Data["url"] = "/my/project"
	} else {
		this.Data["url"] = "/project/manage"
	}

	this.Data["project"] = project
	this.TplName = "projects/project-detail.tpl"
}

//项目统计
type ChartProjectController struct {
	controllers.BaseController
}

func (this *ChartProjectController) Get() {
	//权限检测
	if !strings.Contains(this.GetSession("userPermission").(string), "project-manage") {
		this.Abort("401")
	}
	idstr := this.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idstr)
	longid := int64(id)
	project, err := projects.GetProject(longid)
	if err != nil {
		this.Redirect("/404.html", 302)
	}

	//团队
	chartTeamsNum, _, chartTeams := projects.ChartProjectTeam(longid)
	this.Data["chartTeams"] = chartTeams
	this.Data["chartTeamsNum"] = chartTeamsNum - 1

	//需求
	chartNeedsAcceptNum, _, chartNeedsAccept := projects.ChartProjectNeed("accept", longid)
	this.Data["chartNeedsAccept"] = chartNeedsAccept
	this.Data["chartNeedsAcceptNum"] = chartNeedsAcceptNum - 1

	chartNeedsUserNum, _, chartNeedsUser :=projects.ChartProjectNeed("user", longid)
	this.Data["chartNeedsUser"] = chartNeedsUser
	this.Data["chartNeedsUserNum"] = chartNeedsUserNum - 1

	chartNeedsSourceNum, _, chartNeedsSource := projects.ChartProjectNeedSource(longid)
	this.Data["chartNeedsSource"] = chartNeedsSource
	this.Data["chartNeedsSourceNum"] = chartNeedsSourceNum - 1

	//任务
	chartTasksAcceptNum, _, chartTasksAccept := projects.ChartProjectTask("accept", longid)
	this.Data["chartTasksAccept"] = chartTasksAccept
	this.Data["chartTasksAcceptNum"] = chartTasksAcceptNum - 1

	chartTasksUserNum, _, chartTasksUser := projects.ChartProjectTask("user", longid)
	this.Data["chartTasksUser"] = chartTasksUser
	this.Data["chartTasksUserNum"] = chartTasksUserNum - 1

	chartTasksCompleteNum, _, chartTasksComplete := projects.ChartProjectTask("complete", longid)
	this.Data["chartTasksComplete"] = chartTasksComplete
	this.Data["chartTasksCompleteNum"] = chartTasksCompleteNum - 1

	chartTasksSourceNum, _, chartTasksSource := projects.ChartProjectTaskSource(longid)
	this.Data["chartTasksSource"] = chartTasksSource
	this.Data["chartTasksSourceNum"] = chartTasksSourceNum - 1

	//Bug
	chartTestsAcceptNum, _, chartTestsAccept := projects.ChartProjectTest("accept", longid)
	this.Data["chartTestsAccept"] = chartTestsAccept
	this.Data["chartTestsAcceptNum"] = chartTestsAcceptNum - 1

	chartTestsUserNum, _, chartTestsUser := projects.ChartProjectTest("user", longid)
	this.Data["chartTestsUser"] = chartTestsUser
	this.Data["chartTestsUserNum"] = chartTestsUserNum - 1

	chartTestsCompleteNum, _, chartTestsComplete := projects.ChartProjectTest("complete", longid)
	this.Data["chartTestsComplete"] = chartTestsComplete
	this.Data["chartTestsCompleteNum"] = chartTestsCompleteNum - 1

	this.Data["project"] = project
	this.TplName = "projects/project-chart.tpl"
}
