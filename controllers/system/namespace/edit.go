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

package namespace

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
	c.TplName = "system/namespace/edit.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")

	namespace := c.GetString("namespace")
	if namespace == "" {
		c.Data["actionButtonValue"] = "Create"
		c.Data["pageHeader"] = "Create Namespace"
		c.Data["namespace"] = ""
	} else {
		c.Data["actionButtonValue"] = "Update"
		c.Data["pageHeader"] = "Update Namespace"
		c.Data["namespace"] = namespace
	}

	guimessage.OutputMessage(c.Data)
}

func (c *EditController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	name := c.GetString("name")

	namespace := Namespace{name, false, false, "", "", "", ""}

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort + "/api/v1/namespaces"

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestPostWithStructure(url, namespace, nil, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(guimessagedisplay.GetErrorMessage(err))
	} else {
		guimessage.AddSuccess("Namespace " + name + " is edited")
	}

	c.Ctx.Redirect(302, "/gui/system/namespace/list")

	guimessage.RedirectMessage(c)
}
