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
	"github.com/cloudawan/kubernetes_management_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/kubernetes_management_utility/restclient"
)

type EditController struct {
	beego.Controller
}

func (c *EditController) Get() {
	c.TplNames = "inventory/service/edit.html"

	service := c.GetString("service")
	if service == "" {
		c.Data["actionButtonValue"] = "Create"
		c.Data["pageHeader"] = "Create Service"
		c.Data["serviceName"] = ""
	} else {
		c.Data["actionButtonValue"] = "Update"
		c.Data["pageHeader"] = "Update Service"
		c.Data["serviceName"] = service
	}
}

func (c *EditController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	kubernetesManagementProtocol := beego.AppConfig.String("kubernetesManagementProtocol")
	kubernetesManagementHost := beego.AppConfig.String("kubernetesManagementHost")
	kubernetesManagementPort := beego.AppConfig.String("kubernetesManagementPort")
	kubeapiHost := beego.AppConfig.String("kubeapiHost")
	kubeapiPort := beego.AppConfig.String("kubeapiPort")

	namespace, _ := c.GetSession("namespace").(string)

	name := c.GetString("name")
	selectorName := c.GetString("selectorName")
	//labelName := c.GetString("labelName")
	//portName := c.GetString("portName")
	protocol := c.GetString("protocol")
	port := c.GetString("port")
	targetPort := c.GetString("targetPort")
	nodePort := c.GetString("nodePort")
	sessionAffinity := c.GetString("sessionAffinity")

	labelName := selectorName
	portName := selectorName

	portSlice := make([]ServicePort, 0)
	portSlice = append(portSlice, ServicePort{portName, protocol, port, targetPort, nodePort})
	selectorMap := make(map[string]interface{})
	selectorMap["name"] = selectorName
	labelMap := make(map[string]interface{})
	labelMap["name"] = labelName

	service := Service{name, namespace, portSlice, selectorMap, "", labelMap, sessionAffinity, ""}

	url := kubernetesManagementProtocol + "://" + kubernetesManagementHost + ":" + kubernetesManagementPort +
		"/api/v1/services/" + namespace + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + kubeapiPort

	_, err := restclient.RequestPostWithStructure(url, service, nil)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		guimessage.AddSuccess("Service " + name + " is edited")
	}

	c.Ctx.Redirect(302, "/gui/inventory/service/")

	guimessage.RedirectMessage(c)
}
