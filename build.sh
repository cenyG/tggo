rm -f tggo

unameOut="$(uname -s)"
case "${unameOut}" in
    Darwin*)    goos=darwin;;
    Linux*)     goos=linux;;
    *)          goos=linux
esac

GOOS=${goos} GOARCH=amd64 go build -mod vendor -o tggo ./main.go