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

var DirectMessages = newDirectMessagesTable("public", "direct_messages", "")

type directMessagesTable struct {
	postgres.Table

	// Columns
	ID             postgres.ColumnString
	Content        postgres.ColumnString
	FileURL        postgres.ColumnString
	MemberID       postgres.ColumnString
	ConversationID postgres.ColumnString
	Deleted        postgres.ColumnBool
	CreatedAt      postgres.ColumnTimestampz
	UpdatedAt      postgres.ColumnTimestampz

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type DirectMessagesTable struct {
	directMessagesTable

	EXCLUDED directMessagesTable
}

// AS creates new DirectMessagesTable with assigned alias
func (a DirectMessagesTable) AS(alias string) *DirectMessagesTable {
	return newDirectMessagesTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new DirectMessagesTable with assigned schema name
func (a DirectMessagesTable) FromSchema(schemaName string) *DirectMessagesTable {
	return newDirectMessagesTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new DirectMessagesTable with assigned table prefix
func (a DirectMessagesTable) WithPrefix(prefix string) *DirectMessagesTable {
	return newDirectMessagesTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new DirectMessagesTable with assigned table suffix
func (a DirectMessagesTable) WithSuffix(suffix string) *DirectMessagesTable {
	return newDirectMessagesTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newDirectMessagesTable(schemaName, tableName, alias string) *DirectMessagesTable {
	return &DirectMessagesTable{
		directMessagesTable: newDirectMessagesTableImpl(schemaName, tableName, alias),
		EXCLUDED:            newDirectMessagesTableImpl("", "excluded", ""),
	}
}

func newDirectMessagesTableImpl(schemaName, tableName, alias string) directMessagesTable {
	var (
		IDColumn             = postgres.StringColumn("id")
		ContentColumn        = postgres.StringColumn("content")
		FileURLColumn        = postgres.StringColumn("file_url")
		MemberIDColumn       = postgres.StringColumn("member_id")
		ConversationIDColumn = postgres.StringColumn("conversation_id")
		DeletedColumn        = postgres.BoolColumn("deleted")
		CreatedAtColumn      = postgres.TimestampzColumn("created_at")
		UpdatedAtColumn      = postgres.TimestampzColumn("updated_at")
		allColumns           = postgres.ColumnList{IDColumn, ContentColumn, FileURLColumn, MemberIDColumn, ConversationIDColumn, DeletedColumn, CreatedAtColumn, UpdatedAtColumn}
		mutableColumns       = postgres.ColumnList{ContentColumn, FileURLColumn, MemberIDColumn, ConversationIDColumn, DeletedColumn, CreatedAtColumn, UpdatedAtColumn}
	)

	return directMessagesTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:             IDColumn,
		Content:        ContentColumn,
		FileURL:        FileURLColumn,
		MemberID:       MemberIDColumn,
		ConversationID: ConversationIDColumn,
		Deleted:        DeletedColumn,
		CreatedAt:      CreatedAtColumn,
		UpdatedAt:      UpdatedAtColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
