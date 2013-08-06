package jsondb

import (
	"io/ioutil"
	"os"
)

// Creates a Scanner for iterating all the entries
// in a database.
func (db *Db) NewScanner() (*Scanner, error) {
	files, err := ioutil.ReadDir(db.dir)
	if err != nil {
		return nil, err
	}

	return &Scanner{db, files, -1, len(files)}, nil
}

// Database record scanner. Used for iterating over
// all the records in a database.
type Scanner struct {
	db     *Db
	files  []os.FileInfo
	pos    int
	Length int
}

// Scan to see if there's another record after the
// current one
func (r *Scanner) Scan() bool {
	if r.pos < len(r.files)-1 {
		r.pos += 1
		return true
	}

	return false
}

// Read the current record, see Db.Read for
// error behaviour.
func (r *Scanner) Read(entry Entry) error {
	return r.db.Read(r.files[r.pos].Name(), entry)
}
