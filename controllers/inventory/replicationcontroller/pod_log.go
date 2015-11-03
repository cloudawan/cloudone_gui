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
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/restclient"
)

type PodLogController struct {
	beego.Controller
}

func (c *PodLogController) Get() {
	c.TplNames = "inventory/replicationcontroller/pod_log.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	kubernetesManagementProtocol := beego.AppConfig.String("kubernetesManagementProtocol")
	kubernetesManagementHost := beego.AppConfig.String("kubernetesManagementHost")
	kubernetesManagementPort := beego.AppConfig.String("kubernetesManagementPort")
	kubeapiHost := beego.AppConfig.String("kubeapiHost")
	kubeapiPort := beego.AppConfig.String("kubeapiPort")
	namespace := c.GetString("namespace")
	pod := c.GetString("pod")

	url := kubernetesManagementProtocol + "://" + kubernetesManagementHost + ":" + kubernetesManagementPort +
		"/api/v1/podlogs/" + namespace + "/" + pod + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + kubeapiPort

	result, err := restclient.RequestGet(url, true)
	jsonMap, _ := result.(map[string]interface{})

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		c.Data["log"] = jsonMap["Log"]
	}

	guimessage.OutputMessage(c.Data)
}
