package rsa

import "io/ioutil"

type Reader interface {
	ReadPem() ([]byte, error)
	Read() ([]byte, error)
}

type FileReader struct {
	FileName string
}

func(f *FileReader) ReadPem() ([]byte, error) {
	var pem []byte

	pem,err:=ioutil.ReadFile(f.FileName)
	if err!=nil{
		return pem,err
	}
	return pem,err
}
