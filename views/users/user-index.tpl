<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>{{config "String" "globaltitle" ""}}</title>
{{template "inc/meta.tpl" .}}
<link href="/static/css/table-responsive.css" rel="stylesheet">
</head><body class="sticky-header">
<section> {{template "inc/left.tpl" .}}
  <!-- main content start-->
  <div class="main-content" >
    <!-- header section start-->
    <div class="header-section">
      <!--toggle button start-->
      <a class="toggle-btn"><i class="fa fa-bars"></i></a>
      <!--toggle button end-->
      <!--search start-->
      <form class="searchform" action="/user/manage" method="get">
        <select name="status" class="form-control">
          <option value="">用户状态</option>
          <option value="1" {{if eq "1" .condArr.status}}selected{{end}}>正常</option>
          <option value="2" {{if eq "2" .condArr.status}}selected{{end}}>屏蔽</option>
        </select>
        <input type="text" class="form-control" name="keywords" placeholder="请输入用户名、姓名" value="{{.condArr.keywords}}"/>
        <button type="submit" class="btn btn-primary">搜索</button>
      </form>
      <!--search end-->
      {{template "inc/user-info.tpl" .}} </div>
    <!-- header section end-->
    <!-- page heading start-->
    <div class="page-heading">
      <h3> 组织管理 {{template "users/nav.tpl" .}}</h3>
      <ul class="breadcrumb pull-left">
        <li> <a href="/user/show/{{.LoginUserid}}">OA</a> </li>
        <li> <a href="/user/manage">员工管理</a> </li>
        <li class="active"> 员工 </li>
      </ul>
      <div class="pull-right"><a href="/user/add" class="btn btn-success">+添加新员工</a></div>
    </div>
    <!-- page heading end-->
    <!--body wrapper start-->
    <div class="wrapper">
      <div class="row">
        <div class="col-sm-12">
          <section class="panel">
            <header class="panel-heading"> 员工管理 / 总数：{{.countUser}}<span class="tools pull-right"><a href="javascript:;" class="fa fa-chevron-down"></a>
              <!--a href="javascript:;" class="fa fa-times"></a-->
              </span> </header>
            <div class="panel-body">
              <section id="unseen">
                <form id="user-form-list">
                  <table class="table table-bordered table-striped table-condensed">
                    <thead>
                      <tr>
                        <th>用户名</th>
                        <th>姓名</th>
                        <th>性别</th>
                        <th>手机号</th>
                        <th>紧急电话</th>
                        <th>上次登录</th>
                        <th>状态</th>
                        <th>操作</th>
                      </tr>
                    </thead>
                    <tbody>
                    
                    {{range $k,$v := .user}}
                    <tr>
                      <td>{{$v.Username}}</td>
                      <td><a href="/user/show/{{$v.Id}}">{{$v.Profile.Realname}}</a></td>
                      <td>{{if eq 1 $v.Profile.Sex}}男{{else}}女{{end}}</td>
                      <td>{{$v.Profile.Phone}}</td>
                      <td>{{$v.Profile.Emerphone}}</td>
                      <td>{{getDate $v.Profile.Lasted}}</td>
                      <td>{{if eq 1 $v.Status}}正常{{else}}屏蔽{{end}}</td>
                      <td><div class="btn-group">
                          <button type="button" class="btn btn-primary dropdown-toggle" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false"> 操作<span class="caret"></span> </button>
                          <ul class="dropdown-menu">
                            <li><a href="/user/edit/{{$v.Id}}">编辑</a></li>
                            <!--li role="separator" class="divider"></li>
							<li><a href="/user/permission/{{$v.Id}}">权限</a></li-->
                            <li role="separator" class="divider"></li>
                            {{if eq 1 $v.Status}}
                            <li><a href="javascript:;" class="js-user-single" data-id="{{$v.Id}}" data-status="2">屏蔽</a></li>
                            {{else}}
                            <li><a href="javascript:;" class="js-user-single" data-id="{{$v.Id}}" data-status="1">正常</a></li>
                            {{end}}
							<li role="separator" class="divider"></li>
							<li><a href="/checkwork/all?userid={{$v.Id}}">考勤</a></li>
                          </ul>
                        </div></td>
                    </tr>
                    {{end}}
                    </tbody>
                    
                  </table>
                </form>
                {{template "inc/page.tpl" .}}
				 </section>
            </div>
          </section>
        </div>
      </div>
    </div>
    <!--body wrapper end-->
    <!--footer section start-->
    {{template "inc/foot-info.tpl" .}}
    <!--footer section end-->
  </div>
  <!-- main content end-->
</section>
{{template "inc/foot.tpl" .}}
</body>
</html>
