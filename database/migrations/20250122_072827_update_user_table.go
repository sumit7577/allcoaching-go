package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type UpdateUserTable_20250122_072827 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &UpdateUserTable_20250122_072827{}
	m.Created = "20250122_072827"

	migration.Register("UpdateUserTable_20250122_072827", m)
}

// Run the migrations
func (m *UpdateUserTable_20250122_072827) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update

}

// Reverse the migrations
func (m *UpdateUserTable_20250122_072827) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
