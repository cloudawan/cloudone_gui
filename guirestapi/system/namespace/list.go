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
	"github.com/cloudawan/cloudone_utility/restclient"
	"strconv"
)

type ListController struct {
	beego.Controller
}

type Namespace struct {
	Name     string
	Selected bool
}

// @Title get
// @Description get all namespaces
// @Success 200 {string} []Namespace
// @Failure 404 error reason
// @router / [get]
func (c *ListController) Get() {
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
		c.Data["json"] = make(map[string]interface{})
		c.Data["json"].(map[string]interface{})["error"] = err.Error()
		c.Ctx.Output.Status = 404
		c.ServeJSON()
		return
	} else {
		selectedNamespace := c.GetSession("namespace")

		namespaceSlice := make([]Namespace, 0)
		for _, name := range nameSlice {
			namespace := Namespace{name, false}
			if name == selectedNamespace {
				namespace.Selected = true
			}

			namespaceSlice = append(namespaceSlice, namespace)
		}

		c.Data["json"] = namespaceSlice
		c.ServeJSON()
	}
}
