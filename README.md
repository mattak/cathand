# cathand

cathand is auto play test framework for Android.

## Install

```shell-session
go get github.com/mattak/cathand/pkg/cathand
go get github.com/mattak/cathand/cmd/cathand
go install github.com/mattak/cathand/cmd/cathand
```

## Usage

### help

```shell-session
cathand -h
```

### record

record user playing behaviour.

```shell-session
cathand record sample.record
```

### compose 

compose playable files for auto playing.

```shell-session
cathand compose sample.record sample.play
```

### play

play from playable files

```shell-session
cathand play sample.play sample.result
```

### split

split recorded video into image segments

```shell-session
cathand split sample.record sample.result
```

### verify

verify image differences of auto play result and initial result

```shell-session
cathand verify sample.record sample.result sample.report
```

