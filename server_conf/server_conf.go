package serverconf

import "github.com/AtinAgnihotri/glog-aggregator/internal/database"

type ServerConf struct {
	PORT string
	DB   *database.Queries
}
