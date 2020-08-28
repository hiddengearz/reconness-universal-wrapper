# Reconness Universal Wrapper

Reconness Universal Wrapper is a wrapper/framework that assists with creating [Reconness Agents](https://github.com/reconness/reconness-agents) within the [Reconness Web Application](https://github.com/reconness/reconness). It provides additional functionality that is not currently possibile in reconness and can only be added by modifying the source code of the agent or making a wrapper for it. It does this by allowing you to use "subsitutions" when executing commands.

For example [Massdns](https://github.com/blechschmidt/massdns) requires a file with a list of domains to be used as an argument. In reconness there is no way to create that file thus you need to make your own wrapper that calls reconness's API and creates a file as the argument. Instead you can use the universal wrapper:

`./bin/massdns -r lists/resolvers.txt domains.txt`

`reconness-universal-wrapper exec "./bin/massdns -r lists/resolvers.txt *subdomains -w *outputFile" -a api/targets/exportSubdomains/{{target}}/{{rootDomain}} --silent`

### Subsitutions

Subsitutions can be used to replace an argument within the command you're executing. The current subsitutions are:

| subsitutions  |E xplination   | 
|---|---|
| *subdomains  | Replaces a file input with a list of all subdomains from the specified API endpoint  |
| *outputFile | Replaces the output file with a temp file & displays it's content in stdout, so that the reconness agent can parse it |
|  \*outputDir/*.ext | Replaces the output directory with a temp directory & will display every file with the specified extension to stdout e.g '*outputDir/\*.txt' will print all `txt` files in the temp directory, while `*outputDir/*.*` will print all files in the output directory    |

### Flags
| Flags  | Explination   | 
|---|---|
| -c, --config  | (optional)location of the config file (default is $HOME/.rwrapper.yaml) |
| --debug | Enable debug mode  |
| -h, --help | Display help menu  |
| --silent |  Don't print the commands output |

### Examples

`reconness-universal-wrapper exec "go [redacted]/tools/naabu -hL *subdomains -silent" -a api/targets/exportSubdomains/{{target}}/{{rootDomain}}`

`reconness-universal-wrapper exec "python3.8 [redacted]/tools/OneForAll/oneforall.py --target *subdomains --path *outputDir/*.txt run" --silent`

`reconness-universal-wrapper exec "massdns -r [redacted]/tools/massdns/lists/resolvers.txt *subdomains -w *outputFile -o S" -a api/targets/exportSubdomains/{{target}}/{{rootDomain}}`


### Installation

#### Non-docker

`go get -u github.com/hiddengearz/reconness-universal-wrapper`

`reconness-universal-wrapper setup -u <reconness username> -p <reconness password> -s <https://reconness.mydomain.com:8080>`

#### Docker:

Modify your [Reconness Dockerfile](https://github.com/reconness/reconness/blob/master/src/Dockerfile) by adding the following to the Agents dependencies section.

```
RUN apt-get update && apt-get install -y git
RUN apt-get install -y wget
RUN wget https://dl.google.com/go/go1.13.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.13.linux-amd64.tar.gz
RUN echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.profile
RUN /usr/local/go/bin/go get -u github.com/hiddengearz/reconness-universal-wrapper
RUN cd /root/go/bin/ && ./reconness-universal-wrapper setup -u <reconness username> -p <reconness password> -s <https://reconness.mydomain.com:8080>
```

If you change your reconness username, password or domain name you'll need to update the config aswell. You can do this by typing the following in the cli/docker container `cd ./root/go/bin/reconness-universal-wrapper setup -u <new username> -p <new password> -s <new domain> --force`

### Setup flags
| Flags  | Explination   | 
|---|---|
|-c, --config string |     location to save the config file too (default is $HOME/.rwrapper.yaml) |
|-f, --force |             If an existing config is found, overwite it (default is $HOME/.rwrapper.yaml) |
|-h, --help |              help for setup |
|-p, --password string |   password used to connect to reconness server |
|-s, --server string |     domain name for reconness server e.g <https://reconness.mydomain.com:8080> |
|-u, --username string |   username used to connect to reconness server |
