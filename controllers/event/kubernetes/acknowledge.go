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
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/restclient"
)

type AcknowledgeController struct {
	beego.Controller
}

func (c *AcknowledgeController) Get() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	namespace := c.GetString("namespace")
	id := c.GetString("id")
	acknowledge := c.GetString("acknowledge")

	cloudoneAnalysisProtocol := beego.AppConfig.String("cloudoneAnalysisProtocol")
	cloudoneAnalysisHost := beego.AppConfig.String("cloudoneAnalysisHost")
	cloudoneAnalysisPort := beego.AppConfig.String("cloudoneAnalysisPort")

	url := cloudoneAnalysisProtocol + "://" + cloudoneAnalysisHost + ":" + cloudoneAnalysisPort +
		"/api/v1/historicalevents/" + namespace + "/" + id + "?acknowledge=" + acknowledge

	jsonMapSlice := make([]map[string]interface{}, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestPutWithStructure(url, nil, &jsonMapSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(guimessagedisplay.GetErrorMessage(err))
	} else {
		guimessage.AddSuccess("Acknowledged event")
	}

	if acknowledge == "true" {
		c.Ctx.Redirect(302, "/gui/event/kubernetes/list?acknowledge=false")
	} else {
		c.Ctx.Redirect(302, "/gui/event/kubernetes/list?acknowledge=true")
	}

	guimessage.RedirectMessage(c)
}
