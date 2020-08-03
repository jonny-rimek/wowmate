### build

`CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags='-s -w' -o main . `