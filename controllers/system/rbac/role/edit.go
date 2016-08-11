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

package role

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/rbac"
	"github.com/cloudawan/cloudone_utility/restclient"
)

type EditController struct {
	beego.Controller
}

func (c *EditController) Get() {
	c.TplName = "system/rbac/role/edit.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")

	action := c.GetString("action")
	name := c.GetString("name")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	c.Data["action"] = action

	if action == "create" {
		c.Data["actionButtonValue"] = "Create"
		c.Data["pageHeader"] = "Create Role"
	} else {
		url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
			"/api/v1/authorizations/roles/" + name

		role := rbac.Role{}

		tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

		_, err := restclient.RequestGetWithStructure(url, &role, tokenHeaderMap)

		if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
			return
		}

		if err != nil {
			// Error
			guimessage.AddDanger(guimessagedisplay.GetErrorMessage(err))
			c.Ctx.Redirect(302, "/gui/system/rbac/role/list")
			guimessage.RedirectMessage(c)
			return
		}

		pathMap := make(map[string]bool)
		for _, permission := range role.PermissionSlice {
			if permission.Component == identity.GetConponentName() && permission.Method == "GET" {
				pathMap[permission.Path] = true
			}
		}

		// Dashboard
		setCheckedTag("/gui/dashboard", "checkedTagDashboard", c.Data, pathMap)
		setHiddenTag("/gui/dashboard", "hiddenTagDashboard", c.Data, pathMap)
		setCheckedTag("/gui/dashboard/topology", "checkedTagDashboardTopology", c.Data, pathMap)
		setCheckedTag("/gui/dashboard/healthcheck", "checkedTagDashboardHealthCheck", c.Data, pathMap)
		setCheckedTag("/gui/dashboard/bluegreen", "checkedTagDashboardBlueGreen", c.Data, pathMap)
		setCheckedTag("/gui/dashboard/appservice", "checkedTagDashboardAppService", c.Data, pathMap)
		setCheckedTag("/gui/dashboard/deploy", "checkedTagDashboardDeploy", c.Data, pathMap)

		// Repository
		setCheckedTag("/gui/repository", "checkedTagRepository", c.Data, pathMap)
		setHiddenTag("/gui/repository", "hiddenTagRepository", c.Data, pathMap)
		setCheckedTag("/gui/repository/imageinformation", "checkedTagRepositoryImageInformation", c.Data, pathMap)
		setHiddenTag("/gui/repository/imageinformation", "hiddenTagRepositoryImageInformation", c.Data, pathMap)
		setCheckedTag("/gui/repository/imageinformation/list", "checkedTagRepositoryImageInformationList", c.Data, pathMap)
		setCheckedTag("/gui/repository/imageinformation/create", "checkedTagRepositoryImageInformationCreate", c.Data, pathMap)
		setCheckedTag("/gui/repository/imageinformation/upgrade", "checkedTagRepositoryImageInformationUpgrade", c.Data, pathMap)
		setCheckedTag("/gui/repository/imageinformation/log", "checkedTagRepositoryImageInformationLog", c.Data, pathMap)
		setCheckedTag("/gui/repository/imageinformation/delete", "checkedTagRepositoryImageInformationDelete", c.Data, pathMap)
		setCheckedTag("/gui/repository/imagerecord", "checkedTagRepositoryImageRecord", c.Data, pathMap)
		setHiddenTag("/gui/repository/imagerecord", "hiddenTagRepositoryImageRecord", c.Data, pathMap)
		setCheckedTag("/gui/repository/imagerecord/list", "checkedTagRepositoryImageRecordList", c.Data, pathMap)
		setCheckedTag("/gui/repository/imagerecord/log", "checkedTagRepositoryImageRecordLog", c.Data, pathMap)
		setCheckedTag("/gui/repository/imagerecord/delete", "checkedTagRepositoryImageRecordDelete", c.Data, pathMap)
		setCheckedTag("/gui/repository/thirdparty", "checkedTagRepositoryThirdPartyService", c.Data, pathMap)
		setHiddenTag("/gui/repository/thirdparty", "hiddenTagRepositoryThirdPartyService", c.Data, pathMap)
		setCheckedTag("/gui/repository/thirdparty/list", "checkedTagRepositoryThirdPartyServiceList", c.Data, pathMap)
		setCheckedTag("/gui/repository/thirdparty/edit", "checkedTagRepositoryThirdPartyServiceCreate", c.Data, pathMap)
		setCheckedTag("/gui/repository/thirdparty/launch", "checkedTagRepositoryThirdPartyServiceLaunch", c.Data, pathMap)
		setCheckedTag("/gui/repository/thirdparty/delete", "checkedTagRepositoryThirdPartyServiceDelete", c.Data, pathMap)
		setCheckedTag("/gui/repository/topologytemplate", "checkedTagRepositoryTopologyTemplate", c.Data, pathMap)
		setHiddenTag("/gui/repository/topologytemplate", "hiddenTagRepositoryTopologyTemplate", c.Data, pathMap)
		setCheckedTag("/gui/repository/topologytemplate/list", "checkedTagRepositoryTopologyTemplateList", c.Data, pathMap)
		setCheckedTag("/gui/repository/topologytemplate/clone", "checkedTagRepositoryTopologyTemplateClone", c.Data, pathMap)
		setCheckedTag("/gui/repository/topologytemplate/delete", "checkedTagRepositoryTopologyTemplateDelete", c.Data, pathMap)

		// Deploy
		setCheckedTag("/gui/deploy", "checkedTagDeploy", c.Data, pathMap)
		setHiddenTag("/gui/deploy", "hiddenTagDeploy", c.Data, pathMap)
		setCheckedTag("/gui/deploy/deploy", "checkedTagDeployDeploy", c.Data, pathMap)
		setHiddenTag("/gui/deploy/deploy", "hiddenTagDeployDeploy", c.Data, pathMap)
		setCheckedTag("/gui/deploy/deploy/list", "checkedTagDeployDeployList", c.Data, pathMap)
		setCheckedTag("/gui/deploy/deploy/create", "checkedTagDeployDeployCreate", c.Data, pathMap)
		setCheckedTag("/gui/deploy/deploy/update", "checkedTagDeployDeployUpdate", c.Data, pathMap)
		setCheckedTag("/gui/deploy/deploy/resize", "checkedTagDeployDeployResize", c.Data, pathMap)
		setCheckedTag("/gui/deploy/deploy/delete", "checkedTagDeployDeployDelete", c.Data, pathMap)
		setCheckedTag("/gui/deploy/deploybluegreen", "checkedTagDeployDeployBlueGreen", c.Data, pathMap)
		setCheckedTag("/gui/deploy/clone", "checkedTagDeployClone", c.Data, pathMap)
		setHiddenTag("/gui/deploy/deploybluegreen", "hiddenTagDeployDeployBlueGreen", c.Data, pathMap)
		setCheckedTag("/gui/deploy/deploybluegreen/list", "checkedTagDeployDeployBlueGreenList", c.Data, pathMap)
		setCheckedTag("/gui/deploy/deploybluegreen/select", "checkedTagDeployDeployBlueGreenSelect", c.Data, pathMap)
		setCheckedTag("/gui/deploy/deploybluegreen/delete", "checkedTagDeployDeployBlueGreenDelete", c.Data, pathMap)
		setCheckedTag("/gui/deploy/autoscaler", "checkedTagDeployAutoScaler", c.Data, pathMap)
		setHiddenTag("/gui/deploy/autoscaler", "hiddenTagDeployDeployAutoScaler", c.Data, pathMap)
		setCheckedTag("/gui/deploy/autoscaler/list", "checkedTagDeployAutoScalerList", c.Data, pathMap)
		setCheckedTag("/gui/deploy/autoscaler/edit", "checkedTagDeployAutoScalerCreate", c.Data, pathMap)
		setCheckedTag("/gui/deploy/autoscaler/delete", "checkedTagDeployAutoScalerDelete", c.Data, pathMap)
		setCheckedTag("/gui/deploy/deployclusterapplication", "checkedTagDeployDeployClusterApplication", c.Data, pathMap)
		setHiddenTag("/gui/deploy/deployclusterapplication", "hiddenTagDeployDeployClusterApplication", c.Data, pathMap)
		setCheckedTag("/gui/deploy/deployclusterapplication/list", "checkedTagDeployDeployClusterApplicationList", c.Data, pathMap)
		setCheckedTag("/gui/deploy/deployclusterapplication/size", "checkedTagDeployDeployClusterApplicationSize", c.Data, pathMap)
		setCheckedTag("/gui/deploy/deployclusterapplication/delete", "checkedTagDeployDeployClusterApplicationDelete", c.Data, pathMap)

		// Inventory
		setCheckedTag("/gui/inventory", "checkedTagInventory", c.Data, pathMap)
		setHiddenTag("/gui/inventory", "hiddenTagInventory", c.Data, pathMap)
		setCheckedTag("/gui/inventory/replicationcontroller", "checkedTagInventoryReplicationController", c.Data, pathMap)
		setHiddenTag("/gui/inventory/replicationcontroller", "hiddenTagInventoryReplicationController", c.Data, pathMap)
		setCheckedTag("/gui/inventory/replicationcontroller/list", "checkedTagInventoryReplicationControllerList", c.Data, pathMap)
		setCheckedTag("/gui/inventory/replicationcontroller/edit", "checkedTagInventoryReplicationControllerCreate", c.Data, pathMap)
		setCheckedTag("/gui/inventory/replicationcontroller/size", "checkedTagInventoryReplicationControllerSize", c.Data, pathMap)
		setCheckedTag("/gui/inventory/replicationcontroller/delete", "checkedTagInventoryReplicationControllerDelete", c.Data, pathMap)
		setCheckedTag("/gui/inventory/replicationcontroller/pod/log", "checkedTagInventoryReplicationControllerPodLog", c.Data, pathMap)
		setCheckedTag("/gui/inventory/replicationcontroller/pod/delete", "checkedTagInventoryReplicationControllerPodDelete", c.Data, pathMap)
		setCheckedTag("/gui/inventory/replicationcontroller/dockerterminal", "checkedTagInventoryReplicationControllerDockerTerminal", c.Data, pathMap)
		setCheckedTag("/gui/inventory/service", "checkedTagInventoryService", c.Data, pathMap)
		setHiddenTag("/gui/inventory/service", "hiddenTagInventoryService", c.Data, pathMap)
		setCheckedTag("/gui/inventory/service/list", "checkedTagInventoryServiceList", c.Data, pathMap)
		setCheckedTag("/gui/inventory/service/edit", "checkedTagInventoryServiceCreate", c.Data, pathMap)
		setCheckedTag("/gui/inventory/service/delete", "checkedTagInventoryServiceDelete", c.Data, pathMap)

		// File System
		setCheckedTag("/gui/filesystem", "checkedTagFilesystem", c.Data, pathMap)
		setHiddenTag("/gui/filesystem", "hiddenTagFilesystem", c.Data, pathMap)
		setCheckedTag("/gui/filesystem/glusterfs", "checkedTagFilesystemGlusterfs", c.Data, pathMap)
		setHiddenTag("/gui/filesystem/glusterfs", "hiddenTagFilesystemGlusterfs", c.Data, pathMap)
		setCheckedTag("/gui/filesystem/glusterfs/cluster", "checkedTagFilesystemGlusterfsCluster", c.Data, pathMap)
		setHiddenTag("/gui/filesystem/glusterfs/cluster", "hiddenTagFilesystemGlusterfsCluster", c.Data, pathMap)
		setCheckedTag("/gui/filesystem/glusterfs/cluster/list", "checkedTagFilesystemGlusterfsClusterList", c.Data, pathMap)
		setCheckedTag("/gui/filesystem/glusterfs/cluster/edit", "checkedTagFilesystemGlusterfsClusterCreate", c.Data, pathMap)
		setCheckedTag("/gui/filesystem/glusterfs/cluster/delete", "checkedTagFilesystemGlusterfsClusterDelete", c.Data, pathMap)
		setCheckedTag("/gui/filesystem/glusterfs/volume", "checkedTagFilesystemGlusterfsVolume", c.Data, pathMap)
		setHiddenTag("/gui/filesystem/glusterfs/volume", "hiddenTagFilesystemGlusterfsVolume", c.Data, pathMap)
		setCheckedTag("/gui/filesystem/glusterfs/volume/list", "checkedTagFilesystemGlusterfsVolumeList", c.Data, pathMap)
		setCheckedTag("/gui/filesystem/glusterfs/volume/create", "checkedTagFilesystemGlusterfsVolumeCreate", c.Data, pathMap)
		setCheckedTag("/gui/filesystem/glusterfs/volume/reset", "checkedTagFilesystemGlusterfsVolumeReset", c.Data, pathMap)
		setCheckedTag("/gui/filesystem/glusterfs/volume/delete", "checkedTagFilesystemGlusterfsVolumeDelete", c.Data, pathMap)

		// Monitor
		setCheckedTag("/gui/monitor", "checkedTagMonitor", c.Data, pathMap)
		setHiddenTag("/gui/monitor", "hiddenTagMonitor", c.Data, pathMap)
		setCheckedTag("/gui/monitor/node", "checkedTagMonitorNode", c.Data, pathMap)
		setCheckedTag("/gui/monitor/container", "checkedTagMonitorContainer", c.Data, pathMap)
		setCheckedTag("/gui/monitor/historicalcontainer", "checkedTagMonitorHistoricalContainer", c.Data, pathMap)

		// Event
		setCheckedTag("/gui/event", "checkedTagEvent", c.Data, pathMap)
		setHiddenTag("/gui/event", "hiddenTagEvent", c.Data, pathMap)
		setCheckedTag("/gui/event/audit", "checkedTagEventAudit", c.Data, pathMap)
		setHiddenTag("/gui/event/audit", "hiddenTagEventAudit", c.Data, pathMap)
		setCheckedTag("/gui/event/audit/list", "checkedTagEventAuditList", c.Data, pathMap)
		setCheckedTag("/gui/event/kubernetes", "checkedTagEventKubernetes", c.Data, pathMap)
		setHiddenTag("/gui/event/kubernetes", "hiddenTagEventKubernetes", c.Data, pathMap)
		setCheckedTag("/gui/event/kubernetes/list", "checkedTagEventKubernetesList", c.Data, pathMap)
		setCheckedTag("/gui/event/kubernetes/acknowledge", "checkedTagEventKubernetesAcknowledge", c.Data, pathMap)

		// Notification
		setCheckedTag("/gui/notification", "checkedTagNotification", c.Data, pathMap)
		setHiddenTag("/gui/notification", "hiddenTagNotification", c.Data, pathMap)
		setCheckedTag("/gui/notification/notifier", "checkedTagNotificationNotifier", c.Data, pathMap)
		setHiddenTag("/gui/notification/notifier", "hiddenTagNotificationNotifier", c.Data, pathMap)
		setCheckedTag("/gui/notification/notifier/list", "checkedTagNotificationNotifierList", c.Data, pathMap)
		setCheckedTag("/gui/notification/notifier/edit", "checkedTagNotificationNotifierCreate", c.Data, pathMap)
		setCheckedTag("/gui/notification/notifier/delete", "checkedTagNotificationNotifierDelete", c.Data, pathMap)

		// System
		setCheckedTag("/gui/system", "checkedTagSystem", c.Data, pathMap)
		setHiddenTag("/gui/system", "hiddenTagSystem", c.Data, pathMap)
		setCheckedTag("/gui/system/about", "checkedTagSystemAbout", c.Data, pathMap)
		setCheckedTag("/gui/system/namespace", "checkedTagSystemNamespace", c.Data, pathMap)
		setHiddenTag("/gui/system/namespace", "hiddenTagSystemNamespace", c.Data, pathMap)
		setCheckedTag("/gui/system/namespace/list", "checkedTagSystemNamespaceList", c.Data, pathMap)
		setCheckedTag("/gui/system/namespace/edit", "checkedTagSystemNamespaceCreate", c.Data, pathMap)
		setCheckedTag("/gui/system/namespace/select", "checkedTagSystemNamespaceSelect", c.Data, pathMap)
		setCheckedTag("/gui/system/namespace/bookmark", "checkedTagSystemNamespaceBookmark", c.Data, pathMap)
		setCheckedTag("/gui/system/namespace/delete", "checkedTagSystemNamespaceDelete", c.Data, pathMap)
		setCheckedTag("/gui/system/notification", "checkedTagSystemNotification", c.Data, pathMap)
		setHiddenTag("/gui/system/notification", "hiddenTagSystemNotification", c.Data, pathMap)
		setCheckedTag("/gui/system/notification/emailserver", "checkedTagSystemNotificationEmailServer", c.Data, pathMap)
		setHiddenTag("/gui/system/notification/emailserver", "hiddenTagSystemNotificationEmailServer", c.Data, pathMap)
		setCheckedTag("/gui/system/notification/emailserver/list", "checkedTagSystemNotificationEmailServerList", c.Data, pathMap)
		setCheckedTag("/gui/system/notification/emailserver/create", "checkedTagSystemNotificationEmailServerCreate", c.Data, pathMap)
		setCheckedTag("/gui/system/notification/emailserver/delete", "checkedTagSystemNotificationEmailServerDelete", c.Data, pathMap)
		setCheckedTag("/gui/system/notification/sms", "checkedTagSystemNotificationSMS", c.Data, pathMap)
		setHiddenTag("/gui/system/notification/sms", "hiddenTagSystemNotificationnSMS", c.Data, pathMap)
		setCheckedTag("/gui/system/notification/sms/list", "checkedTagSystemNotificationSMSList", c.Data, pathMap)
		setCheckedTag("/gui/system/notification/sms/create", "checkedTagSystemNotificationSMSCreate", c.Data, pathMap)
		setCheckedTag("/gui/system/notification/sms/delete", "checkedTagSystemNotificationSMSDelete", c.Data, pathMap)
		setCheckedTag("/gui/system/host", "checkedTagSystemHost", c.Data, pathMap)
		setHiddenTag("/gui/system/host", "hiddenTagSystemHost", c.Data, pathMap)
		setCheckedTag("/gui/system/host/credential", "checkedTagSystemHostCredential", c.Data, pathMap)
		setHiddenTag("/gui/system/host/credential", "hiddenTagSystemHostCredential", c.Data, pathMap)
		setCheckedTag("/gui/system/host/credential/list", "checkedTagSystemHostCredentialList", c.Data, pathMap)
		setCheckedTag("/gui/system/host/credential/edit", "checkedTagSystemHostCredentialCreate", c.Data, pathMap)
		setCheckedTag("/gui/system/host/credential/delete", "checkedTagSystemHostCredentialDelete", c.Data, pathMap)
		setCheckedTag("/gui/system/rbac", "checkedTagSystemRBAC", c.Data, pathMap)
		setHiddenTag("/gui/system/rbac", "hiddenTagSystemRBAC", c.Data, pathMap)
		setCheckedTag("/gui/system/rbac/user", "checkedTagSystemRBACUser", c.Data, pathMap)
		setHiddenTag("/gui/system/rbac/user", "hiddenTagSystemRBACUser", c.Data, pathMap)
		setCheckedTag("/gui/system/rbac/user/list", "checkedTagSystemRBACUserList", c.Data, pathMap)
		setCheckedTag("/gui/system/rbac/user/edit", "checkedTagSystemRBACUserCreate", c.Data, pathMap)
		setCheckedTag("/gui/system/rbac/user/delete", "checkedTagSystemRBACUserDelete", c.Data, pathMap)
		setCheckedTag("/gui/system/rbac/role", "checkedTagSystemRBACRole", c.Data, pathMap)
		setHiddenTag("/gui/system/rbac/role", "hiddenTagSystemRBACRole", c.Data, pathMap)
		setCheckedTag("/gui/system/rbac/role/list", "checkedTagSystemRBACRoleList", c.Data, pathMap)
		setCheckedTag("/gui/system/rbac/role/edit", "checkedTagSystemRBACRoleCreate", c.Data, pathMap)
		setCheckedTag("/gui/system/rbac/role/delete", "checkedTagSystemRBACRoleDelete", c.Data, pathMap)
		setCheckedTag("/gui/system/privateregistry", "checkedTagSystemPrivateRegistry", c.Data, pathMap)
		setHiddenTag("/gui/system/privateregistry", "hiddenTagSystemPrivateRegistry", c.Data, pathMap)
		setCheckedTag("/gui/system/privateregistry/server", "checkedTagSystemPrivateRegistryServer", c.Data, pathMap)
		setHiddenTag("/gui/system/privateregistry/server", "hiddenTagSystemPrivateRegistryServer", c.Data, pathMap)
		setCheckedTag("/gui/system/privateregistry/server/list", "checkedTagSystemPrivateRegistryServerList", c.Data, pathMap)
		setCheckedTag("/gui/system/privateregistry/server/edit", "checkedTagSystemPrivateRegistryServerCreate", c.Data, pathMap)
		setCheckedTag("/gui/system/privateregistry/server/delete", "checkedTagSystemPrivateRegistryServerDelete", c.Data, pathMap)
		setCheckedTag("/gui/system/privateregistry/repository", "checkedTagSystemPrivateRegistryRepository", c.Data, pathMap)
		setHiddenTag("/gui/system/privateregistry/repository", "hiddenTagSystemPrivateRegistryRepository", c.Data, pathMap)
		setCheckedTag("/gui/system/privateregistry/repository/list", "checkedTagSystemPrivateRegistryRepositoryList", c.Data, pathMap)
		setCheckedTag("/gui/system/privateregistry/repository/delete", "checkedTagSystemPrivateRegistryRepositoryDelete", c.Data, pathMap)
		setCheckedTag("/gui/system/privateregistry/image", "checkedTagSystemPrivateRegistryImage", c.Data, pathMap)
		setHiddenTag("/gui/system/privateregistry/image", "hiddenTagSystemPrivateRegistryImage", c.Data, pathMap)
		setCheckedTag("/gui/system/privateregistry/image/list", "checkedTagSystemPrivateRegistryImageList", c.Data, pathMap)
		setCheckedTag("/gui/system/privateregistry/image/delete", "checkedTagSystemPrivateRegistryImageDelete", c.Data, pathMap)
		setCheckedTag("/gui/system/slb", "checkedTagSystemSLB", c.Data, pathMap)
		setHiddenTag("/gui/system/slb", "hiddenTagSystemSLB", c.Data, pathMap)
		setCheckedTag("/gui/system/slb/daemon", "checkedTagSystemSLBDaemon", c.Data, pathMap)
		setHiddenTag("/gui/system/slb/daemon", "hiddenTagSystemSLBDaemon", c.Data, pathMap)
		setCheckedTag("/gui/system/slb/daemon/list", "checkedTagSystemSLBDaemonList", c.Data, pathMap)
		setCheckedTag("/gui/system/slb/daemon/edit", "checkedTagSystemSLBDaemonCreate", c.Data, pathMap)
		setCheckedTag("/gui/system/slb/daemon/configure", "checkedTagSystemSLBDaemonConfigure", c.Data, pathMap)
		setCheckedTag("/gui/system/slb/daemon/delete", "checkedTagSystemSLBDaemonDelete", c.Data, pathMap)
		setCheckedTag("/gui/system/upgrade", "checkedTagSystemUpgrade", c.Data, pathMap)

		c.Data["name"] = name
		c.Data["description"] = role.Description
		c.Data["readonly"] = "readonly"

		c.Data["actionButtonValue"] = "Update"
		c.Data["pageHeader"] = "Update Role"
	}

	guimessage.OutputMessage(c.Data)
}

func setCheckedTag(targetPath string, checkedTag string, data map[interface{}]interface{}, pathMap map[string]bool) {
	if pathMap[targetPath] {
		data[checkedTag] = "checked"
	} else {
		data[checkedTag] = ""
	}
}

func setHiddenTag(targetPath string, hiddenTag string, data map[interface{}]interface{}, pathMap map[string]bool) {
	if pathMap[targetPath] {
		data[hiddenTag] = "hidden"
	} else {
		data[hiddenTag] = ""
	}
}

func (c *EditController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	name := c.GetString("name")
	description := c.GetString("description")
	action := c.GetString("action")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	permissionSlice := make([]*rbac.Permission, 0)

	// Dashboard
	if c.GetString("dashboard") == "on" {
		permission := &rbac.Permission{"dashboard", identity.GetConponentName(), "GET", "/gui/dashboard"}
		permissionSlice = append(permissionSlice, permission)
	} else {
		if c.GetString("dashboardTopology") == "on" {
			permission := &rbac.Permission{"dashboardTopology", identity.GetConponentName(), "GET", "/gui/dashboard/topology"}
			permissionSlice = append(permissionSlice, permission)
		}
		if c.GetString("dashboardHealthCheck") == "on" {
			permission := &rbac.Permission{"dashboardHealthCheck", identity.GetConponentName(), "GET", "/gui/dashboard/healthcheck"}
			permissionSlice = append(permissionSlice, permission)
		}
		if c.GetString("dashboardBlueGree") == "on" {
			permission := &rbac.Permission{"dashboardBlueGreen", identity.GetConponentName(), "GET", "/gui/dashboard/bluegreen"}
			permissionSlice = append(permissionSlice, permission)
		}
		if c.GetString("dashboardAppService") == "on" {
			permission := &rbac.Permission{"dashboardAppService", identity.GetConponentName(), "GET", "/gui/dashboard/appservice"}
			permissionSlice = append(permissionSlice, permission)
		}
		if c.GetString("dashboardDeploy") == "on" {
			permission := &rbac.Permission{"dashboardDeploy", identity.GetConponentName(), "GET", "/gui/dashboard/deploy"}
			permissionSlice = append(permissionSlice, permission)
		}
	}

	// Repository
	if c.GetString("repository") == "on" {
		permission := &rbac.Permission{"repository", identity.GetConponentName(), "GET", "/gui/repository"}
		permissionSlice = append(permissionSlice, permission)
	} else {
		if c.GetString("repositoryImageInformation") == "on" {
			permission := &rbac.Permission{"repositoryImageInformation", identity.GetConponentName(), "GET", "/gui/repository/imageinformation"}
			permissionSlice = append(permissionSlice, permission)
		} else {
			if c.GetString("repositoryImageInformationList") == "on" {
				permission := &rbac.Permission{"repositoryImageInformationList", identity.GetConponentName(), "GET", "/gui/repository/imageinformation/list"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("repositoryImageInformationCreate") == "on" {
				permission := &rbac.Permission{"repositoryImageInformationCreate", identity.GetConponentName(), "GET", "/gui/repository/imageinformation/create"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("repositoryImageInformationUpgrade") == "on" {
				permission := &rbac.Permission{"repositoryImageInformationUpgrade", identity.GetConponentName(), "GET", "/gui/repository/imageinformation/upgrade"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("repositoryImageInformationLog") == "on" {
				permission := &rbac.Permission{"repositoryImageInformationLog", identity.GetConponentName(), "GET", "/gui/repository/imageinformation/log"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("repositoryImageInformationDelete") == "on" {
				permission := &rbac.Permission{"repositoryImageInformationDelete", identity.GetConponentName(), "GET", "/gui/repository/imageinformation/delete"}
				permissionSlice = append(permissionSlice, permission)
			}
		}

		if c.GetString("repositoryImageRecord") == "on" {
			permission := &rbac.Permission{"repositoryImageRecord", identity.GetConponentName(), "GET", "/gui/repository/imagerecord"}
			permissionSlice = append(permissionSlice, permission)
		} else {
			if c.GetString("repositoryImageRecordList") == "on" {
				permission := &rbac.Permission{"repositoryImageRecordList", identity.GetConponentName(), "GET", "/gui/repository/imagerecord/list"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("repositoryImageRecordLog") == "on" {
				permission := &rbac.Permission{"repositoryImageRecordLog", identity.GetConponentName(), "GET", "/gui/repository/imagerecord/log"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("repositoryImageRecordDelete") == "on" {
				permission := &rbac.Permission{"repositoryImageRecordDelete", identity.GetConponentName(), "GET", "/gui/repository/imagerecord/delete"}
				permissionSlice = append(permissionSlice, permission)
			}
		}

		if c.GetString("repositoryThirdPartyService") == "on" {
			permission := &rbac.Permission{"repositoryThirdPartyService", identity.GetConponentName(), "GET", "/gui/repository/thirdparty"}
			permissionSlice = append(permissionSlice, permission)
		} else {
			if c.GetString("repositoryThirdPartyServiceList") == "on" {
				permission := &rbac.Permission{"repositoryThirdPartyServiceList", identity.GetConponentName(), "GET", "/gui/repository/thirdparty/list"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("repositoryThirdPartyServiceCreate") == "on" {
				permission := &rbac.Permission{"repositoryThirdPartyServiceCreate", identity.GetConponentName(), "GET", "/gui/repository/thirdparty/edit"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("repositoryThirdPartyServiceLaunch") == "on" {
				permission := &rbac.Permission{"repositoryThirdPartyServiceLaunch", identity.GetConponentName(), "GET", "/gui/repository/thirdparty/launch"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("repositoryThirdPartyServiceDelete") == "on" {
				permission := &rbac.Permission{"repositoryThirdPartyServiceDelete", identity.GetConponentName(), "GET", "/gui/repository/thirdparty/delete"}
				permissionSlice = append(permissionSlice, permission)
			}
		}

		if c.GetString("repositoryTopologyTemplate") == "on" {
			permission := &rbac.Permission{"repositoryTopologyTemplate", identity.GetConponentName(), "GET", "/gui/repository/topologytemplate"}
			permissionSlice = append(permissionSlice, permission)
		} else {
			if c.GetString("repositoryTopologyTemplateList") == "on" {
				permission := &rbac.Permission{"repositoryTopologyTemplateList", identity.GetConponentName(), "GET", "/gui/repository/topologytemplate/list"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("repositoryTopologyTemplateClone") == "on" {
				permission := &rbac.Permission{"repositoryTopologyTemplateClone", identity.GetConponentName(), "GET", "/gui/repository/topologytemplate/clone"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("repositoryTopologyTemplateDelete") == "on" {
				permission := &rbac.Permission{"repositoryTopologyTemplateDelete", identity.GetConponentName(), "GET", "/gui/repository/topologytemplate/delete"}
				permissionSlice = append(permissionSlice, permission)
			}
		}
	}

	// Deploy
	if c.GetString("deploy") == "on" {
		permission := &rbac.Permission{"deploy", identity.GetConponentName(), "GET", "/gui/deploy"}
		permissionSlice = append(permissionSlice, permission)
	} else {
		if c.GetString("deployDeploy") == "on" {
			permission := &rbac.Permission{"deployDeploy", identity.GetConponentName(), "GET", "/gui/deploy/deploy"}
			permissionSlice = append(permissionSlice, permission)
		} else {
			if c.GetString("deployDeployList") == "on" {
				permission := &rbac.Permission{"deployDeployList", identity.GetConponentName(), "GET", "/gui/deploy/deploy/list"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("deployDeployCreate") == "on" {
				permission := &rbac.Permission{"deployDeployCreate", identity.GetConponentName(), "GET", "/gui/deploy/deploy/create"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("deployDeployUpdate") == "on" {
				permission := &rbac.Permission{"deployDeployUpdate", identity.GetConponentName(), "GET", "/gui/deploy/deploy/update"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("deployDeployResize") == "on" {
				permission := &rbac.Permission{"deployDeployResize", identity.GetConponentName(), "GET", "/gui/deploy/deploy/resize"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("deployDeployDelete") == "on" {
				permission := &rbac.Permission{"deployDeployDelete", identity.GetConponentName(), "GET", "/gui/deploy/deploy/delete"}
				permissionSlice = append(permissionSlice, permission)
			}
		}

		if c.GetString("deployDeployBlueGreen") == "on" {
			permission := &rbac.Permission{"deployDeployBlueGreen", identity.GetConponentName(), "GET", "/gui/deploy/deploybluegreen"}
			permissionSlice = append(permissionSlice, permission)
		} else {
			if c.GetString("deployDeployBlueGreenList") == "on" {
				permission := &rbac.Permission{"deployDeployBlueGreenList", identity.GetConponentName(), "GET", "/gui/deploy/deploybluegreen/list"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("deployDeployBlueGreenSelect") == "on" {
				permission := &rbac.Permission{"deployDeployBlueGreenSelect", identity.GetConponentName(), "GET", "/gui/deploy/deploybluegreen/select"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("deployDeployBlueGreenDelete") == "on" {
				permission := &rbac.Permission{"deployDeployBlueGreenDelete", identity.GetConponentName(), "GET", "/gui/deploy/deploybluegreen/delete"}
				permissionSlice = append(permissionSlice, permission)
			}
		}

		if c.GetString("deployAutoScaler") == "on" {
			permission := &rbac.Permission{"deployAutoScaler", identity.GetConponentName(), "GET", "/gui/deploy/autoscaler"}
			permissionSlice = append(permissionSlice, permission)
		} else {
			if c.GetString("deployAutoScalerList") == "on" {
				permission := &rbac.Permission{"deployAutoScalerList", identity.GetConponentName(), "GET", "/gui/deploy/autoscaler/list"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("deployAutoScalerCreate") == "on" {
				permission := &rbac.Permission{"deployAutoScalerCreate", identity.GetConponentName(), "GET", "/gui/deploy/autoscaler/edit"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("deployAutoScalerDelete") == "on" {
				permission := &rbac.Permission{"deployAutoScalerDelete", identity.GetConponentName(), "GET", "/gui/deploy/autoscaler/delete"}
				permissionSlice = append(permissionSlice, permission)
			}
		}

		if c.GetString("deployDeployClusterApplication") == "on" {
			permission := &rbac.Permission{"deployDeployClusterApplication", identity.GetConponentName(), "GET", "/gui/deploy/deployclusterapplication"}
			permissionSlice = append(permissionSlice, permission)
		} else {
			if c.GetString("deployDeployClusterApplicationList") == "on" {
				permission := &rbac.Permission{"deployDeployClusterApplicationList", identity.GetConponentName(), "GET", "/gui/deploy/deployclusterapplication/list"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("deployDeployClusterApplicationSize") == "on" {
				permission := &rbac.Permission{"deployDeployClusterApplicationSize", identity.GetConponentName(), "GET", "/gui/deploy/deployclusterapplication/size"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("deployDeployClusterApplicationDelete") == "on" {
				permission := &rbac.Permission{"deployDeployClusterApplicationDelete", identity.GetConponentName(), "GET", "/gui/deploy/deployclusterapplication/delete"}
				permissionSlice = append(permissionSlice, permission)
			}
		}

		if c.GetString("deployClone") == "on" {
			permission := &rbac.Permission{"deployClone", identity.GetConponentName(), "GET", "/gui/deploy/clone"}
			permissionSlice = append(permissionSlice, permission)
		}
	}

	// Inventory
	if c.GetString("inventory") == "on" {
		permission := &rbac.Permission{"inventory", identity.GetConponentName(), "GET", "/gui/inventory"}
		permissionSlice = append(permissionSlice, permission)
	} else {
		if c.GetString("inventoryReplicationController") == "on" {
			permission := &rbac.Permission{"inventoryReplicationController", identity.GetConponentName(), "GET", "/gui/inventory/replicationcontroller"}
			permissionSlice = append(permissionSlice, permission)
		} else {
			if c.GetString("inventoryReplicationControllerList") == "on" {
				permission := &rbac.Permission{"inventoryReplicationControllerList", identity.GetConponentName(), "GET", "/gui/inventory/replicationcontroller/list"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("inventoryReplicationControllerCreate") == "on" {
				permission := &rbac.Permission{"inventoryReplicationControllerCreate", identity.GetConponentName(), "GET", "/gui/inventory/replicationcontroller/edit"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("inventoryReplicationControllerSize") == "on" {
				permission := &rbac.Permission{"inventoryReplicationControllerSize", identity.GetConponentName(), "GET", "/gui/inventory/replicationcontroller/size"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("inventoryReplicationControllerDelete") == "on" {
				permission := &rbac.Permission{"inventoryReplicationControllerDelete", identity.GetConponentName(), "GET", "/gui/inventory/replicationcontroller/delete"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("inventoryReplicationControllerPodLog") == "on" {
				permission := &rbac.Permission{"inventoryReplicationControllerPodLog", identity.GetConponentName(), "GET", "/gui/inventory/replicationcontroller/pod/log"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("inventoryReplicationControllerPodDelete") == "on" {
				permission := &rbac.Permission{"inventoryReplicationControllerPodDelete", identity.GetConponentName(), "GET", "/gui/inventory/replicationcontroller/pod/delete"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("inventoryReplicationControllerDockerTerminal") == "on" {
				permission := &rbac.Permission{"inventoryReplicationControllerDockerTerminal", identity.GetConponentName(), "GET", "/gui/inventory/replicationcontroller/dockerterminal"}
				permissionSlice = append(permissionSlice, permission)
			}
		}

		if c.GetString("inventoryService") == "on" {
			permission := &rbac.Permission{"inventoryService", identity.GetConponentName(), "GET", "/gui/inventory/service"}
			permissionSlice = append(permissionSlice, permission)
		} else {
			if c.GetString("inventoryServiceList") == "on" {
				permission := &rbac.Permission{"inventoryServiceList", identity.GetConponentName(), "GET", "/gui/inventory/service/list"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("inventoryServiceCreate") == "on" {
				permission := &rbac.Permission{"inventoryServiceCreate", identity.GetConponentName(), "GET", "/gui/inventory/service/edit"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("inventoryServiceDelete") == "on" {
				permission := &rbac.Permission{"inventoryServiceDelete", identity.GetConponentName(), "GET", "/gui/inventory/service/delete"}
				permissionSlice = append(permissionSlice, permission)
			}
		}
	}

	// File System
	if c.GetString("filesystem") == "on" {
		permission := &rbac.Permission{"filesystem", identity.GetConponentName(), "GET", "/gui/filesystem"}
		permissionSlice = append(permissionSlice, permission)
	} else {
		if c.GetString("filesystemGlusterfs") == "on" {
			permission := &rbac.Permission{"filesystemGlusterfs", identity.GetConponentName(), "GET", "/gui/filesystem/glusterfs"}
			permissionSlice = append(permissionSlice, permission)
		} else {
			if c.GetString("filesystemGlusterfsCluster") == "on" {
				permission := &rbac.Permission{"filesystemGlusterfsCluster", identity.GetConponentName(), "GET", "/gui/filesystem/glusterfs/cluster"}
				permissionSlice = append(permissionSlice, permission)
			} else {
				if c.GetString("filesystemGlusterfsClusterList") == "on" {
					permission := &rbac.Permission{"filesystemGlusterfsClusterList", identity.GetConponentName(), "GET", "/gui/filesystem/glusterfs/cluster/list"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("filesystemGlusterfsClusterCreate") == "on" {
					permission := &rbac.Permission{"filesystemGlusterfsClusterCreate", identity.GetConponentName(), "GET", "/gui/filesystem/glusterfs/cluster/edit"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("filesystemGlusterfsClusterDelete") == "on" {
					permission := &rbac.Permission{"filesystemGlusterfsClusterDelete", identity.GetConponentName(), "GET", "/gui/filesystem/glusterfs/cluster/delete"}
					permissionSlice = append(permissionSlice, permission)
				}
			}

			if c.GetString("filesystemGlusterfsVolume") == "on" {
				permission := &rbac.Permission{"filesystemGlusterfsVolume", identity.GetConponentName(), "GET", "/gui/filesystem/glusterfs/volume"}
				permissionSlice = append(permissionSlice, permission)
			} else {
				if c.GetString("filesystemGlusterfsVolumeList") == "on" {
					permission := &rbac.Permission{"filesystemGlusterfsVolumeList", identity.GetConponentName(), "GET", "/gui/filesystem/glusterfs/volume/list"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("filesystemGlusterfsVolumeCreate") == "on" {
					permission := &rbac.Permission{"filesystemGlusterfsVolumeCreate", identity.GetConponentName(), "GET", "/gui/filesystem/glusterfs/volume/create"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("filesystemGlusterfsVolumeReset") == "on" {
					permission := &rbac.Permission{"filesystemGlusterfsVolumeReset", identity.GetConponentName(), "GET", "/gui/filesystem/glusterfs/volume/reset"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("filesystemGlusterfsVolumeDelete") == "on" {
					permission := &rbac.Permission{"filesystemGlusterfsVolumeDelete", identity.GetConponentName(), "GET", "/gui/filesystem/glusterfs/volume/delete"}
					permissionSlice = append(permissionSlice, permission)
				}
			}
		}
	}

	// Monitor
	if c.GetString("monitor") == "on" {
		permission := &rbac.Permission{"monitor", identity.GetConponentName(), "GET", "/gui/monitor"}
		permissionSlice = append(permissionSlice, permission)
	} else {
		if c.GetString("monitorNode") == "on" {
			permission := &rbac.Permission{"monitorNode", identity.GetConponentName(), "GET", "/gui/monitor/node"}
			permissionSlice = append(permissionSlice, permission)
		}
		if c.GetString("monitorContainer") == "on" {
			permission := &rbac.Permission{"monitorContainer", identity.GetConponentName(), "GET", "/gui/monitor/container"}
			permissionSlice = append(permissionSlice, permission)
		}
		if c.GetString("monitorHistoricalContainer") == "on" {
			permission := &rbac.Permission{"monitorHistoricalContainer", identity.GetConponentName(), "GET", "/gui/monitor/historicalcontainer"}
			permissionSlice = append(permissionSlice, permission)
		}
	}

	// Event
	if c.GetString("event") == "on" {
		permission := &rbac.Permission{"event", identity.GetConponentName(), "GET", "/gui/event"}
		permissionSlice = append(permissionSlice, permission)
	} else {
		if c.GetString("eventKubernetes") == "on" {
			permission := &rbac.Permission{"eventKubernetes", identity.GetConponentName(), "GET", "/gui/event/kubernetes"}
			permissionSlice = append(permissionSlice, permission)
		} else {
			if c.GetString("eventKubernetesList") == "on" {
				permission := &rbac.Permission{"eventKubernetesList", identity.GetConponentName(), "GET", "/gui/event/kubernetes/list"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("eventKubernetesAcknowledge") == "on" {
				permission := &rbac.Permission{"eventKubernetesAcknowledge", identity.GetConponentName(), "GET", "/gui/event/kubernetes/acknowledge"}
				permissionSlice = append(permissionSlice, permission)
			}
		}
	}

	// Notification
	if c.GetString("notification") == "on" {
		permission := &rbac.Permission{"notification", identity.GetConponentName(), "GET", "/gui/notification"}
		permissionSlice = append(permissionSlice, permission)
	} else {
		if c.GetString("notificationNotifier") == "on" {
			permission := &rbac.Permission{"notificationNotifier", identity.GetConponentName(), "GET", "/gui/notification/notifier"}
			permissionSlice = append(permissionSlice, permission)
		} else {
			if c.GetString("notificationNotifierList") == "on" {
				permission := &rbac.Permission{"notificationNotifierList", identity.GetConponentName(), "GET", "/gui/notification/notifier/list"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("notificationNotifierCreate") == "on" {
				permission := &rbac.Permission{"notificationNotifierCreate", identity.GetConponentName(), "GET", "/gui/notification/notifier/edit"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("notificationNotifierDelete") == "on" {
				permission := &rbac.Permission{"notificationNotifierDelete", identity.GetConponentName(), "GET", "/gui/notification/notifier/delete"}
				permissionSlice = append(permissionSlice, permission)
			}
		}
	}

	// System
	if c.GetString("system") == "on" {
		permission := &rbac.Permission{"system", identity.GetConponentName(), "GET", "/gui/system"}
		permissionSlice = append(permissionSlice, permission)
	} else {
		if c.GetString("systemAbout") == "on" {
			permission := &rbac.Permission{"systemAbout", identity.GetConponentName(), "GET", "/gui/system/about"}
			permissionSlice = append(permissionSlice, permission)
		}

		if c.GetString("systemNamespace") == "on" {
			permission := &rbac.Permission{"systemNamespace", identity.GetConponentName(), "GET", "/gui/system/namespace"}
			permissionSlice = append(permissionSlice, permission)
		} else {
			if c.GetString("systemNamespaceList") == "on" {
				permission := &rbac.Permission{"systemNamespaceList", identity.GetConponentName(), "GET", "/gui/system/namespace/list"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("systemNamespaceCreate") == "on" {
				permission := &rbac.Permission{"systemNamespaceCreate", identity.GetConponentName(), "GET", "/gui/system/namespace/edit"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("systemNamespaceSelect") == "on" {
				permission := &rbac.Permission{"systemNamespaceSelect", identity.GetConponentName(), "GET", "/gui/system/namespace/select"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("systemNamespaceBookmark") == "on" {
				permission := &rbac.Permission{"systemNamespaceBookmark", identity.GetConponentName(), "GET", "/gui/system/namespace/bookmark"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("systemNamespaceDelete") == "on" {
				permission := &rbac.Permission{"systemNamespaceDelete", identity.GetConponentName(), "GET", "/gui/system/namespace/delete"}
				permissionSlice = append(permissionSlice, permission)
			}
		}

		if c.GetString("systemNotification") == "on" {
			permission := &rbac.Permission{"systemNotification", identity.GetConponentName(), "GET", "/gui/system/notification"}
			permissionSlice = append(permissionSlice, permission)
		} else {
			if c.GetString("systemNotificationEmailServer") == "on" {
				permission := &rbac.Permission{"systemNotificationEmailServer", identity.GetConponentName(), "GET", "/gui/system/notification/emailserver"}
				permissionSlice = append(permissionSlice, permission)
			} else {
				if c.GetString("systemNotificationEmailServerList") == "on" {
					permission := &rbac.Permission{"systemNotificationEmailServerList", identity.GetConponentName(), "GET", "/gui/system/notification/emailserver/list"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("systemNotificationEmailServerCreate") == "on" {
					permission := &rbac.Permission{"systemNotificationEmailServerCreate", identity.GetConponentName(), "GET", "/gui/system/notification/emailserver/create"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("systemNotificationEmailServerDelete") == "on" {
					permission := &rbac.Permission{"systemNotificationEmailServerDelete", identity.GetConponentName(), "GET", "/gui/system/notification/emailserver/delete"}
					permissionSlice = append(permissionSlice, permission)
				}
			}

			if c.GetString("systemNotificationSMS") == "on" {
				permission := &rbac.Permission{"systemNotificationSMS", identity.GetConponentName(), "GET", "/gui/system/notification/sms"}
				permissionSlice = append(permissionSlice, permission)
			} else {
				if c.GetString("systemNotificationSMSList") == "on" {
					permission := &rbac.Permission{"systemNotificationSMSList", identity.GetConponentName(), "GET", "/gui/system/notification/sms/list"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("systemNotificationSMSCreate") == "on" {
					permission := &rbac.Permission{"systemNotificationSMSCreate", identity.GetConponentName(), "GET", "/gui/system/notification/sms/create"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("systemNotificationSMSDelete") == "on" {
					permission := &rbac.Permission{"systemNotificationSMSDelete", identity.GetConponentName(), "GET", "/gui/system/notification/sms/delete"}
					permissionSlice = append(permissionSlice, permission)
				}
			}
		}

		if c.GetString("systemHost") == "on" {
			permission := &rbac.Permission{"systemHost", identity.GetConponentName(), "GET", "/gui/system/host"}
			permissionSlice = append(permissionSlice, permission)
		} else {
			if c.GetString("systemHostCredential") == "on" {
				permission := &rbac.Permission{"systemHostCredential", identity.GetConponentName(), "GET", "/gui/system/host/credential"}
				permissionSlice = append(permissionSlice, permission)
			} else {
				if c.GetString("systemHostCredentialList") == "on" {
					permission := &rbac.Permission{"systemHostCredentialList", identity.GetConponentName(), "GET", "/gui/system/host/credential/list"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("systemHostCredentialCreate") == "on" {
					permission := &rbac.Permission{"systemHostCredentialCreate", identity.GetConponentName(), "GET", "/gui/system/host/credential/edit"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("systemHostCredentialDelete") == "on" {
					permission := &rbac.Permission{"systemHostCredentialDelete", identity.GetConponentName(), "GET", "/gui/system/host/credential/delete"}
					permissionSlice = append(permissionSlice, permission)
				}
			}
		}

		if c.GetString("systemRBAC") == "on" {
			permission := &rbac.Permission{"systemRBAC", identity.GetConponentName(), "GET", "/gui/system/rbac"}
			permissionSlice = append(permissionSlice, permission)
		} else {
			if c.GetString("systemRBACUser") == "on" {
				permission := &rbac.Permission{"systemRBACUser", identity.GetConponentName(), "GET", "/gui/system/rbac/user"}
				permissionSlice = append(permissionSlice, permission)
			} else {
				if c.GetString("systemRBACUserList") == "on" {
					permission := &rbac.Permission{"systemRBACUserList", identity.GetConponentName(), "GET", "/gui/system/rbac/user/list"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("systemRBACUserCreate") == "on" {
					permission := &rbac.Permission{"systemRBACUserCreate", identity.GetConponentName(), "GET", "/gui/system/rbac/user/edit"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("systemRBACUserDelete") == "on" {
					permission := &rbac.Permission{"systemRBACUserDelete", identity.GetConponentName(), "GET", "/gui/system/rbac/user/delete"}
					permissionSlice = append(permissionSlice, permission)
				}
			}

			if c.GetString("systemRBACRole") == "on" {
				permission := &rbac.Permission{"systemRBACRole", identity.GetConponentName(), "GET", "/gui/system/rbac/role"}
				permissionSlice = append(permissionSlice, permission)
			} else {
				if c.GetString("systemRBACRoleList") == "on" {
					permission := &rbac.Permission{"systemRBACRoleList", identity.GetConponentName(), "GET", "/gui/system/rbac/role/list"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("systemRBACRoleCreate") == "on" {
					permission := &rbac.Permission{"systemRBACRoleCreate", identity.GetConponentName(), "GET", "/gui/system/rbac/role/edit"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("systemRBACRoleDelete") == "on" {
					permission := &rbac.Permission{"systemRBACRoleDelete", identity.GetConponentName(), "GET", "/gui/system/rbac/role/delete"}
					permissionSlice = append(permissionSlice, permission)
				}
			}
		}

		if c.GetString("systemPrivateRegistry") == "on" {
			permission := &rbac.Permission{"systemPrivateRegistry", identity.GetConponentName(), "GET", "/gui/system/privateregistry"}
			permissionSlice = append(permissionSlice, permission)
		} else {
			if c.GetString("systemPrivateRegistryServer") == "on" {
				permission := &rbac.Permission{"systemPrivateRegistryServer", identity.GetConponentName(), "GET", "/gui/system/privateregistry/server"}
				permissionSlice = append(permissionSlice, permission)
			} else {
				if c.GetString("systemPrivateRegistryServerList") == "on" {
					permission := &rbac.Permission{"systemPrivateRegistryServerList", identity.GetConponentName(), "GET", "/gui/system/privateregistry/server/list"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("systemPrivateRegistryServerEdit") == "on" {
					permission := &rbac.Permission{"systemPrivateRegistryServerEdit", identity.GetConponentName(), "GET", "/gui/system/privateregistry/server/edit"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("systemPrivateRegistryServerDelete") == "on" {
					permission := &rbac.Permission{"systemPrivateRegistryServerDelete", identity.GetConponentName(), "GET", "/gui/system/privateregistry/server/delete"}
					permissionSlice = append(permissionSlice, permission)
				}
			}

			if c.GetString("systemPrivateRegistryRepository") == "on" {
				permission := &rbac.Permission{"systemPrivateRegistryRepository", identity.GetConponentName(), "GET", "/gui/system/privateregistry/repository"}
				permissionSlice = append(permissionSlice, permission)
			} else {
				if c.GetString("systemPrivateRegistryRepositoryList") == "on" {
					permission := &rbac.Permission{"systemPrivateRegistryRepositoryList", identity.GetConponentName(), "GET", "/gui/system/privateregistry/repository/list"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("systemPrivateRegistryRepositoryDelete") == "on" {
					permission := &rbac.Permission{"systemPrivateRegistryRepositoryDelete", identity.GetConponentName(), "GET", "/gui/system/privateregistry/repository/delete"}
					permissionSlice = append(permissionSlice, permission)
				}
			}

			if c.GetString("systemPrivateRegistryImage") == "on" {
				permission := &rbac.Permission{"systemPrivateRegistryImage", identity.GetConponentName(), "GET", "/gui/system/privateregistry/image"}
				permissionSlice = append(permissionSlice, permission)
			} else {
				if c.GetString("systemPrivateRegistryImageList") == "on" {
					permission := &rbac.Permission{"systemPrivateRegistryImageList", identity.GetConponentName(), "GET", "/gui/system/privateregistry/image/list"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("systemPrivateRegistryRepositoryDelete") == "on" {
					permission := &rbac.Permission{"systemPrivateRegistryImageDelete", identity.GetConponentName(), "GET", "/gui/system/privateregistry/image/delete"}
					permissionSlice = append(permissionSlice, permission)
				}
			}
		}

		if c.GetString("systemSLB") == "on" {
			permission := &rbac.Permission{"systemSLB", identity.GetConponentName(), "GET", "/gui/system/slb"}
			permissionSlice = append(permissionSlice, permission)
		} else {
			if c.GetString("systemSLBDaemon") == "on" {
				permission := &rbac.Permission{"systemSLBDaemon", identity.GetConponentName(), "GET", "/gui/system/slb/daemon"}
				permissionSlice = append(permissionSlice, permission)
			} else {
				if c.GetString("systemSLBDaemonList") == "on" {
					permission := &rbac.Permission{"systemSLBDaemonList", identity.GetConponentName(), "GET", "/gui/system/slb/daemon/list"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("systemSLBDaemonEdit") == "on" {
					permission := &rbac.Permission{"systemSLBDaemonEdit", identity.GetConponentName(), "GET", "/gui/system/slb/daemon/edit"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("systemSLBDaemonConfigure") == "on" {
					permission := &rbac.Permission{"systemSLBDaemonConfigure", identity.GetConponentName(), "GET", "/gui/system/slb/daemon/configure"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("systemSLBDaemonDelete") == "on" {
					permission := &rbac.Permission{"systemSLBDaemonDelete", identity.GetConponentName(), "GET", "/gui/system/slb/daemon/delete"}
					permissionSlice = append(permissionSlice, permission)
				}
			}
		}

		if c.GetString("systemUpgrade") == "on" {
			permission := &rbac.Permission{"systemUpgrade", identity.GetConponentName(), "GET", "/gui/system/upgrade"}
			permissionSlice = append(permissionSlice, permission)
		}
	}

	// For simplified version, only check GUI. The others are all allowed
	permissionSlice = append(permissionSlice, &rbac.Permission{"cloudone-all", "cloudone", "*", "*"})
	permissionSlice = append(permissionSlice, &rbac.Permission{"cloudone_analysis-all", "cloudone_analysis", "*", "*"})
	permissionSlice = append(permissionSlice, &rbac.Permission{"cloudone_gui-POST", identity.GetConponentName(), "POST", "*"})
	// Essentail one
	permissionSlice = append(permissionSlice, &rbac.Permission{"cloudone_gui-logout", identity.GetConponentName(), "GET", "/gui/logout"})

	role := rbac.Role{
		name,
		permissionSlice,
		description,
	}

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	var err error
	if action == "create" {
		url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
			"/api/v1/authorizations/roles"

		_, err = restclient.RequestPostWithStructure(url, role, nil, tokenHeaderMap)

	} else {
		url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
			"/api/v1/authorizations/roles/" + name

		_, err = restclient.RequestPutWithStructure(url, role, nil, tokenHeaderMap)
	}

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(guimessagedisplay.GetErrorMessage(err))
	} else {
		guimessage.AddSuccess("Role " + name + " is edited")
	}

	c.Ctx.Redirect(302, "/gui/system/rbac/role/list")

	guimessage.RedirectMessage(c)
}
