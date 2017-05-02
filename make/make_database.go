package make

// makeDB is a struct which declares information for database settings.
type makeDB struct {
	dbHost string
	dbUser string
	dbPass string
	dbPort int
}

// NewmakeDB will instantiate a new makeDB struct
func NewmakeDB(dbHost, dbUser, dbPass string, dbPort int) *makeDB {
	newDB := &makeDB{}
	newDB.setHost(dbHost)
	newDB.setUser(dbUser)
	newDB.setPass(dbPass)
	newDB.setPort(dbPort)
	return newDB
}

// setHost will set the dbHost field
func (db *makeDB) setHost(dbHost string) {
	db.dbHost = dbHost
}

// getHost will get the dbHost field
func (db *makeDB) getHost() string {
	return db.dbHost
}

// setUser will set the dbUser field
func (db *makeDB) setUser(dbUser string) {
	db.dbUser = dbUser
}

// dbUser will set the dbUser field
func (db *makeDB) getUser() string {
	return db.dbUser
}

// setPass will set the dbPass field
func (db *makeDB) setPass(dbPass string) {
	db.dbPass = dbPass
}

// getPass will get the dbPass field
func (db *makeDB) getPass() string {
	return db.dbPass
}

// setPort will set the dbPort field
func (db *makeDB) setPort(dbPort int) {
	db.dbPort = dbPort
}

// getPort will get the dbPort field
func (db *makeDB) getPort() int {
	return db.dbPort
}

// getInfo returns the database struct
func (db *makeDB) getInfo() *makeDB {
	return db
}
