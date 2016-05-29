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

package clone

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/rbac"
	"github.com/cloudawan/cloudone_utility/restclient"
)

type SelectController struct {
	beego.Controller
}

func (c *SelectController) Get() {
	c.TplName = "deploy/clone/select.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")
	// Authorization for Button
	user, _ := c.GetSession("user").(*rbac.User)

	currentNamespace, _ := c.GetSession("namespace").(string)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort + "/api/v1/namespaces/"

	nameSlice := make([]string, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &nameSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(guimessagedisplay.GetErrorMessage(err))
		guimessage.OutputMessage(c.Data)
		return
	}

	namespaceSlice := make([]string, 0)
	for _, name := range nameSlice {
		if user.HasResource(identity.GetConponentName(), "/namespaces/"+name) {
			namespaceSlice = append(namespaceSlice, name)
		}
	}

	c.Data["currentNamespace"] = currentNamespace
	c.Data["namespaceSlice"] = namespaceSlice

	guimessage.OutputMessage(c.Data)

}

func (c *SelectController) Post() {
	namespace := c.GetString("namespace")
	action := c.GetString("action")

	// Redirect to list
	c.Ctx.Redirect(302, "/gui/deploy/clone/topology?namespace="+namespace+"&action="+action)
}
