/*
Copyright © 2024 Nected<dev@nected.ai>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/nected/sanchaalak/src/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var selectedCmd string

var supportedFlags = map[string]func(){
	"config": generateConfigFile,
	"test": func() {
		fmt.Println("test")
	},
}

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if fn, ok := supportedFlags[selectedCmd]; ok {
			fn()
		} else {
			fmt.Println("Invalid module selected")
		}
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	helpMessage := fmt.Sprintf("Generator for supported modules. Supported modules are:")
	for flag := range supportedFlags {
		helpMessage = fmt.Sprintf("%s \n - %s", helpMessage, flag)
	}
	generateCmd.Flags().StringVar(&selectedCmd, "module", "", helpMessage)
	rootCmd.AddCommand(generateCmd)
}

func generateConfigFile() {
	config := config.NewConfig()
	config.GenerateDefaults()

	data, err := yaml.Marshal(config)
	if err != nil {
		fmt.Println(err)
	}
	f, err := os.Create(viper.ConfigFileUsed())
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		fmt.Println(err)
	}
}
