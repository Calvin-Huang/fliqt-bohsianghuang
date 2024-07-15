package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-gormigrate/gormigrate/v2"

	"fliqt/config"
	"fliqt/internal/model/migration"
	"fliqt/internal/util"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version int
	// rollback to the previous version.
	rollback bool
	// rollback-to is the version to rollback to.
	rollbackTo string

	id, _ = os.Hostname()
)

func init() {
	// Add flag for migration version
	flag.IntVar(&Version, "version", 0, "migration version, eg: -version 0001")
	// Add flag for rolling back migrations to previous version
	flag.BoolVar(&rollback, "rollback", false, "rollback to specified version, eg: -rollback")
	// Add flag for rolling back migrations to specified version
	flag.StringVar(&rollbackTo, "rollback-to", "", "rollback to specified version, eg: -rollback-to 0001")
}

// Check if the database exists.
func checkAndCreateDatabase(cfg *config.Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?parseTime=true&loc=%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, url.QueryEscape(cfg.DBTimezone))
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`;", strings.ReplaceAll(cfg.DBName, "-", "_"))
	_, err = db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	flag.Parse()
	cfg := config.NewConfig()

	err := checkAndCreateDatabase(cfg)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	db, err := util.NewGormDB(cfg)
	if err != nil {
		panic(err)
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	allMigration := migration.AllMigrations()
	// filter the migrations setting smaller than or equal to bc.Data.Database.Version
	var migrations []*gormigrate.Migration
	for idx, m := range allMigration {
		migreateVersion, err := migration.IDToVersion(m.ID)
		if err != nil {
			panic(err)
		}
		if idx+1 != int(migreateVersion) {
			panic(fmt.Errorf("migration version is not continuous"))
		}
		if Version == 0 || migreateVersion <= Version {
			migrations = append(migrations, m)
		}
	}

	migrateOpt := gormigrate.DefaultOptions
	migrateOpt.UseTransaction = true
	m := gormigrate.New(db, migrateOpt, migrations)

	// The flag rollback-to and rollback cannot be used at the same time
	if rollbackTo != "" && rollback {
		log.Fatalf("The flag rollback-to and rollback cannot be used at the same time")
	}

	if rollbackTo != "" {
		if err := m.RollbackTo(rollbackTo); err != nil {
			log.Fatalf("Could not rollback: %v", err)
		}
		log.Printf("Rollback to version %s successful", rollbackTo)
	} else if rollback {
		if err := m.RollbackLast(); err != nil {
			log.Fatalf("Could not rollback: %v", err)
		}
	} else {
		if err = m.Migrate(); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		log.Println("Migration did run successfully")
	}
}
