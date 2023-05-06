# OPCUA Checker
[![Go Report Card](https://goreportcard.com/badge/github.com/obowersa/opcuacheck)](https://goreportcard.com/report/github.com/obowersa/opcuacheck)
## WIP - Simple and light weight way to check OPCUA connections
Golang CLI which allows for retreieving variables over OPCUA as well as checking basic connectivity information.

Still work in progress. Main areas which need improvement:
- Fleshing out the README
- Supporting more OPCUA functions
- Build/release process ( currently relies on folk reading the makefile )
- Expand the testing to cover the CLI component ( similar to what I did for the wfwiki project ) 
- Add in dockerfile support ( with container structure tests/etc )

## To run/build
- Ensure you have make installed
- Ensure you have go 1.20 installed
- Run make-run-all

To build, either run make build-all or specify your platform 

Type make for more help info

Once built, execute the output in bin/ to see the CLI help info
