// Copyright 2018 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package csl

import (
	"github.com/go-kit/kit/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDownload(t *testing.T) {
	if testing.Short() || os.Getenv("TRADEGOV_API_KEY") == "" {
		return
	}

	file, err := Download(log.NewNopLogger(), "")
	if err != nil {
		t.Fatal(err)
	}
	if file == "" {
		t.Fatal("no CSL file")
	}
	defer os.RemoveAll(filepath.Dir(file))

	if !strings.EqualFold("csl.csv", filepath.Base(file)) {
		t.Errorf("unknown file %s", file)
	}
}

func TestDownload_initialDir(t *testing.T) {
	dir, err := ioutil.TempDir("", "iniital-dir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	mk := func(t *testing.T, name string, body string) {
		path := filepath.Join(dir, name)
		if err := ioutil.WriteFile(path, []byte(body), 0644); err != nil {
			t.Fatalf("writing %s: %v", path, err)
		}
	}

	// create each file
	mk(t, "sdn.csv", "file=sdn.csv")
	mk(t, "csl.csv", "file=csl.csv")
	mk(t, "csl.csv", "file=csl.csv")

	file, err := Download(log.NewNopLogger(), dir)
	if err != nil {
		t.Fatal(err)
	}
	if file == "" {
		t.Fatal("no CSL file")
	}

	if strings.EqualFold("csl.csv", filepath.Base(file)) {
		bs, err := ioutil.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		if v := string(bs); v != "file=csl.csv" {
			t.Errorf("csl.csv: %v", v)
		}
	} else {
		t.Fatalf("unknown file: %v", file)
	}
}
