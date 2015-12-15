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

package configuration

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_utility/restclient"
	"strconv"
	"strings"
	"time"
)

var KubeapiHealthCheckTimeoutInMilliSecond = 300

func GetAvailableKubeapiHostAndPort() (returnedHost string, returnedPort int, returnedError error) {
	defer func() {
		if err := recover(); err != nil {
			returnedHost = ""
			returnedPort = 0
			returnedError = err.(error)
		}
	}()

	kubeapiHostAndPortText := beego.AppConfig.String("kubeapiHostAndPort")
	kubeapiHostAndPortSlice := strings.Split(kubeapiHostAndPortText, ",")
	kubeapiHealthCheckTimeoutInMilliSecond, err := beego.AppConfig.Int("kubeapiHealthCheckTimeoutInMilliSecond")
	if err != nil {
		kubeapiHealthCheckTimeoutInMilliSecond = KubeapiHealthCheckTimeoutInMilliSecond
	}
	for _, kubeapiHostAndPort := range kubeapiHostAndPortSlice {
		kubeapiHostAndPort = strings.TrimSpace(kubeapiHostAndPort)
		_, err := restclient.HealthCheck("http://"+kubeapiHostAndPort, time.Duration(kubeapiHealthCheckTimeoutInMilliSecond)*time.Millisecond)
		if err == nil {
			splitSlice := strings.Split(kubeapiHostAndPort, ":")
			host := splitSlice[0]
			port, err := strconv.Atoi(splitSlice[1])
			if err != nil {
				return "", 0, err
			}
			return host, port, nil
		}
	}

	return "", 0, errors.New("No available host and port")
}