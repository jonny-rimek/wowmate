### build

`CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags='-s -w' -o main . `

create new go.mod file `go mod init github.com/jonny-rimek/wowmate/services/converter`