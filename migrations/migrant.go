package main

import (
	"database/sql"
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := flag.String("dsn", "", "MySQL DSN (e.g., user:password@tcp(localhost:3306)/dbname)")
	migrationsDir := flag.String("dir", ".", "Directory containing migration files")
	flag.Parse()

	if *dsn == "" {
		log.Fatal("DSN is required")
	}

	db, err := sql.Open("mysql", *dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	files, err := ioutil.ReadDir(*migrationsDir)
	if err != nil {
		log.Fatal("Failed to read migrations directory:", err)
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		content, err := ioutil.ReadFile(filepath.Join(*migrationsDir, file.Name()))
		if err != nil {
			log.Fatal("Failed to read migration file:", err)
		}

		log.Printf("Applying migration: %s", file.Name())
		_, err = db.Exec(string(content))
		if err != nil {
			log.Fatalf("Failed to apply migration %s: %v", file.Name(), err)
		}
	}

	log.Println("All migrations applied successfully!")
}