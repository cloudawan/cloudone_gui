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

import ()

type GUIMessage struct {
	sessionUtility SessionUtility
	successSlice   []string
	infoSlice      []string
	warningSlice   []string
	dangerSlice    []string
}

const (
	sessionNameGUIMessage = "guiMessage"
)

func GetGUIMessage(sessionUtility SessionUtility) *GUIMessage {
	guiMessage := sessionUtility.GetSession(sessionNameGUIMessage)
	if guiMessage == nil {
		guiMessage = new(GUIMessage)
		guiMessage.(*GUIMessage).sessionUtility = sessionUtility
		guiMessage.(*GUIMessage).successSlice = make([]string, 0)
		guiMessage.(*GUIMessage).infoSlice = make([]string, 0)
		guiMessage.(*GUIMessage).warningSlice = make([]string, 0)
		guiMessage.(*GUIMessage).dangerSlice = make([]string, 0)
	} else {
		sessionUtility.DelSession(sessionNameGUIMessage)
	}
	return guiMessage.(*GUIMessage)
}

func (guiMessage *GUIMessage) AddSuccess(text string) {
	guiMessage.successSlice = append(guiMessage.successSlice, text)
}

func (guiMessage *GUIMessage) AddInfo(text string) {
	guiMessage.infoSlice = append(guiMessage.infoSlice, text)
}

func (guiMessage *GUIMessage) AddWarning(text string) {
	guiMessage.warningSlice = append(guiMessage.warningSlice, text)
}

func (guiMessage *GUIMessage) AddDanger(text string) {
	guiMessage.dangerSlice = append(guiMessage.dangerSlice, text)
}

type SessionUtility interface {
	// SetSession puts value into session.
	SetSession(name interface{}, value interface{})

	// GetSession gets value from session.
	GetSession(name interface{}) interface{}

	// SetSession removes value from session.
	DelSession(name interface{})
}

func (guiMessage *GUIMessage) RedirectMessage(sessionUtility SessionUtility) {
	sessionUtility.SetSession(sessionNameGUIMessage, guiMessage)
}

func (guiMessage *GUIMessage) OutputMessage(data map[interface{}]interface{}) bool {
	// Show global data
	data["currentNamespace"] = guiMessage.sessionUtility.GetSession("namespace")

	if data == nil {
		return false
	} else {
		has := false
		if guiMessage.successSlice != nil && len(guiMessage.successSlice) > 0 {
			data["guiMessageSuccessSlice"] = guiMessage.successSlice
			has = true
		}
		if guiMessage.infoSlice != nil && len(guiMessage.infoSlice) > 0 {
			data["guiMessageInfoSlice"] = guiMessage.infoSlice
			has = true
		}
		if guiMessage.warningSlice != nil && len(guiMessage.warningSlice) > 0 {
			data["guiMessageWarningSlice"] = guiMessage.warningSlice
			has = true
		}
		if guiMessage.dangerSlice != nil && len(guiMessage.dangerSlice) > 0 {
			data["guiMessageDangerSlice"] = guiMessage.dangerSlice
			has = true
		}
		if has == false {
			data["guiMessageDisplay"] = "hidden"
		}
		return true
	}
}
