# Linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./dist/Linux64/ReleaseDelivr main.go

# Windows
#CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./dist/Win64/ReleaseDelivr.exe main.go

# MacOs
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./dist/Mac64/ReleaseDelivr main.go

cd ./dist/Mac64 || exit

./ReleaseDelivr