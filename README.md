# The Things Network Go SDK

[![Build Status](https://travis-ci.org/TheThingsNetwork/go-app-sdk.svg?branch=master)](https://travis-ci.org/TheThingsNetwork/go-app-sdk) [![Coverage Status](https://coveralls.io/repos/github/TheThingsNetwork/go-app-sdk/badge.svg?branch=master)](https://coveralls.io/github/TheThingsNetwork/go-app-sdk?branch=master) [![GoDoc](https://godoc.org/github.com/TheThingsNetwork/go-app-sdk?status.svg)](https://godoc.org/github.com/TheThingsNetwork/go-app-sdk)

![The Things Network](https://thethings.blob.core.windows.net/ttn/logo.svg)

## Usage

참고 : https://www.thethingsnetwork.org/docs/applications/golang/quick-start.html

Assuming you're working on a project `github.com/your-username/your-project`:

```
sudo apt update
sudo apt upgrade
sudo apt-get install golang

mkdir -p $GOPATH/src/github.com/YOUR_USERNAME/ttn-app
cd $GOPATH/src/github.com/YOUR_USERNAME/ttn-app

go get -u github.com/TheThingsNetwork/go-app-sdk/...

sudo apt-get install govendor
cd $GOPATH/src/github.com/TheThingsNetwork/go-app-sdk
govendor init
govendor sync
```

main.go upload!

See the examples [on GoDoc](https://godoc.org/github.com/TheThingsNetwork/go-app-sdk#example-package).

## License

Source code for The Things Network is released under the MIT License, which can be found in the [LICENSE](LICENSE) file. A list of authors can be found in the [AUTHORS](AUTHORS) file.
