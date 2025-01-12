package mydump

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/jamf/go-mysqldump"
	"github.com/stenstromen/s3dbdump/mygzip"
	"github.com/stenstromen/s3dbdump/mys3"
)

var Config mysql.Config

func InitConfig() {
	Config.User = os.Getenv("DB_USER")
	Config.Passwd = os.Getenv("DB_PASSWORD")
	Config.AllowNativePasswords = true
	Config.AllowCleartextPasswords = true
	Config.Net = "tcp"
	db_port := os.Getenv("DB_PORT")
	if db_port == "" {
		db_port = "3306"
	}
	Config.Addr = fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), db_port)
}

func dumpAllDatabases(config mysql.Config) {
	log.Printf("Dumping all databases")

	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return
	}
	defer db.Close()

	rows, err := db.Query("SHOW DATABASES")
	if err != nil {
		log.Printf("Error querying databases: %v", err)
		return
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			log.Printf("Error scanning database name: %v", err)
			return
		}

		if !strings.HasPrefix(dbName, "information_schema") &&
			!strings.HasPrefix(dbName, "performance_schema") &&
			!strings.HasPrefix(dbName, "mysql") &&
			!strings.HasPrefix(dbName, "sys") {
			databases = append(databases, dbName)
		}
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating database names: %v", err)
		return
	}

	for _, database := range databases {
		dumpDatabase(database, config)
	}
}

func dumpDatabase(database string, config mysql.Config) {
	log.Printf("Dumping database %s", database)

	config.DBName = database

	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?)", database).Scan(&exists)
	if err != nil {
		log.Printf("Error checking if database exists: %v", err)
		return
	}
	if !exists {
		log.Fatalf("Database %s does not exist", database)
		return
	}

	dumpDir := os.Getenv("DB_DUMP_PATH")
	if dumpDir == "" {
		dumpDir = "./dumps"
	}
	dumpFilenameFormat := fmt.Sprintf("%s-20060102T150405", database)

	dumper, err := mysqldump.Register(db, dumpDir, dumpFilenameFormat)
	if err != nil {
		log.Printf("Error registering databse: %v", err)
		return
	}

	if err := dumper.Dump(); err != nil {
		log.Fatalf("Error dumping: %v", err)
	}

	if file, ok := dumper.Out.(*os.File); ok {
		if os.Getenv("DB_GZIP") == "0" {
			mys3.UploadToS3(file.Name())
			os.Remove(file.Name())
		} else {
			mygzip.GzipFile(file.Name())
			mys3.UploadToS3(file.Name() + ".gz")
			os.Remove(file.Name() + ".gz")
		}
	} else {
		log.Printf("It's not part of *os.File, but dump is done")
	}

	dumper.Close()
}

func HandleDbDump(config mysql.Config) {
	keepBackups := os.Getenv("DB_DUMP_FILE_KEEP_DAYS")
	if keepBackups == "" {
		keepBackups = "7"
	}

	if os.Getenv("DB_ALL_DATABASES") == "1" {
		dumpAllDatabases(config)
		mys3.KeepOnlyNBackups(keepBackups)
	} else if os.Getenv("DB_NAME") != "" {
		dumpDatabase(os.Getenv("DB_NAME"), config)
		mys3.KeepOnlyNBackups(keepBackups)
	} else {
		log.Printf("No database name provided")
		return
	}
}
