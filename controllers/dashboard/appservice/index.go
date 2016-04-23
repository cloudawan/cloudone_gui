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

package appservice

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/configuration"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/rbac"
	"github.com/cloudawan/cloudone_utility/restclient"
	"strconv"
)

type DeployInformation struct {
	Namespace                 string
	ImageInformationName      string
	CurrentVersion            string
	CurrentVersionDescription string
	Description               string
}

type Service struct {
	Name            string
	Namespace       string
	PortSlice       []ServicePort
	Selector        map[string]interface{}
	ClusterIP       string
	LabelMap        map[string]interface{}
	SessionAffinity string
	Display         string
}

type ServicePort struct {
	Name       string
	Protocol   string
	Port       int
	TargetPort string
	NodePort   int
}

type ReplicationControllerAndRelatedPod struct {
	Name               string
	Namespace          string
	ReplicaAmount      int
	AliveReplicaAmount int
	Selector           map[string]string
	Label              map[string]string
	PodSlice           []Pod
}

type Pod struct {
	Name           string
	Namespace      string
	HostIP         string
	PodIP          string
	ContainerSlice []PodContainer
}

type PodContainer struct {
	Name      string
	Image     string
	PortSlice []PodContainerPort
}

type PodContainerPort struct {
	Name          string
	ContainerPort int
	Protocol      string
}

type DeployClusterApplication struct {
	Name                           string
	Size                           int
	ServiceName                    string
	ReplicationControllerNameSlice []string
}

type IndexController struct {
	beego.Controller
}

func (c *IndexController) Get() {
	guimessage := guimessagedisplay.GetGUIMessage(c)
	c.TplName = "dashboard/appservice/index.html"

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")
	// Dashboard tab menu
	user, _ := c.GetSession("user").(*rbac.User)
	c.Data["dashboardTabMenu"] = identity.GetDashboardTabMenu(user, "appservice")

	cloudoneGUIProtocol := beego.AppConfig.String("cloudoneGUIProtocol")
	cloudoneGUIHost := c.Ctx.Input.Host()
	cloudoneGUIPort := c.Ctx.Input.Port()

	c.Data["cloudoneGUIProtocol"] = cloudoneGUIProtocol
	c.Data["cloudoneGUIHost"] = cloudoneGUIHost
	c.Data["cloudoneGUIPort"] = cloudoneGUIPort

	guimessage.OutputMessage(c.Data)
}

const (
	allKeyword = "All"
)

type DataController struct {
	beego.Controller
}

func (c *DataController) Get() {
	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")
	kubeapiHost, kubeapiPort, err := configuration.GetAvailableKubeapiHostAndPort()
	if err != nil {
		c.Data["json"].(map[string]interface{})["error"] = err.Error()
		c.ServeJSON()
		return
	}

	scope := c.GetString("scope")
	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	namespaceSlice := make([]string, 0)
	if scope == allKeyword {
		url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
			"/api/v1/namespaces/" + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

		_, err := restclient.RequestGetWithStructure(url, &namespaceSlice, tokenHeaderMap)

		if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
			return
		}

		if err != nil {
			c.Data["json"].(map[string]interface{})["error"] = err.Error()
			c.ServeJSON()
			return
		}
	} else {
		namespace, _ := c.GetSession("namespace").(string)
		namespaceSlice = append(namespaceSlice, namespace)
	}

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/deploys/"

	deployInformationSlice := make([]DeployInformation, 0)

	_, err = restclient.RequestGetWithStructure(url, &deployInformationSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		c.Data["json"].(map[string]interface{})["error"] = err.Error()
		c.ServeJSON()
		return
	}

	// Json
	c.Data["json"] = make(map[string]interface{})
	c.Data["json"].(map[string]interface{})["applicationlView"] = make([]interface{}, 0)
	c.Data["json"].(map[string]interface{})["thirdpartyView"] = make([]interface{}, 0)
	c.Data["json"].(map[string]interface{})["errorMap"] = make(map[string]interface{})

	// Application view
	applicationJsonMap := make(map[string]interface{})
	applicationJsonMap["name"] = "App View"
	applicationJsonMap["children"] = make([]interface{}, 0)
	c.Data["json"].(map[string]interface{})["applicationlView"] = append(c.Data["json"].(map[string]interface{})["applicationlView"].([]interface{}), applicationJsonMap)

	// Third-party view
	thirdpartyJsonMap := make(map[string]interface{})
	thirdpartyJsonMap["name"] = "3rd Party View"
	thirdpartyJsonMap["children"] = make([]interface{}, 0)
	c.Data["json"].(map[string]interface{})["thirdpartyView"] = append(c.Data["json"].(map[string]interface{})["thirdpartyView"].([]interface{}), thirdpartyJsonMap)

	applicationViewLeafAmount := 0
	thirdpartyViewLeafAmount := 0
	for _, namespace := range namespaceSlice {
		// Application view
		applicationNamespaceJsonMap := make(map[string]interface{})
		applicationNamespaceJsonMap["name"] = namespace
		applicationNamespaceJsonMap["children"] = make([]interface{}, 0)

		url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
			"/api/v1/services/" + namespace + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

		serviceSlice := make([]Service, 0)

		_, err := restclient.RequestGetWithStructure(url, &serviceSlice, tokenHeaderMap)

		if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
			return
		}

		if err != nil {
			// Error
			c.Data["json"].(map[string]interface{})["error"] = "Get service data error"
			c.Data["json"].(map[string]interface{})["errorMap"].(map[string]interface{})[namespace] = err.Error()
		} else {
			url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
				"/api/v1/replicationcontrollers/" + namespace + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

			replicationControllerAndRelatedPodSlice := make([]ReplicationControllerAndRelatedPod, 0)

			_, err := restclient.RequestGetWithStructure(url, &replicationControllerAndRelatedPodSlice, tokenHeaderMap)

			if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
				return
			}

			if err != nil {
				// Error
				c.Data["json"].(map[string]interface{})["error"] = "Get replication controller data error"
				c.Data["json"].(map[string]interface{})["errorMap"].(map[string]interface{})[namespace] = err.Error()
			} else {
				namespaceOwnAmount := 0
				for _, deployInformation := range deployInformationSlice {
					if namespace == deployInformation.Namespace {
						namespaceOwnAmount++

						serviceName := "NO SERVICE"
						for _, service := range serviceSlice {
							if deployInformation.ImageInformationName == service.Name {
								serviceName = service.Name
								for _, port := range service.PortSlice {
									nodePortText := ""
									if port.NodePort >= 0 {
										nodePortText = strconv.Itoa(port.NodePort)
									}
									serviceName += " " + strconv.Itoa(port.Port) + "/" + nodePortText
								}
							}
						}

						serviceJsonMap := make(map[string]interface{})
						serviceJsonMap["name"] = serviceName
						serviceJsonMap["children"] = make([]interface{}, 0)

						applicationNameJsonMap := make(map[string]interface{})
						applicationNameJsonMap["children"] = make([]interface{}, 0)

						instanceAmount := 0
						for _, replicationControllerAndRelatedPod := range replicationControllerAndRelatedPodSlice {
							if deployInformation.ImageInformationName+deployInformation.CurrentVersion == replicationControllerAndRelatedPod.Name {
								instanceAmount = replicationControllerAndRelatedPod.ReplicaAmount
								replicationControllerJsonMap := make(map[string]interface{})
								replicationControllerJsonMap["name"] = replicationControllerAndRelatedPod.Name + " (" + strconv.Itoa(instanceAmount) + ")"
								replicationControllerJsonMap["children"] = make([]interface{}, 0)
								applicationNameJsonMap["children"] = append(applicationNameJsonMap["children"].([]interface{}), replicationControllerJsonMap)
								applicationViewLeafAmount++
							}
						}

						applicationNameJsonMap["name"] = deployInformation.ImageInformationName + " (" + strconv.Itoa(instanceAmount) + ")"
						serviceJsonMap["children"] = append(serviceJsonMap["children"].([]interface{}), applicationNameJsonMap)
						applicationNamespaceJsonMap["children"] = append(applicationNamespaceJsonMap["children"].([]interface{}), serviceJsonMap)
					}
				}

				if namespaceOwnAmount == 0 {
					applicationViewLeafAmount++
				}

				applicationJsonMap["children"] = append(applicationJsonMap["children"].([]interface{}), applicationNamespaceJsonMap)
			}
		}

		// Third-party view
		thirdpartyNamespaceJsonMap := make(map[string]interface{})
		thirdpartyNamespaceJsonMap["name"] = namespace
		thirdpartyNamespaceJsonMap["children"] = make([]interface{}, 0)

		url = cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
			"/api/v1/deployclusterapplications/" + namespace

		deployClusterApplicationSlice := make([]DeployClusterApplication, 0)

		_, err = restclient.RequestGetWithStructure(url, &deployClusterApplicationSlice, tokenHeaderMap)

		if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
			return
		}

		if err != nil {
			// Error
			c.Data["json"].(map[string]interface{})["error"] = "Get third party application data error"
			c.Data["json"].(map[string]interface{})["errorMap"].(map[string]interface{})[namespace] = err.Error()
		} else {
			for _, deployClusterApplication := range deployClusterApplicationSlice {

				serviceName := deployClusterApplication.ServiceName
				for _, service := range serviceSlice {
					if service.Name == deployClusterApplication.ServiceName {
						for _, port := range service.PortSlice {
							nodePortText := ""
							if port.NodePort >= 0 {
								nodePortText = strconv.Itoa(port.NodePort)
							}
							serviceName += " " + strconv.Itoa(port.Port) + "/" + nodePortText
						}
					}
				}

				serviceJsonMap := make(map[string]interface{})
				serviceJsonMap["name"] = serviceName
				serviceJsonMap["children"] = make([]interface{}, 0)

				thirdpartyNameJsonMap := make(map[string]interface{})
				thirdpartyNameJsonMap["name"] = deployClusterApplication.Name + " (" + strconv.Itoa(deployClusterApplication.Size) + ")"
				thirdpartyNameJsonMap["children"] = make([]interface{}, 0)

				replicationControllerNameAmount := len(deployClusterApplication.ReplicationControllerNameSlice)
				for _, replicationControllerName := range deployClusterApplication.ReplicationControllerNameSlice {
					replicationControllerJsonMap := make(map[string]interface{})
					replicationControllerJsonMap["name"] = replicationControllerName + " (" + strconv.Itoa(deployClusterApplication.Size/replicationControllerNameAmount) + ")"
					replicationControllerJsonMap["children"] = make([]interface{}, 0)
					thirdpartyNameJsonMap["children"] = append(thirdpartyNameJsonMap["children"].([]interface{}), replicationControllerJsonMap)
					thirdpartyViewLeafAmount++
				}

				serviceJsonMap["children"] = append(serviceJsonMap["children"].([]interface{}), thirdpartyNameJsonMap)
				thirdpartyNamespaceJsonMap["children"] = append(thirdpartyNamespaceJsonMap["children"].([]interface{}), serviceJsonMap)
			}
			if len(deployClusterApplicationSlice) == 0 {
				thirdpartyViewLeafAmount++
			}
			thirdpartyJsonMap["children"] = append(thirdpartyJsonMap["children"].([]interface{}), thirdpartyNamespaceJsonMap)
		}
	}
	c.Data["json"].(map[string]interface{})["applicationViewLeafAmount"] = applicationViewLeafAmount
	c.Data["json"].(map[string]interface{})["thirdpartyViewLeafAmount"] = thirdpartyViewLeafAmount

	c.ServeJSON()
}
