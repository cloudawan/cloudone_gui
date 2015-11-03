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
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/restclient"
	"strings"
	"time"
)

type EditController struct {
	beego.Controller
}

func (c *EditController) Get() {
	c.TplNames = "notification/notifier/edit.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	kind := c.GetString("kind")
	name := c.GetString("name")

	c.Data["cpuHidden"] = "hidden"
	c.Data["memoryHidden"] = "hidden"

	if kind == "" || name == "" {
		c.Data["actionButtonValue"] = "Create"
		c.Data["pageHeader"] = "Create Notifier"
		c.Data["kind"] = ""
		c.Data["name"] = ""
		c.Data["readonly"] = ""
	} else {
		namespace, _ := c.GetSession("namespace").(string)

		kubernetesManagementProtocol := beego.AppConfig.String("kubernetesManagementProtocol")
		kubernetesManagementHost := beego.AppConfig.String("kubernetesManagementHost")
		kubernetesManagementPort := beego.AppConfig.String("kubernetesManagementPort")

		url := kubernetesManagementProtocol + "://" + kubernetesManagementHost + ":" + kubernetesManagementPort +
			"/api/v1/notifiers/" + namespace + "/" + kind + "/" + name

		replicationControllerNotifier := ReplicationControllerNotifier{}

		_, err := restclient.RequestGetWithStructure(url, &replicationControllerNotifier)

		if err != nil {
			// Error
			guimessage.AddDanger(err.Error())
		} else {

			for _, notifier := range replicationControllerNotifier.NotifierSlice {
				switch notifier.Kind {
				case "email":
					email := strings.TrimSuffix(notifier.Data, ",")
					c.Data["email"] = email
				case "smsNexmo":
					jsonMap := make(map[string]interface{})
					err := json.Unmarshal([]byte(notifier.Data), &jsonMap)
					if err != nil {
						guimessage.AddDanger(err.Error())
					} else {
						smsNexmoSender := jsonMap["Sender"]
						smsNexmoPhone := jsonMap["ReceiverNumberSlice"].([]interface{})[0]
						c.Data["smsNexmoSender"] = smsNexmoSender
						c.Data["smsNexmoPhone"] = smsNexmoPhone
					}
				}
			}

			for _, indicator := range replicationControllerNotifier.IndicatorSlice {
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
			coolDownDurationInSecond := int(replicationControllerNotifier.CoolDownDuration / time.Second)
			c.Data["coolDownDuration"] = coolDownDurationInSecond
			c.Data["readonly"] = "readonly"
		}

		c.Data["kind"] = kind
		c.Data["name"] = name
		c.Data["actionButtonValue"] = "Update"
		c.Data["pageHeader"] = "Update Notifier"
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

	kubernetesManagementProtocol := beego.AppConfig.String("kubernetesManagementProtocol")
	kubernetesManagementHost := beego.AppConfig.String("kubernetesManagementHost")
	kubernetesManagementPort := beego.AppConfig.String("kubernetesManagementPort")
	kubeapiHost := beego.AppConfig.String("kubeapiHost")
	kubeapiPort, _ := beego.AppConfig.Int("kubeapiPort")

	namespace, _ := c.GetSession("namespace").(string)

	kind := c.GetString("kind")
	name := c.GetString("name")
	coolDownDuration, _ := c.GetInt("coolDownDuration")
	email := c.GetString("email")
	smsNexmoSender := c.GetString("smsNexmoSender")
	smsNexmoPhone := c.GetString("smsNexmoPhone")

	notifierSlice := make([]Notifier, 0)
	if email != "" {
		notifierSlice = append(notifierSlice, Notifier{"email", email})
	}
	if smsNexmoSender != "" && smsNexmoPhone != "" {
		notifierSlice = append(notifierSlice, Notifier{
			"smsNexmo",
			`{"Sender":"` + smsNexmoSender + `","ReceiverNumberSlice":["` + smsNexmoPhone + `"]}`,
		})
	}

	replicationControllerNotifier := ReplicationControllerNotifier{
		true, time.Duration(coolDownDuration) * time.Second, 0, kubeapiHost,
		kubeapiPort, namespace, kind, name, notifierSlice, indicatorSlice}

	url := kubernetesManagementProtocol + "://" + kubernetesManagementHost + ":" + kubernetesManagementPort +
		"/api/v1/notifiers/"

	_, err := restclient.RequestPutWithStructure(url, replicationControllerNotifier, nil)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		guimessage.AddSuccess("Email notifier for " + kind + " " + name + " is edited")
	}

	c.Ctx.Redirect(302, "/gui/notification/notifier/")

	guimessage.RedirectMessage(c)
}
