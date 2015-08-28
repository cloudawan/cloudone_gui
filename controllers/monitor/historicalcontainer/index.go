// Copyright 2015 CloudAwan LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package historicalcontainer

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/cloudawan/kubernetes_management_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/kubernetes_management_utility/restclient"
	"net/url"
	"sort"
	"strconv"
	"time"
)

type IndexController struct {
	beego.Controller
}

const (
	allKeyword = "All"
)

func (c *IndexController) Get() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	c.TplNames = "monitor/historicalcontainer/index.html"

	kubernetesManagementAnalysisProtocol := beego.AppConfig.String("kubernetesManagementAnalysisProtocol")
	kubernetesManagementAnalysisHost := beego.AppConfig.String("kubernetesManagementAnalysisHost")
	kubernetesManagementAnalysisPort := beego.AppConfig.String("kubernetesManagementAnalysisPort")
	kubernetesManagementGUIProtocol := beego.AppConfig.String("kubernetesManagementGUIProtocol")
	kubernetesManagementGUIHost := beego.AppConfig.String("kubernetesManagementGUIHost")
	kubernetesManagementGUIPort := beego.AppConfig.String("kubernetesManagementGUIPort")

	namespaces, _ := c.GetSession("namespace").(string)

	url := kubernetesManagementAnalysisProtocol + "://" + kubernetesManagementAnalysisHost + ":" + kubernetesManagementAnalysisPort +
		"/api/v1/historicalreplicationcontrollers/names/" + namespaces

	jsonMapSlice := make([]interface{}, 0)

	_, err := restclient.RequestGetWithStructure(url, &jsonMapSlice)

	if err != nil {
		// Error
		guimessage.AddDanger("Fail to get all replication controller name with error " + err.Error())
	} else {
		nameSlice := make([]string, 0)
		nameSlice = append(nameSlice, allKeyword)
		for _, value := range jsonMapSlice {
			name, ok := value.(string)
			if ok {
				nameSlice = append(nameSlice, name)
			}
		}

		c.Data["kubernetesManagementGUIProtocol"] = kubernetesManagementGUIProtocol
		c.Data["kubernetesManagementGUIHost"] = kubernetesManagementGUIHost
		c.Data["kubernetesManagementGUIPort"] = kubernetesManagementGUIPort
		c.Data["replicationControllerNameSlice"] = nameSlice
	}

	guimessage.RedirectMessage(c)
}

type DataController struct {
	beego.Controller
}

const (
	aggregationAmount = 60
)

func (c *DataController) Get() {

	kubernetesManagementAnalysisProtocol := beego.AppConfig.String("kubernetesManagementAnalysisProtocol")
	kubernetesManagementAnalysisHost := beego.AppConfig.String("kubernetesManagementAnalysisHost")
	kubernetesManagementAnalysisPort := beego.AppConfig.String("kubernetesManagementAnalysisPort")
	kubeapiHost := beego.AppConfig.String("kubeapiHost")
	kubeapiPort, _ := beego.AppConfig.Int("kubeapiPort")

	namespaces, _ := c.GetSession("namespace").(string)
	timeZoneOffset, _ := c.GetSession("timeZoneOffset").(int)

	replicationControllerName := c.GetString("replicationController")
	fromText := c.GetString("from")
	toText := c.GetString("to")

	//current := time.Now()
	//_, serverTimeZoneOffsetInSecond := current.Zone()

	from, err := time.Parse("01/02/2006 15:04 PM", fromText)
	if err != nil {
		// Error
		errorJsonMap := make(map[string]interface{})
		errorJsonMap["warning"] = "From time is not selected."
		c.Data["json"] = errorJsonMap
		c.ServeJson()
		return
		/*
			from = current.Add(-1 * time.Hour)
			// Convert to UTC
			// Offset server time zone
			from = from.Add(-1 * time.Second * time.Duration(serverTimeZoneOffsetInSecond))
			// Offset browser time zone since time from browser doesn't contain time zone
			from = from.Add(time.Minute * time.Duration(timeZoneOffset))
		*/
	} else {
		// Convert to UTC
		// Offset browser time zone since time from browser doesn't contain time zone
		from = from.Add(time.Minute * time.Duration(timeZoneOffset))
	}

	fromInRFC3339Nano := from.Format(time.RFC3339Nano)

	to, err := time.Parse("01/02/2006 15:04 PM", toText)
	if err != nil {
		// Error
		errorJsonMap := make(map[string]interface{})
		errorJsonMap["warning"] = "To time is not selected."
		c.Data["json"] = errorJsonMap
		c.ServeJson()
		return
		/*
			to = current
			// Convert to UTC
			// Offset server time zone
			to.Add(-1 * time.Second * time.Duration(serverTimeZoneOffsetInSecond))
			// Offset browser time zone since time from browser doesn't contain time zone
			to = to.Add(time.Minute * time.Duration(timeZoneOffset))
		*/
	} else {
		// Convert to UTC
		// Offset browser time zone since time from browser doesn't contain time zone
		to = to.Add(time.Minute * time.Duration(timeZoneOffset))
	}

	toInRFC3339Nano := to.Format(time.RFC3339Nano)

	// Make sure from is before to
	if from.Before(to) == false {
		// Error
		errorJsonMap := make(map[string]interface{})
		errorJsonMap["warning"] = "From need to be before to."
		c.Data["json"] = errorJsonMap
		c.ServeJson()
		return
	}

	allHistoricalReplicationControllerMetricJsonMap := make(map[string]interface{})
	if replicationControllerName != "" && replicationControllerName != allKeyword {

		encodingUrl, _ := url.Parse(kubernetesManagementAnalysisProtocol + "://" + kubernetesManagementAnalysisHost + ":" + kubernetesManagementAnalysisPort +
			"/api/v1/historicalreplicationcontrollermetrics/" + namespaces + "/" + replicationControllerName)
		parameters := url.Values{}
		parameters.Add("kubeapihost", kubeapiHost)
		parameters.Add("kubeapiport", strconv.Itoa(kubeapiPort))
		parameters.Add("from", fromInRFC3339Nano)
		parameters.Add("to", toInRFC3339Nano)
		parameters.Add("aggregationAmount", strconv.Itoa(aggregationAmount))
		encodingUrl.RawQuery = parameters.Encode()

		result, err := restclient.RequestGet(encodingUrl.String(), true)
		historicalReplicationControllerMetricJsonMap, ok := result.(map[string]interface{})
		if err != nil || ok == false {
			// Error
			errorJsonMap := make(map[string]interface{})
			errorJsonMap["error"] = err.Error()
			c.Data["json"] = errorJsonMap
			c.ServeJson()
			return
		}
		allHistoricalReplicationControllerMetricJsonMap[replicationControllerName] = historicalReplicationControllerMetricJsonMap
	} else {
		encodingUrl, _ := url.Parse(kubernetesManagementAnalysisProtocol + "://" + kubernetesManagementAnalysisHost + ":" + kubernetesManagementAnalysisPort +
			"/api/v1/historicalreplicationcontrollermetrics/" + namespaces)
		parameters := url.Values{}
		parameters.Add("kubeapihost", kubeapiHost)
		parameters.Add("kubeapiport", strconv.Itoa(kubeapiPort))
		parameters.Add("from", fromInRFC3339Nano)
		parameters.Add("to", toInRFC3339Nano)
		parameters.Add("aggregationAmount", strconv.Itoa(aggregationAmount))
		encodingUrl.RawQuery = parameters.Encode()

		result, err := restclient.RequestGet(encodingUrl.String(), true)
		var ok bool
		allHistoricalReplicationControllerMetricJsonMap, ok = result.(map[string]interface{})
		if err != nil || ok == false {
			// Error
			errorJsonMap := make(map[string]interface{})
			errorJsonMap["error"] = err.Error()
			c.Data["json"] = errorJsonMap
			c.ServeJson()
			return
		}
	}

	// At least two points to draw a line
	// Filtered out those less than 2
	filterHistoricalReplicationControllerMetricJsonMap := make(map[string]interface{})
	for key, historicalReplicationControllerMetricJsonMap := range allHistoricalReplicationControllerMetricJsonMap {
		timestampSlice, _ := historicalReplicationControllerMetricJsonMap.(map[string]interface{})["timestamp"].([]interface{})
		resultAmount := len(timestampSlice)
		if resultAmount >= 2 {
			filterHistoricalReplicationControllerMetricJsonMap[key] = historicalReplicationControllerMetricJsonMap
		}
	}

	// No data to show
	if len(filterHistoricalReplicationControllerMetricJsonMap) == 0 {
		// Error
		errorJsonMap := make(map[string]interface{})
		errorJsonMap["error"] = "Insufficient data"
		c.Data["json"] = errorJsonMap
		c.ServeJson()
		return
	}

	// Due to boundary (start + end), it may be aggregationAmount or aggregationAmount + 1 for different replication controller metrics.
	// If the covered range has no data, it may be less. Force to use the smallest one to align.
	smallestResultAmount := aggregationAmount
	for _, historicalReplicationControllerMetricJsonMap := range filterHistoricalReplicationControllerMetricJsonMap {
		timestampSlice, _ := historicalReplicationControllerMetricJsonMap.(map[string]interface{})["timestamp"].([]interface{})
		resultAmount := len(timestampSlice)
		if resultAmount < smallestResultAmount {
			smallestResultAmount = resultAmount
		}
	}

	resultAmount := smallestResultAmount
	differenceAmount := resultAmount - 1

	// CPU usage total
	cpuUsageTotalJsonMap := make(map[string]interface{})
	cpuUsageTotalJsonMap["metadata"] = make(map[string]interface{})
	cpuUsageTotalJsonMap["metadata"].(map[string]interface{})["title"] = "CPU (ms/1s)"
	cpuUsageTotalJsonMap["metadata"].(map[string]interface{})["lineName"] = make([]string, 0)
	cpuUsageTotalJsonMap["data"] = make([]interface{}, differenceAmount)
	for i := 0; i < differenceAmount; i++ {
		cpuUsageTotalJsonMap["data"].([]interface{})[i] = make(map[string]interface{})
		cpuUsageTotalJsonMap["data"].([]interface{})[i].(map[string]interface{})["y"] = make([]int64, 0)
	}

	// Memory usage
	memoryUsageJsonMap := make(map[string]interface{})
	memoryUsageJsonMap["metadata"] = make(map[string]interface{})
	memoryUsageJsonMap["metadata"].(map[string]interface{})["title"] = "Memory(MB)"
	memoryUsageJsonMap["metadata"].(map[string]interface{})["lineName"] = make([]string, 0)
	memoryUsageJsonMap["data"] = make([]interface{}, resultAmount)
	for i := 0; i < resultAmount; i++ {
		memoryUsageJsonMap["data"].([]interface{})[i] = make(map[string]interface{})
		memoryUsageJsonMap["data"].([]interface{})[i].(map[string]interface{})["y"] = make([]int64, 0)
	}

	// Disk I/O byte
	diskIOServiceBytesStatsJsonMap := make(map[string]interface{})
	diskIOServiceBytesStatsJsonMap["metadata"] = make(map[string]interface{})
	diskIOServiceBytesStatsJsonMap["metadata"].(map[string]interface{})["title"] = "Disk I/O (Byte/s)"
	diskIOServiceBytesStatsJsonMap["metadata"].(map[string]interface{})["lineName"] = make([]string, 0)
	diskIOServiceBytesStatsJsonMap["data"] = make([]interface{}, differenceAmount)
	for i := 0; i < differenceAmount; i++ {
		diskIOServiceBytesStatsJsonMap["data"].([]interface{})[i] = make(map[string]interface{})
		diskIOServiceBytesStatsJsonMap["data"].([]interface{})[i].(map[string]interface{})["y"] = make([]int64, 0)
	}

	// Disk I/O count
	diskIOServicedStatsJsonMap := make(map[string]interface{})
	diskIOServicedStatsJsonMap["metadata"] = make(map[string]interface{})
	diskIOServicedStatsJsonMap["metadata"].(map[string]interface{})["title"] = "Disk I/O (count/s)"
	diskIOServicedStatsJsonMap["metadata"].(map[string]interface{})["lineName"] = make([]string, 0)
	diskIOServicedStatsJsonMap["data"] = make([]interface{}, differenceAmount)
	for i := 0; i < differenceAmount; i++ {
		diskIOServicedStatsJsonMap["data"].([]interface{})[i] = make(map[string]interface{})
		diskIOServicedStatsJsonMap["data"].([]interface{})[i].(map[string]interface{})["y"] = make([]int64, 0)
	}

	// Network RX Byte
	networkRXBytesJsonMap := make(map[string]interface{})
	networkRXBytesJsonMap["metadata"] = make(map[string]interface{})
	networkRXBytesJsonMap["metadata"].(map[string]interface{})["title"] = "Network RX (Bytes/s)"
	networkRXBytesJsonMap["metadata"].(map[string]interface{})["lineName"] = make([]string, 0)
	networkRXBytesJsonMap["data"] = make([]interface{}, differenceAmount)
	for i := 0; i < differenceAmount; i++ {
		networkRXBytesJsonMap["data"].([]interface{})[i] = make(map[string]interface{})
		networkRXBytesJsonMap["data"].([]interface{})[i].(map[string]interface{})["y"] = make([]int64, 0)
	}

	// Network TX Byte
	networkTXBytesJsonMap := make(map[string]interface{})
	networkTXBytesJsonMap["metadata"] = make(map[string]interface{})
	networkTXBytesJsonMap["metadata"].(map[string]interface{})["title"] = "Network TX (Bytes/s)"
	networkTXBytesJsonMap["metadata"].(map[string]interface{})["lineName"] = make([]string, 0)
	networkTXBytesJsonMap["data"] = make([]interface{}, differenceAmount)
	for i := 0; i < differenceAmount; i++ {
		networkTXBytesJsonMap["data"].([]interface{})[i] = make(map[string]interface{})
		networkTXBytesJsonMap["data"].([]interface{})[i].(map[string]interface{})["y"] = make([]int64, 0)
	}

	// Network RX Packet
	networkRXPacketsJsonMap := make(map[string]interface{})
	networkRXPacketsJsonMap["metadata"] = make(map[string]interface{})
	networkRXPacketsJsonMap["metadata"].(map[string]interface{})["title"] = "Network RX (packet/s)"
	networkRXPacketsJsonMap["metadata"].(map[string]interface{})["lineName"] = make([]string, 0)
	networkRXPacketsJsonMap["data"] = make([]interface{}, differenceAmount)
	for i := 0; i < differenceAmount; i++ {
		networkRXPacketsJsonMap["data"].([]interface{})[i] = make(map[string]interface{})
		networkRXPacketsJsonMap["data"].([]interface{})[i].(map[string]interface{})["y"] = make([]int64, 0)
	}

	// Network TX Packet
	networkTXPacketsJsonMap := make(map[string]interface{})
	networkTXPacketsJsonMap["metadata"] = make(map[string]interface{})
	networkTXPacketsJsonMap["metadata"].(map[string]interface{})["title"] = "Network TX (packet/s)"
	networkTXPacketsJsonMap["metadata"].(map[string]interface{})["lineName"] = make([]string, 0)
	networkTXPacketsJsonMap["data"] = make([]interface{}, differenceAmount)
	for i := 0; i < differenceAmount; i++ {
		networkTXPacketsJsonMap["data"].([]interface{})[i] = make(map[string]interface{})
		networkTXPacketsJsonMap["data"].([]interface{})[i].(map[string]interface{})["y"] = make([]int64, 0)
	}

	allHistoricalReplicationControllerMetricKeySlice := getAllKeyFromJsonMap(filterHistoricalReplicationControllerMetricJsonMap)
	for _, replicationControllerName := range allHistoricalReplicationControllerMetricKeySlice {
		differenceTimestampSlice := make([]string, 0)
		originalTimestampSlice := make([]string, 0)
		timestampSlice, _ := filterHistoricalReplicationControllerMetricJsonMap[replicationControllerName].(map[string]interface{})["timestamp"].([]interface{})
		for index, value := range timestampSlice {
			timestamp, _ := time.Parse(time.RFC3339Nano, value.(string))
			// Convert from UTC to local time zone
			timestamp = timestamp.Add(-1 * time.Minute * time.Duration(timeZoneOffset))
			timestampText := timestamp.Format("2006-01-02 15:04:05")
			if index == 0 {
				originalTimestampSlice = append(originalTimestampSlice, timestampText)
			} else {
				originalTimestampSlice = append(originalTimestampSlice, timestampText)
				differenceTimestampSlice = append(differenceTimestampSlice, timestampText)
			}
		}

		allHistoricalPodJsonMap, _ := filterHistoricalReplicationControllerMetricJsonMap[replicationControllerName].(map[string]interface{})
		allHistoricalPodKeySlice := getAllKeyFromJsonMap(allHistoricalPodJsonMap)
		for _, podName := range allHistoricalPodKeySlice {
			if podName != "timestamp" {
				allHistoricalContainerJsonMap, _ := allHistoricalPodJsonMap[podName].(map[string]interface{})
				allHistoricalContainerKeySlice := getAllKeyFromJsonMap(allHistoricalContainerJsonMap)
				for _, containerName := range allHistoricalContainerKeySlice {
					historicalContainerJsonMap, _ := allHistoricalContainerJsonMap[containerName].(map[string]interface{})

					// CPU usage total
					cpuUsageTotalJsonMap["metadata"].(map[string]interface{})["lineName"] = append(
						cpuUsageTotalJsonMap["metadata"].(map[string]interface{})["lineName"].([]string),
						replicationControllerName+"/"+podName+"/"+containerName)
					// Memory usage
					memoryUsageJsonMap["metadata"].(map[string]interface{})["lineName"] = append(
						memoryUsageJsonMap["metadata"].(map[string]interface{})["lineName"].([]string),
						replicationControllerName+"/"+podName+"/"+containerName)
					// Disk I/O byte
					diskIOServiceBytesStatsJsonMap["metadata"].(map[string]interface{})["lineName"] = append(
						diskIOServiceBytesStatsJsonMap["metadata"].(map[string]interface{})["lineName"].([]string),
						replicationControllerName+"/"+podName+"/"+containerName)
					// Disk I/O count
					diskIOServicedStatsJsonMap["metadata"].(map[string]interface{})["lineName"] = append(
						diskIOServicedStatsJsonMap["metadata"].(map[string]interface{})["lineName"].([]string),
						replicationControllerName+"/"+podName+"/"+containerName)
					// Network RX Byte
					networkRXBytesJsonMap["metadata"].(map[string]interface{})["lineName"] = append(
						networkRXBytesJsonMap["metadata"].(map[string]interface{})["lineName"].([]string),
						replicationControllerName+"/"+podName+"/"+containerName)
					// Network TX Byte
					networkTXBytesJsonMap["metadata"].(map[string]interface{})["lineName"] = append(
						networkTXBytesJsonMap["metadata"].(map[string]interface{})["lineName"].([]string),
						replicationControllerName+"/"+podName+"/"+containerName)
					// Network RX Packet
					networkRXPacketsJsonMap["metadata"].(map[string]interface{})["lineName"] = append(
						networkRXPacketsJsonMap["metadata"].(map[string]interface{})["lineName"].([]string),
						replicationControllerName+"/"+podName+"/"+containerName)
					// Network TX Packet
					networkTXPacketsJsonMap["metadata"].(map[string]interface{})["lineName"] = append(
						networkTXPacketsJsonMap["metadata"].(map[string]interface{})["lineName"].([]string),
						replicationControllerName+"/"+podName+"/"+containerName)

					for i := 0; i < differenceAmount; i++ {
						count, _ := historicalContainerJsonMap["documentCountSlice"].([]interface{})[i].(json.Number).Int64()
						// CPU usage total
						second, _ := historicalContainerJsonMap["minimumCpuUsageTotalSlice"].([]interface{})[i+1].(json.Number).Int64()
						first, _ := historicalContainerJsonMap["minimumCpuUsageTotalSlice"].([]interface{})[i].(json.Number).Int64()
						value := (second - first) / count
						// Convert from ns to ms
						value /= 1000000
						cpuUsageTotalJsonMap["data"].([]interface{})[i].(map[string]interface{})["y"] = append(
							cpuUsageTotalJsonMap["data"].([]interface{})[i].(map[string]interface{})["y"].([]int64), value)
						cpuUsageTotalJsonMap["data"].([]interface{})[i].(map[string]interface{})["x"] = differenceTimestampSlice[i]

						// Disk I/O byte
						second, _ = historicalContainerJsonMap["minimumDiskioIoServiceBytesStatsTotalSlice"].([]interface{})[i+1].(json.Number).Int64()
						first, _ = historicalContainerJsonMap["minimumDiskioIoServiceBytesStatsTotalSlice"].([]interface{})[i].(json.Number).Int64()
						value = (second - first) / count
						diskIOServiceBytesStatsJsonMap["data"].([]interface{})[i].(map[string]interface{})["y"] = append(
							diskIOServiceBytesStatsJsonMap["data"].([]interface{})[i].(map[string]interface{})["y"].([]int64), value)
						diskIOServiceBytesStatsJsonMap["data"].([]interface{})[i].(map[string]interface{})["x"] = differenceTimestampSlice[i]

						// Disk I/O byte
						second, _ = historicalContainerJsonMap["minimumDiskioIoServicedStatsTotalSlice"].([]interface{})[i+1].(json.Number).Int64()
						first, _ = historicalContainerJsonMap["minimumDiskioIoServicedStatsTotalSlice"].([]interface{})[i].(json.Number).Int64()
						value = (second - first) / count
						diskIOServicedStatsJsonMap["data"].([]interface{})[i].(map[string]interface{})["y"] = append(
							diskIOServicedStatsJsonMap["data"].([]interface{})[i].(map[string]interface{})["y"].([]int64), value)
						diskIOServicedStatsJsonMap["data"].([]interface{})[i].(map[string]interface{})["x"] = differenceTimestampSlice[i]

						// Network RX Byte
						second, _ = historicalContainerJsonMap["minimumNetworkRxBytesSlice"].([]interface{})[i+1].(json.Number).Int64()
						first, _ = historicalContainerJsonMap["minimumNetworkRxBytesSlice"].([]interface{})[i].(json.Number).Int64()
						value = (second - first) / count
						networkRXBytesJsonMap["data"].([]interface{})[i].(map[string]interface{})["y"] = append(
							networkRXBytesJsonMap["data"].([]interface{})[i].(map[string]interface{})["y"].([]int64), value)
						networkRXBytesJsonMap["data"].([]interface{})[i].(map[string]interface{})["x"] = differenceTimestampSlice[i]

						// Network TX Byte
						second, _ = historicalContainerJsonMap["minimumNetworkTxBytesSlice"].([]interface{})[i+1].(json.Number).Int64()
						first, _ = historicalContainerJsonMap["minimumNetworkTxBytesSlice"].([]interface{})[i].(json.Number).Int64()
						value = (second - first) / count
						networkTXBytesJsonMap["data"].([]interface{})[i].(map[string]interface{})["y"] = append(
							networkTXBytesJsonMap["data"].([]interface{})[i].(map[string]interface{})["y"].([]int64), value)
						networkTXBytesJsonMap["data"].([]interface{})[i].(map[string]interface{})["x"] = differenceTimestampSlice[i]

						// Network RX Packet
						second, _ = historicalContainerJsonMap["minimumNetworkRxPacketsSlice"].([]interface{})[i+1].(json.Number).Int64()
						first, _ = historicalContainerJsonMap["minimumNetworkRxPacketsSlice"].([]interface{})[i].(json.Number).Int64()
						value = (second - first) / count
						networkRXPacketsJsonMap["data"].([]interface{})[i].(map[string]interface{})["y"] = append(
							networkRXPacketsJsonMap["data"].([]interface{})[i].(map[string]interface{})["y"].([]int64), value)
						networkRXPacketsJsonMap["data"].([]interface{})[i].(map[string]interface{})["x"] = differenceTimestampSlice[i]

						// Network TX Packet
						second, _ = historicalContainerJsonMap["minimumNetworkTxPacketsSlice"].([]interface{})[i+1].(json.Number).Int64()
						first, _ = historicalContainerJsonMap["minimumNetworkTxPacketsSlice"].([]interface{})[i].(json.Number).Int64()
						value = (second - first) / count
						networkTXPacketsJsonMap["data"].([]interface{})[i].(map[string]interface{})["y"] = append(
							networkTXPacketsJsonMap["data"].([]interface{})[i].(map[string]interface{})["y"].([]int64), value)
						networkTXPacketsJsonMap["data"].([]interface{})[i].(map[string]interface{})["x"] = differenceTimestampSlice[i]
					}

					for i := 0; i < resultAmount; i++ {
						// Memory usage
						value, _ := historicalContainerJsonMap["averageMemoryUsageSlice"].([]interface{})[i].(json.Number).Int64()
						// Convert from byte to MB
						value /= (1024 * 1024)
						memoryUsageJsonMap["data"].([]interface{})[i].(map[string]interface{})["y"] = append(
							memoryUsageJsonMap["data"].([]interface{})[i].(map[string]interface{})["y"].([]int64), value)
						memoryUsageJsonMap["data"].([]interface{})[i].(map[string]interface{})["x"] = originalTimestampSlice[i]
					}
				}
			}
		}
	}

	chartSlice := make(map[string]interface{})

	// CPU usage total
	chartSlice["cpuUsageTotal"] = cpuUsageTotalJsonMap
	// Memory usage
	chartSlice["memoryUsage"] = memoryUsageJsonMap
	// Disk I/O byte
	chartSlice["diskIOServiceBytesStats"] = diskIOServiceBytesStatsJsonMap
	// Disk I/O count
	chartSlice["diskIOServicedStats"] = diskIOServicedStatsJsonMap
	// Network RX Byte
	chartSlice["networkRXBytes"] = networkRXBytesJsonMap
	// Network TX Byte
	chartSlice["networkTXBytes"] = networkTXBytesJsonMap

	// Network RX Packet
	chartSlice["networkRXPackets"] = networkRXPacketsJsonMap
	// Network TX Packet
	chartSlice["networkTXPackets"] = networkTXPacketsJsonMap

	c.Data["json"] = chartSlice

	c.ServeJson()
}

func getAllKeyFromJsonMap(jsonMap map[string]interface{}) []string {
	keySlice := make([]string, 0)
	for key, _ := range jsonMap {
		keySlice = append(keySlice, key)
	}
	sort.Strings(keySlice)
	return keySlice
}
