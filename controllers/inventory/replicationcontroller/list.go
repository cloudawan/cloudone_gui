package replicationcontroller

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/kubernetes_management_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/kubernetes_management_utility/restclient"
)

type ListController struct {
	beego.Controller
}

type ReplicationControllerAndRelatedPod struct {
	Name               string
	Namespace          string
	ReplicaAmount      int
	AliveReplicaAmount int
	Selector           map[string]string
	Label              map[string]string
	PodSlice           []Pod
	Display            string
}

type Pod struct {
	Name           string
	Namespace      string
	HostIP         string
	PodIP          string
	ContainerSlice []PodContainer
}

type PodContainer struct {
	Name      string
	Image     string
	PortSlice []PodContainerPort
}

type PodContainerPort struct {
	Name          string
	ContainerPort int
	Protocol      string
}

var displayMap map[string]string = map[string]string{
	"kube-dns-v6":           "disabled",
	"private-registry":      "disabled",
	"kubernetes-management": "disabled",
}

func (c *ListController) Get() {
	c.TplNames = "inventory/replicationcontroller/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	kubernetesManagementProtocol := beego.AppConfig.String("kubernetesManagementProtocol")
	kubernetesManagementHost := beego.AppConfig.String("kubernetesManagementHost")
	kubernetesManagementPort := beego.AppConfig.String("kubernetesManagementPort")
	kubeapiHost := beego.AppConfig.String("kubeapiHost")
	kubeapiPort := beego.AppConfig.String("kubeapiPort")
	namespace, _ := c.GetSession("namespace").(string)

	url := kubernetesManagementProtocol + "://" + kubernetesManagementHost + ":" + kubernetesManagementPort +
		"/api/v1/replicationcontrollers/" + namespace + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + kubeapiPort

	replicationControllerAndRelatedPodSlice := make([]ReplicationControllerAndRelatedPod, 0)
	_, err := restclient.RequestGetWithStructure(url, &replicationControllerAndRelatedPodSlice)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		for i := 0; i < len(replicationControllerAndRelatedPodSlice); i++ {
			replicationControllerAndRelatedPodSlice[i].Display =
				displayMap[replicationControllerAndRelatedPodSlice[i].Name]
		}
		c.Data["replicationControllerAndRelatedPodSlice"] = replicationControllerAndRelatedPodSlice
	}

	guimessage.OutputMessage(c.Data)
}
