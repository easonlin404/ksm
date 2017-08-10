package rsa

import (
	"io/ioutil"
)

type Reader interface {
	ReadPem() ([]byte, error)
}

var (
	_ Reader = FileReader{}
	_ Reader = TextReader{}
)


type FileReader struct {
	FileName string
}

func(f FileReader) ReadPem() ([]byte, error) {
	var pem []byte
	pem,err:=ioutil.ReadFile(f.FileName)
	if err!=nil{
		return pem,err
	}
	return pem,err
}

type TextReader struct {
	Pem []byte
}

func(t TextReader) ReadPem() ([]byte, error) {
	if len(t.Pem) > 0 {
		return t.Pem,nil
	}
	return nil,nil

}
