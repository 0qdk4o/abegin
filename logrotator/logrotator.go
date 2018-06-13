package logrotator

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// RotateWriter represents custom log rotator
type RotateWriter struct {
	filename     string
	fp           *os.File
	backupCount  int
	nextRotateTs int64
}

// New return new instance of Rotator
func New(filename string) *RotateWriter {
	r := &RotateWriter{
		filename:    filename,
		backupCount: 8,
	}
	r.updateNextRotateTs()

	if _, err := os.Stat(r.filename); err == nil {
		r.fp, err = os.OpenFile(r.filename, os.O_RDWR|os.O_APPEND, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "OpenFile error: %s\n", err)
			return nil
		}
	} else {
		var err error
		r.fp, err = os.Create(r.filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Create error: %s\n", err)
			return nil
		}
	}
	return r
}

func (r *RotateWriter) updateNextRotateTs() {
	timeString := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	timeString += "T00:00:00+08:00"
	tm, _ := time.Parse(time.RFC3339, timeString)
	r.nextRotateTs = tm.Unix()
}

func (r *RotateWriter) Write(p []byte) (n int, err error) {
	if r.isNeedRotation() {
		r.updateNextRotateTs()
		err := r.rotate()
		if err != nil {
			return 0, err
		}
	}
	return r.fp.Write(p)
}

func (r *RotateWriter) isNeedRotation() bool {
	if time.Now().Unix() > r.nextRotateTs {
		fmt.Println("NeedRotation")
		return true
	}
	return false
}

func (r *RotateWriter) deleteLogFiles(delTm int64) error {
	if len(r.filename) < 0 {
		return errors.New("Uninitialized instance of logrotator")
	}
	abPath, err := filepath.Abs(filepath.Dir(r.filename))
	if err != nil {
		return err
	}
	f, err := os.Open(abPath)
	if err != nil {
		return err
	}
	names, err := f.Readdirnames(-1)
	if err != nil {
		return err
	}
	var delErr error
	for _, file := range names {
		dateSuffix := strings.Split(file, ".")
		timeString := dateSuffix[len(dateSuffix)-1]
		timeString += "T00:00:00+08:00"
		tm, _ := time.Parse(time.RFC3339, timeString)
		if tm.Unix() < delTm {
			delErr = os.Remove(filepath.Join(abPath, file))
		}
	}
	return delErr
}

func (r *RotateWriter) rotate() error {

	suffix := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	if r.fp != nil {
		err := r.fp.Close()
		r.fp = nil
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat(r.filename); err == nil {
		err = os.Rename(r.filename, r.filename+"."+suffix)
		if err != nil {
			return err
		}
	}
	var err error
	r.fp, err = os.Create(r.filename)
	if err != nil {
		return err
	}

	return r.deleteLogFiles(time.Now().AddDate(0, 0, -r.backupCount).Unix())
}
