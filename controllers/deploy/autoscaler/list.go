package autoscaler

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/kubernetes_management_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/kubernetes_management_utility/restclient"
	"time"
)

type ListController struct {
	beego.Controller
}

type ReplicationControllerAutoScaler struct {
	Check             bool
	CoolDownDuration  time.Duration
	RemainingCoolDown time.Duration
	KubeapiHost       string
	KubeapiPort       int
	Namespace         string
	Kind              string
	Name              string
	MaximumReplica    int
	MinimumReplica    int
	IndicatorSlice    []Indicator
}

type Indicator struct {
	Type                  string
	AboveAllOrOne         bool
	AbovePercentageOfData float64
	AboveThreshold        int64
	BelowAllOrOne         bool
	BelowPercentageOfData float64
	BelowThreshold        int64
}

func (c *ListController) Get() {
	c.TplNames = "deploy/autoscaler/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	kubernetesManagementProtocol := beego.AppConfig.String("kubernetesManagementProtocol")
	kubernetesManagementHost := beego.AppConfig.String("kubernetesManagementHost")
	kubernetesManagementPort := beego.AppConfig.String("kubernetesManagementPort")

	url := kubernetesManagementProtocol + "://" + kubernetesManagementHost + ":" + kubernetesManagementPort +
		"/api/v1/autoscalers/"

	replicationControllerAutoScalerSlice := make([]ReplicationControllerAutoScaler, 0)

	returnedReplicationControllerAutoScalerSlice, err := restclient.RequestGetWithStructure(url, &replicationControllerAutoScalerSlice)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		c.Data["replicationControllerAutoScalerSlice"] = returnedReplicationControllerAutoScalerSlice
	}

	guimessage.OutputMessage(c.Data)
}
