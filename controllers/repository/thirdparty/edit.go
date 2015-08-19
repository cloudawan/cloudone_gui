package thirdparty

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/cloudawan/kubernetes_management_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/kubernetes_management_utility/restclient"
)

type EditController struct {
	beego.Controller
}

func (c *EditController) Get() {
	c.TplNames = "repository/thirdparty/edit.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	kubernetesManagementProtocol := beego.AppConfig.String("kubernetesManagementProtocol")
	kubernetesManagementHost := beego.AppConfig.String("kubernetesManagementHost")
	kubernetesManagementPort := beego.AppConfig.String("kubernetesManagementPort")

	name := c.GetString("name")
	if name == "" {
		c.Data["actionButtonValue"] = "Create"
		c.Data["pageHeader"] = "Create third party service"
		c.Data["name"] = ""

		guimessage.OutputMessage(c.Data)
	} else {
		url := kubernetesManagementProtocol + "://" + kubernetesManagementHost + ":" + kubernetesManagementPort +
			"/api/v1/statelessapplications/" + name
		statelessSerializable := StatelessSerializable{}
		_, err := restclient.RequestGetWithStructure(url, &statelessSerializable)
		if err != nil {
			guimessage.AddDanger("Fail to get with error" + err.Error())
			// Redirect to list
			c.Ctx.Redirect(302, "/gui/repository/thirdparty/")

			guimessage.RedirectMessage(c)
		} else {
			replicationControllerByteSlice, err := json.MarshalIndent(statelessSerializable.ReplicationControllerJson, "", "    ")
			if err != nil {
				guimessage.AddDanger("Fail to get with error" + err.Error())
				// Redirect to list
				c.Ctx.Redirect(302, "/gui/repository/thirdparty/")

				guimessage.RedirectMessage(c)
			}

			serviceByteSlice, err := json.MarshalIndent(statelessSerializable.ServiceJson, "", "    ")
			if err != nil {
				guimessage.AddDanger("Fail to get with error" + err.Error())
				// Redirect to list
				c.Ctx.Redirect(302, "/gui/repository/thirdparty/")

				guimessage.RedirectMessage(c)
			}

			environmentByteSlice, err := json.MarshalIndent(statelessSerializable.Environment, "", "    ")
			if err != nil {
				guimessage.AddDanger("Fail to get with error" + err.Error())
				// Redirect to list
				c.Ctx.Redirect(302, "/gui/repository/thirdparty/")

				guimessage.RedirectMessage(c)
			}

			c.Data["actionButtonValue"] = "Update"
			c.Data["pageHeader"] = "Update third party service"
			c.Data["name"] = statelessSerializable.Name
			c.Data["description"] = statelessSerializable.Description
			c.Data["replicationControllerJson"] = string(replicationControllerByteSlice)
			c.Data["serviceJson"] = string(serviceByteSlice)
			c.Data["environment"] = string(environmentByteSlice)

			guimessage.OutputMessage(c.Data)
		}
	}
}

func (c *EditController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	kubernetesManagementProtocol := beego.AppConfig.String("kubernetesManagementProtocol")
	kubernetesManagementHost := beego.AppConfig.String("kubernetesManagementHost")
	kubernetesManagementPort := beego.AppConfig.String("kubernetesManagementPort")

	name := c.GetString("name")
	description := c.GetString("description")
	replicationControllerJsonText := c.GetString("replicationControllerJson")
	serviceJsonText := c.GetString("serviceJson")
	environmentText := c.GetString("environment")
	if environmentText == "" {
		environmentText = "{}"
	}

	replicationControllerJsonMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(replicationControllerJsonText), &replicationControllerJsonMap)
	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
		c.Ctx.Redirect(302, "/gui/repository/thirdparty/")
		guimessage.RedirectMessage(c)
		return
	}

	serviceJsonMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(serviceJsonText), &serviceJsonMap)
	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
		c.Ctx.Redirect(302, "/gui/repository/thirdparty/")
		guimessage.RedirectMessage(c)
		return
	}

	environmentJsonMap := make(map[string]string)
	err = json.Unmarshal([]byte(environmentText), &environmentJsonMap)
	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
		c.Ctx.Redirect(302, "/gui/repository/thirdparty/")
		guimessage.RedirectMessage(c)
		return
	}

	statelessApplicationJsonMap := make(map[string]interface{})
	statelessApplicationJsonMap["name"] = name
	statelessApplicationJsonMap["description"] = description
	statelessApplicationJsonMap["replicationControllerJson"] = replicationControllerJsonMap
	statelessApplicationJsonMap["serviceJson"] = serviceJsonMap
	statelessApplicationJsonMap["environment"] = environmentJsonMap

	url := kubernetesManagementProtocol + "://" + kubernetesManagementHost + ":" + kubernetesManagementPort +
		"/api/v1/statelessapplications/"

	_, err = restclient.RequestPost(url, statelessApplicationJsonMap, true)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		guimessage.AddSuccess("Third party application " + name + " is edited")
	}

	c.Ctx.Redirect(302, "/gui/repository/thirdparty/")

	guimessage.RedirectMessage(c)
}
