package imageinformation

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/cloudawan/kubernetes_management_utility/restclient"
)

type UpdateController struct {
	beego.Controller
}

type DeployUpgradeInput struct {
	ImageInformationName string
	Description          string
}

func (c *UpdateController) Post() {

	kubernetesManagementProtocol := beego.AppConfig.String("kubernetesManagementProtocol")
	kubernetesManagementHost := beego.AppConfig.String("kubernetesManagementHost")
	kubernetesManagementPort := beego.AppConfig.String("kubernetesManagementPort")

	deployUpgradeInput := DeployUpgradeInput{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &deployUpgradeInput)
	if err != nil {
		// Error
		errorJsonMap := make(map[string]interface{})
		errorJsonMap["error"] = err.Error()
		c.Data["json"] = errorJsonMap
		c.ServeJson()
		c.Abort("401")
		return
	}

	url := kubernetesManagementProtocol + "://" + kubernetesManagementHost + ":" + kubernetesManagementPort +
		"/api/v1/imageinformations/upgrade/"

	_, err = restclient.RequestPutWithStructure(url, deployUpgradeInput, nil)
	if err != nil {
		// Error
		errorJsonMap := make(map[string]interface{})
		errorJsonMap["error"] = err.Error()
		c.Data["json"] = errorJsonMap
		c.ServeJson()
		c.Abort("401")
		return
	}

	jsonMap := make(map[string]interface{})
	c.Data["json"] = jsonMap
	c.ServeJson()
}
