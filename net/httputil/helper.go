// Copyright 2020 The Ledger Authors
//
// Licensed under the AGPL, Version 3.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.gnu.org/licenses/agpl-3.0.en.html
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package httputil

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (er *ErrorResponse) Error() string {
	return fmt.Sprintf("%d - %s", er.Code, er.Message)
}

// Abort is a helper function that calls `Abort()` and then `JSON` internally.
// This method stops the chain, writes the status code and return a JSON body with HTTP status code and error message.
// It also sets the Content-Type as "application/json".
func Abort(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, &ErrorResponse{
		Code:    code,
		Message: message,
	})
}
