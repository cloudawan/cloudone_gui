package imagerecord

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/kubernetes_management_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/kubernetes_management_utility/restclient"
)

type DeleteController struct {
	beego.Controller
}

func (c *DeleteController) Get() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	kubernetesManagementProtocol := beego.AppConfig.String("kubernetesManagementProtocol")
	kubernetesManagementHost := beego.AppConfig.String("kubernetesManagementHost")
	kubernetesManagementPort := beego.AppConfig.String("kubernetesManagementPort")

	imageInformationName := c.GetString("name")
	imageRecordVersion := c.GetString("version")

	url := kubernetesManagementProtocol + "://" + kubernetesManagementHost + ":" + kubernetesManagementPort +
		"/api/v1/imagerecords/" + imageInformationName + "/" + imageRecordVersion
	_, err := restclient.RequestDelete(url, nil, true)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		guimessage.AddSuccess("Image record " + imageRecordVersion + " belonging to " + imageInformationName + " is deleted")
	}

	// Redirect to list
	c.Ctx.Redirect(302, "/gui/repository/imagerecord/?name="+imageInformationName)

	guimessage.RedirectMessage(c)
}
