# SQLite Store Package

This package provides primitives for working with SQLite databases stored in a directory, with game-specific schemas applied at runtime.

## Functions

### `New(base string) *Store`
Creates a `Store` instance bound to a base directory for database files.  
The directory will be created if it does not exist.

**Parameters:**
- `base`: string — Directory path where `.db` files will be stored.

**Returns:**
- `*Store`: A pointer to a `Store` instance.

---

### `OpenFor(gameID string, schemaPath string) (*sql.DB, error)`
Opens (or creates) a SQLite database file for a given `gameID` and ensures the schema from a `.sql` file is applied.

**Parameters:**
- `gameID`: string — Unique identifier for the game. The DB filename will be `<gameID>.db`.
- `schemaPath`: string — Filesystem path to the `.sql` schema file to execute.

**Returns:**
- `*sql.DB`: An open database handle.
- `error`: Non-nil if the database could not be opened or the schema could not be applied.

---

## Example Usage

```go
package main

import (
    "log"
    "github.com/yourname/yourmodule/internal/store/sqlite"
)

func main() {
    st := sqlite.New("Storage/games")
    db, err := st.OpenFor("demo", "internal/tictactoe/store/schema.sql")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Database is ready with schema applied.
}
