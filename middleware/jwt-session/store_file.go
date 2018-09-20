package jwt

import (
	"encoding/binary"
	"os"
	"path"
	"strconv"
)

// FileStore FileStore
type FileStore struct {
	Path string
}

// NewFileStore NewFileStore
func NewFileStore() *FileStore {
	return &FileStore{}
}

// Init config store dir path
func (p *FileStore) Init(config string) error {
	p.Path = config

	err := os.MkdirAll(p.Path, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// GetCounter GetCounter
func (p *FileStore) GetCounter(uid int64) (int64, error) {
	return p.readFile(uid, CAATypeCounter)
}

// SetCounter SetCounter
func (p *FileStore) SetCounter(uid int64, n int64) error {
	return p.writeFile(uid, CAATypeCounter, n)
}

// GetTimeout GetTimeout
func (p *FileStore) GetTimeout(uid int64) (int64, error) {
	return p.readFile(uid, CAATypeTimeout)
}

// SetTimeout SetTimeout
func (p *FileStore) SetTimeout(uid int64, t int64) error {
	return p.writeFile(uid, CAATypeTimeout, t)
}

func (p *FileStore) readFile(uid int64, ctype CAAType) (val int64, err error) {
	var f *os.File
	f, err = p.uidFile(uid)
	if err != nil {
		return
	}
	defer f.Close()

	if ctype == CAATypeCounter {
		f.Seek(0, os.SEEK_SET)
	} else {
		f.Seek(8, os.SEEK_SET)
	}

	b := make([]byte, 8)
	_, err = f.Read(b)
	if err != nil {
		return 0, err
	}

	val = int64(binary.LittleEndian.Uint64(b))

	return
}

func (p *FileStore) writeFile(uid int64, ctype CAAType, val int64) (err error) {
	var f *os.File
	f, err = p.uidFile(uid)
	if err != nil {
		return
	}
	defer f.Close()

	if ctype == CAATypeCounter {
		f.Seek(0, os.SEEK_SET)
	} else {
		f.Seek(8, os.SEEK_SET)
	}

	return binary.Write(f, binary.LittleEndian, val)
}

func (p *FileStore) uidFile(uid int64) (f *os.File, err error) {
	ustr := strconv.FormatInt(uid, 10)

	if len(ustr) < 4 {
		for i := 0; i < 4-len(ustr); i++ {
			ustr = "0" + ustr
		}
	}

	fpath := path.Join(p.Path, string(ustr[0]), string(ustr[1]), string(ustr[2]), ustr)
	err = os.MkdirAll(path.Dir(fpath), 0777)
	if err != nil {
		return
	}

	_, err = os.Stat(fpath)
	if err == nil {
		f, err = os.OpenFile(fpath, os.O_RDWR, os.ModePerm)
	} else if os.IsNotExist(err) {
		f, err = os.Create(fpath)
	} else {
		return
	}

	var stat os.FileInfo
	stat, err = f.Stat()
	if err != nil {
		return
	}

	if stat.Size() != 16 {
		f.Write(make([]byte, 16))
		f.Sync()
	}

	f.Seek(0, os.SEEK_SET)

	return
}

func init() {
	RegisterStore("file", NewFileStore())
}
