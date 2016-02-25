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
	"github.com/cloudawan/cloudone_utility/restclient"
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
}

type SMSNexmo struct {
	Name      string
	Url       string
	APIKey    string
	APISecret string
}

// @Title get
// @Description get the notifier
// @Param kind path string true "The type of target notifier configured for"
// @Param name path string true "The name of target notifier configured for"
// @Success 200 {object} guirestapi.notification.notifier.ReplicationControllerNotifier
// @Failure 404 error reason
// @router /:kind/:name [get]
func (c *EditController) Get() {
	kind := c.GetString("kind")
	name := c.GetString("name")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/notifiers/emailserversmtp"
	emailServerSMTPSlice := make([]EmailServerSMTP, 0)
	_, err := restclient.RequestGetWithStructure(url, &emailServerSMTPSlice)
	if err != nil {
		// Error
		c.Data["json"] = make(map[string]interface{})
		c.Data["json"].(map[string]interface{})["error"] = err.Error()
		c.Ctx.Output.Status = 404
		c.ServeJson()
		return
	}

	url = cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/notifiers/smsnexmo"
	smsNexmoSlice := make([]SMSNexmo, 0)
	_, err = restclient.RequestGetWithStructure(url, &smsNexmoSlice)
	if err != nil {
		// Error
		c.Data["json"] = make(map[string]interface{})
		c.Data["json"].(map[string]interface{})["error"] = err.Error()
		c.Ctx.Output.Status = 404
		c.ServeJson()
		return
	}

	if len(emailServerSMTPSlice) == 0 {
		// Error
		c.Data["json"] = make(map[string]interface{})
		c.Data["json"].(map[string]interface{})["error"] = "No Email server is configured"
		c.Ctx.Output.Status = 404
		c.ServeJson()
		return
	}

	if len(smsNexmoSlice) == 0 {
		// Error
		c.Data["json"] = make(map[string]interface{})
		c.Data["json"].(map[string]interface{})["error"] = "No SMS server is configured"
		c.Ctx.Output.Status = 404
		c.ServeJson()
		return
	}

	namespace, _ := c.GetSession("namespace").(string)

	url = cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/notifiers/" + namespace + "/" + kind + "/" + name

	replicationControllerNotifier := ReplicationControllerNotifier{}

	_, err = restclient.RequestGetWithStructure(url, &replicationControllerNotifier)

	if err != nil {
		// Error
		c.Data["json"] = make(map[string]interface{})
		c.Data["json"].(map[string]interface{})["error"] = err.Error()
		c.Ctx.Output.Status = 404
		c.ServeJson()
		return
	} else {
		c.Data["json"] = make(map[string]interface{})
		c.Data["json"].(map[string]interface{})["emailServerSMTPSlice"] = emailServerSMTPSlice
		c.Data["json"].(map[string]interface{})["smsNexmoSlice"] = smsNexmoSlice
		c.Data["json"].(map[string]interface{})["replicationControllerNotifier"] = replicationControllerNotifier
		c.ServeJson()
	}
}

// @Title update
// @Description update the notifier
// @Param body body guirestapi.notification.notifier.ReplicationControllerNotifier true "body for notifier"
// @Success 200 {string} {}
// @Failure 404 error reason
// @router / [put]
func (c *EditController) Put() {
	inputBody := c.Ctx.Input.RequestBody
	replicationControllerNotifier := ReplicationControllerNotifier{}
	err := json.Unmarshal(inputBody, &replicationControllerNotifier)
	if err != nil {
		// Error
		c.Data["json"] = make(map[string]interface{})
		c.Data["json"].(map[string]interface{})["error"] = err.Error()
		c.Ctx.Output.Status = 404
		c.ServeJson()
		return
	}

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")
	kubeapiHost, kubeapiPort, _ := configuration.GetAvailableKubeapiHostAndPort()

	namespace, _ := c.GetSession("namespace").(string)

	replicationControllerNotifier.KubeapiHost = kubeapiHost
	replicationControllerNotifier.KubeapiPort = kubeapiPort
	replicationControllerNotifier.Namespace = namespace

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/notifiers/"

	_, err = restclient.RequestPutWithStructure(url, replicationControllerNotifier, nil)

	if err != nil {
		// Error
		c.Data["json"] = make(map[string]interface{})
		c.Data["json"].(map[string]interface{})["error"] = err.Error()
		c.Ctx.Output.Status = 404
		c.ServeJson()
		return
	} else {
		c.Data["json"] = make(map[string]interface{})
		c.ServeJson()
	}
}
