package test

import (
	"KittyStager/internal/task/ps"
	"github.com/go-playground/assert/v2"
	"testing"
)

var (
	psPpid int
	psPid  int
	psName string
	pList  []ps.Process
)

func init() {
	psPpid = 12
	psPid = 13
	psName = "test"
}

func TestGetPs(t *testing.T) {
	t.Parallel()
	p := ps.NewProcess(psPpid, psPid, psName)
	assert.Equal(t, p.Ppid, psPpid)
	assert.Equal(t, p.Pid, psPid)
	assert.Equal(t, p.Name, psName)
}

func TestMarshallUnmashallPs(t *testing.T) {
	t.Parallel()
	// given
	p := ps.NewProcess(psPpid, psPid, psName)
	pList = append(pList, *p)
	pl := ps.NewProcessList(pList)
	// when
	b, err := pl.MarshallProcessList()
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	err = pl.UnmarshallProcessList(b)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
}
