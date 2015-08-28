package thirdparty

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/kubernetes_management_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/kubernetes_management_utility/restclient"
)

type ListController struct {
	beego.Controller
}

type ThirdPartyApplication struct {
	Name        string
	Description string
}

func (c *ListController) Get() {
	c.TplNames = "repository/thirdparty/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	kubernetesManagementProtocol := beego.AppConfig.String("kubernetesManagementProtocol")
	kubernetesManagementHost := beego.AppConfig.String("kubernetesManagementHost")
	kubernetesManagementPort := beego.AppConfig.String("kubernetesManagementPort")

	url := kubernetesManagementProtocol + "://" + kubernetesManagementHost + ":" + kubernetesManagementPort +
		"/api/v1/clusterapplications/"

	thirdPartyApplicationSlice := make([]ThirdPartyApplication, 0)

	_, err := restclient.RequestGetWithStructure(url, &thirdPartyApplicationSlice)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		c.Data["thirdPartyApplicationSlice"] = thirdPartyApplicationSlice
	}

	guimessage.OutputMessage(c.Data)
}
