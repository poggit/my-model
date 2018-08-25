/*
 * Poggit
 *
 * Copyright (C) 2018 Poggit
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package myModel

import (
	"errors"
	"github.com/philopon/go-toposort"
	"reflect"
	"sort"
)

type Schema struct {
	Tables        map[string]*MainTable
	sortedList    []*MainTable
	graphOutdated bool
}

func (schema *Schema) getTable(typ reflect.Type) *MainTable {
	if _, isNew := schema.Tables[typ.Name()]; !isNew {
		schema.graphOutdated = true
		schema.Tables[typ.Name()] = NewMainTable(typ)
	}
	return schema.Tables[typ.Name()]
}
func (schema *Schema) mustGetTable(name string) *MainTable {
	table, exists := schema.Tables[name]
	if !exists {
		panic("table " + name + " does not exist")
	}
	return table
}
func (schema *Schema) getSortedTables() []*MainTable {
	if schema.graphOutdated {
		graph := toposort.NewGraph(len(schema.Tables))
		tables := make([]string, 0, len(schema.Tables))
		for _, table := range schema.Tables {
			tables = append(tables, table.Name)
		}
		sort.Strings(tables)
		for _, table := range tables {
			graph.AddNode(table)
		}
		for _, dependent := range schema.Tables {
			for _, dependency := range schema.Tables {
				if dependent != dependency && dependent.Depends(dependency) {
					graph.AddEdge(dependency.Name, dependent.Name)
				}
			}
		}

		schema.sortedList = []*MainTable{}

		names, ok := graph.Toposort()
		if !ok {
			panic(errors.New("cyclic dependency detected"))
		}
		for _, name := range names {
			schema.sortedList = append(schema.sortedList, schema.Tables[name])
		}
	}

	return schema.sortedList
}

type Table struct {
	Name          string
	SimpleFields  []*MysqlField
	PrimaryKeys   []string
	UniqueKeys    map[string][]string
	CompositeKeys map[string][]string
	ForeignKeys   []ForeignKey
}

func NewTable(name string) *Table {
	return &Table{
		Name:          name,
		SimpleFields:  []*MysqlField{},
		PrimaryKeys:   []string{},
		UniqueKeys:    map[string][]string{},
		CompositeKeys: map[string][]string{},
	}
}

func (table *Table) FindField(name string) *MysqlField {
	for _, field := range table.SimpleFields {
		if field.Name == name {
			return field
		}
	}

	panic("field " + name + " not found in " + table.Name)
}

type MainTable struct {
	*Table
	AuxTables []*Table

	Edges       []*Edge
	Type        reflect.Type
	knownParent *MainTable // set from the parent type, to be validated if there is an EdgeTypeMultiOneParent
	yielded     bool
}

func NewMainTable(typ reflect.Type) *MainTable {
	return &MainTable{
		Table:     NewTable(typ.Name()),
		Edges:     []*Edge{},
		AuxTables: []*Table{},
		Type:      typ,
	}
}

func (table *MainTable) Depends(dependency *MainTable) bool {
	if table.knownParent == dependency {
		return true
	}
	for _, edge := range table.Edges {
		if edge.Type != EdgeTypeOneOne && edge.Type != EdgeTypeOneMulti && edge.PeerTable == dependency.Name {
			return true
		}
	}
	return false
}

func (table *MainTable) FindEdgeByName(name string) *Edge {
	for _, edge := range table.Edges {
		if edge.Name == name {
			return edge
		}
	}
	return nil
}

func (table *MainTable) FindEdgeByPeerTable(peerTable string) *Edge {
	for _, edge := range table.Edges {
		if edge.PeerTable == peerTable {
			return edge
		}
	}
	return nil
}

type MysqlField struct {
	Name          string
	Type          string
	Nullable      bool
	AutoIncrement bool
}

type ForeignKey struct {
	SourceColumns []string
	RefTable      string
	RefColumns    []string
	OnUpdate      ReferenceOption
	OnDelete      ReferenceOption
}

type ReferenceOption string

const (
	ReferenceOptionRestrict ReferenceOption = "RESTRICT"
	ReferenceOptionCascade  ReferenceOption = "CASCADE"
	ReferenceOptionSetNull  ReferenceOption = "SET NULL"
)

func MakeForeignKey(refTable string) ForeignKey {
	return ForeignKey{
		SourceColumns: []string{},
		RefTable:      refTable,
		RefColumns:    []string{},
		OnUpdate:      ReferenceOptionRestrict,
	}
}
