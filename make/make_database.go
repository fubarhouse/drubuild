package make

type makeDB struct {
	dbHost string
	dbUser string
	dbPass string
	dbPort int
}

func NewmakeDB(dbHost, dbUser, dbPass string, dbPort int) *makeDB {
	newDB := &makeDB{}
	newDB.setHost(dbHost)
	newDB.setUser(dbUser)
	newDB.setPass(dbPass)
	newDB.setPort(dbPort)
	return newDB
}

func (db *makeDB) setHost(dbHost string) {
	db.dbHost = dbHost
}

func (db *makeDB) getHost() string {
	return db.dbHost
}

func (db *makeDB) setUser(dbUser string) {
	db.dbUser = dbUser
}

func (db *makeDB) getUser() string {
	return db.dbUser
}

func (db *makeDB) setPass(dbPass string) {
	db.dbPass = dbPass
}

func (db *makeDB) getPass() string {
	return db.dbPass
}

func (db *makeDB) setPort(dbPort int) {
	db.dbPort = dbPort
}

func (db *makeDB) getPort() int {
	return db.dbPort
}

func (db *makeDB) getInfo() *makeDB {
	return db
}
