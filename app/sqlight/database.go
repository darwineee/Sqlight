package sqlight

import (
	"com.sentry.dev/app/sqlight/cell"
	"com.sentry.dev/app/sqlight/cmd"
	"com.sentry.dev/app/sqlight/table"
	"com.sentry.dev/app/util"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"syscall"
)

type Database struct {
	file     *os.File
	fileInfo *os.FileInfo
	schema   *table.SqliteSchema
	data     []byte
	//lock       sync.RWMutex used for concurrent access when write to the database
}

var (
	database *Database
	once     sync.Once
)

func GetInstance(dbPath string) *Database {
	once.Do(func() {
		handleIfErr := func(file *os.File, err error) {
			if err != nil {
				_ = file.Close()
				log.Fatal(err)
			}
		}
		//file, err := os.OpenFile(dbPath, os.O_RDWR|os.O_CREATE, 0644)
		file, err := os.OpenFile(dbPath, os.O_RDWR, 0644)
		handleIfErr(file, err)
		fileInfo, err := file.Stat()
		handleIfErr(file, err)
		data, err := syscall.Mmap(
			int(file.Fd()),
			0,
			int(fileInfo.Size()),
			syscall.PROT_READ|syscall.PROT_WRITE,
			syscall.MAP_SHARED,
		)
		handleIfErr(file, err)
		schema, err := table.ParseSqliteSchema(data)
		handleIfErr(file, err)
		database = &Database{
			file:     file,
			fileInfo: &fileInfo,
			schema:   &schema,
			data:     data,
		}
	})
	return database
}

func (database *Database) Close() {
	if err := syscall.Munmap(database.data); err != nil {
		log.Printf("Failed to unmap database file: %v", err)
	}
	if err := database.file.Close(); err != nil {
		log.Printf("Failed to close database file: %v", err)
	}
}

func (database *Database) Execute(query string) {
	switch query[0] {
	case '.':
		database.executeDotCmd(query[1:])
	default:
		log.Println("Invalid or unsupported command")
	}
}

func (database *Database) executeDotCmd(query string) {
	switch query {
	case cmd.DbInfo:
		fmt.Println("database page size: ", database.schema.DbHeader.GetRealPageSize())
		fmt.Println("number of tables: ", database.schema.PageHeader.CellCount)
	case cmd.Tables:
		names := util.Map(database.schema.CellContent, func(record cell.SchemaRecord) string {
			return record.Name
		})
		fmt.Println(strings.Join(names, " "))

	default:
		log.Println("Invalid or unsupported command")
	}
}
