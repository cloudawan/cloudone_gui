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

type ClusterLaunch struct {
	Size                              int
	EnvironmentSlice                  []interface{}
	ReplicationControllerExtraJsonMap map[string]interface{}
}

type Region struct {
	Name           string
	LocationTagged bool
	ZoneSlice      []Zone
}

type Zone struct {
	Name           string
	LocationTagged bool
	NodeSlice      []Node
}

type Node struct {
	Name     string
	Address  string
	Capacity Capacity
}

type Capacity struct {
	Cpu    string
	Memory string
}

type LaunchController struct {
	beego.Controller
}

func (c *LaunchController) Get() {
	c.TplName = "repository/thirdparty/launch.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")

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
		guimessage.AddDanger("Fail to get data with error" + err.Error())
		// Redirect to list
		c.Ctx.Redirect(302, "/gui/repository/thirdparty/list")

		guimessage.RedirectMessage(c)
		return
	}

	kubeapiHost, kubeapiPort, err := configuration.GetAvailableKubeapiHostAndPort()
	if err != nil {
		// Error
		guimessage.AddDanger("No availabe host and port with error " + err.Error())
		// Redirect to list
		c.Ctx.Redirect(302, "/gui/repository/thirdparty/list")

		guimessage.RedirectMessage(c)
		return
	}

	url = cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/nodes/topology?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

	regionSlice := make([]Region, 0)

	_, err = restclient.RequestGetWithStructure(url, &regionSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		guimessage.AddDanger("Fail to get node topology with error" + err.Error())
		// Redirect to list
		c.Ctx.Redirect(302, "/gui/repository/thirdparty/list")

		guimessage.RedirectMessage(c)
		return
	}

	filteredRegionSlice := make([]Region, 0)
	for _, region := range regionSlice {
		if region.LocationTagged {
			filteredRegionSlice = append(filteredRegionSlice, region)
		}
	}

	if cluster.Environment != nil {
		namespace, _ := c.GetSession("namespace").(string)

		// Try to set the known common parameter
		for key, _ := range cluster.Environment {
			if key == "SERVICE_NAME" {
				cluster.Environment[key] = name
			}
			if key == "NAMESPACE" {
				cluster.Environment[key] = namespace
			}
		}
	}

	c.Data["actionButtonValue"] = "Launch"
	c.Data["pageHeader"] = "Launch third party service"
	c.Data["thirdPartyApplicationName"] = name
	c.Data["environment"] = cluster.Environment
	c.Data["regionSlice"] = filteredRegionSlice

	guimessage.OutputMessage(c.Data)
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
	size, _ := c.GetInt("size")

	region := c.GetString("region")
	zone := c.GetString("zone")

	if region == "Any" {
		region = ""
	}
	if zone == "Any" {
		zone = ""
	}

	keySlice := make([]string, 0)
	inputMap := c.Input()
	if inputMap != nil {
		for key, _ := range inputMap {
			// Ignore the non environment field
			if strings.HasPrefix(key, "environment_") {
				keySlice = append(keySlice, key)
			}
		}
	}

	environmentSlice := make([]interface{}, 0)
	for _, key := range keySlice {
		value := c.GetString(key)
		if len(value) > 0 {
			environmentMap := make(map[string]string)
			environmentMap["name"] = key[len("environment_"):]
			environmentMap["value"] = value
			environmentSlice = append(environmentSlice, environmentMap)
		}
	}

	extraJsonMap := make(map[string]interface{})
	if len(region) > 0 {
		extraJsonMap["spec"] = make(map[string]interface{})
		extraJsonMap["spec"].(map[string]interface{})["template"] = make(map[string]interface{})
		extraJsonMap["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"] = make(map[string]interface{})
		extraJsonMap["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["nodeSelector"] = make(map[string]interface{})

		extraJsonMap["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["nodeSelector"].(map[string]interface{})["region"] = region
		if len(zone) > 0 {
			extraJsonMap["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["nodeSelector"].(map[string]interface{})["zone"] = zone
		}
	} else {
		extraJsonMap = nil
	}

	clusterLaunch := ClusterLaunch{
		size,
		environmentSlice,
		extraJsonMap,
	}

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/clusterapplications/launch/" + namespace + "/" + name +
		"?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)
	jsonMap := make(map[string]interface{})

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err = restclient.RequestPostWithStructure(url, clusterLaunch, &jsonMap, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		errorMessage, _ := jsonMap["Error"].(string)
		if errorMessage == "The cluster application already exists" {
			guimessage.AddDanger("Cluster application " + name + " already exists")
		} else {
			guimessage.AddDanger(err.Error())
		}
	} else {
		guimessage.AddSuccess("Cluster application " + name + " is launched")
	}

	// Redirect to list
	c.Ctx.Redirect(302, "/gui/deploy/deployclusterapplication/list")

	guimessage.RedirectMessage(c)
}
