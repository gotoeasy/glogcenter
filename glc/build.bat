SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
echo now the CGO_ENABLED:
 go env CGO_ENABLED

echo now the GOOS:
 go env GOOS

echo now the GOARCH:
 set h=%time:~0,2%
 set h=%h: =0%
 set version=%date:~0,4%%date:~5,2%%date:~8,2%%h%%time:~3,2%%time:~6,2%
 go env GOARCH
 go build -mod=vendor -o glc_amd64 -ldflags  "-s -w" main.go
