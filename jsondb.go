package jsondb

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"time"
)

type JsonDb interface {
	// Create a new entry and write it to the database.
	Create(entry Entry) error

	// Delete an existing entry. Will return NotFoundError if the
	// record isn't in the database.
	Delete(id string) error

	// Read an entry from the database. Returns a NotFoundError
	// if an entry with the supplied id can't be found. Will
	// also return any serialisation errors if they occur.
	Read(id string, entry Entry) error

	// Update an existing entry. Will return NotFoundError if the
	// record isn't in the database.
	Update(id string, entry Entry) error

	// Creates a Scanner for iterating all the entries
	// in a database.
	NewScanner() (Scanner, error)
}

// func for generating a new identifier. Use your uuid
// of choice
type IdGenerator func() string

// Create a new database client pointed at a directory. Will
// create the directory if it doesn't exist.
func New(dir string, newid IdGenerator) JsonDb {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			panic("Unable to create directory " + dir)
		}
	}

	return &jsondatabase{dir, newid}
}

// A database entry representation needs to implement
// this interface.
type Entry interface {
	// Assigns a generated id for a new entry, or assigns
	// the existing id when an entry is read
	AssignId(id string)

	// Sets the created time for a new entry
	Created(at time.Time)

	// Sets the modified time when an entry is updated
	Modified(at time.Time)
}

func (db *jsondatabase) Read(id string, entry Entry) error {
	outputPath := db.path(id)

	if _, err := os.Stat(outputPath); err != nil {
		if os.IsNotExist(err) {
			return &NotFoundError{id}
		} else {
			return err
		}
	}

	data, err := ioutil.ReadFile(outputPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, entry); err != nil {
		return err
	}

	entry.AssignId(id)

	return nil
}

func (db *jsondatabase) Create(entry Entry) error {
	id := db.newid()
	entry.AssignId(id)
	entry.Created(time.Now())

	b, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(db.path(id), b, 0644); err != nil {
		return err
	}

	return nil
}

func (db *jsondatabase) Delete(id string) error {
	outputPath := db.path(id)

	if _, err := os.Stat(outputPath); err != nil {
		if os.IsNotExist(err) {
			return &NotFoundError{id}
		} else {
			return err
		}
	}

	return os.Remove(outputPath)
}

func (db *jsondatabase) Update(id string, entry Entry) error {
	outputPath := db.path(id)

	if _, err := os.Stat(outputPath); err != nil {
		if os.IsNotExist(err) {
			return &NotFoundError{id}
		} else {
			return err
		}
	}

	entry.Modified(time.Now())

	b, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path.Join(db.dir, id), b, 0644)
}

func (db *jsondatabase) NewScanner() (Scanner, error) {
	files, err := ioutil.ReadDir(db.dir)
	if err != nil {
		return nil, err
	}

	return &scanner{db, files, -1, len(files)}, nil
}

type jsondatabase struct {
	dir   string
	newid IdGenerator
}

func (db *jsondatabase) path(id string) string {
	return path.Join(db.dir, id)
}
