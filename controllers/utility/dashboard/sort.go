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

package dashboard

import (
	"sort"
)

type ByJsonMap []interface{}

func (b ByJsonMap) Len() int      { return len(b) }
func (b ByJsonMap) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b ByJsonMap) Less(i, j int) bool {
	iJsonMap, _ := b[i].(map[string]interface{})
	jJsonMap, _ := b[j].(map[string]interface{})
	iName, _ := iJsonMap["name"].(string)
	jName, _ := jJsonMap["name"].(string)
	return iName < jName
}

func RecursiveSortTheDataInGraphJsonMap(jsonMap map[string]interface{}) {
	if jsonMap == nil {
		return
	}

	childrenJsonSlice, ok := jsonMap["children"].([]interface{})

	if ok {
		recursiveSortTheDataInGraphJsonSlice(childrenJsonSlice)
	} else {
		for key, _ := range jsonMap {
			_, ok := jsonMap[key].(map[string]interface{})
			if ok {
				RecursiveSortTheDataInGraphJsonMap(jsonMap[key].(map[string]interface{}))
			}
			_, ok = jsonMap[key].([]interface{})
			if ok {
				recursiveSortTheDataInGraphJsonSlice(jsonMap[key].([]interface{}))
			}
		}
	}
}

func recursiveSortTheDataInGraphJsonSlice(jsonSlice []interface{}) {
	if jsonSlice == nil {
		return
	}

	for i := 0; i < len(jsonSlice); i++ {
		_, ok := jsonSlice[i].(map[string]interface{})
		if ok {
			RecursiveSortTheDataInGraphJsonMap(jsonSlice[i].(map[string]interface{}))
		}
	}

	sort.Sort(ByJsonMap(jsonSlice))
}
