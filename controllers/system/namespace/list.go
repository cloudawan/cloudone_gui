package namespace

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/kubernetes_management_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/kubernetes_management_utility/restclient"
)

type ListController struct {
	beego.Controller
}

type Namespace struct {
	Name     string
	Selected bool
	Display  string
}

var displayMap map[string]string = map[string]string{
	"default": "disabled",
}

func (c *ListController) Get() {
	c.TplNames = "system/namespace/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	kubernetesManagementProtocol := beego.AppConfig.String("kubernetesManagementProtocol")
	kubernetesManagementHost := beego.AppConfig.String("kubernetesManagementHost")
	kubernetesManagementPort := beego.AppConfig.String("kubernetesManagementPort")
	kubeapiHost := beego.AppConfig.String("kubeapiHost")
	kubeapiPort := beego.AppConfig.String("kubeapiPort")

	url := kubernetesManagementProtocol + "://" + kubernetesManagementHost + ":" + kubernetesManagementPort +
		"/api/v1/namespaces/" + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + kubeapiPort

	nameSlice := make([]string, 0)
	_, err := restclient.RequestGetWithStructure(url, &nameSlice)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		selectedNamespace := c.GetSession("namespace")

		namespaceSlice := make([]Namespace, 0)
		for _, name := range nameSlice {

			namespace := Namespace{name, false, ""}
			if name == selectedNamespace {
				namespace.Selected = true
			}
			namespace.Display = displayMap[name]

			namespaceSlice = append(namespaceSlice, namespace)
		}

		c.Data["namespaceSlice"] = namespaceSlice
	}

	guimessage.OutputMessage(c.Data)
}
