cls
go env -w GOOS=windows
go env -w GOARCH=amd64
go build -o movie.exe
.\movie.exe