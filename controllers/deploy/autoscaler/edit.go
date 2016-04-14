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

package autoscaler

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/configuration"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/restclient"
	"time"
)

type EditController struct {
	beego.Controller
}

func (c *EditController) Get() {
	c.TplName = "deploy/autoscaler/edit.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")

	kind := c.GetString("kind")
	name := c.GetString("name")

	c.Data["cpuHidden"] = "hidden"
	c.Data["memoryHidden"] = "hidden"

	if kind == "" || name == "" {
		c.Data["actionButtonValue"] = "Create"
		c.Data["pageHeader"] = "Create Autoscaler"
		c.Data["kind"] = ""
		c.Data["name"] = ""
		c.Data["readonly"] = ""
		c.Data["maximumReplica"] = 1
		c.Data["minimumReplica"] = 1
	} else {
		namespace, _ := c.GetSession("namespace").(string)

		cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
		cloudoneHost := beego.AppConfig.String("cloudoneHost")
		cloudonePort := beego.AppConfig.String("cloudonePort")

		url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
			"/api/v1/autoscalers/" + namespace + "/" + kind + "/" + name

		replicationControllerAutoScaler := ReplicationControllerAutoScaler{}

		tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

		_, err := restclient.RequestGetWithStructure(url, &replicationControllerAutoScaler, tokenHeaderMap)

		if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
			return
		}

		if err != nil {
			// Error
			guimessage.AddDanger(err.Error())
		} else {
			c.Data["maximumReplica"] = replicationControllerAutoScaler.MaximumReplica
			c.Data["minimumReplica"] = replicationControllerAutoScaler.MinimumReplica

			for _, indicator := range replicationControllerAutoScaler.IndicatorSlice {
				switch indicator.Type {
				case "cpu":
					c.Data["cpuChecked"] = "checked"
					delete(c.Data, "cpuHidden")

					if indicator.AboveAllOrOne {
						c.Data["cpuAboveAllOrOneChecked"] = "checked"
					}

					c.Data["cpuAbovePercentageOfData"] = int(indicator.AbovePercentageOfData * 100)
					c.Data["cpuAboveThreshold"] = indicator.AboveThreshold / 1000000

					if indicator.BelowAllOrOne {
						c.Data["cpuBelowAllOrOneChecked"] = "checked"
					}

					c.Data["cpuBelowPercentageOfData"] = int(indicator.BelowPercentageOfData * 100)
					c.Data["cpuBelowThreshold"] = indicator.BelowThreshold / 1000000
				case "memory":
					c.Data["memoryChecked"] = "checked"
					delete(c.Data, "memoryHidden")

					if indicator.AboveAllOrOne {
						c.Data["memoryAboveAllOrOneChecked"] = "checked"
					}

					c.Data["memoryAbovePercentageOfData"] = int(indicator.AbovePercentageOfData * 100)
					c.Data["memoryAboveThreshold"] = indicator.AboveThreshold / (1024 * 1024)

					if indicator.BelowAllOrOne {
						c.Data["memoryBelowAllOrOneChecked"] = "checked"
					}

					c.Data["memoryBelowPercentageOfData"] = int(indicator.BelowPercentageOfData * 100)
					c.Data["memoryBelowThreshold"] = indicator.BelowThreshold / (1024 * 1024)
				}
			}
			coolDownDurationInSecond := int(replicationControllerAutoScaler.CoolDownDuration / time.Second)
			c.Data["coolDownDuration"] = coolDownDurationInSecond
			c.Data["readonly"] = "readonly"

			if kind == "selector" {
				c.Data["kindSelectorSelected"] = "selected"
				c.Data["kindReplicationControllerSelected"] = ""
			} else if kind == "replicationController" {
				c.Data["kindSelectorSelected"] = ""
				c.Data["kindReplicationControllerSelected"] = "selected"
			}
		}

		c.Data["kind"] = kind
		c.Data["name"] = name
		c.Data["actionButtonValue"] = "Update"
		c.Data["pageHeader"] = "Update Autoscaler"
	}

	guimessage.OutputMessage(c.Data)
}

func (c *EditController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	indicatorSlice := make([]Indicator, 0)
	cpu := c.GetString("cpu")
	if cpu == "on" {
		cpuAboveAllOrOneText := c.GetString("cpuAboveAllOrOne")
		var cpuAboveAllOrOne bool
		if cpuAboveAllOrOneText == "on" {
			cpuAboveAllOrOne = true
		} else {
			cpuAboveAllOrOne = false
		}
		cpuAbovePercentageOfData, _ := c.GetFloat("cpuAbovePercentageOfData")
		cpuAboveThreshold, _ := c.GetInt64("cpuAboveThreshold")
		cpuBelowAllOrOneText := c.GetString("cpuBelowAllOrOne")
		var cpuBelowAllOrOne bool
		if cpuBelowAllOrOneText == "on" {
			cpuBelowAllOrOne = true
		} else {
			cpuBelowAllOrOne = false
		}
		cpuBelowPercentageOfData, _ := c.GetFloat("cpuBelowPercentageOfData")
		cpuBelowThreshold, _ := c.GetInt64("cpuBelowThreshold")
		indicatorSlice = append(indicatorSlice, Indicator{"cpu",
			cpuAboveAllOrOne, cpuAbovePercentageOfData / 100.0, cpuAboveThreshold * 1000000,
			cpuBelowAllOrOne, cpuBelowPercentageOfData / 100.0, cpuBelowThreshold * 1000000})
	}
	memory := c.GetString("memory")
	if memory == "on" {
		memoryAboveAllOrOneText := c.GetString("memoryAboveAllOrOne")
		var memoryAboveAllOrOne bool
		if memoryAboveAllOrOneText == "on" {
			memoryAboveAllOrOne = true
		} else {
			memoryAboveAllOrOne = false
		}
		memoryAbovePercentageOfData, _ := c.GetFloat("memoryAbovePercentageOfData")
		memoryAboveThreshold, _ := c.GetInt64("memoryAboveThreshold")
		memoryBelowAllOrOneText := c.GetString("memoryBelowAllOrOne")
		var memoryBelowAllOrOne bool
		if memoryBelowAllOrOneText == "on" {
			memoryBelowAllOrOne = true
		} else {
			memoryBelowAllOrOne = false
		}
		memoryBelowPercentageOfData, _ := c.GetFloat("memoryBelowPercentageOfData")
		memoryBelowThreshold, _ := c.GetInt64("memoryBelowThreshold")
		indicatorSlice = append(indicatorSlice, Indicator{"memory",
			memoryAboveAllOrOne, memoryAbovePercentageOfData / 100.0, memoryAboveThreshold * 1024 * 1024,
			memoryBelowAllOrOne, memoryBelowPercentageOfData / 100.0, memoryBelowThreshold * 1024 * 1024})
	}

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")
	kubeapiHost, kubeapiPort, err := configuration.GetAvailableKubeapiHostAndPort()
	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
		guimessage.RedirectMessage(c)
		c.Ctx.Redirect(302, "/gui/deploy/autoscaler/list")
		return
	}

	namespace, _ := c.GetSession("namespace").(string)

	kind := c.GetString("kind")
	name := c.GetString("name")
	coolDownDuration, _ := c.GetInt("coolDownDuration")
	maximumReplica, _ := c.GetInt("maximumReplica")
	minimumReplica, _ := c.GetInt("minimumReplica")

	replicationControllerAutoScaler := ReplicationControllerAutoScaler{
		true, time.Duration(coolDownDuration) * time.Second, 0, kubeapiHost, kubeapiPort, namespace, kind, name,
		maximumReplica, minimumReplica, indicatorSlice, "", ""}

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/autoscalers/"

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err = restclient.RequestPutWithStructure(url, replicationControllerAutoScaler, nil, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		guimessage.AddSuccess("Auto scaler for " + kind + " " + name + " is edited")
	}

	c.Ctx.Redirect(302, "/gui/deploy/autoscaler/list")

	guimessage.RedirectMessage(c)
}
