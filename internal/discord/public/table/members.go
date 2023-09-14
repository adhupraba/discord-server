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

var Members = newMembersTable("public", "members", "")

type membersTable struct {
	postgres.Table

	// Columns
	ID        postgres.ColumnString
	Role      postgres.ColumnString
	ProfileID postgres.ColumnString
	ServerID  postgres.ColumnString
	CreatedAt postgres.ColumnTimestamp
	UpdatedAt postgres.ColumnTimestamp

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type MembersTable struct {
	membersTable

	EXCLUDED membersTable
}

// AS creates new MembersTable with assigned alias
func (a MembersTable) AS(alias string) *MembersTable {
	return newMembersTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new MembersTable with assigned schema name
func (a MembersTable) FromSchema(schemaName string) *MembersTable {
	return newMembersTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new MembersTable with assigned table prefix
func (a MembersTable) WithPrefix(prefix string) *MembersTable {
	return newMembersTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new MembersTable with assigned table suffix
func (a MembersTable) WithSuffix(suffix string) *MembersTable {
	return newMembersTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newMembersTable(schemaName, tableName, alias string) *MembersTable {
	return &MembersTable{
		membersTable: newMembersTableImpl(schemaName, tableName, alias),
		EXCLUDED:     newMembersTableImpl("", "excluded", ""),
	}
}

func newMembersTableImpl(schemaName, tableName, alias string) membersTable {
	var (
		IDColumn        = postgres.StringColumn("id")
		RoleColumn      = postgres.StringColumn("role")
		ProfileIDColumn = postgres.StringColumn("profile_id")
		ServerIDColumn  = postgres.StringColumn("server_id")
		CreatedAtColumn = postgres.TimestampColumn("created_at")
		UpdatedAtColumn = postgres.TimestampColumn("updated_at")
		allColumns      = postgres.ColumnList{IDColumn, RoleColumn, ProfileIDColumn, ServerIDColumn, CreatedAtColumn, UpdatedAtColumn}
		mutableColumns  = postgres.ColumnList{RoleColumn, ProfileIDColumn, ServerIDColumn, CreatedAtColumn, UpdatedAtColumn}
	)

	return membersTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:        IDColumn,
		Role:      RoleColumn,
		ProfileID: ProfileIDColumn,
		ServerID:  ServerIDColumn,
		CreatedAt: CreatedAtColumn,
		UpdatedAt: UpdatedAtColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
