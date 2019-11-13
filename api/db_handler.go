package api

import (
	"fmt"

	"gitlab.odds.team/plus1/backend-go/config"
	"gopkg.in/mgo.v2"
)

// MongoDB holds metadata about session database and collections name.
type (
	MongoDB struct {
		Conn   *mgo.Session
		UCol   *mgo.Collection
		SCol   *mgo.Collection
		ICol   *mgo.Collection
		TCol   *mgo.Collection
		SumCol *mgo.Collection
		RCol   *mgo.Collection
	}
)

// NewMongoDB creates a new macOddsTeamDB backed by a given Mongo server.
func NewMongoDB() (*MongoDB, error) {
	s := config.Spec()
	conn, err := mgo.Dial(s.DBHost)

	if err != nil {
		return nil, fmt.Errorf("mongo: could not dial: %v", err)
	}

	return &MongoDB{
		Conn:   conn,
		UCol:   conn.DB(s.DBName).C(s.DBUsersCol),
		SCol:   conn.DB(s.DBName).C(s.DBStatusCol),
		ICol:   conn.DB(s.DBName).C(s.DBItemCol),
		TCol:   conn.DB(s.DBName).C(s.DBTypeCol),
		RCol:   conn.DB(s.DBName).C(s.DBReserveCol),
		SumCol: conn.DB(s.DBName).C(s.DBSummaryCol),
	}, nil
}

// Close closes the database.
func (db *MongoDB) Close() {
	db.Conn.Close()
}
