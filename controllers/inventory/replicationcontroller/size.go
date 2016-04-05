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

package replicationcontroller

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/configuration"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/restclient"
	"strconv"
)

type SizeController struct {
	beego.Controller
}

func (c *SizeController) Get() {
	c.TplName = "inventory/replicationcontroller/size.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")

	name := c.GetString("name")
	size := c.GetString("size")
	c.Data["name"] = name
	c.Data["size"] = size

	guimessage.OutputMessage(c.Data)
}

func (c *SizeController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")
	kubeapiHost, kubeapiPort, err := configuration.GetAvailableKubeapiHostAndPort()
	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
		guimessage.RedirectMessage(c)
		c.Ctx.Redirect(302, "/gui/inventory/replicationcontroller/list")
		return
	}

	namespace, _ := c.GetSession("namespace").(string)

	name := c.GetString("name")
	size, _ := c.GetInt("size")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/replicationcontrollers/size/" + namespace + "/" + name + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

	putBodyJsonMap := make(map[string]interface{})
	putBodyJsonMap["Size"] = size

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err = restclient.RequestPut(url, putBodyJsonMap, tokenHeaderMap, true)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		guimessage.AddSuccess("Replication Controller " + name + " is resized")
	}

	c.Ctx.Redirect(302, "/gui/inventory/replicationcontroller/list")

	guimessage.RedirectMessage(c)
}
