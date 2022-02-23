package fisy

import (
	"os"
	"reflect"
	"testing"
)

func TestOpen(t *testing.T)      {

	tmp, err := os.CreateTemp(os.TempDir(), "zrm_test_file")

	if err != nil {
		t.Fatalf("failed to create temporary file, reason %v", err)
	}

	contents := []byte("The zrm tool is quite powerful")
	filepath := tmp.Name()

	tmp.Write(contents)
	tmp.Close()

	f := open(filepath)

	if f == nil {
		t.FailNow()
	}

	buff := make([]byte, len(contents))

	f.Read(buff)

	if !reflect.DeepEqual(contents, buff) {
		t.Fatalf("wants %v, got %v", string(contents), string(buff))
	}
}

func TestWriteFill(t *testing.T) {

	tmp, err := os.CreateTemp(os.TempDir(), "zrm_test_file")

	if err != nil {
		t.Fatalf("failed to create temporary file, reason %v", err)
	}

	contents := []byte("The zrm tool is quite powerful, zeroing files")
	filepath := tmp.Name()

	tmp.Write(contents)
	tmp.Close()

	WriteFill(filepath, []byte{0})

	f := open(filepath)

	if f == nil {
		t.FailNow()
	}

	buff := make([]byte, len(contents))
	zero := make([]byte, len(contents))

	f.Read(buff)

	if !reflect.DeepEqual(zero, buff) {
		t.Fatalf("wants %v, got %v", string(contents), string(buff))
	}
}

func TestDelete(t *testing.T) {

	tmp, err := os.CreateTemp(os.TempDir(), "zrm_test_file")

	if err != nil {
		t.Fatalf("failed to create temporary file, reason %v", err)
	}

	filepath := tmp.Name()

	tmp.Write([]byte("The zrm tool is quite powerful"))
	tmp.Close()

	if err = Delete(filepath); err != nil {
		t.Fatalf("failed to delete file, reason %v", err)
	}

}
