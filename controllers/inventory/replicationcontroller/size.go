package replicationcontroller

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/kubernetes_management_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/kubernetes_management_utility/restclient"
)

type SizeController struct {
	beego.Controller
}

func (c *SizeController) Get() {
	c.TplNames = "inventory/replicationcontroller/size.html"

	name := c.GetString("name")
	size := c.GetString("size")
	c.Data["name"] = name
	c.Data["size"] = size
}

func (c *SizeController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	kubernetesManagementProtocol := beego.AppConfig.String("kubernetesManagementProtocol")
	kubernetesManagementHost := beego.AppConfig.String("kubernetesManagementHost")
	kubernetesManagementPort := beego.AppConfig.String("kubernetesManagementPort")
	kubeapiHost := beego.AppConfig.String("kubeapiHost")
	kubeapiPort := beego.AppConfig.String("kubeapiPort")

	namespace, _ := c.GetSession("namespace").(string)

	name := c.GetString("name")
	size, _ := c.GetInt("size")

	url := kubernetesManagementProtocol + "://" + kubernetesManagementHost + ":" + kubernetesManagementPort +
		"/api/v1/replicationcontrollers/size/" + namespace + "/" + name + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + kubeapiPort

	putBodyJsonMap := make(map[string]interface{})
	putBodyJsonMap["size"] = size

	_, err := restclient.RequestPut(url, putBodyJsonMap, true)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		guimessage.AddSuccess("Replication Controller " + name + " is resized")
	}

	c.Ctx.Redirect(302, "/gui/inventory/replicationcontroller/")

	guimessage.RedirectMessage(c)
}
