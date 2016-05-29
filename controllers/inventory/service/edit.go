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
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/restclient"
	"regexp"
)

type EditController struct {
	beego.Controller
}

func (c *EditController) Get() {
	c.TplName = "inventory/service/edit.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")

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

	guimessage.OutputMessage(c.Data)
}

func (c *EditController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	namespace, _ := c.GetSession("namespace").(string)

	name := c.GetString("name")
	selectorName := c.GetString("selectorName")
	//labelName := c.GetString("labelName")
	//portName := c.GetString("portName")
	protocol := c.GetString("protocol")
	port, _ := c.GetInt("port")
	targetPort := c.GetString("targetPort")
	useNodePort := c.GetString("useNodePort")
	autoGeneratedNodePort := c.GetString("autoGeneratedNodePort")
	nodePort, _ := c.GetInt("nodePort")
	sessionAffinity := c.GetString("sessionAffinity")

	// Name need to be a DNS 952 label
	match, _ := regexp.MatchString("^[a-z]{1}[a-z0-9-]{1,23}$", name)
	if match == false {
		guimessage.AddDanger("The name need to be a DNS 952 label ^[a-z]{1}[a-z0-9-]{1,23}$")
		c.Ctx.Redirect(302, "/gui/inventory/service/list")
		guimessage.RedirectMessage(c)
		return
	}

	labelName := selectorName
	portName := selectorName

	if useNodePort == "on" {
		if autoGeneratedNodePort == "on" {
			nodePort = 0
		}
	} else {
		nodePort = -1
	}

	portSlice := make([]ServicePort, 0)
	portSlice = append(portSlice, ServicePort{portName, protocol, port, targetPort, nodePort, "", "", ""})
	selectorMap := make(map[string]interface{})
	selectorMap["name"] = selectorName
	labelMap := make(map[string]interface{})
	labelMap["name"] = labelName

	service := Service{name, namespace, portSlice, selectorMap, "", labelMap, sessionAffinity, "", ""}

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/services/" + namespace

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestPostWithStructure(url, service, nil, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(guimessagedisplay.GetErrorMessage(err))
	} else {
		guimessage.AddSuccess("Service " + name + " is edited")
	}

	c.Ctx.Redirect(302, "/gui/inventory/service/list")

	guimessage.RedirectMessage(c)
}
