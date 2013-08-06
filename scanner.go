package jsondb

// Database record scanner. Used for iterating over
// all the records in a database.
type Scanner interface {
	// Scan to see if there's another record after the
	// current one
	Scan() bool

	// Read the current record, see Db.Read for
	// error behaviour.
	Read(entry Entry) error

	// The length of the underlying dataset
	Length() int
}

func NewScanner(db JsonDb) (Scanner, error) {
	ids, err := db.ids()
	if err != nil {
		return nil, err
	}

	return &scanner{db, ids, 0, len(ids)}, nil
}

type scanner struct {
	db     JsonDb
	ids    []string
	pos    int
	length int
}

func (s *scanner) Scan() bool {
	if s.pos < len(s.ids)-1 {
		s.pos += 1
		return true
	}

	return false
}

func (s *scanner) Read(entry Entry) error {
	return s.db.Read(s.ids[s.pos], entry)
}

func (s *scanner) Length() int {
	return s.length
}
