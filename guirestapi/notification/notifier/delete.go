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
	"github.com/cloudawan/cloudone_utility/restclient"
)

type DeleteController struct {
	beego.Controller
}

// @Title delete
// @Description delete the notifier
// @Param kind path string true "The type of target notifier configured for"
// @Param name path string true "The name of target notifier configured for"
// @Success 200 {string} {}
// @Failure 404 error reason
// @router /:kind/:name [delete]
func (c *DeleteController) Delete() {
	kind := c.GetString("kind")
	name := c.GetString("name")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	namespace, _ := c.GetSession("namespace").(string)

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/notifiers/" + namespace + "/" + kind + "/" + name
	_, err := restclient.RequestDelete(url, nil, true)

	if err != nil {
		// Error
		c.Data["json"] = make(map[string]interface{})
		c.Data["json"].(map[string]interface{})["error"] = err.Error()
		c.Ctx.Output.Status = 404
		c.ServeJSON()
		return
	} else {
		c.Data["json"] = make(map[string]interface{})
		c.ServeJSON()
	}
}
