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

package service

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/utility/configuration"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/restclient"
	"strconv"
)

type ListController struct {
	beego.Controller
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
	Port       string
	TargetPort string
	NodePort   string
}

var displayMap map[string]string = map[string]string{
	"kube-dns":              "disabled",
	"kubernetes":            "disabled",
	"private-registry":      "disabled",
	"kubernetes-management": "disabled",
}

func (c *ListController) Get() {
	c.TplName = "inventory/service/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")
	kubeapiHost, kubeapiPort, err := configuration.GetAvailableKubeapiHostAndPort()
	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
		guimessage.OutputMessage(c.Data)
		return
	}

	namespace := c.GetSession("namespace").(string)

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/services/" + namespace + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

	serviceSlice := make([]Service, 0)
	_, err = restclient.RequestGetWithStructure(url, &serviceSlice)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		for i := 0; i < len(serviceSlice); i++ {
			serviceSlice[i].Display = displayMap[serviceSlice[i].Name]
		}
		c.Data["serviceSlice"] = serviceSlice
	}

	guimessage.OutputMessage(c.Data)
}
