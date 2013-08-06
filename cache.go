package jsondb

import "reflect"

// Create a memory caching wrapper around a JsonDb.
// Unlike an uncached JsonDb, a cache won't be able
// to deal with other processes changing the underlying
// files. Don't use if you've got that scenario.
func Cache(db JsonDb) JsonDb {
	return &cachingdb{
		inner:         db,
		cachedentries: make(map[string]Entry),
		cachedids:     make(map[string]bool),
	}
}

type cachingdb struct {
	inner         JsonDb
	cachedentries map[string]Entry
	cachedids     map[string]bool
	initd         bool
}

func (db *cachingdb) Create(entry Entry) (string, error) {
	id, err := db.inner.Create(entry)

	if err != nil {
		return id, err
	}

	db.cachedentries[id] = entry
	db.cachedids[id] = true
	return id, nil
}

func (db *cachingdb) Delete(id string) error {
	if err := db.inner.Delete(id); err != nil {
		return err
	}

	delete(db.cachedentries, id)
	delete(db.cachedids, id)
	return nil
}

func (db *cachingdb) Read(id string, entry Entry) error {
	cached := db.cachedentries[id]

	if cached != nil {
		// A bit of nutty unsafe pointer manipulation here.
		// Due to the api being designed around the way the
		// json unmarshaler takes a pointer, it prevents us
		// from just passing out our cached copy. We have to
		// poke the cached value into the pointer using reflect.
		entryptr := reflect.ValueOf(entry)
		entrystruct := entryptr.Elem()

		cachedptr := reflect.ValueOf(cached)
		cachedstruct := cachedptr.Elem()

		entrystruct.Set(cachedstruct)
		return nil
	}

	if err := db.inner.Read(id, entry); err != nil {
		return err
	}

	db.cachedentries[id] = entry
	db.cachedids[id] = true
	return nil
}

func (db *cachingdb) Update(id string, entry Entry) error {
	if err := db.inner.Update(id, entry); err != nil {
		return err
	}

	db.cachedentries[id] = entry
	db.cachedids[id] = true
	return nil
}

func (db *cachingdb) ids() ([]string, error) {
	if !db.initd {
		ids, err := db.inner.ids()
		if err != nil {
			return []string{}, err
		}
		for _, id := range ids {
			db.cachedids[id] = true
		}
		db.initd = true
	}

	ids := make([]string, 0, len(db.cachedentries))
	for id, _ := range db.cachedids {
		ids = append(ids, id)
	}

	return ids, nil
}
