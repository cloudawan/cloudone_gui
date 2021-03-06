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

package credential

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
	c.TplName = "system/host/credential/edit.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")

	ip := c.GetString("ip")

	if ip == "" {
		c.Data["actionButtonValue"] = "Create"
		c.Data["pageHeader"] = "Create Host Credential"
		c.Data["createOrUpdate"] = "create"
		c.Data["fieldDisabledChecked"] = ""

		c.Data["sshPort"] = 22
	} else {
		cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
		cloudoneHost := beego.AppConfig.String("cloudoneHost")
		cloudonePort := beego.AppConfig.String("cloudonePort")

		url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
			"/api/v1/hosts/credentials/" + ip

		credential := Credential{}

		tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

		_, err := restclient.RequestGetWithStructure(url, &credential, tokenHeaderMap)

		if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
			return
		}

		if err != nil {
			// Error
			guimessage.AddDanger(guimessagedisplay.GetErrorMessage(err))
			guimessage.OutputMessage(c.Data)
			return
		}

		c.Data["actionButtonValue"] = "Update"
		c.Data["pageHeader"] = "Update Host Credential"
		c.Data["createOrUpdate"] = "update"

		c.Data["ip"] = credential.IP
		c.Data["sshPort"] = credential.SSH.Port
		c.Data["sshUser"] = credential.SSH.User
		c.Data["sshPassword"] = credential.SSH.Password
		if credential.Disabled {
			c.Data["fieldDisabledChecked"] = "checked"
		} else {
			c.Data["fieldDisabledChecked"] = ""
		}

		c.Data["ipFieldReadOnly"] = "readonly"
	}

	guimessage.OutputMessage(c.Data)
}

func (c *EditController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	createOrUpdate := c.GetString("createOrUpdate")
	ip := c.GetString("ip")
	sshPort, _ := c.GetInt("sshPort")
	disabledText := c.GetString("disabled")
	sshUser := c.GetString("sshUser")
	sshPassword := c.GetString("sshPassword")

	disabled := false
	if disabledText == "on" {
		disabled = true
	}

	credential := Credential{
		ip,
		SSH{
			sshPort,
			sshUser,
			sshPassword,
		},
		disabled,
		"",
		"",
	}

	if createOrUpdate == "create" {
		url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
			"/api/v1/hosts/credentials/"

		tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

		_, err := restclient.RequestPostWithStructure(url, credential, nil, tokenHeaderMap)

		if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
			return
		}

		if err != nil {
			// Error
			guimessage.AddDanger(guimessagedisplay.GetErrorMessage(err))
		} else {
			guimessage.AddSuccess("Host credential " + ip + " is created")
		}
	} else {
		url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
			"/api/v1/hosts/credentials/" + ip

		tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

		_, err := restclient.RequestPutWithStructure(url, credential, nil, tokenHeaderMap)

		if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
			return
		}

		if err != nil {
			// Error
			guimessage.AddDanger(guimessagedisplay.GetErrorMessage(err))
		} else {
			guimessage.AddSuccess("Host credential " + ip + " is updated")
		}
	}

	c.Ctx.Redirect(302, "/gui/system/host/credential/list")

	guimessage.RedirectMessage(c)
}
