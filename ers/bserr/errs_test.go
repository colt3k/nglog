package bserr

import (
	"errors"
	"testing"
)

func TestBSErr_Err(t *testing.T) {
	stders.Err(errors.New("test error"), false, nil...)
}
func TestBSErr_NotErr(t *testing.T) {
	stders.NotErr(errors.New("test error"), nil...)
}
func TestBSErr_NoPrintErr(t *testing.T) {
	stders.NoPrintErr(errors.New("test error"))
}
func TestBSErr_StopErr(t *testing.T) {
	stders.StopErr(errors.New("test error"), nil...)
}
