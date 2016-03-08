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

package deployclusterapplication

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/utility/configuration"
	"github.com/cloudawan/cloudone_utility/restclient"
	"strconv"
)

type ListController struct {
	beego.Controller
}

type DeployClusterApplication struct {
	Name                           string
	Size                           int
	ServiceName                    string
	ReplicationControllerNameSlice []string
}

// @Title get
// @Description get all cluster application deployments
// @Success 200 {string} []DeployClusterApplication
// @Failure 404 error reason
// @router / [get]
func (c *ListController) Get() {
	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")
	kubeapiHost, kubeapiPort, _ := configuration.GetAvailableKubeapiHostAndPort()

	namespace := c.GetSession("namespace").(string)

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/deployclusterapplications/" + namespace + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

	deployClusterApplicationSlice := make([]DeployClusterApplication, 0)

	_, err := restclient.RequestGetWithStructure(url, &deployClusterApplicationSlice)

	if err != nil {
		// Error
		c.Data["json"] = make(map[string]interface{})
		c.Data["json"].(map[string]interface{})["error"] = err.Error()
		c.Ctx.Output.Status = 404
		c.ServeJSON()
		return
	} else {
		c.Data["json"] = deployClusterApplicationSlice
		c.ServeJSON()
	}
}
