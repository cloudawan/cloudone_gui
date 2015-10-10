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
			"/api/v1/clusterapplications/" + name
		cluster := Cluster{}
		_, err := restclient.RequestGetWithStructure(url, &cluster)
		if err != nil {
			guimessage.AddDanger("Fail to get with error" + err.Error())
			// Redirect to list
			c.Ctx.Redirect(302, "/gui/repository/thirdparty/")

			guimessage.RedirectMessage(c)
		} else {
			environmentByteSlice, err := json.MarshalIndent(cluster.Environment, "", "    ")
			if err != nil {
				guimessage.AddDanger("Fail to get with error" + err.Error())
				// Redirect to list
				c.Ctx.Redirect(302, "/gui/repository/thirdparty/")

				guimessage.RedirectMessage(c)
			}

			c.Data["actionButtonValue"] = "Update"
			c.Data["pageHeader"] = "Update third party service"
			c.Data["name"] = cluster.Name
			c.Data["description"] = cluster.Description
			c.Data["replicationControllerJson"] = string(cluster.ReplicationControllerJson)
			c.Data["serviceJson"] = cluster.ServiceJson
			c.Data["environment"] = string(environmentByteSlice)
			c.Data["scriptType"] = cluster.ScriptType
			c.Data["scriptContent"] = cluster.ScriptContent

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
	replicationControllerJson := c.GetString("replicationControllerJson")
	if replicationControllerJson == "" {
		replicationControllerJson = "{}"
	}
	serviceJson := c.GetString("serviceJson")
	if serviceJson == "" {
		serviceJson = "{}"
	}
	environmentText := c.GetString("environment")
	if environmentText == "" {
		environmentText = "{}"
	}
	scriptType := c.GetString("scriptType")
	scriptContent := c.GetString("scriptContent")

	environmentJsonMap := make(map[string]string)
	err := json.Unmarshal([]byte(environmentText), &environmentJsonMap)
	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
		c.Ctx.Redirect(302, "/gui/repository/thirdparty/")
		guimessage.RedirectMessage(c)
		return
	}

	cluster := Cluster{
		name,
		description,
		replicationControllerJson,
		serviceJson,
		environmentJsonMap,
		scriptType,
		scriptContent,
	}

	url := kubernetesManagementProtocol + "://" + kubernetesManagementHost + ":" + kubernetesManagementPort +
		"/api/v1/clusterapplications/"

	_, err = restclient.RequestPost(url, cluster, true)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		guimessage.AddSuccess("Third party application " + name + " is edited")
	}

	c.Ctx.Redirect(302, "/gui/repository/thirdparty/")

	guimessage.RedirectMessage(c)
}
