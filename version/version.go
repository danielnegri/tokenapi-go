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

// Package version contains version information for this app.
package version

import (
	"fmt"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

// Version is set by the build scripts.
var (
	BuildTime  = time.Now().In(time.UTC).Format(time.Stamp + " 2006 UTC")
	CommitHash = ""
	Version    = ""
)

func NewCommand(description string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version and exit",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(description)
			fmt.Printf("Go Version: %s\n", runtime.Version())
			fmt.Printf("Go OS/ARCH: %s %s\n", runtime.GOOS, runtime.GOARCH)
			fmt.Printf("Build Time: %s\n", BuildTime)
			fmt.Printf("Commit: %s\n", CommitHash)
			fmt.Printf("Version: %s\n", Version)
		},
	}
}
