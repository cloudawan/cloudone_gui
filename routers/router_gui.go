// Copyright 2015 CloudAwan LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package routers

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers"
	"github.com/cloudawan/cloudone_gui/controllers/dashboard/appservice"
	"github.com/cloudawan/cloudone_gui/controllers/dashboard/bluegreen"
	dashboarddeploy "github.com/cloudawan/cloudone_gui/controllers/dashboard/deploy"
	"github.com/cloudawan/cloudone_gui/controllers/dashboard/healthcheck"
	"github.com/cloudawan/cloudone_gui/controllers/dashboard/topology"
	"github.com/cloudawan/cloudone_gui/controllers/deploy/autoscaler"
	"github.com/cloudawan/cloudone_gui/controllers/deploy/clone"
	"github.com/cloudawan/cloudone_gui/controllers/deploy/deploy"
	"github.com/cloudawan/cloudone_gui/controllers/deploy/deploybluegreen"
	"github.com/cloudawan/cloudone_gui/controllers/deploy/deployclusterapplication"
	"github.com/cloudawan/cloudone_gui/controllers/event/audit"
	"github.com/cloudawan/cloudone_gui/controllers/event/kubernetes"
	"github.com/cloudawan/cloudone_gui/controllers/filesystem/glusterfs/cluster"
	"github.com/cloudawan/cloudone_gui/controllers/filesystem/glusterfs/volume"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/inventory/replicationcontroller"
	"github.com/cloudawan/cloudone_gui/controllers/inventory/service"
	"github.com/cloudawan/cloudone_gui/controllers/monitor/container"
	"github.com/cloudawan/cloudone_gui/controllers/monitor/historicalcontainer"
	"github.com/cloudawan/cloudone_gui/controllers/monitor/node"
	"github.com/cloudawan/cloudone_gui/controllers/notification/notifier"
	"github.com/cloudawan/cloudone_gui/controllers/repository/imageinformation"
	"github.com/cloudawan/cloudone_gui/controllers/repository/imagerecord"
	"github.com/cloudawan/cloudone_gui/controllers/repository/thirdparty"
	"github.com/cloudawan/cloudone_gui/controllers/repository/topologytemplate"
	"github.com/cloudawan/cloudone_gui/controllers/system/about"
	"github.com/cloudawan/cloudone_gui/controllers/system/host/credential"
	"github.com/cloudawan/cloudone_gui/controllers/system/namespace"
	"github.com/cloudawan/cloudone_gui/controllers/system/notification/emailserver"
	"github.com/cloudawan/cloudone_gui/controllers/system/notification/sms"
	"github.com/cloudawan/cloudone_gui/controllers/system/rbac/role"
	"github.com/cloudawan/cloudone_gui/controllers/system/rbac/user"
	"github.com/cloudawan/cloudone_gui/controllers/system/upgrade"
)

func init() {

	beego.Router("/", &controllers.MainController{})
	beego.Router("/gui/login", &identity.LoginController{})
	beego.Router("/gui/logout", &identity.LogoutController{})

	beego.Router("/gui/dashboard/topology", &topology.IndexController{})
	beego.Router("/gui/dashboard/topology/data", &topology.DataController{})
	beego.Router("/gui/dashboard/healthcheck/list", &healthcheck.ListController{})
	beego.Router("/gui/dashboard/bluegreen", &bluegreen.IndexController{})
	beego.Router("/gui/dashboard/bluegreen/data", &bluegreen.DataController{})
	beego.Router("/gui/dashboard/appservice", &appservice.IndexController{})
	beego.Router("/gui/dashboard/appservice/data", &appservice.DataController{})
	beego.Router("/gui/dashboard/deploy", &dashboarddeploy.IndexController{})
	beego.Router("/gui/dashboard/deploy/data", &dashboarddeploy.DataController{})
	beego.Router("/gui/repository/imageinformation/list", &imageinformation.ListController{})
	beego.Router("/gui/repository/imageinformation/create", &imageinformation.CreateController{})
	beego.Router("/gui/repository/imageinformation/upgrade", &imageinformation.UpgradeController{})
	beego.Router("/gui/repository/imageinformation/delete", &imageinformation.DeleteController{})
	beego.Router("/gui/repository/imageinformation/log", &imageinformation.LogController{})
	beego.Router("/gui/repository/imageinformation/websocket", &imageinformation.WebSocketController{})
	beego.Router("/gui/repository/imagerecord/list", &imagerecord.ListController{})
	beego.Router("/gui/repository/imagerecord/log", &imagerecord.LogController{})
	beego.Router("/gui/repository/imagerecord/delete", &imagerecord.DeleteController{})
	beego.Router("/gui/repository/thirdparty/list", &thirdparty.ListController{})
	beego.Router("/gui/repository/thirdparty/edit", &thirdparty.EditController{})
	beego.Router("/gui/repository/thirdparty/launch", &thirdparty.LaunchController{})
	beego.Router("/gui/repository/thirdparty/delete", &thirdparty.DeleteController{})
	beego.Router("/gui/repository/topologytemplate/list", &topologytemplate.ListController{})
	beego.Router("/gui/repository/topologytemplate/delete", &topologytemplate.DeleteController{})
	beego.Router("/gui/repository/topologytemplate/clone", &topologytemplate.CloneController{})
	beego.Router("/gui/deploy/deploy/list", &deploy.ListController{})
	beego.Router("/gui/deploy/deploy/create", &deploy.CreateController{})
	beego.Router("/gui/deploy/deploy/update", &deploy.UpdateController{})
	beego.Router("/gui/deploy/deploy/resize", &deploy.ResizeController{})
	beego.Router("/gui/deploy/deploy/delete", &deploy.DeleteController{})
	beego.Router("/gui/deploy/deploybluegreen/list", &deploybluegreen.ListController{})
	beego.Router("/gui/deploy/deploybluegreen/select", &deploybluegreen.SelectController{})
	beego.Router("/gui/deploy/deploybluegreen/delete", &deploybluegreen.DeleteController{})
	beego.Router("/gui/deploy/autoscaler/list", &autoscaler.ListController{})
	beego.Router("/gui/deploy/autoscaler/edit", &autoscaler.EditController{})
	beego.Router("/gui/deploy/autoscaler/delete", &autoscaler.DeleteController{})
	beego.Router("/gui/deploy/deployclusterapplication/list", &deployclusterapplication.ListController{})
	beego.Router("/gui/deploy/deployclusterapplication/size", &deployclusterapplication.SizeController{})
	beego.Router("/gui/deploy/deployclusterapplication/delete", &deployclusterapplication.DeleteController{})
	beego.Router("/gui/deploy/clone/select", &clone.SelectController{})
	beego.Router("/gui/deploy/clone/topology", &clone.TopologyController{})
	beego.Router("/gui/inventory/replicationcontroller/list", &replicationcontroller.ListController{})
	beego.Router("/gui/inventory/replicationcontroller/edit", &replicationcontroller.EditController{})
	beego.Router("/gui/inventory/replicationcontroller/size", &replicationcontroller.SizeController{})
	beego.Router("/gui/inventory/replicationcontroller/delete", &replicationcontroller.DeleteController{})
	beego.Router("/gui/inventory/replicationcontroller/pod/log", &replicationcontroller.PodLogController{})
	beego.Router("/gui/inventory/replicationcontroller/pod/delete", &replicationcontroller.PodDeleteController{})
	beego.Router("/gui/inventory/replicationcontroller/dockerterminal", &replicationcontroller.TerminalController{})
	beego.Router("/gui/inventory/replicationcontroller/dockerterminal/websocket", &replicationcontroller.WebSocketController{})
	beego.Router("/gui/inventory/service/list", &service.ListController{})
	beego.Router("/gui/inventory/service/edit", &service.EditController{})
	beego.Router("/gui/inventory/service/delete", &service.DeleteController{})
	beego.Router("/gui/filesystem/glusterfs/cluster/list", &cluster.ListController{})
	beego.Router("/gui/filesystem/glusterfs/cluster/edit", &cluster.EditController{})
	beego.Router("/gui/filesystem/glusterfs/cluster/delete", &cluster.DeleteController{})
	beego.Router("/gui/filesystem/glusterfs/volume/list", &volume.ListController{})
	beego.Router("/gui/filesystem/glusterfs/volume/create", &volume.CreateController{})
	beego.Router("/gui/filesystem/glusterfs/volume/reset", &volume.ResetController{})
	beego.Router("/gui/filesystem/glusterfs/volume/delete", &volume.DeleteController{})
	beego.Router("/gui/monitor/node", &node.IndexController{})
	beego.Router("/gui/monitor/node/data", &node.DataController{})
	beego.Router("/gui/monitor/container", &container.IndexController{})
	beego.Router("/gui/monitor/container/data", &container.DataController{})
	beego.Router("/gui/monitor/historicalcontainer", &historicalcontainer.IndexController{})
	beego.Router("/gui/monitor/historicalcontainer/data", &historicalcontainer.DataController{})
	beego.Router("/gui/event/audit/list", &audit.ListController{})
	beego.Router("/gui/event/kubernetes/list", &kubernetes.ListController{})
	beego.Router("/gui/event/kubernetes/acknowledge", &kubernetes.AcknowledgeController{})
	beego.Router("/gui/notification/notifier/list", &notifier.ListController{})
	beego.Router("/gui/notification/notifier/edit", &notifier.EditController{})
	beego.Router("/gui/notification/notifier/delete", &notifier.DeleteController{})
	beego.Router("/gui/system/about", &about.IndexController{})
	beego.Router("/gui/system/namespace/list", &namespace.ListController{})
	beego.Router("/gui/system/namespace/edit", &namespace.EditController{})
	beego.Router("/gui/system/namespace/select", &namespace.SelectController{})
	beego.Router("/gui/system/namespace/bookmark", &namespace.BookmarkController{})
	beego.Router("/gui/system/namespace/delete", &namespace.DeleteController{})
	beego.Router("/gui/system/notification/emailserver/list", &emailserver.ListController{})
	beego.Router("/gui/system/notification/emailserver/create", &emailserver.CreateController{})
	beego.Router("/gui/system/notification/emailserver/delete", &emailserver.DeleteController{})
	beego.Router("/gui/system/notification/sms/list", &sms.ListController{})
	beego.Router("/gui/system/notification/sms/create", &sms.CreateController{})
	beego.Router("/gui/system/notification/sms/delete", &sms.DeleteController{})
	beego.Router("/gui/system/host/credential/list", &credential.ListController{})
	beego.Router("/gui/system/host/credential/edit", &credential.EditController{})
	beego.Router("/gui/system/host/credential/delete", &credential.DeleteController{})
	beego.Router("/gui/system/rbac/user/list", &user.ListController{})
	beego.Router("/gui/system/rbac/user/delete", &user.DeleteController{})
	beego.Router("/gui/system/rbac/user/edit", &user.EditController{})
	beego.Router("/gui/system/rbac/role/list", &role.ListController{})
	beego.Router("/gui/system/rbac/role/delete", &role.DeleteController{})
	beego.Router("/gui/system/rbac/role/edit", &role.EditController{})
	beego.Router("/gui/system/upgrade", &upgrade.IndexController{})
	beego.Router("/gui/system/upgrade/websocket", &upgrade.WebSocketController{})
}
