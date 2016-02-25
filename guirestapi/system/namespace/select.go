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
)

type SelectController struct {
	beego.Controller
}

// @Title select
// @Description select the current namespace
// @Param name path string true "The name of namespace"
// @Success 200 {string} {}
// @Failure 404 error reason
// @router /select/:name [put]
func (c *SelectController) Put() {
	name := c.GetString(":name")

	c.SetSession("namespace", name)

	c.Data["json"] = make(map[string]interface{})
	c.ServeJson()
}
