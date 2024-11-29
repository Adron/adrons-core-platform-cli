/*
Copyright © 2019 Adron Hall <adron@thrashingcode.com>

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
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "The 'config' subcommand is for use in management of configuration.",
	Long: func() string {
		baseDesc := `The 'config' subcommand is for use in management of configuration. It can be used, in combination with the
other subcommands 'add', 'update', 'view', and 'delete'.

Configuration should be provided via config.json file. See config.sample.json for an example.`

		// Only append debug info if DEBUG is true
		if viper.GetBool("DEBUG") {
			debugInfo := fmt.Sprintf(`

Debug Information:
----------------
Postgres URL: %s
Username: %s
Debug Mode: %v`,
				viper.GetString("POSTGRES_URL"),
				viper.GetString("USERNAME"),
				viper.GetBool("DEBUG"))

			return baseDesc + debugInfo
		}

		return baseDesc
	}(),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Current Configuration:")
		fmt.Println("---------------------")

		// Get all settings from viper
		allSettings := viper.AllSettings()

		// Print each key-value pair
		for key, value := range allSettings {
			fmt.Printf("%s: %v\n", key, value)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.PersistentFlags().StringP("key", "k", "", "The key for the key value set to add to the configuration.")
	configCmd.PersistentFlags().StringP("value", "v", "", "The value for the key value set to add to the configuration.")

	// Setup Viper for JSON config
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("json")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // look for config in the working directory

	// Read the JSON config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No config.json file found. Please copy config.sample.json to config.json and modify as needed.")
		} else {
			fmt.Printf("Error reading config file: %s\n", err)
		}
	}
}
