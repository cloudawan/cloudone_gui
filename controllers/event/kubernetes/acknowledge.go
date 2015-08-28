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

package kubernetes

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/kubernetes_management_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/kubernetes_management_utility/restclient"
)

type AcknowledgeController struct {
	beego.Controller
}

func (c *AcknowledgeController) Get() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	namespace := c.GetString("namespace")
	id := c.GetString("id")
	acknowledge := c.GetString("acknowledge")

	kubernetesManagementAnalysisProtocol := beego.AppConfig.String("kubernetesManagementAnalysisProtocol")
	kubernetesManagementAnalysisHost := beego.AppConfig.String("kubernetesManagementAnalysisHost")
	kubernetesManagementAnalysisPort := beego.AppConfig.String("kubernetesManagementAnalysisPort")

	url := kubernetesManagementAnalysisProtocol + "://" + kubernetesManagementAnalysisHost + ":" + kubernetesManagementAnalysisPort +
		"/api/v1/historicalevents/" + namespace + "/" + id + "?acknowledge=" + acknowledge

	jsonMapSlice := make([]map[string]interface{}, 0)
	_, err := restclient.RequestPutWithStructure(url, nil, &jsonMapSlice)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		guimessage.AddSuccess("Acknowledged event")
	}

	if acknowledge == "true" {
		c.Ctx.Redirect(302, "/gui/event/kubernetes/?acknowledge=false")
	} else {
		c.Ctx.Redirect(302, "/gui/event/kubernetes/?acknowledge=true")
	}

	guimessage.RedirectMessage(c)
}
