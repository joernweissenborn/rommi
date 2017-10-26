package main

import (
	"io"
	"os"
	"time"
)

type cacheFile struct {
	*os.File
	path string
}

func (c *cacheFile) exists() bool {
	_, err := os.Stat(c.path)
	return err == nil
}

func (c *cacheFile) valid() bool {
	f, err := os.Stat(c.path)
	if err != nil {
		return false
	}
	modTime := f.ModTime()
	return time.Since(modTime) < cacheValidTime
}

func (c *cacheFile) open() (err error) {
	c.File, err = os.OpenFile(c.path, os.O_RDWR|os.O_CREATE, 0777)
	return
}

func (c *cacheFile) write(r io.Reader) (err error) {
	log.Info("Wrting Cache")
	err = c.open()
	if err != nil {
		log.Error("Error writing cache: ", err)
		return
	}
	defer c.Close()
	_, err = io.Copy(c, r)
	if err != nil {
		log.Error("Error writing cache: ", err)
	}
	return
}
