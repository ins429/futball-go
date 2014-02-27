package main

import (
	"database/sql"
	"encoding/base32"
	"github.com/coopernurse/gorp"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"net/http"
	"strings"
	"time"
)

type PGStore struct {
	Codecs  []securecookie.Codec
	Options *sessions.Options
	path    string
	DbMap   *gorp.DbMap
}

// Session type
type Session struct {
	Id         int64     `db: "id"`
	Key        string    `db: "key"`
	Data       string    `db: "data"`
	CreatedOn  time.Time `db: "created_on"`
	ModifiedOn time.Time `db: "modified_on"`
	ExpiresOn  time.Time `db: "expires_on"`
}

func NewPGStore(dbUrl string, keyPairs ...[]byte) *PGStore {
	db, err := sql.Open("postgres", dbUrl)
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	dbStore := &PGStore{
		Codecs: securecookie.CodecsFromPairs(keyPairs...),
		Options: &sessions.Options{
			Path:   "/",
			MaxAge: 86400 * 30,
		},
		DbMap: dbmap,
	}

	if err != nil {
		// Ignore and return nil
		//fmt.Printf("SQL connection error: %s\n", err)
		return nil
	}

	// Create table if it doesn't exist
	dbmap.AddTableWithName(Session{}, "http_sessions").SetKeys(true, "Id")
	err = dbmap.CreateTablesIfNotExists()

	if err != nil {
		// Ignore and return nil
		//fmt.Printf("Failed to create table http_sessions. Error: %s\n", err)
		return nil
	}

	return dbStore
}

func (db *PGStore) Close() {
	db.DbMap.Db.Close()
}

// Fetches a session for a given name after it has been added to the registry.
func (db *PGStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(db, name)
}

// New returns a new session for the given name w/o adding it to the registry.
func (db *PGStore) New(r *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(db, name)
	if session == nil {
		return session, nil
	}
	session.Options = &(*db.Options)
	session.IsNew = true

	var err error
	if c, errCookie := r.Cookie(name); errCookie == nil {
		err = securecookie.DecodeMulti(name, c.Value, &session.ID, db.Codecs...)
		if err == nil {
			err = db.load(session)
			if err == nil {
				session.IsNew = false
			}
		}
	}

	return session, err
}

func (db *PGStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	// Set delete if max-age is < 0
	if session.Options.MaxAge < 0 {
		if err := db.destroy(session); err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), "", session.Options))
		return nil
	}

	if session.ID == "" {
		// Generate a random session ID key suitable for storage in the DB
		session.ID = string(securecookie.GenerateRandomKey(32))
		session.ID = strings.TrimRight(
			base32.StdEncoding.EncodeToString(
				securecookie.GenerateRandomKey(32)), "=")
	}

	if err := db.save(session); err != nil {
		return err
	}

	// Keep the session ID key in a cookie so it can be looked up in DB later.
	encoded, err := securecookie.EncodeMulti(session.Name(), session.ID, db.Codecs...)
	if err != nil {
		return err
	}

	http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))
	return nil
}

//load fetches a session by ID from the database and decodes its content into session.Values
func (db *PGStore) load(session *sessions.Session) error {
	var s Session
	err := db.DbMap.SelectOne(&s, "SELECT * FROM http_sessions WHERE key = $1", session.ID)

	if err := securecookie.DecodeMulti(session.Name(), string(s.Data),
		&session.Values, db.Codecs...); err != nil {
		return err
	}

	return err
}

// save writes encoded session.Values to a database record.
// writes to http_sessions table by default.
func (db *PGStore) save(session *sessions.Session) error {
	encoded, err := securecookie.EncodeMulti(session.Name(), session.Values,
		db.Codecs...)

	if err != nil {
		return err
	}

	var createdOn time.Time
	var expiresOn time.Time

	crOn := session.Values["created_on"]
	exOn := session.Values["expires_on"]

	if crOn == nil {
		createdOn = time.Now()
	} else {
		createdOn = crOn.(time.Time)
	}

	if exOn == nil {
		expiresOn = time.Now().Add(time.Second * time.Duration(session.Options.MaxAge))
	} else {
		expiresOn = exOn.(time.Time)
		if expiresOn.Sub(time.Now().Add(time.Second*time.Duration(session.Options.MaxAge))) < 0 {
			expiresOn = time.Now().Add(time.Second * time.Duration(session.Options.MaxAge))
		}
	}

	s := Session{
		Key:        session.ID,
		Data:       encoded,
		CreatedOn:  createdOn,
		ExpiresOn:  expiresOn,
		ModifiedOn: time.Now(),
	}

	if session.IsNew {
		err = db.DbMap.Insert(&s)
	} else {
		_, err = db.DbMap.Update(&s)
	}

	return err
}

// Delete session
func (db *PGStore) destroy(session *sessions.Session) error {
	_, err := db.DbMap.Db.Exec("DELETE FROM http_sessions WHERE key = $1", session.ID)
	return err
}
