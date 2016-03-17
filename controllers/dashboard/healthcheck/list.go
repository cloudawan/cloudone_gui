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
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
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

type GlusterfsStatus struct {
	IP        string
	Active    bool
	Glusterfs bool
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

type SortGlusterfsStatusByIP []GlusterfsStatus

func (sortGlusterfsStatusByIP SortGlusterfsStatusByIP) Len() int {
	return len(sortGlusterfsStatusByIP)
}
func (sortGlusterfsStatusByIP SortGlusterfsStatusByIP) Swap(i, j int) {
	sortGlusterfsStatusByIP[i], sortGlusterfsStatusByIP[j] = sortGlusterfsStatusByIP[j], sortGlusterfsStatusByIP[i]
}
func (sortGlusterfsStatusByIP SortGlusterfsStatusByIP) Less(i, j int) bool {
	return sortGlusterfsStatusByIP[i].IP < sortGlusterfsStatusByIP[j].IP
}

func (c *ListController) Get() {
	c.TplName = "dashboard/healthcheck/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

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
	_, err := restclient.RequestGetWithStructure(url, &cloudoneJsonMap)

	if err != nil {
		// Error
		allErrorMessageSlice = append(allErrorMessageSlice, err.Error())
	}

	url = cloudoneAnalysisProtocol + "://" + cloudoneAnalysisHost + ":" + cloudoneAnalysisPort +
		"/api/v1/healthchecks/"

	cloudoneAnalysisJsonMap := make(map[string]interface{}, 0)
	_, err = restclient.RequestGetWithStructure(url, &cloudoneAnalysisJsonMap)

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

	sort.Sort(SortKubernetesStatusByIP(kubernetesStatusSlice))

	c.Data["componentStatusSlice"] = componentStatusSlice
	c.Data["kubernetesStatusSlice"] = kubernetesStatusSlice

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
				if kubernetesStatus.Docker, _ = serviceJsonMap["docker"].(bool); ok == false {
					errorText := fmt.Sprintf("Fail to convert field docker in serviceJsonMap[docker] %v cloudoneJsonMap %v", serviceJsonMap["docker"], cloudoneJsonMap)
					errorMessageSlice = append(errorMessageSlice, errorText)
				}
				if kubernetesStatus.Flannel, _ = serviceJsonMap["flanneld"].(bool); ok == false {
					errorText := fmt.Sprintf("Fail to convert field flanneld in serviceJsonMap[docker] %v cloudoneJsonMap %v", serviceJsonMap["flanneld"], cloudoneJsonMap)
					errorMessageSlice = append(errorMessageSlice, errorText)
				}
				if kubernetesStatus.KubeApiserver, _ = serviceJsonMap["kube-apiserver"].(bool); ok == false {
					errorText := fmt.Sprintf("Fail to convert field kube-apiserver in serviceJsonMap[docker] %v cloudoneJsonMap %v", serviceJsonMap["kube-apiserver"], cloudoneJsonMap)
					errorMessageSlice = append(errorMessageSlice, errorText)
				}
				if kubernetesStatus.KubeControllerManager, _ = serviceJsonMap["kube-controller-manager"].(bool); ok == false {
					errorText := fmt.Sprintf("Fail to convert field kube-controller-manager in serviceJsonMap[docker] %v cloudoneJsonMap %v", serviceJsonMap["kube-controller-manager"], cloudoneJsonMap)
					errorMessageSlice = append(errorMessageSlice, errorText)
				}
				if kubernetesStatus.KubeProxy, _ = serviceJsonMap["kube-proxy"].(bool); ok == false {
					errorText := fmt.Sprintf("Fail to convert field kube-proxy in serviceJsonMap[docker] %v cloudoneJsonMap %v", serviceJsonMap["kube-proxy"], cloudoneJsonMap)
					errorMessageSlice = append(errorMessageSlice, errorText)
				}
				if kubernetesStatus.KubeScheduler, _ = serviceJsonMap["kube-scheduler"].(bool); ok == false {
					errorText := fmt.Sprintf("Fail to convert field kube-scheduler in serviceJsonMap[docker] %v cloudoneJsonMap %v", serviceJsonMap["kube-scheduler"], cloudoneJsonMap)
					errorMessageSlice = append(errorMessageSlice, errorText)
				}
				if kubernetesStatus.Kubelet, _ = serviceJsonMap["kubelet"].(bool); ok == false {
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
