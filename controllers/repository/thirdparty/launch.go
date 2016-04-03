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

package thirdparty

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/configuration"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/restclient"
	"strconv"
	"strings"
)

type Cluster struct {
	Name                      string
	Description               string
	ReplicationControllerJson string
	ServiceJson               string
	Environment               map[string]string
	ScriptType                string
	ScriptContent             string
}

type LaunchController struct {
	beego.Controller
}

func (c *LaunchController) Get() {
	c.TplName = "repository/thirdparty/launch.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	name := c.GetString("name")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/clusterapplications/" + name
	cluster := Cluster{}

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &cluster, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		guimessage.AddDanger("Fail to get with error" + err.Error())
		// Redirect to list
		c.Ctx.Redirect(302, "/gui/repository/thirdparty/")

		guimessage.RedirectMessage(c)
	} else {
		c.Data["actionButtonValue"] = "Launch"
		c.Data["pageHeader"] = "Launch third party service"
		c.Data["thirdPartyApplicationName"] = name
		c.Data["environment"] = cluster.Environment

		guimessage.OutputMessage(c.Data)
	}
}

func (c *LaunchController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")
	kubeapiHost, kubeapiPort, err := configuration.GetAvailableKubeapiHostAndPort()
	if err != nil {
		// Error
		errorJsonMap := make(map[string]interface{})
		errorJsonMap["error"] = err.Error()
		c.Data["json"] = errorJsonMap
		c.ServeJSON()
		return
	}

	namespace, _ := c.GetSession("namespace").(string)
	name := c.GetString("name")
	size := c.GetString("size")

	keySlice := make([]string, 0)
	inputMap := c.Input()
	if inputMap != nil {
		for key, _ := range inputMap {
			// Ignore the non environment field
			if key != "name" && key != "size" {
				keySlice = append(keySlice, key)
			}
		}
	}

	environmentSlice := make([]interface{}, 0)
	for _, key := range keySlice {
		value := c.GetString(key)
		if len(value) > 0 {
			environmentMap := make(map[string]string)
			environmentMap["name"] = key
			environmentMap["value"] = value
			environmentSlice = append(environmentSlice, environmentMap)
		}
	}

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/clusterapplications/launch/" + namespace + "/" + name +
		"?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort) + "&size=" + size
	jsonMap := make(map[string]interface{})

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err = restclient.RequestPostWithStructure(url, environmentSlice, &jsonMap, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		errorMessage, _ := jsonMap["Error"].(string)
		if strings.HasPrefix(errorMessage, "Replication controller already exists") {
			guimessage.AddDanger("Replication controller " + name + " already exists")
		} else {
			guimessage.AddDanger(err.Error())
		}
	} else {
		guimessage.AddSuccess("Cluster application " + name + " is launched")
	}

	// Redirect to list
	c.Ctx.Redirect(302, "/gui/deploy/deployclusterapplication/")

	guimessage.RedirectMessage(c)
}
