package deploy

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/kubernetes_management_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/kubernetes_management_utility/restclient"
)

type ListController struct {
	beego.Controller
}

type DeployInformation struct {
	Namespace            string
	ImageInformationName string
	CurrentVersion       string
	Description          string
}

func (c *ListController) Get() {
	c.TplNames = "deploy/deploy/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	kubernetesManagementProtocol := beego.AppConfig.String("kubernetesManagementProtocol")
	kubernetesManagementHost := beego.AppConfig.String("kubernetesManagementHost")
	kubernetesManagementPort := beego.AppConfig.String("kubernetesManagementPort")

	namespaces, _ := c.GetSession("namespace").(string)

	url := kubernetesManagementProtocol + "://" + kubernetesManagementHost + ":" + kubernetesManagementPort +
		"/api/v1/deploys/"

	deployInformationSlice := make([]DeployInformation, 0)

	_, err := restclient.RequestGetWithStructure(url, &deployInformationSlice)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		// Only show those belonging to this namespace
		filteredDeployInformationSlice := make([]DeployInformation, 0)
		for _, deployInformation := range deployInformationSlice {
			if deployInformation.Namespace == namespaces {
				filteredDeployInformationSlice = append(filteredDeployInformationSlice, deployInformation)
			}
		}

		c.Data["deployInformationSlice"] = filteredDeployInformationSlice
	}

	guimessage.OutputMessage(c.Data)
}
