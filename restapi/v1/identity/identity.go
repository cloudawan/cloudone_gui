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

package identity

import (
	"github.com/astaxie/beego/context"
)

const (
	loginURL         = "/api/v1/identity/login"
	webhookGithubURL = "/api/v1/webhook/github" // The webhook has its own verification
)

func FilterToken(ctx *context.Context) {
	if (ctx.Input.IsGet() || ctx.Input.IsPost()) && (ctx.Input.URL() == loginURL || ctx.Input.URL() == webhookGithubURL) {
		// Don't redirect itself to prevent the circle
	} else {
		// No need to verify here since it is just relay the data.
		// The authorization happens in the real processing places
	}
}
