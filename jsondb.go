package jsondb

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"time"
)

// func for generating a new identifier. Use your uuid
// of choice
type IdGenerator func() string

// Create a new database client pointed at a directory. Will
// create the directory if it doesn't exist.
func New(dir string, newid IdGenerator) *Db {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			panic("Unable to create directory " + dir)
		}
	}

	return &Db{dir, newid}
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

// Read an entry from the database. Returns a NotFoundError
// if an entry with the supplied id can't be found. Will
// also return any serialisation errors if they occur.
func (db *Db) Read(id string, entry Entry) error {
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

// Create a new entry and write it to the database.
func (db *Db) Create(entry Entry) error {
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

// Delete an existing entry. Will return NotFoundError if the
// record isn't in the database.
func (db *Db) Delete(id string) error {
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

// Update an existing entry. Will return NotFoundError if the
// record isn't in the database.
func (db *Db) Update(id string, entry Entry) error {
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

type Db struct {
	dir   string
	newid IdGenerator
}

func (db *Db) path(id string) string {
	return path.Join(db.dir, id)
}
