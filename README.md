# jsondb

Very dumb, stupid, crappy JSON-based file database. No gaurentees for performance, correctness, robustness or usability given.

## Install

    go get github.com/jagregory/jsondb

## Usage

    db := jsondb.New("./data", func() string {
      return generateNextUuid()
    })
    
    var entry Foo
    db.Read("abcd", &foo)
