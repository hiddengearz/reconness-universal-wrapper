package cmd

import (
	"bufio"
	"fmt"
	"os"

	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type config struct {
	username string `yaml:"username"`
	password string `yaml:"password"`
	server   string `yaml:"server"`
}

var (
	conf = &config{}
	// Used for flags.
	cfgFile       string
	authApi       = "api/Auth/Login"
	api           string
	force         bool = false
	debug         bool = false
	silent        bool = false
	outputFile    string
	outputDir     string
	fileExt       string
	subdomainFile string

	//maxtasks int
	rootCmd = &cobra.Command{Use: "app"}
)

func Execute() error {
	return rootCmd.Execute()
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func GetPathandFile(fullPath string) (theDir string, thefile string) {
	dir, file := filepath.Split(fullPath)
	return dir, file
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		viper.SetConfigType("yaml")

	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			er(err)
		}
		cfgFile = home + "/" + ".rwrapper.yaml"
		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".rwrapper")
		viper.SetConfigType("yaml")
	}

	//viper.AutomaticEnv()
	/*
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		} else {
			er(err)
		}
	*/
}

func writeConfig() {
	if force {
		err := viper.WriteConfigAs(cfgFile)
		if err != nil {
			er(err)
		}
	} else {
		err := viper.SafeWriteConfigAs(cfgFile)
		if err != nil {
			er(err)
		}
	}
	fmt.Println("Config written to " + cfgFile)

}

func readConfig() *config {

	err := viper.ReadInConfig()

	if err != nil {
		er(err)
	}

	conf := &config{}
	err = viper.Unmarshal(&conf)
	if err != nil {
		er(err)
	}
	conf.server = viper.GetString("server")
	conf.username = viper.GetString("username")
	conf.password = viper.GetString("password")

	return conf
}

//ReadFile reads the file and saves to an array
func ReadFile(filePath string) []string {
	file, err := os.Open(filePath)
	var content []string

	if err != nil {
		er(err)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content = append(content, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		er(err)
		return nil
	}

	return content
}
