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

package notifier

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/kubernetes_management_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/kubernetes_management_utility/restclient"
	"time"
)

type ListController struct {
	beego.Controller
}

type ReplicationControllerNotifier struct {
	Check             bool
	CoolDownDuration  time.Duration
	RemainingCoolDown time.Duration
	KubeapiHost       string
	KubeapiPort       int
	Namespace         string
	Kind              string
	Name              string
	NotifierSlice     []Notifier
	IndicatorSlice    []Indicator
}

type Notifier struct {
	Kind string
	Data string
}

type Indicator struct {
	Type                  string
	AboveAllOrOne         bool
	AbovePercentageOfData float64
	AboveThreshold        int64
	BelowAllOrOne         bool
	BelowPercentageOfData float64
	BelowThreshold        int64
}

func (c *ListController) Get() {
	c.TplNames = "notification/notifier/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	kubernetesManagementProtocol := beego.AppConfig.String("kubernetesManagementProtocol")
	kubernetesManagementHost := beego.AppConfig.String("kubernetesManagementHost")
	kubernetesManagementPort := beego.AppConfig.String("kubernetesManagementPort")

	url := kubernetesManagementProtocol + "://" + kubernetesManagementHost + ":" + kubernetesManagementPort +
		"/api/v1/notifiers/"

	replicationControllerNotifierSlice := make([]ReplicationControllerNotifier, 0)

	returnedReplicationControllerNotifierSlice, err := restclient.RequestGetWithStructure(url, &replicationControllerNotifierSlice)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		c.Data["replicationControllerNotifierSlice"] = returnedReplicationControllerNotifierSlice
	}

	guimessage.OutputMessage(c.Data)
}
