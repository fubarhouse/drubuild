package make

// Struct which declares information for database settings.
type makeDB struct {
	dbHost string
	dbUser string
	dbPass string
	dbPort int
}

// Instantiate a new makeDB struct
func NewmakeDB(dbHost, dbUser, dbPass string, dbPort int) *makeDB {
	newDB := &makeDB{}
	newDB.setHost(dbHost)
	newDB.setUser(dbUser)
	newDB.setPass(dbPass)
	newDB.setPort(dbPort)
	return newDB
}

// Set the dbHost field
func (db *makeDB) setHost(dbHost string) {
	db.dbHost = dbHost
}

// Get the dbHost field
func (db *makeDB) getHost() string {
	return db.dbHost
}

// Set the dbUser field
func (db *makeDB) setUser(dbUser string) {
	db.dbUser = dbUser
}

// Get the dbUser field
func (db *makeDB) getUser() string {
	return db.dbUser
}

// Set the dbPass field
func (db *makeDB) setPass(dbPass string) {
	db.dbPass = dbPass
}

// Get the dbPass field
func (db *makeDB) getPass() string {
	return db.dbPass
}

// Set the dbPort field
func (db *makeDB) setPort(dbPort int) {
	db.dbPort = dbPort
}

// Get the dbPort field
func (db *makeDB) getPort() int {
	return db.dbPort
}

// Returns the database struct
func (db *makeDB) getInfo() *makeDB {
	return db
}
