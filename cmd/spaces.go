/*
Copyright © 2026 Amanda Hager Lopes de Andrade Katz amandahla@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/amandahla/syncli/internal"
	"github.com/amandahla/syncli/internal/synapse"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// spacesCmd represents the spaces command
var spacesCmd = &cobra.Command{
	Use:   "spaces",
	Short: "Retrieve a list of public spaces from the Synapse Matrix homeserver.",
	Long:  `The list contains name, members, child count and child rooms ids`,
	Run: func(cmd *cobra.Command, args []string) {
		err := getSpaces(config)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"event": "get_spaces_error",
				"error": err,
			}).Error("Error occurred while getting spaces")
		}
	},
}

func init() {
	getCmd.AddCommand(spacesCmd)
}

func getSpaces(config internal.Config) error {
	spaces, err := synapse.GetSpaces(config, logger)
	if err != nil {
		return err
	}

	internal.Print(spaces, false)
	return nil
}
