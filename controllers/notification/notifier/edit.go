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
	"github.com/cloudawan/cloudone_gui/controllers/utility/configuration"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/restclient"
	"strings"
	"time"
)

type EditController struct {
	beego.Controller
}

type EmailServerSMTP struct {
	Name     string
	Account  string
	Password string
	Host     string
	Port     int
	Selected string
}

type SMSNexmo struct {
	Name      string
	Url       string
	APIKey    string
	APISecret string
	Selected  string
}

func (c *EditController) Get() {
	c.TplNames = "notification/notifier/edit.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	kind := c.GetString("kind")
	name := c.GetString("name")

	c.Data["cpuHidden"] = "hidden"
	c.Data["memoryHidden"] = "hidden"

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/notifiers/emailserversmtp"
	emailServerSMTPSlice := make([]EmailServerSMTP, 0)
	_, err := restclient.RequestGetWithStructure(url, &emailServerSMTPSlice)
	if err != nil {
		guimessage.AddDanger(err.Error())
	}

	url = cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/notifiers/smsnexmo"
	smsNexmoSlice := make([]SMSNexmo, 0)
	_, err = restclient.RequestGetWithStructure(url, &smsNexmoSlice)
	if err != nil {
		guimessage.AddDanger(err.Error())
	}

	c.Data["emailServerSMTPSlice"] = emailServerSMTPSlice
	c.Data["smsNexmoSlice"] = smsNexmoSlice

	if len(emailServerSMTPSlice) == 0 {
		guimessage.AddDanger("No Email server is configured")
	}

	if len(smsNexmoSlice) == 0 {
		guimessage.AddDanger("No SMS server is configured")
	}

	if kind == "" || name == "" {
		c.Data["actionButtonValue"] = "Create"
		c.Data["pageHeader"] = "Create Notifier"
		c.Data["kind"] = ""
		c.Data["name"] = ""
		c.Data["readonly"] = ""
		c.Data["selectorSelected"] = ""
		c.Data["replicationControllerSelected"] = ""
	} else {
		namespace, _ := c.GetSession("namespace").(string)

		url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
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
					notifierEmail := NotifierEmail{}
					err := json.Unmarshal([]byte(notifier.Data), &notifierEmail)
					if err != nil {
						guimessage.AddDanger(err.Error())
					} else {
						receiverAccountList := ""
						length := len(notifierEmail.ReceiverAccountSlice)
						for i := 0; i < length; i++ {
							if i == length-1 {
								receiverAccountList += notifierEmail.ReceiverAccountSlice[i]
							} else {
								receiverAccountList += notifierEmail.ReceiverAccountSlice[i] + ", "
							}
						}
						c.Data["email"] = receiverAccountList
						c.Data["emailServerName"] = notifierEmail.Destination

						for i := 0; i < len(emailServerSMTPSlice); i++ {
							if emailServerSMTPSlice[i].Name == notifierEmail.Destination {
								emailServerSMTPSlice[i].Selected = "selected"
							}
						}
					}
				case "smsNexmo":
					notifierSMSNexmo := NotifierSMSNexmo{}
					if err != nil {
						guimessage.AddDanger(err.Error())
					} else {
						receiverNumberList := ""
						length := len(notifierSMSNexmo.ReceiverNumberSlice)
						for i := 0; i < length; i++ {
							if i == length-1 {
								receiverNumberList += notifierSMSNexmo.ReceiverNumberSlice[i]
							} else {
								receiverNumberList += notifierSMSNexmo.ReceiverNumberSlice[i] + ", "
							}
						}
						c.Data["smsNexmoSender"] = notifierSMSNexmo.Sender
						c.Data["smsNexmoPhone"] = receiverNumberList
						c.Data["smsNexmoName"] = notifierSMSNexmo.Destination

						for i := 0; i < len(smsNexmoSlice); i++ {
							if smsNexmoSlice[i].Name == notifierSMSNexmo.Destination {
								smsNexmoSlice[i].Selected = "selected"
							}
						}
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
		c.Data["selectorSelected"] = ""
		c.Data["replicationControllerSelected"] = ""
		switch kind {
		case "selector":
			c.Data["selectorSelected"] = "selected"
		case "replicationController":
			c.Data["replicationControllerSelected"] = "selected"
		}
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
	kubeapiHost, kubeapiPort, _ := configuration.GetAvailableKubeapiHostAndPort()

	namespace, _ := c.GetSession("namespace").(string)

	kind := c.GetString("kind")
	name := c.GetString("name")
	coolDownDuration, _ := c.GetInt("coolDownDuration")
	emailField := c.GetString("email")
	emailServerName := c.GetString("emailServerName")
	smsNexmoSender := c.GetString("smsNexmoSender")
	smsNexmoPhoneField := c.GetString("smsNexmoPhone")
	smsNexmoName := c.GetString("smsNexmoName")

	if len(emailServerName) == 0 {
		guimessage.AddDanger("Email server configuration name can't be empty")
		guimessage.OutputMessage(c.Data)
		return
	}

	if len(smsNexmoName) == 0 {
		guimessage.AddDanger("SMS Nexom configuration name can't be empty")
		guimessage.OutputMessage(c.Data)
		return
	}

	notifierSlice := make([]Notifier, 0)
	if emailField != "" {
		emailSlice := make([]string, 0)
		for _, email := range strings.Split(emailField, ",") {
			value := strings.TrimSpace(email)
			if len(value) > 0 {
				emailSlice = append(emailSlice, value)
			}
		}
		notifierEmail := NotifierEmail{
			emailServerName,
			emailSlice,
		}
		byteSlice, err := json.Marshal(notifierEmail)
		if err != nil {
			guimessage.AddDanger(err.Error())
			guimessage.OutputMessage(c.Data)
			return
		}

		notifierSlice = append(notifierSlice, Notifier{"email", string(byteSlice)})
	}
	if smsNexmoSender != "" && smsNexmoPhoneField != "" {
		smsNexmoPhoneSlice := make([]string, 0)
		for _, smsNexmoPhone := range strings.Split(smsNexmoPhoneField, ",") {
			value := strings.TrimSpace(smsNexmoPhone)
			if len(value) > 0 {
				smsNexmoPhoneSlice = append(smsNexmoPhoneSlice, value)
			}
		}
		notifierSMSNexmo := NotifierSMSNexmo{
			smsNexmoName,
			smsNexmoSender,
			smsNexmoPhoneSlice,
		}
		byteSlice, err := json.Marshal(notifierSMSNexmo)
		if err != nil {
			guimessage.AddDanger(err.Error())
			guimessage.OutputMessage(c.Data)
			return
		}

		notifierSlice = append(notifierSlice, Notifier{"smsNexmo", string(byteSlice)})
	}

	replicationControllerNotifier := ReplicationControllerNotifier{
		true, time.Duration(coolDownDuration) * time.Second, 0, kubeapiHost,
		kubeapiPort, namespace, kind, name, notifierSlice, indicatorSlice}

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
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
