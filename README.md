# Reconness Universal Wrapper

Too tired, finish this tomorrow...

Example

`go run wrapper.go exec "python3.8 [redacted]/tools/OneForAll/oneforall.py --target domains.txt --path *outputDir/*.txt run" --silent`

`go run wrapper.go exec "go [redacted]/tools/naabu -hL *subdomains -silent" -a api/targets/exportSubdomains/test/test`

`go run main.go exec "massdns -r [redacted]/tools/massdns/lists/resolvers.txt *subdomains -w *output -o S" -a api/targets/exportSubdomains/.mil/.mil`
