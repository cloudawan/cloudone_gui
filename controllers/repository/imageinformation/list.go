package imageinformation

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/kubernetes_management_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/kubernetes_management_utility/restclient"
)

type ListController struct {
	beego.Controller
}

type ImageInformation struct {
	Name           string
	Kind           string
	Description    string
	CurrentVersion string
	BuildParameter map[string]string
}

func (c *ListController) Get() {
	c.TplNames = "repository/imageinformation/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	kubernetesManagementProtocol := beego.AppConfig.String("kubernetesManagementProtocol")
	kubernetesManagementHost := beego.AppConfig.String("kubernetesManagementHost")
	kubernetesManagementPort := beego.AppConfig.String("kubernetesManagementPort")

	url := kubernetesManagementProtocol + "://" + kubernetesManagementHost + ":" + kubernetesManagementPort +
		"/api/v1/imageinformations/"

	imageInformationSlice := make([]ImageInformation, 0)

	returnedImageInformationSlice, err := restclient.RequestGetWithStructure(url, &imageInformationSlice)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		c.Data["imageInformationSlice"] = returnedImageInformationSlice
	}

	guimessage.OutputMessage(c.Data)
}
