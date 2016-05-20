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

package github

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/utility/configuration"
	"github.com/cloudawan/cloudone_gui/controllers/utility/limit"
	"github.com/cloudawan/cloudone_utility/restclient"
	"strconv"
)

type GithubPost struct {
	User             string
	ImageInformation string
	Signature        string
	Payload          string
}

type PushController struct {
	beego.Controller
}

func (c *PushController) Post() {
	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")
	kubeapiHost, kubeapiPort, err := configuration.GetAvailableKubeapiHostAndPort()
	if err != nil {
		errorJsonMap := make(map[string]interface{})
		errorJsonMap["error"] = err.Error()
		c.Data["json"] = errorJsonMap
		c.Ctx.Output.Status = 401
		c.ServeJSON()
		return
	}

	user := c.GetString("user")
	imageInformation := c.GetString("imageInformation")
	signature := c.Ctx.Input.Header("X-Hub-Signature")
	payload := string(c.Ctx.Input.CopyBody(limit.InputPostBodyMaximum))

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/webhooks/github/?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

	githubPost := GithubPost{
		user,
		imageInformation,
		signature,
		payload,
	}

	errorJsonMap := make(map[string]interface{})

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err = restclient.RequestPostWithStructure(url, githubPost, &errorJsonMap, tokenHeaderMap)

	if err != nil {
		// Error
		errorJsonMap := make(map[string]interface{})
		errorJsonMap["error"] = err.Error()
		c.Data["json"] = errorJsonMap
		c.Ctx.Output.Status = 401
		c.ServeJSON()
		return
	}

	jsonMap := make(map[string]interface{})
	c.Data["json"] = jsonMap
	c.ServeJSON()
}
