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
}

type ServicePort struct {
	Name       string
	Protocol   string
	Port       string
	TargetPort string
	NodePort   string
}

// @Title get
// @Description get all services
// @Success 200 {string} []Service
// @Failure 404 error reason
// @router / [get]
func (c *ListController) Get() {
	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	kubeapiHost, kubeapiPort, _ := configuration.GetAvailableKubeapiHostAndPort()

	namespace := c.GetSession("namespace").(string)

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/services/" + namespace + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

	serviceSlice := make([]Service, 0)
	_, err := restclient.RequestGetWithStructure(url, &serviceSlice)

	if err != nil {
		// Error
		c.Data["json"] = make(map[string]interface{})
		c.Data["json"].(map[string]interface{})["error"] = err.Error()
		c.Ctx.Output.Status = 404
		c.ServeJson()
		return
	} else {
		c.Data["json"] = serviceSlice
		c.ServeJson()
	}
}
