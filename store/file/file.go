package file

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/peterhellberg/kop/list"
)

func Store(name string, store list.Store) (list.Store, error) {
	if name == "" {
		return store, nil
	}

	dir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(dir, "kop", name+".kop")

	r, err := open(path, os.O_RDONLY)
	if err != nil {
		return nil, err
	}

	for s := bufio.NewScanner(r); s.Scan(); {
		store.Add(s.Text())
	}

	return &cache{
		Store: store,
		path:  path,
	}, nil
}

type cache struct {
	list.Store
	path string
}

func (c *cache) Add(vals ...string) {
	c.Store.Add(vals...)
	c.write()
}

func (c *cache) Remove(vals ...string) {
	c.Store.Remove(vals...)
	c.write()
}

func (c *cache) Clear() {
	c.Store.Clear()
	c.clear()
}

func (c *cache) clear() {
	if err := os.Truncate(c.path, 0); err != nil {
		fmt.Printf("Unable to truncate file, %v\n", err)
	}
}

func (c *cache) write() (int, error) {
	w, err := open(c.path, os.O_WRONLY|os.O_TRUNC)
	if err != nil {
		return 0, err
	}

	return w.WriteString(strings.Join(c.Members(), "\n") + "\n")
}

func open(path string, flag int) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, err
	}

	return os.OpenFile(path, flag|os.O_CREATE, 0600)
}
