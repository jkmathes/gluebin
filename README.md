# gluebin

_*This works only on ELF binaries*_

To build:

```
CGO_ENABLED=0 go build -o gluebin cmd/gluebin/main.go
```

To run:

```
gluebin dynamic_binary static_binary
```

For example:
```
gluebin /usr/bin/curl ./mycurl
```

