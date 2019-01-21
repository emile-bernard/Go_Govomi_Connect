# go_govomi_connect
This project contains a basic Go program to connect Govomi client to vcenter.

### Install Go

Download from: https://golang.org/dl/

```
$ tar -C /usr/local -xzf go1.11.2.linux-amd64.tar.gz
```

```
$ export PATH=$PATH:/usr/local/go/bin
```

### Install govomi package
```
$ go get -u github.com/vmware/govmomi
```

### List all environnement variables
```
$ printenv
```

### Set environnement variables
```
$ export GOVOMIHOST=vcenter.ulaval.ca/sdk
```

```
$ export GOVOMIUID=userid
```

```
$ export GOVOMIPWD='password'
```

### Build and run the program
```
$ go build
```

```
$ ./go_govomi_connect
```

### Links
[Govomi Library](https://github.com/vmware/govmomi)

[Go Downloads](https://golang.org/dl/)

[Getting started with Go](https://golang.org/doc/install?download=go1.11.2.linux-amd64.tar.gz)
