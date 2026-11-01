package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	p := "/mnt/us/kron-fs-probe.txt"
	fmt.Println("FS probe start")

	f, err := os.OpenFile(
		p,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		fmt.Println("OPEN ERROR:", err)
		os.Exit(1)
	}

	fmt.Println("OPEN OK")

	info, err := f.Stat()
	if err != nil {
		fmt.Println("STAT ERROR:", err)
		f.Close()
		os.Exit(1)
	}

	fmt.Println("STAT OK:", info.Size())

	n, err := f.Write([]byte("Kindle FS probe\n"))
	fmt.Println("WRITE:", n, err)

	err = f.Close()
	fmt.Println("CLOSE:", err)

	abs, err := filepath.Abs(p)
	fmt.Println("ABS:", abs, err)

	fmt.Println("FS probe done")
}
