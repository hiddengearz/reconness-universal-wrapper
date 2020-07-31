package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var cmdExec = &cobra.Command{
	Use:   "exec",
	Short: "Execute a command",
	Long:  `Execute the command specified`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, arguments []string) {

		home, err := homedir.Dir() //get home directory windows/linux
		if err != nil {
			er(err)
		}

		if cfgFile != "" { //if user specified a config use the config
			if debug {
				fmt.Println("Using config located at: " + cfgFile)
			}

			initConfig()
		} else if fileExists(home + "/.rwrapper.yaml") { //else try the config in the default location
			if debug {
				fmt.Println("Using default config located at: " + home + "/.rwrapper.yaml")
			}
			initConfig()
		} else { //otherwise no config specified/found
			er("No config found at " + home + "/.rwrapper.yaml. Please specify a config location")
		}

		conf = readConfig() //read the config
		conf.server = "https://" + conf.server

		args := strings.Fields(arguments[0]) //split the argument given

		for i, arg := range args { //get JWT and subdomains
			if arg == "*subdomains" {
				if api == "" {
					er("You must provide an api for -a")
					os.Exit(1)
				}
				// Do the authentication and obtain the jwt
				jwt := Auth(conf.server, authApi, conf.username, conf.password)
				// Get the token to allow us send auth request
				token := GetToken(jwt)
				subdomainFile = ExportSubdomains(conf.server, api, token)
				args[i] = subdomainFile
			}
		}

		for i, arg := range args { //this for loop is done seperately so we don't create temp files if the above errors out
			if arg == "*outputFile" { //Replace *outputFile with a tempfile
				outputFile = CreateTmpFile()
				args[i] = outputFile
			}
			if strings.Contains(arg, "*outputDir") { //replace *outputFolder with a temp folder
				outputDir = CreateTmpDir()
				if strings.Contains(arg, "*.") { //if a wildcard for the file extension to read is found, save it
					tmp := strings.Split(arg, "*.")
					fileExt = tmp[1]
				}
				args[i] = outputDir
				//time.Sleep(60 * time.Second)
			}
		}

		command := exec.Command(args[0], args[1:]...) //execute command

		stderr, _ := command.StderrPipe()
		if err := command.Start(); err != nil {
			log.Fatal(err)
		}

		if !silent { //Print the commands output
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}
		}
		command.Wait()

		if outputFile != "" { //if *outputFile was used print the content of the temp file
			output := ReadFile(outputFile)
			for _, line := range output {
				fmt.Println(line)
			}
		}

		if outputDir != "" { //if *outputDir was used, print the content of the directory if....
			files, err := ioutil.ReadDir(outputDir)
			if err != nil {
				er(err)
			}

			if fileExt == "*" { //if the file ext is * print all files
				for _, file := range files {
					content := ReadFile(outputDir + "/" + file.Name())
					for _, line := range content {
						fmt.Println(line)
					}

				}
			} else { // print the files that only contain the spcified file extension
				for _, file := range files {
					if strings.Contains(file.Name(), fileExt) {
						content := ReadFile(outputDir + "/" + file.Name())
						for _, line := range content {
							fmt.Println(line)
						}
					}
				}
			}

		}

	},
}

func init() {
	rootCmd.AddCommand(cmdExec)
	cmdExec.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "location of the config file (default is $HOME/.rwrapper.yaml)")
	cmdExec.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug mode")
	cmdExec.PersistentFlags().BoolVar(&silent, "silent", false, "Don't print the commands output")

	cmdExec.PersistentFlags().StringVarP(&api, "api", "a", "", "api endpoint")
	//viper.BindPFlag("api", cmdSetup.PersistentFlags().Lookup("api"))

}

func Auth(url string, authApi string, username string, password string) string {

	values := map[string]string{"UserName": username, "Password": password}
	jsonValue, _ := json.Marshal(values)

	req, err := http.NewRequest("POST", url+"/"+authApi, bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		er(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return string(body)
}

//GetToken Retrieves the JWT for authentication
func GetToken(jwt string) string {

	in := []byte(jwt)
	var raw map[string]interface{}
	if err := json.Unmarshal(in, &raw); err != nil {
		panic(err)
	}
	return raw["auth_token"].(string)
}

//ExportSubdomains exports subdomains to a temporary file
func ExportSubdomains(url string, subdomainApi string, token string) (filepath string) {

	req, err := http.NewRequest("GET", url+"/"+subdomainApi, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	subDomains := strings.Split(string(bodyBytes), ",") //csv format to txt
	tmpFile, err := ioutil.TempFile(os.TempDir(), "tmp-")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}
	//defer os.Remove(tmpFile.Name())

	for _, data := range subDomains {
		text := []byte(data + "\n")
		if _, err = tmpFile.Write(text); err != nil {
			log.Fatal("Failed to write to temporary file", err)
		}
	}

	return tmpFile.Name()
}

func GenerateRandomString(length int) (randomString string) {
	rand.Seed(time.Now().UnixNano())
	chars := []rune(
		"abcdefghijklmnopqrstuvwxyz" +
			"0123456789")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String() // E.g. "ExcbsVQs"

	return str
}

func CreateTmpFile() string {
	tmpFile, err := ioutil.TempFile(os.TempDir(), GenerateRandomString(4))
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}
	return tmpFile.Name()
}

func CreateTmpDir() string {
	dname, err := ioutil.TempDir("", GenerateRandomString(4))
	if err != nil {
		er(err)
	}
	return dname
}
