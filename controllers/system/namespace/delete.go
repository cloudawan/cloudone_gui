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
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/restclient"
	"time"
)

type DeleteController struct {
	beego.Controller
}

func (c *DeleteController) Get() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	name := c.GetString("name")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")
	kubeapiHost := beego.AppConfig.String("kubeapiHost")
	kubeapiPort := beego.AppConfig.String("kubeapiPort")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/namespaces/" + name + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + kubeapiPort

	_, err := restclient.RequestDelete(url, nil, true)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		guimessage.AddSuccess("Namespace " + name + " is deleted")

		selectedNamespace := c.GetSession("namespace")
		if selectedNamespace.(string) == name {
			c.SetSession("namespace", "default")
		}
	}
	time.Sleep(1000 * time.Millisecond)
	// Redirect to list
	c.Ctx.Redirect(302, "/gui/system/namespace/")

	guimessage.RedirectMessage(c)
}
