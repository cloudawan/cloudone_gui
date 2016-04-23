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
	"github.com/cloudawan/cloudone_gui/controllers/deploy/autoscaler"
	"github.com/cloudawan/cloudone_gui/controllers/deploy/deploy"
	"github.com/cloudawan/cloudone_gui/controllers/deploy/deployclusterapplication"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/notification/notifier"
	"github.com/cloudawan/cloudone_gui/controllers/utility/configuration"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/restclient"
	"strconv"
	"time"
)

type DeleteController struct {
	beego.Controller
}

func (c *DeleteController) Get() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	name := c.GetString("name")

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

	// Delete deploy
	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/deploys/" + name

	deployInformationSlice := make([]deploy.DeployInformation, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err = restclient.RequestGetWithStructure(url, &deployInformationSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		for _, deployInformation := range deployInformationSlice {
			url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
				"/api/v1/deploys/" + name + "/" + deployInformation.ImageInformationName + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

			_, err := restclient.RequestDelete(url, nil, tokenHeaderMap, true)

			if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
				return
			}

			if err != nil {
				guimessage.AddDanger(err.Error())
			}
		}
	}

	// Delete third party service
	url = cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/deployclusterapplications/" + name

	deployClusterApplicationSlice := make([]deployclusterapplication.DeployClusterApplication, 0)
	_, err = restclient.RequestGetWithStructure(url, &deployClusterApplicationSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		for _, deployClusterApplication := range deployClusterApplicationSlice {
			url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
				"/api/v1/deployclusterapplications/" + name + "/" + deployClusterApplication.Name + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

			_, err := restclient.RequestDelete(url, nil, tokenHeaderMap, true)

			if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
				return
			}

			if err != nil {
				guimessage.AddDanger(err.Error())
			}
		}
	}

	// Delete namespace
	url = cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/namespaces/" + name + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

	_, err = restclient.RequestDelete(url, nil, tokenHeaderMap, true)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		guimessage.AddSuccess("Namespace " + name + " is deleted")

		selectedNamespace := c.GetSession("namespace")
		if selectedNamespace.(string) == name {
			c.SetSession("namespace", "default")
		}
	}

	// Delete autoscaler
	url = cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/autoscalers/"

	replicationControllerAutoScalerSlice := make([]autoscaler.ReplicationControllerAutoScaler, 0)
	_, err = restclient.RequestGetWithStructure(url, &replicationControllerAutoScalerSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		for _, replicationControllerAutoScaler := range replicationControllerAutoScalerSlice {
			if replicationControllerAutoScaler.Namespace == name {
				url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
					"/api/v1/autoscalers/" + name + "/" + replicationControllerAutoScaler.Kind + "/" + replicationControllerAutoScaler.Name

				_, err := restclient.RequestDelete(url, nil, tokenHeaderMap, true)

				if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
					return
				}

				if err != nil {
					guimessage.AddDanger(err.Error())
				}
			}
		}
	}

	// Delete notifier
	url = cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/notifiers/"

	replicationControllerNotifierSlice := make([]notifier.ReplicationControllerNotifier, 0)
	_, err = restclient.RequestGetWithStructure(url, &replicationControllerNotifierSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		for _, replicationControllerNotifier := range replicationControllerNotifierSlice {
			if replicationControllerNotifier.Namespace == name {
				url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
					"/api/v1/notifiers/" + name + "/" + replicationControllerNotifier.Kind + "/" + replicationControllerNotifier.Name

				_, err := restclient.RequestDelete(url, nil, tokenHeaderMap, true)

				if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
					return
				}

				if err != nil {
					guimessage.AddDanger(err.Error())
				}
			}
		}
	}

	// Since the name deletion is asynchronized, wait for some time
	time.Sleep(time.Millisecond * 500)

	// Redirect to list
	c.Ctx.Redirect(302, "/gui/system/namespace/list")

	guimessage.RedirectMessage(c)
}
