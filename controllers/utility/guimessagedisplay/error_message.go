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

package guimessagedisplay

import (
	"fmt"
	"github.com/cloudawan/cloudone_utility/restclient"
)

func GetErrorMessage(err error) string {
	requestError, ok := err.(restclient.RequestError)
	if ok {
		if requestError.ResponseData != nil {
			responseDataJsonMap, ok := requestError.ResponseData.(map[string]interface{})
			if ok {
				errorField, ok := responseDataJsonMap["Error"].(string)
				if ok {
					errorMessageField, ok := responseDataJsonMap["ErrorMessage"].(string)
					if ok {
						return errorField + ": " + errorMessageField
					} else {
						return errorField
					}
				} else {
					return fmt.Sprintf("%v", requestError.ResponseData)
				}
			} else {
				return fmt.Sprintf("%v", responseDataJsonMap)
			}
		} else {
			return requestError.Error()
		}
	} else {
		return err.Error()
	}
}
