package deploybluegreen

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/kubernetes_management_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/kubernetes_management_utility/restclient"
)

type ListController struct {
	beego.Controller
}

type DeployBlueGreen struct {
	ImageInformation string
	Namespace        string
	NodePort         int
	Description      string
	SessionAffinity  string
}

func (c *ListController) Get() {
	c.TplNames = "deploy/deploybluegreen/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	kubernetesManagementProtocol := beego.AppConfig.String("kubernetesManagementProtocol")
	kubernetesManagementHost := beego.AppConfig.String("kubernetesManagementHost")
	kubernetesManagementPort := beego.AppConfig.String("kubernetesManagementPort")

	url := kubernetesManagementProtocol + "://" + kubernetesManagementHost + ":" + kubernetesManagementPort +
		"/api/v1/deploybluegreens/"

	deployBlueGreenSlice := make([]DeployBlueGreen, 0)

	_, err := restclient.RequestGetWithStructure(url, &deployBlueGreenSlice)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		c.Data["deployBlueGreenSlice"] = deployBlueGreenSlice
	}

	guimessage.OutputMessage(c.Data)
}
