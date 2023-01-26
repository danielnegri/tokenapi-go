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

package main

import (
	"fmt"
	"os"

	"github.com/danielnegri/tokenapi-go/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	envPrefix        = "AdhereTech"
	longDescription  = "AdhereTech Ledger"
	shortDescription = "ledger"
)

var environment string = os.Getenv("ENVIRONMENT")

func commandRoot() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   shortDescription,
		Short: shortDescription,
		Long:  longDescription,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			os.Exit(2)
		},
	}

	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()

	rootCmd.AddCommand(commandServe())
	rootCmd.AddCommand(version.NewCommand(longDescription))

	return rootCmd
}

func main() {
	if err := commandRoot().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}
}
