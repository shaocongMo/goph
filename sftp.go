// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package goph

import (
	"fmt"
	"io"
	"os"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// Upload local file to remote.
func Upload(c *ssh.Client, src string, dest string) (err error) {
	if session, err := c.NewSession(); err == nil {
		defer session.Close()
		go func() {
			fmt.Println("start upload file:", src)
			Buf := make([]byte, 1024)
			w, _ := session.StdinPipe()
			defer w.Close()
			File, _ := os.Open(src)
			info, _ := File.Stat()
			fmt.Fprintln(w, "C0644", info.Size(), info.Name())
			for {
				n, err := File.Read(Buf)
				fmt.Fprint(w, string(Buf[:n]))
				if err != nil {
					if err == io.EOF {
						return
					} else {
						fmt.Println("Upload File error:", err)
						return
					}
				}
			}
		}()
		if err := session.Run(fmt.Sprintf("/usr/bin/scp -qrt %s", dest)); err != nil {
			if err != nil {
				if err.Error() != "Process exited with: 1. Reason was: ()" {
					fmt.Println(err.Error())
				}
			}
			fmt.Printf("upload file %s finish\n", src)
		}
	}
	return
}

// Download remote file to local.
func Download(c *ssh.Client, src string, dest string) (err error) {

	client, err := sftp.NewClient(c)

	if err != nil {
		return
	}

	defer client.Close()

	destFile, err := os.Create(dest)

	if err != nil {
		return
	}

	defer destFile.Close()

	srcFile, err := client.Open(src)

	if err != nil {
		return
	}

	defer srcFile.Close()

	if _, err = io.Copy(destFile, srcFile); err != nil {
		return
	}

	return destFile.Sync()
}
