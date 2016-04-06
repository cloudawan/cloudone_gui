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
			guimessage.AddDanger(err.Error())
			c.Ctx.Redirect(302, "/gui/system/rbac/role/list")
			guimessage.RedirectMessage(c)
			return
		}

		c.Data["name"] = name
		c.Data["description"] = role.Description

		c.Data["actionButtonValue"] = "Update"
		c.Data["pageHeader"] = "Update Role"
	}

	guimessage.OutputMessage(c.Data)
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
		permission := &rbac.Permission{"dashboard", "cloudone_gui", "GET", "/gui/dashboard"}
		permissionSlice = append(permissionSlice, permission)
	} else {
		if c.GetString("dashboardTopology") == "on" {
			permission := &rbac.Permission{"dashboardTopology", "cloudone_gui", "GET", "/gui/dashboard/topology"}
			permissionSlice = append(permissionSlice, permission)
		}
		if c.GetString("dashboardHealthCheck") == "on" {
			permission := &rbac.Permission{"dashboardHealthCheck", "cloudone_gui", "GET", "/gui/dashboard/healthcheck"}
			permissionSlice = append(permissionSlice, permission)
		}
		if c.GetString("dashboardBlueGree") == "on" {
			permission := &rbac.Permission{"dashboardBlueGreen", "cloudone_gui", "GET", "/gui/dashboard/bluegreen"}
			permissionSlice = append(permissionSlice, permission)
		}
		if c.GetString("dashboardAppService") == "on" {
			permission := &rbac.Permission{"dashboardAppService", "cloudone_gui", "GET", "/gui/dashboard/appservice"}
			permissionSlice = append(permissionSlice, permission)
		}
	}

	// Repository
	if c.GetString("repository") == "on" {
		permission := &rbac.Permission{"repository", "cloudone_gui", "GET", "/gui/repository"}
		permissionSlice = append(permissionSlice, permission)
	} else {
		if c.GetString("repositoryImageInformation") == "on" {
			permission := &rbac.Permission{"repositoryImageInformation", "cloudone_gui", "GET", "/gui/repository/imageinformation"}
			permissionSlice = append(permissionSlice, permission)
		} else {
			if c.GetString("repositoryImageInformationList") == "on" {
				permission := &rbac.Permission{"repositoryImageInformationList", "cloudone_gui", "GET", "/gui/repository/imageinformation/list"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("repositoryImageInformationCreate") == "on" {
				permission := &rbac.Permission{"repositoryImageInformationCreate", "cloudone_gui", "GET", "/gui/repository/imageinformation/create"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("repositoryImageInformationUpgrade") == "on" {
				permission := &rbac.Permission{"repositoryImageInformationUpgrade", "cloudone_gui", "GET", "/gui/repository/imageinformation/upgrade"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("repositoryImageInformationDeploy") == "on" {
				permission := &rbac.Permission{"repositoryImageInformationDeploy", "cloudone_gui", "GET", "/gui/deploy/deploy/create"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("repositoryImageInformationDelete") == "on" {
				permission := &rbac.Permission{"repositoryImageInformationDelete", "cloudone_gui", "GET", "/gui/repository/imageinformation/delete"}
				permissionSlice = append(permissionSlice, permission)
			}
		}

		if c.GetString("repositoryImageRecord") == "on" {
			permission := &rbac.Permission{"repositoryImageRecord", "cloudone_gui", "GET", "/gui/repository/imagerecord"}
			permissionSlice = append(permissionSlice, permission)
		} else {
			if c.GetString("repositoryImageRecordList") == "on" {
				permission := &rbac.Permission{"repositoryImageRecordList", "cloudone_gui", "GET", "/gui/repository/imagerecord/list"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("repositoryImageRecordDelete") == "on" {
				permission := &rbac.Permission{"repositoryImageRecordDelete", "cloudone_gui", "GET", "/gui/repository/imagerecord/delete"}
				permissionSlice = append(permissionSlice, permission)
			}
		}

		if c.GetString("repositoryThirdPartyService") == "on" {
			permission := &rbac.Permission{"repositoryThirdPartyService", "cloudone_gui", "GET", "/gui/repository/thirdparty"}
			permissionSlice = append(permissionSlice, permission)
		} else {
			if c.GetString("repositoryThirdPartyServiceList") == "on" {
				permission := &rbac.Permission{"repositoryThirdPartyServiceList", "cloudone_gui", "GET", "/gui/repository/thirdparty/list"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("repositoryThirdPartyServiceCreate") == "on" {
				permission := &rbac.Permission{"repositoryThirdPartyServiceCreate", "cloudone_gui", "GET", "/gui/repository/thirdparty/edit"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("repositoryThirdPartyServiceLaunch") == "on" {
				permission := &rbac.Permission{"repositoryThirdPartyServiceLaunch", "cloudone_gui", "GET", "/gui/repository/thirdparty/launch"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("repositoryImageInformationDeploy") == "on" {
				permission := &rbac.Permission{"repositoryImageInformationDeploy", "cloudone_gui", "GET", "/gui/repository/thirdparty/delete"}
				permissionSlice = append(permissionSlice, permission)
			}
		}
	}

	// Deploy
	if c.GetString("deploy") == "on" {
		permission := &rbac.Permission{"deploy", "cloudone_gui", "GET", "/gui/deploy"}
		permissionSlice = append(permissionSlice, permission)
	} else {
		if c.GetString("repositoryDeployDeploy") == "on" {
			permission := &rbac.Permission{"repositoryDeployDeploy", "cloudone_gui", "GET", "/gui/deploy/deploy"}
			permissionSlice = append(permissionSlice, permission)
		} else {
			if c.GetString("repositoryDeployDeployList") == "on" {
				permission := &rbac.Permission{"repositoryDeployDeployList", "cloudone_gui", "GET", "/gui/deploy/deploy/list"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("repositoryDeployDeployCreate") == "on" {
				permission := &rbac.Permission{"repositoryDeployDeployCreate", "cloudone_gui", "GET", "/gui/deploy/deploy/create"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("repositoryDeployDeployUpdate") == "on" {
				permission := &rbac.Permission{"repositoryDeployDeployUpdate", "cloudone_gui", "GET", "/gui/deploy/deploy/update"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("repositoryDeployDeployDelete") == "on" {
				permission := &rbac.Permission{"repositoryDeployDeployDelete", "cloudone_gui", "GET", "/gui/deploy/deploy/delete"}
				permissionSlice = append(permissionSlice, permission)
			}
		}

		if c.GetString("deployDeployBlueGreen") == "on" {
			permission := &rbac.Permission{"deployDeployBlueGreen", "cloudone_gui", "GET", "/gui/deploy/deploybluegreen"}
			permissionSlice = append(permissionSlice, permission)
		} else {
			if c.GetString("deployDeployBlueGreenList") == "on" {
				permission := &rbac.Permission{"deployDeployBlueGreenList", "cloudone_gui", "GET", "/gui/deploy/deploybluegreen/list"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("deployDeployBlueGreenSelect") == "on" {
				permission := &rbac.Permission{"deployDeployBlueGreenSelect", "cloudone_gui", "GET", "/gui/deploy/deploybluegreen/select"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("deployDeployBlueGreenDelete") == "on" {
				permission := &rbac.Permission{"deployDeployBlueGreenDelete", "cloudone_gui", "GET", "/gui/deploy/deploybluegreen/delete"}
				permissionSlice = append(permissionSlice, permission)
			}
		}

		if c.GetString("deployAutoScaler") == "on" {
			permission := &rbac.Permission{"deployAutoScaler", "cloudone_gui", "GET", "/gui/deploy/autoscaler"}
			permissionSlice = append(permissionSlice, permission)
		} else {
			if c.GetString("deployAutoScalerList") == "on" {
				permission := &rbac.Permission{"deployAutoScalerList", "cloudone_gui", "GET", "/gui/deploy/autoscaler/list"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("deployAutoScalerCreate") == "on" {
				permission := &rbac.Permission{"deployAutoScalerCreate", "cloudone_gui", "GET", "/gui/deploy/autoscaler/edit"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("deployAutoScalerDelete") == "on" {
				permission := &rbac.Permission{"deployAutoScalerDelete", "cloudone_gui", "GET", "/gui/deploy/autoscaler/delete"}
				permissionSlice = append(permissionSlice, permission)
			}
		}

		if c.GetString("deployDeployClusterApplication") == "on" {
			permission := &rbac.Permission{"deployDeployClusterApplication", "cloudone_gui", "GET", "/gui/deploy/deployclusterapplication"}
			permissionSlice = append(permissionSlice, permission)
		} else {
			if c.GetString("deployDeployClusterApplicationList") == "on" {
				permission := &rbac.Permission{"deployDeployClusterApplicationList", "cloudone_gui", "GET", "/gui/deploy/deployclusterapplication/list"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("deployDeployClusterApplicationSize") == "on" {
				permission := &rbac.Permission{"deployDeployClusterApplicationSize", "cloudone_gui", "GET", "/gui/deploy/deployclusterapplication/size"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("deployDeployClusterApplicationDelete") == "on" {
				permission := &rbac.Permission{"deployDeployClusterApplicationDelete", "cloudone_gui", "GET", "/gui/deploy/deployclusterapplication/delete"}
				permissionSlice = append(permissionSlice, permission)
			}
		}
	}

	// Inventory
	if c.GetString("inventory") == "on" {
		permission := &rbac.Permission{"inventory", "cloudone_gui", "GET", "/gui/inventory"}
		permissionSlice = append(permissionSlice, permission)
	} else {
		if c.GetString("inventoryReplicationController") == "on" {
			permission := &rbac.Permission{"inventoryReplicationController", "cloudone_gui", "GET", "/gui/inventory/replicationcontroller"}
			permissionSlice = append(permissionSlice, permission)
		} else {
			if c.GetString("inventoryReplicationControllerList") == "on" {
				permission := &rbac.Permission{"inventoryReplicationControllerList", "cloudone_gui", "GET", "/gui/inventory/replicationcontroller/list"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("inventoryReplicationControllerCreate") == "on" {
				permission := &rbac.Permission{"inventoryReplicationControllerCreate", "cloudone_gui", "GET", "/gui/inventory/replicationcontroller/edit"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("inventoryReplicationControllerSize") == "on" {
				permission := &rbac.Permission{"inventoryReplicationControllerSize", "cloudone_gui", "GET", "/gui/inventory/replicationcontroller/size"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("inventoryReplicationControllerDelete") == "on" {
				permission := &rbac.Permission{"inventoryReplicationControllerDelete", "cloudone_gui", "GET", "/gui/inventory/replicationcontroller/delete"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("inventoryReplicationControllerPodLog") == "on" {
				permission := &rbac.Permission{"inventoryReplicationControllerPodLog", "cloudone_gui", "GET", "/gui/inventory/replicationcontroller/podlog"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("inventoryReplicationControllerDockerTerminal") == "on" {
				permission := &rbac.Permission{"inventoryReplicationControllerDockerTerminal", "cloudone_gui", "GET", "/gui/inventory/replicationcontroller/dockerterminal"}
				permissionSlice = append(permissionSlice, permission)
			}
		}

		if c.GetString("inventoryService") == "on" {
			permission := &rbac.Permission{"inventoryService", "cloudone_gui", "GET", "/gui/inventory/service"}
			permissionSlice = append(permissionSlice, permission)
		} else {
			if c.GetString("inventoryServiceList") == "on" {
				permission := &rbac.Permission{"inventoryServiceList", "cloudone_gui", "GET", "/gui/inventory/service/list"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("inventoryServiceCreate") == "on" {
				permission := &rbac.Permission{"inventoryServiceCreate", "cloudone_gui", "GET", "/gui/inventory/service/edit"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("inventoryServiceDelete") == "on" {
				permission := &rbac.Permission{"inventoryServiceDelete", "cloudone_gui", "GET", "/gui/inventory/service/delete"}
				permissionSlice = append(permissionSlice, permission)
			}
		}
	}

	// File System
	if c.GetString("filesystem") == "on" {
		permission := &rbac.Permission{"filesystem", "cloudone_gui", "GET", "/gui/filesystem"}
		permissionSlice = append(permissionSlice, permission)
	} else {
		if c.GetString("filesystemGlusterfs") == "on" {
			permission := &rbac.Permission{"filesystemGlusterfs", "cloudone_gui", "GET", "/gui/filesystem/glusterfs"}
			permissionSlice = append(permissionSlice, permission)
		} else {
			if c.GetString("filesystemGlusterfsCluster") == "on" {
				permission := &rbac.Permission{"filesystemGlusterfsCluster", "cloudone_gui", "GET", "/gui/filesystem/glusterfs/cluster"}
				permissionSlice = append(permissionSlice, permission)
			} else {
				if c.GetString("filesystemGlusterfsClusterList") == "on" {
					permission := &rbac.Permission{"filesystemGlusterfsClusterList", "cloudone_gui", "GET", "/gui/filesystem/glusterfs/cluster/list"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("filesystemGlusterfsClusterCreate") == "on" {
					permission := &rbac.Permission{"filesystemGlusterfsClusterCreate", "cloudone_gui", "GET", "/gui/filesystem/glusterfs/cluster/edit"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("filesystemGlusterfsClusterDelete") == "on" {
					permission := &rbac.Permission{"filesystemGlusterfsClusterDelete", "cloudone_gui", "GET", "/gui/filesystem/glusterfs/cluster/delete"}
					permissionSlice = append(permissionSlice, permission)
				}
			}

			if c.GetString("filesystemGlusterfsVolume") == "on" {
				permission := &rbac.Permission{"filesystemGlusterfsVolume", "cloudone_gui", "GET", "/gui/filesystem/glusterfs/volume"}
				permissionSlice = append(permissionSlice, permission)
			} else {
				if c.GetString("filesystemGlusterfsVolumeList") == "on" {
					permission := &rbac.Permission{"filesystemGlusterfsVolumeList", "cloudone_gui", "GET", "/gui/filesystem/glusterfs/volume/list"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("filesystemGlusterfsVolumeCreate") == "on" {
					permission := &rbac.Permission{"filesystemGlusterfsVolumeCreate", "cloudone_gui", "GET", "/gui/filesystem/glusterfs/volume/create"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("filesystemGlusterfsVolumeDelete") == "on" {
					permission := &rbac.Permission{"filesystemGlusterfsVolumeDelete", "cloudone_gui", "GET", "/gui/filesystem/glusterfs/volume/delete"}
					permissionSlice = append(permissionSlice, permission)
				}
			}
		}
	}

	// Monitor
	if c.GetString("monitor") == "on" {
		permission := &rbac.Permission{"monitor", "cloudone_gui", "GET", "/gui/monitor"}
		permissionSlice = append(permissionSlice, permission)
	} else {
		if c.GetString("monitorNode") == "on" {
			permission := &rbac.Permission{"monitorNode", "cloudone_gui", "GET", "/gui/monitor/node"}
			permissionSlice = append(permissionSlice, permission)
		}
		if c.GetString("monitorContainer") == "on" {
			permission := &rbac.Permission{"monitorContainer", "cloudone_gui", "GET", "/gui/monitor/container"}
			permissionSlice = append(permissionSlice, permission)
		}
		if c.GetString("monitorHistoricalContainer") == "on" {
			permission := &rbac.Permission{"monitorHistoricalContainer", "cloudone_gui", "GET", "/gui/monitor/historicalcontainer"}
			permissionSlice = append(permissionSlice, permission)
		}
	}

	// Event
	if c.GetString("event") == "on" {
		permission := &rbac.Permission{"event", "cloudone_gui", "GET", "/gui/event"}
		permissionSlice = append(permissionSlice, permission)
	} else {
		if c.GetString("eventKubernetes") == "on" {
			permission := &rbac.Permission{"eventKubernetes", "cloudone_gui", "GET", "/gui/event/kubernetes"}
			permissionSlice = append(permissionSlice, permission)
		} else {
			if c.GetString("eventKubernetesList") == "on" {
				permission := &rbac.Permission{"eventKubernetesList", "cloudone_gui", "GET", "/gui/event/kubernetes/list"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("eventKubernetesAcknowledge") == "on" {
				permission := &rbac.Permission{"eventKubernetesAcknowledge", "cloudone_gui", "GET", "/gui/event/kubernetes/acknowledge"}
				permissionSlice = append(permissionSlice, permission)
			}
		}
	}

	// Notification
	if c.GetString("notification") == "on" {
		permission := &rbac.Permission{"notification", "cloudone_gui", "GET", "/gui/notification"}
		permissionSlice = append(permissionSlice, permission)
	} else {
		if c.GetString("notificationNotifier") == "on" {
			permission := &rbac.Permission{"notificationNotifier", "cloudone_gui", "GET", "/gui/notification/notifier"}
			permissionSlice = append(permissionSlice, permission)
		} else {
			if c.GetString("notificationNotifierList") == "on" {
				permission := &rbac.Permission{"notificationNotifierList", "cloudone_gui", "GET", "/gui/notification/notifier/list"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("notificationNotifierCreate") == "on" {
				permission := &rbac.Permission{"notificationNotifierCreate", "cloudone_gui", "GET", "/gui/notification/notifier/edit"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("notificationNotifierDelete") == "on" {
				permission := &rbac.Permission{"notificationNotifierDelete", "cloudone_gui", "GET", "/gui/notification/notifier/delete"}
				permissionSlice = append(permissionSlice, permission)
			}
		}
	}

	// System
	if c.GetString("system") == "on" {
		permission := &rbac.Permission{"system", "cloudone_gui", "GET", "/gui/system"}
		permissionSlice = append(permissionSlice, permission)
	} else {
		if c.GetString("systemNamespace") == "on" {
			permission := &rbac.Permission{"systemNamespace", "cloudone_gui", "GET", "/gui/system/namespace"}
			permissionSlice = append(permissionSlice, permission)
		} else {
			if c.GetString("systemNamespaceList") == "on" {
				permission := &rbac.Permission{"systemNamespaceList", "cloudone_gui", "GET", "/gui/system/namespace/list"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("systemNamespaceCreate") == "on" {
				permission := &rbac.Permission{"systemNamespaceCreate", "cloudone_gui", "GET", "/gui/system/namespace/edit"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("systemNamespaceSelect") == "on" {
				permission := &rbac.Permission{"systemNamespaceSelect", "cloudone_gui", "GET", "/gui/system/namespace/select"}
				permissionSlice = append(permissionSlice, permission)
			}
			if c.GetString("systemNamespaceDelete") == "on" {
				permission := &rbac.Permission{"systemNamespaceDelete", "cloudone_gui", "GET", "/gui/system/namespace/delete"}
				permissionSlice = append(permissionSlice, permission)
			}
		}

		if c.GetString("systemNotification") == "on" {
			permission := &rbac.Permission{"systemNotification", "cloudone_gui", "GET", "/gui/system/notification"}
			permissionSlice = append(permissionSlice, permission)
		} else {
			if c.GetString("systemNotificationEmailServer") == "on" {
				permission := &rbac.Permission{"systemNotificationEmailServer", "cloudone_gui", "GET", "/gui/system/notification/emailserver"}
				permissionSlice = append(permissionSlice, permission)
			} else {
				if c.GetString("systemNotificationEmailServerList") == "on" {
					permission := &rbac.Permission{"systemNotificationEmailServerList", "cloudone_gui", "GET", "/gui/system/notification/emailserver/list"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("systemNotificationEmailServerCreate") == "on" {
					permission := &rbac.Permission{"systemNotificationEmailServerCreate", "cloudone_gui", "GET", "/gui/system/notification/emailserver/create"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("systemNotificationEmailServerDelete") == "on" {
					permission := &rbac.Permission{"systemNotificationEmailServerDelete", "cloudone_gui", "GET", "/gui/system/notification/emailserver/delete"}
					permissionSlice = append(permissionSlice, permission)
				}
			}

			if c.GetString("systemNotificationSMS") == "on" {
				permission := &rbac.Permission{"systemNotificationSMS", "cloudone_gui", "GET", "/gui/system/notification/emailserver"}
				permissionSlice = append(permissionSlice, permission)
			} else {
				if c.GetString("systemNotificationSMSList") == "on" {
					permission := &rbac.Permission{"systemNotificationSMSList", "cloudone_gui", "GET", "/gui/system/notification/sms/list"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("systemNotificationSMSCreate") == "on" {
					permission := &rbac.Permission{"systemNotificationSMSCreate", "cloudone_gui", "GET", "/gui/system/notification/sms/create"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("systemNotificationSMSDelete") == "on" {
					permission := &rbac.Permission{"systemNotificationSMSDelete", "cloudone_gui", "GET", "/gui/system/notification/sms/delete"}
					permissionSlice = append(permissionSlice, permission)
				}
			}
		}

		if c.GetString("systemHost") == "on" {
			permission := &rbac.Permission{"systemHost", "cloudone_gui", "GET", "/gui/system/host"}
			permissionSlice = append(permissionSlice, permission)
		} else {
			if c.GetString("systemHostCredential") == "on" {
				permission := &rbac.Permission{"systemHostCredential", "cloudone_gui", "GET", "/gui/system/host/credential"}
				permissionSlice = append(permissionSlice, permission)
			} else {
				if c.GetString("systemHostCredentialList") == "on" {
					permission := &rbac.Permission{"systemHostCredentialList", "cloudone_gui", "GET", "/gui/system/host/credential/list"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("systemHostCredentialCreate") == "on" {
					permission := &rbac.Permission{"systemHostCredentialCreate", "cloudone_gui", "GET", "/gui/system/host/credential/edit"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("systemHostCredentialDelete") == "on" {
					permission := &rbac.Permission{"systemHostCredentialDelete", "cloudone_gui", "GET", "/gui/system/host/credential/delete"}
					permissionSlice = append(permissionSlice, permission)
				}
			}
		}

		if c.GetString("systemRBAC") == "on" {
			permission := &rbac.Permission{"systemRBAC", "cloudone_gui", "GET", "/gui/system/rbac"}
			permissionSlice = append(permissionSlice, permission)
		} else {
			if c.GetString("systemRBACUser") == "on" {
				permission := &rbac.Permission{"systemRBACUser", "cloudone_gui", "GET", "/gui/system/rbac/user"}
				permissionSlice = append(permissionSlice, permission)
			} else {
				if c.GetString("systemRBACUserList") == "on" {
					permission := &rbac.Permission{"systemRBACUserList", "cloudone_gui", "GET", "/gui/system/rbac/user/list"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("systemRBACUserCreate") == "on" {
					permission := &rbac.Permission{"systemRBACUserCreate", "cloudone_gui", "GET", "/gui/system/rbac/user/edit"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("systemRBACUserDelete") == "on" {
					permission := &rbac.Permission{"systemRBACUserDelete", "cloudone_gui", "GET", "/gui/system/rbac/user/delete"}
					permissionSlice = append(permissionSlice, permission)
				}
			}

			if c.GetString("systemRBACRole") == "on" {
				permission := &rbac.Permission{"systemRBACRole", "cloudone_gui", "GET", "/gui/system/rbac/role"}
				permissionSlice = append(permissionSlice, permission)
			} else {
				if c.GetString("systemRBACRoleList") == "on" {
					permission := &rbac.Permission{"systemRBACRoleList", "cloudone_gui", "GET", "/gui/system/rbac/role/list"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("systemRBACRoleCreate") == "on" {
					permission := &rbac.Permission{"systemRBACRoleCreate", "cloudone_gui", "GET", "/gui/system/rbac/role/edit"}
					permissionSlice = append(permissionSlice, permission)
				}
				if c.GetString("systemRBACRoleDelete") == "on" {
					permission := &rbac.Permission{"systemRBACRoleDelete", "cloudone_gui", "GET", "/gui/system/rbac/role/delete"}
					permissionSlice = append(permissionSlice, permission)
				}
			}
		}

		if c.GetString("systemUpgrade") == "on" {
			permission := &rbac.Permission{"systemUpgrade", "cloudone_gui", "GET", "/gui/system/upgrade"}
			permissionSlice = append(permissionSlice, permission)
		}
	}

	// For simplified version, only check GUI. The others are all allowed
	permissionSlice = append(permissionSlice, &rbac.Permission{"cloudone-all", "cloudone", "*", "*"})
	permissionSlice = append(permissionSlice, &rbac.Permission{"cloudone_analysis-all", "cloudone_analysis", "*", "*"})
	// Essentail one
	permissionSlice = append(permissionSlice, &rbac.Permission{"cloudone_gui_logout", "cloudone_gui", "GET", "/gui/logout"})

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
		guimessage.AddDanger(err.Error())
	} else {
		guimessage.AddSuccess("Role " + name + " is created")
	}

	c.Ctx.Redirect(302, "/gui/system/rbac/role/list")

	guimessage.RedirectMessage(c)
}
