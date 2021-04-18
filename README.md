# Vega wallet GUI
Multiplatform visual tool to simplify the usage of Vega wallet. 
It is a wrapper for the [official Vega wallet](https://github.com/vegaprotocol/go-wallet/)
The application is written in Golang and takes advantage of astielectron.
## Install
Get the latest version for your OS in the release page.
## Build
If you want to build it from source follow this procedure:
- Make sure you installed Golang
- Clone the repository locally
- Install astielectron
```shell
go get -u github.com/asticode/go-astilectron
```
- Install astielectron bundler:
```shell
go get -u github.com/asticode/go-astilectron-bundler/...
```
- Build the project
```shell
astilectron-bundler
```
## Debug mode
If you want to run the application in debug mode you can run it using the -d flag.
Once the application starts you can open dev tools by pressing CTRL+D or selecting debug in the file menu.

# About Vega
[Vega](https://vega.xyz) is a protocol for creating and trading derivatives on a fully decentralised network. The network, secured with proof-of-stake, will facilitate fully automated, end-to-end margin trading and execution of complex financial products. Anyone will be able to build decentralised markets using the protocol.

Read more at [https://vega.xyz](https://vega.xyz).
