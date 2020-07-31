package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmdSetup = &cobra.Command{
	Use:   "setup [OPTIONS] [COMMANDS]",
	Short: "Setup the remote agent & save config to specific location",
	Long:  `Initialize the setup process for the remote agent and save the config to a specified location (default is $HOME/.rragent.yaml)`,
	Args:  cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		initConfig()
		writeConfig()

	},
}

func init() {
	cmdSetup.PersistentFlags().StringVarP(&conf.server, "server", "s", "", "domain name for reconness server e.g mydomain.com")
	cmdSetup.MarkFlagRequired("server")
	viper.BindPFlag("server", cmdSetup.PersistentFlags().Lookup("server"))

	cmdSetup.PersistentFlags().StringVarP(&conf.username, "username", "u", "", "username used to connect to reconness server")
	cmdSetup.MarkFlagRequired("username")
	viper.BindPFlag("username", cmdSetup.PersistentFlags().Lookup("username"))

	cmdSetup.PersistentFlags().StringVarP(&conf.password, "password", "p", "", "password used to connect to reconness server ")
	cmdSetup.MarkFlagRequired("password")
	viper.BindPFlag("password", cmdSetup.PersistentFlags().Lookup("password"))

	cmdSetup.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "location to save the config file too (default is $HOME/.rragent.yaml)")

	cmdSetup.PersistentFlags().BoolVarP(&force, "force", "f", false, "If an existing config is found, overwite it (default is $HOME/.rragent.yaml)")

	//cmdSetup.PersistentFlags().IntVar(&maxtasks, "maxtasks", "", "how many tasks to run")

	rootCmd.AddCommand(cmdSetup)

}
