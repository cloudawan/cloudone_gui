package namespace

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/kubernetes_management_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/kubernetes_management_utility/restclient"
)

type EditController struct {
	beego.Controller
}

func (c *EditController) Get() {
	c.TplNames = "system/namespace/edit.html"

	namespace := c.GetString("namespace")
	if namespace == "" {
		c.Data["actionButtonValue"] = "Create"
		c.Data["pageHeader"] = "Create Namespace"
		c.Data["namespace"] = ""
	} else {
		c.Data["actionButtonValue"] = "Update"
		c.Data["pageHeader"] = "Update Namespace"
		c.Data["namespace"] = namespace
	}
}

func (c *EditController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	kubernetesManagementProtocol := beego.AppConfig.String("kubernetesManagementProtocol")
	kubernetesManagementHost := beego.AppConfig.String("kubernetesManagementHost")
	kubernetesManagementPort := beego.AppConfig.String("kubernetesManagementPort")
	kubeapiHost := beego.AppConfig.String("kubeapiHost")
	kubeapiPort := beego.AppConfig.String("kubeapiPort")

	name := c.GetString("name")

	namespace := Namespace{name, false, ""}

	url := kubernetesManagementProtocol + "://" + kubernetesManagementHost + ":" + kubernetesManagementPort +
		"/api/v1/namespaces/" + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + kubeapiPort

	_, err := restclient.RequestPostWithStructure(url, namespace, nil)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		guimessage.AddSuccess("Namespace " + name + " is edited")
	}

	c.Ctx.Redirect(302, "/gui/system/namespace/")

	guimessage.RedirectMessage(c)
}
