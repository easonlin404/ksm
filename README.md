# ksm


[![Build Status](https://travis-ci.org/easonlin404/ksm.svg)](https://travis-ci.org/easonlin404/ksm)
[![codecov](https://codecov.io/gh/easonlin404/ksm/branch/master/graph/badge.svg)](https://codecov.io/gh/easonlin404/ksm)
[![Go Report Card](https://goreportcard.com/badge/github.com/easonlin404/ksm)](https://goreportcard.com/report/github.com/easonlin404/ksm)
[![GoDoc](https://godoc.org/github.com/easonlin404/ksm?status.svg)](https://godoc.org/github.com/easonlin404/ksm)

FairPlay Streaming Key Security Module written in Go (Golang).

## Usage

### Start using it

Download and install it:

```bash
$ go get github.com/easonlin404/ksm
```

### Verify ckc
Perform verification utility `verify_ckc` to test KSM implementation.
```
testdata/verify_ckc -s testdata/FPS/spc1.bin -c testdata/FPS/ckc1.bin
```
