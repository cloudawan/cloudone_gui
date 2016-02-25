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

// @APIVersion 1.0.0
// @Title Cloudone GUI REST API
// @Description This is used for client side GUI only. For functional API, use each component's API directly.
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/guirestapi/dashboard/healthcheck"
	"github.com/cloudawan/cloudone_gui/guirestapi/dashboard/topology"
	"github.com/cloudawan/cloudone_gui/guirestapi/deploy/autoscaler"
	"github.com/cloudawan/cloudone_gui/guirestapi/deploy/deploy"
	"github.com/cloudawan/cloudone_gui/guirestapi/deploy/deploybluegreen"
	"github.com/cloudawan/cloudone_gui/guirestapi/deploy/deployclusterapplication"
	"github.com/cloudawan/cloudone_gui/guirestapi/event/kubernetes"
	"github.com/cloudawan/cloudone_gui/guirestapi/filesystem/glusterfs/cluster"
	"github.com/cloudawan/cloudone_gui/guirestapi/filesystem/glusterfs/volume"
	"github.com/cloudawan/cloudone_gui/guirestapi/inventory/replicationcontroller"
	"github.com/cloudawan/cloudone_gui/guirestapi/inventory/service"
	"github.com/cloudawan/cloudone_gui/guirestapi/monitor/container"
	"github.com/cloudawan/cloudone_gui/guirestapi/monitor/historicalcontainer"
	"github.com/cloudawan/cloudone_gui/guirestapi/monitor/node"
	"github.com/cloudawan/cloudone_gui/guirestapi/notification/notifier"
	"github.com/cloudawan/cloudone_gui/guirestapi/repository/imageinformation"
	"github.com/cloudawan/cloudone_gui/guirestapi/repository/imagerecord"
	"github.com/cloudawan/cloudone_gui/guirestapi/repository/thirdparty"
	"github.com/cloudawan/cloudone_gui/guirestapi/system/namespace"
	"github.com/cloudawan/cloudone_gui/guirestapi/system/notification/emailserver"
	"github.com/cloudawan/cloudone_gui/guirestapi/system/notification/sms"
)

// Beego Rest API generation will use this file router.go
func init() {
	ns := beego.NewNamespace("/guirestapi/v1",
		beego.NSNamespace("/dashboardhealthcheck",
			beego.NSInclude(
				&healthcheck.ListController{},
			),
		),
		beego.NSNamespace("/dashboardtopology",
			beego.NSInclude(
				&topology.IndexController{},
			),
		),
		beego.NSNamespace("/deployautoscaler",
			beego.NSInclude(
				&autoscaler.DeleteController{},
				&autoscaler.EditController{},
				&autoscaler.ListController{},
			),
		),
		beego.NSNamespace("/deploydeploy",
			beego.NSInclude(
				&deploy.CreateController{},
				&deploy.DeleteController{},
				&deploy.ListController{},
				&deploy.UpdateController{},
			),
		),
		beego.NSNamespace("/deploydeploybluegreen",
			beego.NSInclude(
				&deploybluegreen.DeleteController{},
				&deploybluegreen.ListController{},
				&deploybluegreen.SelectController{},
			),
		),
		beego.NSNamespace("/deploydeployclusterapplication",
			beego.NSInclude(
				&deployclusterapplication.DeleteController{},
				&deployclusterapplication.ListController{},
				&deployclusterapplication.SizeController{},
			),
		),
		beego.NSNamespace("/eventkubernetes",
			beego.NSInclude(
				&kubernetes.AcknowledgeController{},
				&kubernetes.ListController{},
			),
		),
		beego.NSNamespace("/filesystemglusterfscluster",
			beego.NSInclude(
				&cluster.DeleteController{},
				&cluster.EditController{},
				&cluster.ListController{},
			),
		),
		beego.NSNamespace("/filesystemglusterfsvolume",
			beego.NSInclude(
				&volume.CreateController{},
				&volume.DeleteController{},
				&volume.ListController{},
			),
		),
		beego.NSNamespace("/inventoryreplicationcontroller",
			beego.NSInclude(
				&replicationcontroller.DeleteController{},
				&replicationcontroller.EditController{},
				&replicationcontroller.ListController{},
				&replicationcontroller.PodLogController{},
				&replicationcontroller.SizeController{},
			),
		),
		beego.NSNamespace("/inventoryservice",
			beego.NSInclude(
				&service.DeleteController{},
				&service.EditController{},
				&service.ListController{},
			),
		),
		beego.NSNamespace("/monitorcontainer",
			beego.NSInclude(
				&container.IndexController{},
				&container.DataController{},
			),
		),
		beego.NSNamespace("/monitorchistoricalcontainer",
			beego.NSInclude(
				&historicalcontainer.IndexController{},
				&historicalcontainer.DataController{},
			),
		),
		beego.NSNamespace("/monitornode",
			beego.NSInclude(
				&node.IndexController{},
				&node.DataController{},
			),
		),
		beego.NSNamespace("/notificationnotifier",
			beego.NSInclude(
				&notifier.DeleteController{},
				&notifier.EditController{},
				&notifier.ListController{},
			),
		),
		beego.NSNamespace("/repositoryimageinformation",
			beego.NSInclude(
				&imageinformation.CreateController{},
				&imageinformation.DeleteController{},
				&imageinformation.ListController{},
				&imageinformation.UpgradeController{},
			),
		),
		beego.NSNamespace("/repositoryimagerecord",
			beego.NSInclude(
				&imagerecord.DeleteController{},
				&imagerecord.ListController{},
			),
		),
		beego.NSNamespace("/repositorythirdparty",
			beego.NSInclude(
				&thirdparty.DeleteController{},
				&thirdparty.EditController{},
				&thirdparty.LaunchController{},
				&thirdparty.ListController{},
			),
		),
		beego.NSNamespace("/systemnamespace",
			beego.NSInclude(
				&namespace.DeleteController{},
				&namespace.EditController{},
				&namespace.ListController{},
				&namespace.SelectController{},
			),
		),
		beego.NSNamespace("/systemnotificationemailserver",
			beego.NSInclude(
				&emailserver.CreateController{},
				&emailserver.DeleteController{},
				&emailserver.ListController{},
			),
		),
		beego.NSNamespace("/systemnotificationsms",
			beego.NSInclude(
				&sms.CreateController{},
				&sms.DeleteController{},
				&sms.ListController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
