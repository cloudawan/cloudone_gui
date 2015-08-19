package kubernetes

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/kubernetes_management_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/kubernetes_management_utility/restclient"
	"strconv"
	"time"
)

type ListController struct {
	beego.Controller
}

type KubernetesEvent struct {
	Namespace      string
	Id             string
	FirstTimestamp string
	LastTimestamp  string
	Count          int
	Message        string
	Reason         string
	Acknowledge    bool
	Action         string
	Button         string
}

const (
	amountPerPage = 10
)

func (c *ListController) Get() {
	c.TplNames = "event/kubernetes/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	kubernetesManagementAnalysisProtocol := beego.AppConfig.String("kubernetesManagementAnalysisProtocol")
	kubernetesManagementAnalysisHost := beego.AppConfig.String("kubernetesManagementAnalysisHost")
	kubernetesManagementAnalysisPort := beego.AppConfig.String("kubernetesManagementAnalysisPort")

	acknowledge := c.GetString("acknowledge")
	if acknowledge == "" {
		acknowledge = "false"
	}

	offset, _ := c.GetInt("offset")

	timeZoneOffset, _ := c.GetSession("timeZoneOffset").(int)

	url := kubernetesManagementAnalysisProtocol + "://" + kubernetesManagementAnalysisHost + ":" + kubernetesManagementAnalysisPort +
		"/api/v1/historicalevents?acknowledge=" + acknowledge + "&size=" + strconv.Itoa(amountPerPage) + "&offset=" + strconv.Itoa(offset)

	jsonMapSlice := make([]map[string]interface{}, 0)
	_, err := restclient.RequestGetWithStructure(url, &jsonMapSlice)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		var action string
		var button string
		if acknowledge == "true" {
			action = "false"
			button = "Unacknowledge"
		} else {
			action = "true"
			button = "Acknowledge"
		}

		kubernetesEventSlice := make([]KubernetesEvent, 0)
		for _, jsonMap := range jsonMapSlice {
			sourceJsonMap, _ := jsonMap["_source"].(map[string]interface{})

			namespace, _ := sourceJsonMap["metadata"].(map[string]interface{})["namespace"].(string)
			id, _ := jsonMap["_id"].(string)
			firstTimestamp, _ := sourceJsonMap["firstTimestamp"].(string)
			lastTimestamp, _ := sourceJsonMap["lastTimestamp"].(string)
			count, _ := sourceJsonMap["count"].(float64)
			message, _ := sourceJsonMap["message"].(string)
			reason, _ := sourceJsonMap["reason"].(string)
			acknowledge, _ := sourceJsonMap["searchMetaData"].(map[string]interface{})["acknowledge"].(bool)

			firstTime, err := time.Parse(time.RFC3339, firstTimestamp)
			if err != nil {
				// Fail to parse, show original one
			} else {
				firstTimestamp = firstTime.Add(time.Minute * time.Duration(timeZoneOffset) * -1).Format(time.RFC3339)
			}

			lastTime, err := time.Parse(time.RFC3339, lastTimestamp)
			if err != nil {
				// Fail to parse, show original one
			} else {
				lastTimestamp = lastTime.Add(time.Minute * time.Duration(timeZoneOffset) * -1).Format(time.RFC3339)
			}

			kubernetesEvent := KubernetesEvent{
				namespace,
				id,
				firstTimestamp,
				lastTimestamp,
				int(count),
				message,
				reason,
				acknowledge,
				action,
				button,
			}

			kubernetesEventSlice = append(kubernetesEventSlice, kubernetesEvent)
		}

		previousOffset := offset - amountPerPage
		if previousOffset < 0 {
			previousOffset = 0
		}
		nextOffset := offset + amountPerPage

		previousFrom := previousOffset
		if previousFrom < 0 {
			previousFrom = 0
		}
		previousFrom += 1
		previousTo := previousOffset + amountPerPage
		c.Data["previousLabel"] = strconv.Itoa(previousFrom) + "~" + strconv.Itoa(previousTo)
		if offset == 0 {
			c.Data["previousButtonHidden"] = "hidden"
		} else {
			c.Data["previousButtonHidden"] = ""
		}

		nextFrom := nextOffset + 1
		nextTo := nextOffset + amountPerPage
		c.Data["nextLabel"] = strconv.Itoa(nextFrom) + "~" + strconv.Itoa(nextTo)

		if acknowledge == "true" {
			c.Data["acknowledgeActive"] = "active"
			c.Data["paginationUrlPrevious"] = "/gui/event/kubernetes/?acknowledge=true&offset=" + strconv.Itoa(previousOffset)
			c.Data["paginationUrlNext"] = "/gui/event/kubernetes/?acknowledge=true&offset=" + strconv.Itoa(nextOffset)
		} else {
			c.Data["unacknowledgeActive"] = "active"
			c.Data["paginationUrlPrevious"] = "/gui/event/kubernetes/?acknowledge=false&offset=" + strconv.Itoa(previousOffset)
			c.Data["paginationUrlNext"] = "/gui/event/kubernetes/?acknowledge=false&offset=" + strconv.Itoa(nextOffset)
		}

		c.Data["kubernetesEventSlice"] = kubernetesEventSlice
	}

	guimessage.OutputMessage(c.Data)
}