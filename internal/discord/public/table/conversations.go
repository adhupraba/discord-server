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

var Conversations = newConversationsTable("public", "conversations", "")

type conversationsTable struct {
	postgres.Table

	// Columns
	ID          postgres.ColumnString
	MemberOneID postgres.ColumnString
	MemberTwoID postgres.ColumnString
	CreatedAt   postgres.ColumnTimestampz
	UpdatedAt   postgres.ColumnTimestampz

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type ConversationsTable struct {
	conversationsTable

	EXCLUDED conversationsTable
}

// AS creates new ConversationsTable with assigned alias
func (a ConversationsTable) AS(alias string) *ConversationsTable {
	return newConversationsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new ConversationsTable with assigned schema name
func (a ConversationsTable) FromSchema(schemaName string) *ConversationsTable {
	return newConversationsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new ConversationsTable with assigned table prefix
func (a ConversationsTable) WithPrefix(prefix string) *ConversationsTable {
	return newConversationsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new ConversationsTable with assigned table suffix
func (a ConversationsTable) WithSuffix(suffix string) *ConversationsTable {
	return newConversationsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newConversationsTable(schemaName, tableName, alias string) *ConversationsTable {
	return &ConversationsTable{
		conversationsTable: newConversationsTableImpl(schemaName, tableName, alias),
		EXCLUDED:           newConversationsTableImpl("", "excluded", ""),
	}
}

func newConversationsTableImpl(schemaName, tableName, alias string) conversationsTable {
	var (
		IDColumn          = postgres.StringColumn("id")
		MemberOneIDColumn = postgres.StringColumn("member_one_id")
		MemberTwoIDColumn = postgres.StringColumn("member_two_id")
		CreatedAtColumn   = postgres.TimestampzColumn("created_at")
		UpdatedAtColumn   = postgres.TimestampzColumn("updated_at")
		allColumns        = postgres.ColumnList{IDColumn, MemberOneIDColumn, MemberTwoIDColumn, CreatedAtColumn, UpdatedAtColumn}
		mutableColumns    = postgres.ColumnList{MemberOneIDColumn, MemberTwoIDColumn, CreatedAtColumn, UpdatedAtColumn}
	)

	return conversationsTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:          IDColumn,
		MemberOneID: MemberOneIDColumn,
		MemberTwoID: MemberTwoIDColumn,
		CreatedAt:   CreatedAtColumn,
		UpdatedAt:   UpdatedAtColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
