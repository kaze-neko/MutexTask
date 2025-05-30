package main

import (
	"os"
	"os/exec"
	"second/file"
	"sync"
	"testing"
)

func ConcurrentCopy() {
	var file1 = file.NewFile("path1")
	file1.Write([]byte("content1"))
	var file2 = file.NewFile("path2")
	file2.Write([]byte("content2"))
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		file.Copy(file1, file2)
	}()
	go func() {
		defer wg.Done()
		file.Copy(file2, file1)
	}()
	wg.Wait()
}

func TestCopyDeadlock(t *testing.T) {
	if os.Getenv("CRASHED") == "1" {
		ConcurrentCopy()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestCopy")
	cmd.Env = append(os.Environ(), "CRASHED=1")
	var err = cmd.Run()
	if err != nil {
		t.Errorf("Возникло состояние deadlock, когда ожидалось его отсутствие.")
	}
}
