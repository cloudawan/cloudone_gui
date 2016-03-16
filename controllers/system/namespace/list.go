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
	"github.com/cloudawan/cloudone_gui/controllers/utility/configuration"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/restclient"
	"sort"
	"strconv"
)

type ListController struct {
	beego.Controller
}

type Namespace struct {
	Name     string
	Selected bool
	Display  string
}

var displayMap map[string]string = map[string]string{
	"default":     "disabled",
	"kube-system": "disabled",
}

type ByNamespace []Namespace

func (b ByNamespace) Len() int           { return len(b) }
func (b ByNamespace) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByNamespace) Less(i, j int) bool { return b[i].Name < b[j].Name }

func (c *ListController) Get() {
	c.TplName = "system/namespace/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")
	kubeapiHost, kubeapiPort, _ := configuration.GetAvailableKubeapiHostAndPort()

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/namespaces/" + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

	nameSlice := make([]string, 0)
	_, err := restclient.RequestGetWithStructure(url, &nameSlice)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		selectedNamespace := c.GetSession("namespace")

		namespaceSlice := make([]Namespace, 0)
		for _, name := range nameSlice {

			namespace := Namespace{name, false, ""}
			if name == selectedNamespace {
				namespace.Selected = true
			}
			namespace.Display = displayMap[name]

			namespaceSlice = append(namespaceSlice, namespace)
		}

		sort.Sort(ByNamespace(namespaceSlice))
		c.Data["namespaceSlice"] = namespaceSlice
	}

	guimessage.OutputMessage(c.Data)
}
