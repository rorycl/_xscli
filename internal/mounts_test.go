package internal

import (
	"embed"
	_ "embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

//go:embed testdata
var testdata embed.FS

//go:embed testdata/dirA
var testdataDirA embed.FS

func TestMounts(t *testing.T) {

	tests := []struct {
		name       string
		mountName  string
		embeddedFS fs.FS
		dirPath    string
		dirToStat  string
		err        error
	}{
		{
			name:       "embedded fs mount",
			mountName:  "testdata",
			embeddedFS: testdata,
			dirPath:    "",
			dirToStat:  "dirA/dirB",
		},
		{
			name:       "directory fs mount",
			mountName:  "testdata",
			embeddedFS: testdata,
			dirPath:    "./testdata",
			dirToStat:  "dirA/dirB",
		},
		{
			name:       "directory fs mount fail",
			mountName:  "testdata",
			embeddedFS: testdata,
			dirPath:    "./doesNotExist",
			err:        errors.New(`new mount at "./doesNotExist"`),
		},
		{
			name:       "embedded fs mount for dirA",
			mountName:  "testdata/dirA",
			embeddedFS: testdataDirA,
			dirPath:    "",
			dirToStat:  "dirB",
		},
		{
			name:       "directory fs mount for dirA",
			mountName:  "testdata/dirA",
			embeddedFS: testdataDirA,
			dirPath:    "testdata/dirA",
			dirToStat:  "dirB",
		},
	}

	for ii, tt := range tests {
		t.Run(fmt.Sprintf("%d_%s", ii, tt.name), func(t *testing.T) {

			testDir := t.TempDir()
			// uncomment to inspect directory
			// testDir, err := os.MkdirTemp("", "mount_*")
			// if err != nil {
			// 	t.Fatal(err)
			// }

			fm, err := NewFileMount(tt.mountName, tt.embeddedFS, tt.dirPath)
			if err != nil {
				if tt.err != nil {
					if got, want := err.Error(), tt.err.Error(); !strings.Contains(got, want) {
						fmt.Errorf("error got %q want substring %q", got, want)
					}
					return
				}
				t.Fatalf("unexpected New error %v", err)
			}

			stat, err := fs.Stat(fm.FS, tt.dirToStat)
			if err != nil {
				t.Fatalf("could not find dir 'dirB' in 'dirA' at top level of fs: %v", err)
			}
			if !stat.IsDir() {
				t.Errorf("dir 'dirB' in 'dirA' of fs is not a dir: %v", stat.Name())
			}
			err = fm.Materialize(testDir)
			if err != nil {
				t.Errorf("unexpected error %v", err)
			}

			// Given a target of /tmp the materialized output is put in (for example)
			// /tmp/testdata/. To compensate for this the top level of the materialized
			// output is popped.

			matFS := os.DirFS(testDir)
			materializedFS, err := fs.Sub(matFS, tt.mountName)
			if err != nil {
				t.Fatalf("could not submount materialized dir")
			}
			materializedFSAsString, err := PrintFS(materializedFS)
			if err != nil {
				t.Fatal(err)
			}

			mountFSAsString, err := PrintFS(fm.FS)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(materializedFSAsString, mountFSAsString); diff != "" {
				t.Errorf("unexpected difference between materialization and mount:\n%s", diff)
			}

		})
	}
}
