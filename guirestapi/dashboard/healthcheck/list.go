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
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
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
	Cassandra        bool
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
	KubeScherduler        bool
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

func (sortKubernetesStatusByIP *SortKubernetesStatusByIP) Len() int {
	return len(*sortKubernetesStatusByIP)
}
func (sortKubernetesStatusByIP *SortKubernetesStatusByIP) Swap(i, j int) {
	(*sortKubernetesStatusByIP)[i], (*sortKubernetesStatusByIP)[j] = (*sortKubernetesStatusByIP)[j], (*sortKubernetesStatusByIP)[i]
}
func (sortKubernetesStatusByIP *SortKubernetesStatusByIP) Less(i, j int) bool {
	return (*sortKubernetesStatusByIP)[i].IP < (*sortKubernetesStatusByIP)[j].IP
}

// @Title Get
// @Description get health check
// @Success 200 (string) {}
// @router / [get]
func (c *ListController) Get() {
	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")
	cloudoneAnalysisProtocol := beego.AppConfig.String("cloudoneAnalysisProtocol")
	cloudoneAnalysisHost := beego.AppConfig.String("cloudoneAnalysisHost")
	cloudoneAnalysisPort := beego.AppConfig.String("cloudoneAnalysisPort")

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
		c.Data["json"] = make(map[string]interface{})
		c.Data["json"].(map[string]interface{})["error"] = err.Error()
		c.Ctx.Output.Status = 404
		c.ServeJSON()
		return
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
		c.Data["json"] = make(map[string]interface{})
		c.Data["json"].(map[string]interface{})["error"] = err.Error()
		c.Ctx.Output.Status = 404
		c.ServeJSON()
		return
	}

	componentStatus, err := parseComponentStatus(cloudoneJsonMap, cloudoneAnalysisJsonMap)
	if err != nil {
		// Error
		c.Data["json"] = make(map[string]interface{})
		c.Data["json"].(map[string]interface{})["error"] = err.Error()
		c.Ctx.Output.Status = 404
		c.ServeJSON()
		return
	}
	componentStatusSlice := make([]ComponentStatus, 0)
	componentStatusSlice = append(componentStatusSlice, *componentStatus)

	kubernetesStatusSlice, err := parseKubernetesStatusSlice(cloudoneJsonMap)
	if err != nil {
		// Error
		c.Data["json"] = make(map[string]interface{})
		c.Data["json"].(map[string]interface{})["error"] = err.Error()
		c.Ctx.Output.Status = 404
		c.ServeJSON()
		return
	}
	sortKubernetesStatusByIP := SortKubernetesStatusByIP(kubernetesStatusSlice)
	sort.Sort(&sortKubernetesStatusByIP)

	c.Data["json"] = make(map[string]interface{})
	c.Data["json"].(map[string]interface{})["componentStatusSlice"] = componentStatusSlice
	c.Data["json"].(map[string]interface{})["kubernetesStatusSlice"] = kubernetesStatusSlice
	c.ServeJSON()
}

func parseComponentStatus(cloudoneJsonMap map[string]interface{}, cloudoneAnalysisJsonMap map[string]interface{}) (*ComponentStatus, error) {
	cloudoneStatusJsonMap, ok := cloudoneJsonMap["cloudone"].(map[string]interface{})
	if ok == false {
		errorText := fmt.Sprintf("Fail to parse cloudoneJsonMap[cloudone] %v cloudoneJsonMap %v", cloudoneJsonMap["cloudone"], cloudoneJsonMap)
		return nil, errors.New(errorText)
	}
	storage, ok := cloudoneStatusJsonMap["storage"].(bool)
	docker, ok := cloudoneStatusJsonMap["docker"].(bool)
	cloudone, ok := cloudoneStatusJsonMap["restapi"].(bool)

	cloudoneAnalysisStatusJsonMap, ok := cloudoneAnalysisJsonMap["cloudone_analysis"].(map[string]interface{})
	if ok == false {
		errorText := fmt.Sprintf("Fail to parse cloudoneAnalysisJsonMap[cloudone_analysis] %v cloudoneJsonMap %v", cloudoneAnalysisJsonMap["cloudone_analysis"], cloudoneAnalysisJsonMap)
		return nil, errors.New(errorText)
	}
	elasticSearch, ok := cloudoneAnalysisStatusJsonMap["elasticsearch"].(bool)
	cloudoneAnalysis, ok := cloudoneAnalysisStatusJsonMap["restapi"].(bool)

	componentStatus := &ComponentStatus{
		cloudone,
		cloudoneAnalysis,
		true,
		storage,
		elasticSearch,
		docker,
	}
	return componentStatus, nil
}

func parseKubernetesStatusSlice(cloudoneJsonMap map[string]interface{}) ([]KubernetesStatus, error) {
	kubernetesStatusJsonMap, ok := cloudoneJsonMap["kubernetes"].(map[string]interface{})
	if ok == false {
		errorText := fmt.Sprintf("Fail to parse cloudoneJsonMap[kubernetes] %v cloudoneJsonMap %v", cloudoneJsonMap["kubernetes"], cloudoneJsonMap)
		return nil, errors.New(errorText)
	}

	kubernetesStatusSlice := make([]KubernetesStatus, 0)
	for key, value := range kubernetesStatusJsonMap {
		nodeStatusJsonMap, ok := value.(map[string]interface{})
		if ok == false {
			errorText := fmt.Sprintf("Fail to parse value %v cloudoneJsonMap %v", value, cloudoneJsonMap)
			return nil, errors.New(errorText)
		}
		active, ok := nodeStatusJsonMap["active"].(bool)
		if ok == false {
			errorText := fmt.Sprintf("Fail to parse nodeStatusJsonMap[active] %v cloudoneJsonMap %v", nodeStatusJsonMap["active"], cloudoneJsonMap)
			return nil, errors.New(errorText)
		}
		if active {
			serviceJsonMap, ok := nodeStatusJsonMap["service"].(map[string]interface{})
			if ok == false {
				errorText := fmt.Sprintf("Fail to parse nodeStatusJsonMap[service] %v cloudoneJsonMap %v", nodeStatusJsonMap["service"], cloudoneJsonMap)
				return nil, errors.New(errorText)
			}
			docker, _ := serviceJsonMap["docker"].(bool)
			flanneld, _ := serviceJsonMap["flanneld"].(bool)
			kubeApiserver, _ := serviceJsonMap["kube-apiserver"].(bool)
			kubeControllerManager, _ := serviceJsonMap["kube-controller-manager"].(bool)
			kubeProxy, _ := serviceJsonMap["kube-proxy"].(bool)
			kubeScheduler, _ := serviceJsonMap["kube-scheduler"].(bool)
			kubelet, _ := serviceJsonMap["kubelet"].(bool)

			dockerJsonMap, ok := nodeStatusJsonMap["docker"].(map[string]interface{})
			if ok == false {
				errorText := fmt.Sprintf("Fail to parse nodeStatusJsonMap[docker] %v cloudoneJsonMap %v", nodeStatusJsonMap["docker"], cloudoneJsonMap)
				return nil, errors.New(errorText)
			}
			dockerIP, dockerIPOk := dockerJsonMap["ip"].(string)

			flannelJsonMap, ok := nodeStatusJsonMap["flannel"].(map[string]interface{})
			if ok == false {
				errorText := fmt.Sprintf("Fail to parse nodeStatusJsonMap[flannel] %v cloudoneJsonMap %v", nodeStatusJsonMap["flannel"], cloudoneJsonMap)
				return nil, errors.New(errorText)
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

			kubernetesStatus := KubernetesStatus{
				key,
				active,
				docker,
				flanneld,
				kubeProxy,
				kubelet,
				kubeApiserver,
				kubeScheduler,
				kubeControllerManager,
				dockerIP,
				flannelIP,
				dockerIPValid,
			}
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
	return kubernetesStatusSlice, nil
}
