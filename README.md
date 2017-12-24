# Simple Go Playground like.

* Simpel go play ground like on cli.
* Can use your own go environment External/3rd Party libraries.
* Create play ground page like src file.
* Watch changes and run src file.

## Usage

* Go to play-dir
```
cd ${PLAY_DIR_NAME}
```
* New
```
# default make play.n.go(n:1, 2, 3...)
godo new
> create play.go/play.n.go in ${PLAY_DIR_NAME}
```
```
# ${NAME}.n.go
godo new -- ${NAME}
> create ${NAME}.go/${NAME}.n.go in ${PLAY_DIR_NAME}
```
* Run
```
godo run
> search play.${MAX_OF_N}.go in ${PLAY_DIR_NAME}
> goimports -w ${SEARCHED}.go
> go run ${SEARCHED}.go
```
```
godo run -- ${NAME}
> search ${NAME}.${MAX_OF_N}.go in ${PLAY_DIR_NAME}
...
> go run ${SEARCHED}.go
```
```
godo run -- ${GO_FILENAME}
...
> run ${GO_FILENAME}
```
* Watch
```
godo -w
# change ${CHANGED}.go file
> goimports -w ${CHANGED}.go
> go run ${CHANGED}.go
```

## Setup

* go get
```
go get -u gopkg.in/godo.v2/cmd/godo
go get golang.org/x/tools/cmd/goimports
```
* clone
```
git clone https://github.com/ynishi/simple-go-playground.git ${PLAY_DIR_NAME}
```

## Credit and License

* Copyright (c) 2017, Yutaka Nishimura. Licensed under Apache2.0, see LICENSE.
