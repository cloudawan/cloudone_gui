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
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	"time"
)

const (
	loginExpired = time.Minute * 10
	signingKey   = "tokenKey"
)

var tokenMap map[string]*jwt.Token = make(map[string]*jwt.Token)

func createToken(username string, password string) (string, error) {
	// TODO verify
	if username != "admin" || password != "password" {
		return "", nil
	} else {
		// Create the token
		token := jwt.New(jwt.SigningMethodHS512)
		// Set some claims
		token.Claims["username"] = username
		token.Claims["exp"] = time.Now().Add(loginExpired).Unix()

		tokenMap[username] = token

		// Sign and get the complete encoded token as a string
		return token.SignedString([]byte(signingKey))
	}
}

func verifyToken(token string) error {
	_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		// TODO verify
		username, _ := token.Claims["username"].(string)
		if tokenMap[username] == nil {
			return nil, errors.New("Unauthorized user")
		}

		return []byte(signingKey), nil
	})

	return err
}
