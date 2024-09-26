# server
the Walkyria server is a key value database.

# documentation
https://walkyria.io

# rules
1. The rule number one is, since now, we should NEVER push a commit that didn't passed the test.

# requirements to releases
- v0.3.0? Logfile for reconstruct the mem hashmap in case of reboot
- v0.2.0? Real in mem hashmap
- v0.1.1-alpha standardizing responses using REST standard.
- v0.1.0-alpha better testing, default port and port selection using -p or --port on executing.

# Building Walkyria from source
## Windows
```
go build -o 'YourBinaryName' .\main.go .\adm_routes.go .\consume_routes.go .\db_crud.go .\general_funcions.go
```
## Linux
```
go build -o 'YourBinaryName' ./main.go ./adm_routes.go ./consume_routes.go ./db_crud.go ./general_funcions.go
```
