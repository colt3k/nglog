package ng

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//******************* File APPENDER ********************

type FileAppender struct {
	*OutAppender
	file           *os.File
	append         bool
	buffered       bool
	bufferSize     int
	buf            *bufio.Writer
	fileName       string
	immediateFlush bool
	FilePerm       os.FileMode //default 0777
	UserID         int
	GroupID        int
}

func NewFileAppender(filter, fileName, name string, bufferSize int) (*FileAppender, error) {
	var err error
	oa := newOutAppender(filter, name)
	t := new(FileAppender)
	t.OutAppender = oa
	t.disableColor = true
	t.append = true
	t.FilePerm = os.ModePerm

	t.fileName = fileName
	//Ensure directory structure exists or create it
	if strings.Index(t.fileName, string(filepath.Separator)) > -1 {
		dir, _ := filepath.Split(t.fileName)
		err = os.MkdirAll(dir, t.FilePerm)
		if err != nil {
			return nil, err
		}
	}
	if t.append {
		t.file, err = os.OpenFile(t.fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			return nil, err
		}
	} else {
		t.file, err = os.OpenFile(t.fileName, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return nil, err
		}
	}

	if t.UserID > 0 && t.GroupID > 0 {
		err = os.Chown(t.fileName, t.UserID, t.GroupID)
		if err != nil {
			//os.Stdout.Write([]byte(fmt.Sprint(err)))
			return nil, err
		}
	}
	t.immediateFlush = true
	t.buffered = true
	t.bufferSize = func() int {
		if bufferSize == 0 {
			return 8192
		} else {
			return bufferSize
		}
	}()
	if t.buffered {
		t.buf = bufio.NewWriterSize(t.file, t.bufferSize)
	}

	return t, nil
}
func (f *FileAppender) Name() string {
	if len(f.name) > 0 {
		return f.name
	}
	return fmt.Sprintf("%T", f)
}
func (f *FileAppender) DisableColor() bool {
	return f.disableColor
}
func (f *FileAppender) Applicable(msg string) bool {
	if f.filter == "*" {
		return true
	}
	if strings.Index(msg, f.filter) > -1 {
		return true
	}
	return false
}

func (f *FileAppender) Process(msg []byte) {
	if f.buffered {
		f.buf.Write(msg)
		if f.immediateFlush {
			f.buf.Flush()
		}
	} else {
		f.Out.Write(msg)
	}
}
