package replicationcontroller

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/kubernetes_management_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/kubernetes_management_utility/restclient"
)

type PodLogController struct {
	beego.Controller
}

func (c *PodLogController) Get() {
	c.TplNames = "inventory/replicationcontroller/pod_log.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	kubernetesManagementProtocol := beego.AppConfig.String("kubernetesManagementProtocol")
	kubernetesManagementHost := beego.AppConfig.String("kubernetesManagementHost")
	kubernetesManagementPort := beego.AppConfig.String("kubernetesManagementPort")
	kubeapiHost := beego.AppConfig.String("kubeapiHost")
	kubeapiPort := beego.AppConfig.String("kubeapiPort")
	namespace := c.GetString("namespace")
	pod := c.GetString("pod")

	url := kubernetesManagementProtocol + "://" + kubernetesManagementHost + ":" + kubernetesManagementPort +
		"/api/v1/podlogs/" + namespace + "/" + pod + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + kubeapiPort

	result, err := restclient.RequestGet(url, true)
	jsonMap, _ := result.(map[string]interface{})

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		c.Data["log"] = jsonMap["Log"]
	}

	guimessage.OutputMessage(c.Data)
}
