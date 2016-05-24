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

package server

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/restclient"
)

type EditController struct {
	beego.Controller
}

func (c *EditController) Get() {
	c.TplName = "system/privateregistry/server/edit.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")

	serverName := c.GetString("serverName")

	if serverName == "" {
		c.Data["actionButtonValue"] = "Create"
		c.Data["pageHeader"] = "Create Server Configuration"
		c.Data["createOrUpdate"] = "create"
	} else {
		cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
		cloudoneHost := beego.AppConfig.String("cloudoneHost")
		cloudonePort := beego.AppConfig.String("cloudonePort")

		url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
			"/api/v1/privateregistries/servers/" + serverName

		privateRegistry := PrivateRegistry{}

		tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

		_, err := restclient.RequestGetWithStructure(url, &privateRegistry, tokenHeaderMap)

		if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
			return
		}

		if err != nil {
			// Error
			guimessage.AddDanger("Fail to get the server configuration with error " + err.Error())
			guimessage.OutputMessage(c.Data)
			return
		}

		c.Data["actionButtonValue"] = "Update"
		c.Data["pageHeader"] = "Update Server Configuration"
		c.Data["createOrUpdate"] = "update"

		c.Data["name"] = privateRegistry.Name
		c.Data["host"] = privateRegistry.Host
		c.Data["port"] = privateRegistry.Port

		c.Data["nameFieldReadOnly"] = "readonly"
	}

	guimessage.OutputMessage(c.Data)
}

func (c *EditController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	name := c.GetString("name")
	host := c.GetString("host")
	port, _ := c.GetInt("port")
	createOrUpdate := c.GetString("createOrUpdate")

	privateRegistry := PrivateRegistry{
		name,
		host,
		port,
		"",
		"",
		"",
	}

	if createOrUpdate == "create" {
		url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
			"/api/v1/privateregistries/servers/"

		tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

		_, err := restclient.RequestPostWithStructure(url, privateRegistry, nil, tokenHeaderMap)

		if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
			return
		}

		if err != nil {
			// Error
			guimessage.AddDanger(err.Error())
		} else {
			guimessage.AddSuccess("Private register server configuration " + name + " is created")
		}
	} else {
		url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
			"/api/v1/privateregistries/servers/" + name

		tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

		_, err := restclient.RequestPutWithStructure(url, privateRegistry, nil, tokenHeaderMap)

		if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
			return
		}

		if err != nil {
			// Error
			guimessage.AddDanger(err.Error())
		} else {
			guimessage.AddSuccess("Private register server configuration " + name + " is updated")
		}
	}

	c.Ctx.Redirect(302, "/gui/system/privateregistry/server/list")

	guimessage.RedirectMessage(c)
}
