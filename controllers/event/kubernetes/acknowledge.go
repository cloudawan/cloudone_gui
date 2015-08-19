package kubernetes

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/kubernetes_management_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/kubernetes_management_utility/restclient"
)

type AcknowledgeController struct {
	beego.Controller
}

func (c *AcknowledgeController) Get() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	namespace := c.GetString("namespace")
	id := c.GetString("id")
	acknowledge := c.GetString("acknowledge")

	kubernetesManagementAnalysisProtocol := beego.AppConfig.String("kubernetesManagementAnalysisProtocol")
	kubernetesManagementAnalysisHost := beego.AppConfig.String("kubernetesManagementAnalysisHost")
	kubernetesManagementAnalysisPort := beego.AppConfig.String("kubernetesManagementAnalysisPort")

	url := kubernetesManagementAnalysisProtocol + "://" + kubernetesManagementAnalysisHost + ":" + kubernetesManagementAnalysisPort +
		"/api/v1/historicalevents/" + namespace + "/" + id + "?acknowledge=" + acknowledge

	jsonMapSlice := make([]map[string]interface{}, 0)
	_, err := restclient.RequestPutWithStructure(url, nil, &jsonMapSlice)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		guimessage.AddSuccess("Acknowledged event")
	}

	if acknowledge == "true" {
		c.Ctx.Redirect(302, "/gui/event/kubernetes/?acknowledge=false")
	} else {
		c.Ctx.Redirect(302, "/gui/event/kubernetes/?acknowledge=true")
	}

	guimessage.RedirectMessage(c)
}
