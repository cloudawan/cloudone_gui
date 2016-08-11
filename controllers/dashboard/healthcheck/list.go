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

package healthcheck

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/rbac"
	"github.com/cloudawan/cloudone_utility/restclient"
	"sort"
	"strings"
)

type ListController struct {
	beego.Controller
}

type ComponentStatus struct {
	Cloudone         bool
	CloudoneAnalysis bool
	CloudoneGUI      bool
	Storage          bool
	ElasticSearch    bool
	Docker           bool
}

type KubernetesStatus struct {
	IP                    string
	Active                bool
	Docker                bool
	Flannel               bool
	KubeProxy             bool
	Kubelet               bool
	KubeApiserver         bool
	KubeScheduler         bool
	KubeControllerManager bool
	DockerIP              string
	FlannelIP             string
	DockerIPValid         bool
}

type GlusterfsClusterStatus struct {
	Name        string
	ServerSlice []GlusterfsServerStatus
}

type GlusterfsServerStatus struct {
	IP              string
	Active          bool
	GlusterfsServer bool
}

type SLBSetStatus struct {
	Name        string
	ServerSlice []SLBServerStatus
}

type SLBServerStatus struct {
	IP              string
	Active          bool
	CloudOneSLB     bool
	Keepalived      bool
	HAProxy         bool
	LastCommandTime string
}

type SortKubernetesStatusByIP []KubernetesStatus

func (sortKubernetesStatusByIP SortKubernetesStatusByIP) Len() int {
	return len(sortKubernetesStatusByIP)
}
func (sortKubernetesStatusByIP SortKubernetesStatusByIP) Swap(i, j int) {
	sortKubernetesStatusByIP[i], sortKubernetesStatusByIP[j] = sortKubernetesStatusByIP[j], sortKubernetesStatusByIP[i]
}
func (sortKubernetesStatusByIP SortKubernetesStatusByIP) Less(i, j int) bool {
	return sortKubernetesStatusByIP[i].IP < sortKubernetesStatusByIP[j].IP
}

type SortGlusterfsClusterStatusByName []GlusterfsClusterStatus

func (sortGlusterfsClusterStatusByName SortGlusterfsClusterStatusByName) Len() int {
	return len(sortGlusterfsClusterStatusByName)
}
func (sortGlusterfsClusterStatusByName SortGlusterfsClusterStatusByName) Swap(i, j int) {
	sortGlusterfsClusterStatusByName[i], sortGlusterfsClusterStatusByName[j] = sortGlusterfsClusterStatusByName[j], sortGlusterfsClusterStatusByName[i]
}
func (sortGlusterfsClusterStatusByName SortGlusterfsClusterStatusByName) Less(i, j int) bool {
	return sortGlusterfsClusterStatusByName[i].Name < sortGlusterfsClusterStatusByName[j].Name
}

type SortGlusterfsServerStatusByIP []GlusterfsServerStatus

func (sortGlusterfsServerStatusByIP SortGlusterfsServerStatusByIP) Len() int {
	return len(sortGlusterfsServerStatusByIP)
}
func (sortGlusterfsServerStatusByIP SortGlusterfsServerStatusByIP) Swap(i, j int) {
	sortGlusterfsServerStatusByIP[i], sortGlusterfsServerStatusByIP[j] = sortGlusterfsServerStatusByIP[j], sortGlusterfsServerStatusByIP[i]
}
func (sortGlusterfsServerStatusByIP SortGlusterfsServerStatusByIP) Less(i, j int) bool {
	return sortGlusterfsServerStatusByIP[i].IP < sortGlusterfsServerStatusByIP[j].IP
}

type SortSLBSetStatusByName []SLBSetStatus

func (sortSLBSetStatusByName SortSLBSetStatusByName) Len() int {
	return len(sortSLBSetStatusByName)
}
func (sortSLBSetStatusByName SortSLBSetStatusByName) Swap(i, j int) {
	sortSLBSetStatusByName[i], sortSLBSetStatusByName[j] = sortSLBSetStatusByName[j], sortSLBSetStatusByName[i]
}
func (sortSLBSetStatusByName SortSLBSetStatusByName) Less(i, j int) bool {
	return sortSLBSetStatusByName[i].Name < sortSLBSetStatusByName[j].Name
}

type SortSLBServerStatusByIP []SLBServerStatus

func (sortSLBServerStatusByIP SortSLBServerStatusByIP) Len() int {
	return len(sortSLBServerStatusByIP)
}
func (sortSLBServerStatusByIP SortSLBServerStatusByIP) Swap(i, j int) {
	sortSLBServerStatusByIP[i], sortSLBServerStatusByIP[j] = sortSLBServerStatusByIP[j], sortSLBServerStatusByIP[i]
}
func (sortSLBServerStatusByIP SortSLBServerStatusByIP) Less(i, j int) bool {
	return sortSLBServerStatusByIP[i].IP < sortSLBServerStatusByIP[j].IP
}

func (c *ListController) Get() {
	c.TplName = "dashboard/healthcheck/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")
	// Dashboard tab menu
	user, _ := c.GetSession("user").(*rbac.User)
	c.Data["dashboardTabMenu"] = identity.GetDashboardTabMenu(user, "healthcheck")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")
	cloudoneAnalysisProtocol := beego.AppConfig.String("cloudoneAnalysisProtocol")
	cloudoneAnalysisHost := beego.AppConfig.String("cloudoneAnalysisHost")
	cloudoneAnalysisPort := beego.AppConfig.String("cloudoneAnalysisPort")

	allErrorMessageSlice := make([]string, 0)

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/healthchecks/"

	cloudoneJsonMap := make(map[string]interface{}, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &cloudoneJsonMap, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		allErrorMessageSlice = append(allErrorMessageSlice, err.Error())
	}

	url = cloudoneAnalysisProtocol + "://" + cloudoneAnalysisHost + ":" + cloudoneAnalysisPort +
		"/api/v1/healthchecks/"

	cloudoneAnalysisJsonMap := make(map[string]interface{}, 0)

	_, err = restclient.RequestGetWithStructure(url, &cloudoneAnalysisJsonMap, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		allErrorMessageSlice = append(allErrorMessageSlice, err.Error())
	}

	componentStatus, errorMessageSlice := parseComponentStatus(cloudoneJsonMap, cloudoneAnalysisJsonMap)
	if len(errorMessageSlice) > 0 {
		// Error
		allErrorMessageSlice = append(allErrorMessageSlice, errorMessageSlice...)
	}

	componentStatusSlice := make([]ComponentStatus, 0)
	componentStatusSlice = append(componentStatusSlice, *componentStatus)

	kubernetesStatusSlice, errorMessageSlice := parseKubernetesStatusSlice(cloudoneJsonMap)
	if len(errorMessageSlice) > 0 {
		// Error
		allErrorMessageSlice = append(allErrorMessageSlice, errorMessageSlice...)
	}

	glusterfsClusterStatusSlice, errorMessageSlice := parseGlusterfsClusterStatusSlice(cloudoneJsonMap)
	if len(errorMessageSlice) > 0 {
		// Error
		allErrorMessageSlice = append(allErrorMessageSlice, errorMessageSlice...)
	}

	slbSetStatusSlice, errorMessageSlice := parseSLBSetStatusSlice(cloudoneJsonMap)
	if len(errorMessageSlice) > 0 {
		// Error
		allErrorMessageSlice = append(allErrorMessageSlice, errorMessageSlice...)
	}

	sort.Sort(SortKubernetesStatusByIP(kubernetesStatusSlice))
	sort.Sort(SortGlusterfsClusterStatusByName(glusterfsClusterStatusSlice))
	sort.Sort(SortSLBSetStatusByName(slbSetStatusSlice))

	c.Data["componentStatusSlice"] = componentStatusSlice
	c.Data["kubernetesStatusSlice"] = kubernetesStatusSlice
	c.Data["glusterfsClusterStatusSlice"] = glusterfsClusterStatusSlice
	c.Data["slbSetStatusSlice"] = slbSetStatusSlice

	if len(allErrorMessageSlice) > 0 {
		errorText := fmt.Sprintf("%v", allErrorMessageSlice)
		guimessage.AddDanger(errorText)
	}

	guimessage.OutputMessage(c.Data)
}

func parseComponentStatus(cloudoneJsonMap map[string]interface{}, cloudoneAnalysisJsonMap map[string]interface{}) (*ComponentStatus, []string) {
	componentStatus := &ComponentStatus{}
	componentStatus.CloudoneGUI = true

	errorMessageSlice := make([]string, 0)

	cloudoneStatusJsonMap, ok := cloudoneJsonMap["cloudone"].(map[string]interface{})
	if ok == false {
		errorText := fmt.Sprintf("Fail to parse cloudoneJsonMap[cloudone] %v cloudoneJsonMap %v", cloudoneJsonMap["cloudone"], cloudoneJsonMap)
		errorMessageSlice = append(errorMessageSlice, errorText)
	} else {
		if componentStatus.Storage, ok = cloudoneStatusJsonMap["storage"].(bool); ok == false {
			errorText := fmt.Sprintf("Fail to convert field storage in cloudoneJsonMap[cloudone] %v cloudoneJsonMap %v", cloudoneJsonMap["cloudone"], cloudoneJsonMap)
			errorMessageSlice = append(errorMessageSlice, errorText)
		}
		if componentStatus.Docker, ok = cloudoneStatusJsonMap["docker"].(bool); ok == false {
			errorText := fmt.Sprintf("Fail to convert field docker in cloudoneJsonMap[cloudone] %v cloudoneJsonMap %v", cloudoneJsonMap["cloudone"], cloudoneJsonMap)
			errorMessageSlice = append(errorMessageSlice, errorText)
		}
		if componentStatus.Cloudone, ok = cloudoneStatusJsonMap["restapi"].(bool); ok == false {
			errorText := fmt.Sprintf("Fail to convert field restapi in cloudoneJsonMap[cloudone] %v cloudoneJsonMap %v", cloudoneJsonMap["cloudone"], cloudoneJsonMap)
			errorMessageSlice = append(errorMessageSlice, errorText)
		}
	}

	cloudoneAnalysisStatusJsonMap, ok := cloudoneAnalysisJsonMap["cloudone_analysis"].(map[string]interface{})
	if ok == false {
		errorText := fmt.Sprintf("Fail to parse cloudoneAnalysisJsonMap[cloudone_analysis] %v cloudoneJsonMap %v", cloudoneAnalysisJsonMap["cloudone_analysis"], cloudoneAnalysisJsonMap)
		errorMessageSlice = append(errorMessageSlice, errorText)
	} else {
		if componentStatus.ElasticSearch, ok = cloudoneAnalysisStatusJsonMap["elasticsearch"].(bool); ok == false {
			errorText := fmt.Sprintf("Fail to convert field elasticsearch in cloudoneAnalysisJsonMap[cloudone_analysis] %v cloudoneAnalysisJsonMap %v", cloudoneAnalysisJsonMap["cloudone_analysis"], cloudoneAnalysisJsonMap)
			errorMessageSlice = append(errorMessageSlice, errorText)
		}
		if componentStatus.CloudoneAnalysis, ok = cloudoneAnalysisStatusJsonMap["restapi"].(bool); ok == false {
			errorText := fmt.Sprintf("Fail to convert field restapi in cloudoneAnalysisJsonMap[cloudone_analysis] %v cloudoneAnalysisJsonMap %v", cloudoneAnalysisJsonMap["cloudone_analysis"], cloudoneAnalysisJsonMap)
			errorMessageSlice = append(errorMessageSlice, errorText)
		}
	}

	return componentStatus, errorMessageSlice
}

func parseKubernetesStatusSlice(cloudoneJsonMap map[string]interface{}) ([]KubernetesStatus, []string) {
	kubernetesStatusSlice := make([]KubernetesStatus, 0)

	errorMessageSlice := make([]string, 0)

	kubernetesStatusJsonMap, ok := cloudoneJsonMap["kubernetes"].(map[string]interface{})
	if ok == false {
		errorText := fmt.Sprintf("Fail to parse cloudoneJsonMap[kubernetes] %v cloudoneJsonMap %v", cloudoneJsonMap["kubernetes"], cloudoneJsonMap)
		errorMessageSlice = append(errorMessageSlice, errorText)
	}

	for key, value := range kubernetesStatusJsonMap {
		nodeStatusJsonMap, ok := value.(map[string]interface{})
		if ok == false {
			errorText := fmt.Sprintf("Fail to parse value %v cloudoneJsonMap %v", value, cloudoneJsonMap)
			errorMessageSlice = append(errorMessageSlice, errorText)
		}
		active, ok := nodeStatusJsonMap["active"].(bool)
		if ok == false {
			errorText := fmt.Sprintf("Fail to parse nodeStatusJsonMap[active] %v cloudoneJsonMap %v", nodeStatusJsonMap["active"], cloudoneJsonMap)
			errorMessageSlice = append(errorMessageSlice, errorText)
		}
		if active {
			kubernetesStatus := KubernetesStatus{}
			kubernetesStatus.IP = key
			kubernetesStatus.Active = active

			serviceJsonMap, ok := nodeStatusJsonMap["service"].(map[string]interface{})
			if ok == false {
				errorText := fmt.Sprintf("Fail to parse nodeStatusJsonMap[service] %v cloudoneJsonMap %v", nodeStatusJsonMap["service"], cloudoneJsonMap)
				errorMessageSlice = append(errorMessageSlice, errorText)
			} else {
				if kubernetesStatus.Docker, ok = serviceJsonMap["docker"].(bool); ok == false {
					errorText := fmt.Sprintf("Fail to convert field docker in serviceJsonMap[docker] %v cloudoneJsonMap %v", serviceJsonMap["docker"], cloudoneJsonMap)
					errorMessageSlice = append(errorMessageSlice, errorText)
				}
				if kubernetesStatus.Flannel, ok = serviceJsonMap["flanneld"].(bool); ok == false {
					errorText := fmt.Sprintf("Fail to convert field flanneld in serviceJsonMap[docker] %v cloudoneJsonMap %v", serviceJsonMap["flanneld"], cloudoneJsonMap)
					errorMessageSlice = append(errorMessageSlice, errorText)
				}
				if kubernetesStatus.KubeApiserver, ok = serviceJsonMap["kube-apiserver"].(bool); ok == false {
					errorText := fmt.Sprintf("Fail to convert field kube-apiserver in serviceJsonMap[docker] %v cloudoneJsonMap %v", serviceJsonMap["kube-apiserver"], cloudoneJsonMap)
					errorMessageSlice = append(errorMessageSlice, errorText)
				}
				if kubernetesStatus.KubeControllerManager, ok = serviceJsonMap["kube-controller-manager"].(bool); ok == false {
					errorText := fmt.Sprintf("Fail to convert field kube-controller-manager in serviceJsonMap[docker] %v cloudoneJsonMap %v", serviceJsonMap["kube-controller-manager"], cloudoneJsonMap)
					errorMessageSlice = append(errorMessageSlice, errorText)
				}
				if kubernetesStatus.KubeProxy, ok = serviceJsonMap["kube-proxy"].(bool); ok == false {
					errorText := fmt.Sprintf("Fail to convert field kube-proxy in serviceJsonMap[docker] %v cloudoneJsonMap %v", serviceJsonMap["kube-proxy"], cloudoneJsonMap)
					errorMessageSlice = append(errorMessageSlice, errorText)
				}
				if kubernetesStatus.KubeScheduler, ok = serviceJsonMap["kube-scheduler"].(bool); ok == false {
					errorText := fmt.Sprintf("Fail to convert field kube-scheduler in serviceJsonMap[docker] %v cloudoneJsonMap %v", serviceJsonMap["kube-scheduler"], cloudoneJsonMap)
					errorMessageSlice = append(errorMessageSlice, errorText)
				}
				if kubernetesStatus.Kubelet, ok = serviceJsonMap["kubelet"].(bool); ok == false {
					errorText := fmt.Sprintf("Fail to convert field kubelet in serviceJsonMap[docker] %v cloudoneJsonMap %v", serviceJsonMap["kubelet"], cloudoneJsonMap)
					errorMessageSlice = append(errorMessageSlice, errorText)
				}
			}

			dockerJsonMap, ok := nodeStatusJsonMap["docker"].(map[string]interface{})
			if ok == false {
				errorText := fmt.Sprintf("Fail to parse nodeStatusJsonMap[docker] %v cloudoneJsonMap %v", nodeStatusJsonMap["docker"], cloudoneJsonMap)
				errorMessageSlice = append(errorMessageSlice, errorText)
			}
			dockerIP, dockerIPOk := dockerJsonMap["ip"].(string)

			flannelJsonMap, ok := nodeStatusJsonMap["flannel"].(map[string]interface{})
			if ok == false {
				errorText := fmt.Sprintf("Fail to parse nodeStatusJsonMap[flannel] %v cloudoneJsonMap %v", nodeStatusJsonMap["flannel"], cloudoneJsonMap)
				errorMessageSlice = append(errorMessageSlice, errorText)
			}
			flannelIP, flannelIPOk := flannelJsonMap["ip"].(string)

			dockerIPValid := true
			if dockerIPOk == false || flannelIPOk == false {
				dockerIPValid = false
			} else {
				dockerIPSplits := strings.Split(dockerIP, ".")
				flannelIPSplits := strings.Split(flannelIP, ".")
				if len(dockerIPSplits) == 4 && len(flannelIPSplits) == 4 {
					if dockerIPSplits[0] != flannelIPSplits[0] ||
						dockerIPSplits[1] != flannelIPSplits[1] ||
						dockerIPSplits[2] != flannelIPSplits[2] {
						dockerIPValid = false
					}
				} else {
					dockerIPValid = false
				}
			}

			kubernetesStatus.DockerIP = dockerIP
			kubernetesStatus.FlannelIP = flannelIP
			kubernetesStatus.DockerIPValid = dockerIPValid

			kubernetesStatusSlice = append(kubernetesStatusSlice, kubernetesStatus)
		} else {
			kubernetesStatus := KubernetesStatus{
				key,
				active,
				false,
				false,
				false,
				false,
				false,
				false,
				false,
				"",
				"",
				false,
			}
			kubernetesStatusSlice = append(kubernetesStatusSlice, kubernetesStatus)
		}
	}

	return kubernetesStatusSlice, errorMessageSlice
}

func parseGlusterfsClusterStatusSlice(cloudoneJsonMap map[string]interface{}) ([]GlusterfsClusterStatus, []string) {
	glusterfsClusterStatusSlice := make([]GlusterfsClusterStatus, 0)

	errorMessageSlice := make([]string, 0)

	glusterfsClusterStatusJsonMap, ok := cloudoneJsonMap["glusterfs"].(map[string]interface{})
	if ok == false {
		errorText := fmt.Sprintf("Fail to parse cloudoneJsonMap[glusterfs] %v cloudoneJsonMap %v", cloudoneJsonMap["glusterfs"], cloudoneJsonMap)
		errorMessageSlice = append(errorMessageSlice, errorText)
	}

	for clusterKey, clusterValue := range glusterfsClusterStatusJsonMap {
		glusterfsClusterStatus := GlusterfsClusterStatus{}
		glusterfsClusterStatus.Name = clusterKey
		glusterfsClusterStatus.ServerSlice = make([]GlusterfsServerStatus, 0)

		clusterStatusJsonMap, ok := clusterValue.(map[string]interface{})
		if ok == false {
			errorText := fmt.Sprintf("Fail to parse clusterValue %v cloudoneJsonMap %v", clusterValue, cloudoneJsonMap)
			errorMessageSlice = append(errorMessageSlice, errorText)
		} else {
			for serverKey, serverValue := range clusterStatusJsonMap {
				serverStatusJsonMap, ok := serverValue.(map[string]interface{})
				if ok == false {
					errorText := fmt.Sprintf("Fail to parse serverValue %v cloudoneJsonMap %v", serverValue, cloudoneJsonMap)
					errorMessageSlice = append(errorMessageSlice, errorText)
				} else {
					glusterfsServerStatus := GlusterfsServerStatus{}
					glusterfsServerStatus.IP = serverKey

					glusterfsServerStatus.Active, ok = serverStatusJsonMap["active"].(bool)
					if ok == false {
						errorText := fmt.Sprintf("Fail to parse serverStatusJsonMap[active] %v cloudoneJsonMap %v", serverStatusJsonMap["active"], cloudoneJsonMap)
						errorMessageSlice = append(errorMessageSlice, errorText)
					}

					if glusterfsServerStatus.Active {
						serviceJsonMap, ok := serverStatusJsonMap["service"].(map[string]interface{})
						if ok == false {
							errorText := fmt.Sprintf("Fail to parse serverStatusJsonMap[service] %v cloudoneJsonMap %v", serverStatusJsonMap["service"], cloudoneJsonMap)
							errorMessageSlice = append(errorMessageSlice, errorText)
						} else {
							if glusterfsServerStatus.GlusterfsServer, ok = serviceJsonMap["glusterfs-server"].(bool); ok == false {
								errorText := fmt.Sprintf("Fail to convert field glusterfs-server in serviceJsonMap[glusterfs-server] %v cloudoneJsonMap %v", serviceJsonMap["glusterfs-server"], cloudoneJsonMap)
								errorMessageSlice = append(errorMessageSlice, errorText)
							}
						}
					}

					glusterfsClusterStatus.ServerSlice = append(glusterfsClusterStatus.ServerSlice, glusterfsServerStatus)
				}
			}
		}
		sort.Sort(SortGlusterfsServerStatusByIP(glusterfsClusterStatus.ServerSlice))
		glusterfsClusterStatusSlice = append(glusterfsClusterStatusSlice, glusterfsClusterStatus)
	}

	return glusterfsClusterStatusSlice, errorMessageSlice
}

func parseSLBSetStatusSlice(cloudoneJsonMap map[string]interface{}) ([]SLBSetStatus, []string) {
	slbSetStatusSlice := make([]SLBSetStatus, 0)

	errorMessageSlice := make([]string, 0)

	slbSetStatusJsonMap, ok := cloudoneJsonMap["slb"].(map[string]interface{})
	if ok == false {
		errorText := fmt.Sprintf("Fail to parse cloudoneJsonMap[slb] %v cloudoneJsonMap %v", cloudoneJsonMap["slb"], cloudoneJsonMap)
		errorMessageSlice = append(errorMessageSlice, errorText)
	}

	for clusterKey, clusterValue := range slbSetStatusJsonMap {
		slbSetStatus := SLBSetStatus{}
		slbSetStatus.Name = clusterKey
		slbSetStatus.ServerSlice = make([]SLBServerStatus, 0)

		clusterStatusJsonMap, ok := clusterValue.(map[string]interface{})
		if ok == false {
			errorText := fmt.Sprintf("Fail to parse clusterValue %v cloudoneJsonMap %v", clusterValue, cloudoneJsonMap)
			errorMessageSlice = append(errorMessageSlice, errorText)
		} else {
			for serverKey, serverValue := range clusterStatusJsonMap {
				serverStatusJsonMap, ok := serverValue.(map[string]interface{})
				if ok == false {
					errorText := fmt.Sprintf("Fail to parse serverValue %v cloudoneJsonMap %v", serverValue, cloudoneJsonMap)
					errorMessageSlice = append(errorMessageSlice, errorText)
				} else {
					slbServerStatus := SLBServerStatus{}
					slbServerStatus.IP = serverKey

					slbServerStatus.Active, ok = serverStatusJsonMap["active"].(bool)
					if ok == false {
						errorText := fmt.Sprintf("Fail to parse serverStatusJsonMap[active] %v cloudoneJsonMap %v", serverStatusJsonMap["active"], cloudoneJsonMap)
						errorMessageSlice = append(errorMessageSlice, errorText)
					}

					if slbServerStatus.Active {
						serviceJsonMap, ok := serverStatusJsonMap["service"].(map[string]interface{})
						if ok == false {
							errorText := fmt.Sprintf("Fail to parse serverStatusJsonMap[service] %v cloudoneJsonMap %v", serverStatusJsonMap["service"], cloudoneJsonMap)
							errorMessageSlice = append(errorMessageSlice, errorText)
						} else {
							if slbServerStatus.CloudOneSLB, ok = serviceJsonMap["cloudone_slb"].(bool); ok == false {
								errorText := fmt.Sprintf("Fail to convert field cloudone_slb in serviceJsonMap[cloudone_slb] %v cloudoneJsonMap %v", serviceJsonMap["cloudone_slb"], cloudoneJsonMap)
								errorMessageSlice = append(errorMessageSlice, errorText)
							}
							if slbServerStatus.HAProxy, ok = serviceJsonMap["haproxy"].(bool); ok == false {
								errorText := fmt.Sprintf("Fail to convert field haproxy in serviceJsonMap[haproxy] %v cloudoneJsonMap %v", serviceJsonMap["haproxy"], cloudoneJsonMap)
								errorMessageSlice = append(errorMessageSlice, errorText)
							}
							if slbServerStatus.Keepalived, ok = serviceJsonMap["keepalived"].(bool); ok == false {
								errorText := fmt.Sprintf("Fail to convert field keepalived in serviceJsonMap[keepalived] %v cloudoneJsonMap %v", serviceJsonMap["keepalived"], cloudoneJsonMap)
								errorMessageSlice = append(errorMessageSlice, errorText)
							}
						}

						slbDaemonJsonMap, ok := serverStatusJsonMap["slb_daemon"].(map[string]interface{})
						if ok == false {
							errorText := fmt.Sprintf("Fail to parse serverStatusJsonMap[slb_daemon] %v cloudoneJsonMap %v", serverStatusJsonMap["slb_daemon"], cloudoneJsonMap)
							errorMessageSlice = append(errorMessageSlice, errorText)
						} else {
							if slbServerStatus.LastCommandTime, ok = slbDaemonJsonMap["last_command_created_time"].(string); ok == false {
								errorText := fmt.Sprintf("Fail to convert field last_command_created_time in slbDaemonJsonMap[last_command_created_time] %v cloudoneJsonMap %v", slbDaemonJsonMap["last_command_created_time"], cloudoneJsonMap)
								errorMessageSlice = append(errorMessageSlice, errorText)
							}
						}
					}

					slbSetStatus.ServerSlice = append(slbSetStatus.ServerSlice, slbServerStatus)
				}
			}
		}
		sort.Sort(SortSLBServerStatusByIP(slbSetStatus.ServerSlice))
		slbSetStatusSlice = append(slbSetStatusSlice, slbSetStatus)
	}

	return slbSetStatusSlice, errorMessageSlice
}
