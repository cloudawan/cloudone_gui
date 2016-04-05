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

type PodLogController struct {
	beego.Controller
}

func (c *PodLogController) Get() {
	c.TplName = "inventory/replicationcontroller/pod_log.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")

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

	namespace := c.GetString("namespace")
	pod := c.GetString("pod")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/podlogs/" + namespace + "/" + pod + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	result, err := restclient.RequestGet(url, tokenHeaderMap, true)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	jsonMap, _ := result.(map[string]interface{})

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		c.Data["logJsonMap"] = jsonMap
	}

	guimessage.OutputMessage(c.Data)
}
