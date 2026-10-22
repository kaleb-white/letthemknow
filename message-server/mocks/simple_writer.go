package mocks

import "sync"

var store []byte = []byte{}
var m sync.Mutex = sync.Mutex{}

func setSimpleReaderWriterValue(p []byte) {
	m.Lock()
	store = p	
	m.Unlock()
}

func readSimpleReadWriterValueIntoStore(p *[]byte) {
	*p = store
}

type SimpleReaderWriter struct {
}

func (s SimpleReaderWriter) Read(p []byte) (n int, err error) {
	readSimpleReadWriterValueIntoStore(&p)
	return len(p), nil
}

func (s SimpleReaderWriter) Write(p []byte) (n int, err error) {
	setSimpleReaderWriterValue(p)
	return len(p), nil
}

func (s *SimpleReaderWriter) GetStore() []byte {
	return store
}
