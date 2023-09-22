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

var People = newPeopleTable("public", "people", "")

type peopleTable struct {
	postgres.Table

	// Columns
	ID        postgres.ColumnString
	Nickname  postgres.ColumnString
	Name      postgres.ColumnString
	Birthdate postgres.ColumnDate
	Stack     postgres.ColumnString
	Search    postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type PeopleTable struct {
	peopleTable

	EXCLUDED peopleTable
}

// AS creates new PeopleTable with assigned alias
func (a PeopleTable) AS(alias string) *PeopleTable {
	return newPeopleTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new PeopleTable with assigned schema name
func (a PeopleTable) FromSchema(schemaName string) *PeopleTable {
	return newPeopleTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new PeopleTable with assigned table prefix
func (a PeopleTable) WithPrefix(prefix string) *PeopleTable {
	return newPeopleTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new PeopleTable with assigned table suffix
func (a PeopleTable) WithSuffix(suffix string) *PeopleTable {
	return newPeopleTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newPeopleTable(schemaName, tableName, alias string) *PeopleTable {
	return &PeopleTable{
		peopleTable: newPeopleTableImpl(schemaName, tableName, alias),
		EXCLUDED:    newPeopleTableImpl("", "excluded", ""),
	}
}

func newPeopleTableImpl(schemaName, tableName, alias string) peopleTable {
	var (
		IDColumn        = postgres.StringColumn("id")
		NicknameColumn  = postgres.StringColumn("nickname")
		NameColumn      = postgres.StringColumn("name")
		BirthdateColumn = postgres.DateColumn("birthdate")
		StackColumn     = postgres.StringColumn("stack")
		SearchColumn    = postgres.StringColumn("search")
		allColumns      = postgres.ColumnList{IDColumn, NicknameColumn, NameColumn, BirthdateColumn, StackColumn, SearchColumn}
		mutableColumns  = postgres.ColumnList{NicknameColumn, NameColumn, BirthdateColumn, StackColumn, SearchColumn}
	)

	return peopleTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:        IDColumn,
		Nickname:  NicknameColumn,
		Name:      NameColumn,
		Birthdate: BirthdateColumn,
		Stack:     StackColumn,
		Search:    SearchColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
