package logrotator

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDeleteLogFiles(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "logrotator_")
	if err != nil {
		t.Fatalf("tmp dir creation failed. %s\n", err)
	}
	t.Logf("tmp dir: %s\n", tmpDir)
	filename := filepath.Join(tmpDir, "test.log")
	savedCount := 4
	rotator := New(filename, savedCount-1)
	if rotator == nil {
		t.Fatal("New rotator failed.\n")
	}

	for i := 0; i < savedCount; i++ {
		tmSuffix := time.Now().AddDate(0, 0, -savedCount+i).Format("2006-01-02")
		backupFile := filename + "." + tmSuffix
		f, err := os.OpenFile(backupFile, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			t.Fatalf("create test backup file error: %s\n", err)
		}
		t.Logf("create test backup file: %s\n", backupFile)
		f.Close()
	}

	delTime := time.Now().AddDate(0, 0, -savedCount+1).Format("2006-01-02")
	delTime += "T00:00:00+08:00"
	tm, _ := time.Parse(time.RFC3339, delTime)
	err = rotator.deleteLogFiles(tm.Unix())
	t.Logf("delete backup file before: %s\n", delTime)
	if err != nil {
		t.Fatalf("deleteLogFiles error: %s\n", err)
	} else {
		for i := 1; i < savedCount; i++ {
			suffix := time.Now().AddDate(0, 0, -savedCount+i).Format("2006-01-02")
			if _, err := os.Stat(filename + "." + suffix); err != nil {
				t.Errorf("file %s should exist, but there is error: %s", filename+"."+suffix, err)
			}
		}
		suffix := time.Now().AddDate(0, 0, -savedCount).Format("2006-01-02")
		if _, err := os.Stat(filename + "." + suffix); err == nil {
			t.Errorf("file %s should del, but it is alive", filename+"."+suffix)
		}
	}

	// clean temp dir
	rotator.fp.Close()
	rotator.fp = nil
	os.RemoveAll(tmpDir)
}
