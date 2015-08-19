package imageinformation

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/kubernetes_management_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/kubernetes_management_utility/restclient"
)

type UpgradeController struct {
	beego.Controller
}

type DeployUpgradeInput struct {
	ImageInformationName string
	Description          string
}

func (c *UpgradeController) Get() {
	c.TplNames = "repository/imageinformation/upgrade.html"
	name := c.GetString("name")
	c.Data["name"] = name
}

func (c *UpgradeController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	kubernetesManagementProtocol := beego.AppConfig.String("kubernetesManagementProtocol")
	kubernetesManagementHost := beego.AppConfig.String("kubernetesManagementHost")
	kubernetesManagementPort := beego.AppConfig.String("kubernetesManagementPort")

	imageInformationName := c.GetString("name")
	description := c.GetString("description")

	url := kubernetesManagementProtocol + "://" + kubernetesManagementHost + ":" + kubernetesManagementPort +
		"/api/v1/imageinformations/upgrade/"

	deployUpgradeInput := DeployUpgradeInput{imageInformationName, description}

	_, err := restclient.RequestPutWithStructure(url, deployUpgradeInput, nil)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		guimessage.AddSuccess(imageInformationName + " upgrade success")
	}

	// Redirect to list
	c.Ctx.Redirect(302, "/gui/repository/imageinformation/")

	guimessage.RedirectMessage(c)
}
