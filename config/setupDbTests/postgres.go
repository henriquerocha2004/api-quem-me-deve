package setupdbtests

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func TruncateTables(db *gorm.DB) error {
	var tables []string

	err := db.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public' AND tablename != 'schema_migrations'").Scan(&tables).Error

	if err != nil {
		return err
	}

	if len(tables) == 0 {
		return nil
	}

	sql := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", strings.Join(wrapTableNames(tables), ", "))
	return db.Exec(sql).Error
}

func wrapTableNames(names []string) []string {
	out := make([]string, len(names))
	for i, name := range names {
		out[i] = `"` + name + `"`
	}
	return out
}
