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
func New(filename string, backCount int) *RotateWriter {
	r := &RotateWriter{
		filename:    filename,
		backupCount: backCount,
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
		if !isYMDFormat(timeString) {
			continue
		}
		basename := strings.Join(dateSuffix[:len(dateSuffix)-1], ".")
		if strings.Compare(basename, filepath.Base(r.filename)) != 0 {
			continue
		}

		timeString += "T00:00:00+08:00"
		tm, _ := time.Parse(time.RFC3339, timeString)
		if tm.Unix() < delTm {
			delErr = os.Remove(filepath.Join(abPath, file))
		}
	}
	return delErr
}

// isYMDForamt tell a strings such as 2018-06-14 is date format or not
func isYMDFormat(ymd string) bool {
	source := []byte(ymd)
	if len(source) != 10 {
		return false
	}
	isNumberic := func(c byte) bool {
		if c < 48 || c > 57 {
			return false
		}
		return true
	}
	for i, c := range source {
		if i == 4 || i == 7 {
			if c != 45 {
				return false
			}
			continue
		}
		if !isNumberic(c) {
			return false
		}
	}
	return true
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

	delBefore := time.Now().AddDate(0, 0, -r.backupCount).Format("2006-01-02")
	delBefore += "T00:00:00+08:00"
	tm, _ := time.Parse(time.RFC3339, delBefore)
	return r.deleteLogFiles(tm.Unix())
}
