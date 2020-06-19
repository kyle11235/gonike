# linux
GOOS=linux GOARCH=amd64 go build -o gonike_linux

# windows
GOOS=windows GOARCH=amd64 go build -o gonike_win

# adrwin
GOOS=darwin GOARCH=amd64 go build -o gonike_mac
