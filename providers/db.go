package providers

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strconv"
	"time"
)

type IDbClient interface {
	GetDBInstance(config AppConfig) (*gorm.DB, error)
	Close() error
}

type DBClient struct {
	db *gorm.DB
}

type PQDBInfo struct {
	SchemaName          string
	SingularTable       bool
	TimeZone            string
	DBHost              string
	DBPort              string
	DBUser              string
	DBPassword          string
	DBName              string
	DBMaxOpenConn       string
	DBMaxIdleConn       string
	DBConnectionTimeout string
}

func (dbClient *DBClient) GetDBInstance(config AppConfig) (*gorm.DB, error) {
	if dbClient.db == nil {
		port := strconv.Itoa(config.DbConfig.Port)
		timeout := strconv.Itoa(config.DbConfig.Timeout)
		maxOpenConn := strconv.Itoa(config.DbConfig.MaxOpenConn)
		maxIdleConn := strconv.Itoa(config.DbConfig.MaxIdleConn)
		pdbInfo := PQDBInfo{
			SchemaName:          "url_shortener.",
			SingularTable:       true,
			TimeZone:            "UTC",
			DBHost:              config.DbConfig.Host,
			DBPort:              port,
			DBUser:              config.DbConfig.User,
			DBPassword:          config.DbConfig.Password,
			DBName:              config.DbConfig.Name,
			DBConnectionTimeout: timeout,
			DBMaxOpenConn:       maxOpenConn,
			DBMaxIdleConn:       maxIdleConn,
		}
		return GetGormSqlClient(&pdbInfo), nil
	} else {
		return dbClient.db, nil
	}
}

func (dbClient *DBClient) Close() error {
	if dbClient.db != nil {
		sqlDB, err := dbClient.db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

func GetGormSqlClient(pqdbinfo *PQDBInfo) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		pqdbinfo.DBHost, pqdbinfo.DBPort, pqdbinfo.DBUser, pqdbinfo.DBPassword, pqdbinfo.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: false,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   pqdbinfo.SchemaName,
			SingularTable: pqdbinfo.SingularTable,
		},
		NowFunc: func() time.Time {
			tz, err := time.LoadLocation(pqdbinfo.TimeZone)
			if err != nil {
				panic(err)
			}
			return time.Now().In(tz)
		},
	})
	if err != nil {
		panic(err)
	}
	return db
}
