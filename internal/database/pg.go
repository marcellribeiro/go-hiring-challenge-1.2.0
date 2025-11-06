package database

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database interface {
	Find(dest interface{}, conds ...interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) *gorm.DB
	Create(value interface{}) *gorm.DB
	Save(value interface{}) *gorm.DB
	Delete(value interface{}, conds ...interface{}) *gorm.DB
	Where(query interface{}, args ...interface{}) *gorm.DB
	Limit(limit int) *gorm.DB
	Offset(offset int) *gorm.DB
	Preload(query string, args ...interface{}) *gorm.DB
	Count(count *int64) *gorm.DB
	Exec(query string, args ...interface{}) *gorm.DB
	Model(value interface{}) *gorm.DB
}

func New(user, password, dbname, port string) (Database, func() error) {
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", user, password, port, dbname)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %s", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection: %s", err)
	}

	return &GormDB{DB: db}, sqlDB.Close
}

type GormDB struct {
	*gorm.DB
}

func (g *GormDB) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	return g.DB.Find(dest, conds...)
}

func (g *GormDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	return g.DB.First(dest, conds...)
}

func (g *GormDB) Create(value interface{}) *gorm.DB {
	return g.DB.Create(value)
}

func (g *GormDB) Save(value interface{}) *gorm.DB {
	return g.DB.Save(value)
}

func (g *GormDB) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	return g.DB.Delete(value, conds...)
}

func (g *GormDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	return g.DB.Where(query, args...)
}

func (g *GormDB) Limit(limit int) *gorm.DB {
	return g.DB.Limit(limit)
}

func (g *GormDB) Offset(offset int) *gorm.DB {
	return g.DB.Offset(offset)
}

func (g *GormDB) Preload(query string, args ...interface{}) *gorm.DB {
	return g.DB.Preload(query, args...)
}

func (g *GormDB) Count(count *int64) *gorm.DB {
	return g.DB.Count(count)
}

func (g *GormDB) Exec(query string, args ...interface{}) *gorm.DB {
	return g.DB.Exec(query, args...)
}
func (g *GormDB) Model(value interface{}) *gorm.DB {
	return g.DB.Model(value)
}
