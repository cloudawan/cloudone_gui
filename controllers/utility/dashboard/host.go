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
	"github.com/astaxie/beego/context"
	"strings"
)

func GetServerHostAndPortFromUserRequest(beegoInput *context.BeegoInput) (string, int) {
	host := beegoInput.Host()
	port := beegoInput.Port()
	// Due to beego's bug, when the port is skipped for 80(http) and 443(https), it will set to 80
	// So we need to judge whether it is 80 or 443
	if port == 80 {
		if strings.Contains(beegoInput.Context.Request.Host, ":") == false {
			if strings.HasPrefix(beegoInput.Referer(), "https") {
				port = 443
			}
		}
	}

	return host, port
}
