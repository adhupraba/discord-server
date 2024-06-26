//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var Servers = newServersTable("public", "servers", "")

type serversTable struct {
	postgres.Table

	// Columns
	ID         postgres.ColumnString
	Name       postgres.ColumnString
	ImageURL   postgres.ColumnString
	InviteCode postgres.ColumnString
	ProfileID  postgres.ColumnString
	CreatedAt  postgres.ColumnTimestampz
	UpdatedAt  postgres.ColumnTimestampz

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type ServersTable struct {
	serversTable

	EXCLUDED serversTable
}

// AS creates new ServersTable with assigned alias
func (a ServersTable) AS(alias string) *ServersTable {
	return newServersTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new ServersTable with assigned schema name
func (a ServersTable) FromSchema(schemaName string) *ServersTable {
	return newServersTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new ServersTable with assigned table prefix
func (a ServersTable) WithPrefix(prefix string) *ServersTable {
	return newServersTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new ServersTable with assigned table suffix
func (a ServersTable) WithSuffix(suffix string) *ServersTable {
	return newServersTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newServersTable(schemaName, tableName, alias string) *ServersTable {
	return &ServersTable{
		serversTable: newServersTableImpl(schemaName, tableName, alias),
		EXCLUDED:     newServersTableImpl("", "excluded", ""),
	}
}

func newServersTableImpl(schemaName, tableName, alias string) serversTable {
	var (
		IDColumn         = postgres.StringColumn("id")
		NameColumn       = postgres.StringColumn("name")
		ImageURLColumn   = postgres.StringColumn("image_url")
		InviteCodeColumn = postgres.StringColumn("invite_code")
		ProfileIDColumn  = postgres.StringColumn("profile_id")
		CreatedAtColumn  = postgres.TimestampzColumn("created_at")
		UpdatedAtColumn  = postgres.TimestampzColumn("updated_at")
		allColumns       = postgres.ColumnList{IDColumn, NameColumn, ImageURLColumn, InviteCodeColumn, ProfileIDColumn, CreatedAtColumn, UpdatedAtColumn}
		mutableColumns   = postgres.ColumnList{NameColumn, ImageURLColumn, InviteCodeColumn, ProfileIDColumn, CreatedAtColumn, UpdatedAtColumn}
	)

	return serversTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:         IDColumn,
		Name:       NameColumn,
		ImageURL:   ImageURLColumn,
		InviteCode: InviteCodeColumn,
		ProfileID:  ProfileIDColumn,
		CreatedAt:  CreatedAtColumn,
		UpdatedAt:  UpdatedAtColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
