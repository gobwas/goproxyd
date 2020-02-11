package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"hash"
	"io"
	"io/ioutil"
	"mime"
	"os"
	"path"
	"strings"
	"time"

	"github.com/goproxy/goproxy"
)

type cacher struct {
	root     string
	readonly bool
}

// NewHash returns a new instance of the `hash.Hash` used to compute the
// checksums of the caches in the underlying cacher.
func (c *cacher) NewHash() hash.Hash {
	return md5.New()
}

// SetCache sets the c to the underlying cacher.
//
// It is the caller's responsibility to close the c.
func (c *cacher) SetCache(ctx context.Context, item goproxy.Cache) (err error) {
	if c.readonly {
		return nil
	}
	filename := path.Join(c.root, item.Name())
	if err := os.MkdirAll(path.Dir(filename), 0755); err != nil {
		return err
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	defer item.Close()

	_, err = io.Copy(file, item)

	return err
}

// Cache returns the matched `Cache` for the name from the underlying
// cacher. It returns the `ErrCacheNotFound` if not found.
//
// It is the caller's responsibility to close the returned `Cache`.
func (c *cacher) Cache(ctx context.Context, name string) (item goproxy.Cache, err error) {
	filename := path.Join(c.root, name)
	stat, err := os.Stat(filename)
	if err != nil {
		return nil, goproxy.ErrCacheNotFound
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	md5sum := md5.Sum(data)

	return &cache{
		r:        bytes.NewReader(data),
		size:     int64(len(data)),
		name:     name,
		mime:     mimeType(filename),
		modtime:  stat.ModTime(),
		checksum: md5sum[:],
	}, nil
}

func mimeType(name string) string {
	switch ext := strings.ToLower(path.Ext(name)); ext {
	case ".info":
		return "application/json; charset=utf-8"
	case ".mod":
		return "text/plain; charset=utf-8"
	case ".zip":
		return "application/zip"
	default:
		return mime.TypeByExtension(ext)
	}
}

type cache struct {
	r        *bytes.Reader
	size     int64
	name     string
	mime     string
	modtime  time.Time
	checksum []byte
}

func (c *cache) Read(b []byte) (int, error) {
	return c.r.Read(b)
}
func (c *cache) Seek(offset int64, whence int) (int64, error) {
	return c.r.Seek(offset, whence)
}
func (c *cache) Close() error {
	return nil
}
func (c *cache) Name() string {
	return c.name
}
func (c *cache) MIMEType() string {
	return c.mime
}
func (c *cache) Size() int64 {
	return c.size
}
func (c *cache) ModTime() time.Time {
	return c.modtime
}
func (c *cache) Checksum() []byte {
	return c.checksum
}
