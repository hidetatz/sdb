package sdb

import (
	"fmt"
)

type SDB struct {
	parser      Parser
	planner     Planner
	catalog     Catalog
	executor    Executor
	engine      Engine
	diskManager DiskManager
}

type Parameter struct {
	Query string // Raw SQL
	// TODO: define transaction context etc.
}

type Result struct {
	Code  string     // "OK" or "NG" for now
	RS    *ResultSet // filled when "OK"
	Error *Error     // filled when "NG"
}

type Error struct {
	Message string
}

type ResultSet struct {
	Message string
	Columns []string   // empty when insert, update, delete
	Values  [][]string // empty when insert, update, delete
	Count   int        // empty when insert
}

func New(parser Parser, planner Planner, catalog Catalog, executor Executor, engine Engine, diskManager DiskManager) *SDB {
	return &SDB{
		parser:      parser,
		planner:     planner,
		catalog:     catalog,
		executor:    executor,
		engine:      engine,
		diskManager: diskManager,
	}
}

func (sdb *SDB) ExecuteQuery(param *Parameter) *Result {
	stmt, err := sdb.parser.Parse(param.Query)
	if err != nil {
		return &Result{
			Code: "NG",
			// FUTURE WORK: should define and return &ParserError to return
			// which part of the query looks wrong
			Error: &Error{Message: fmt.Sprintf("failed to parse query: %s", err)},
		}
	}

	plan, err := sdb.planner.Plan(stmt)
	if err != nil {
		return &Result{
			Code: "NG",
			// FUTURE WORK: should define and return &ParserError to return
			// which part of the query looks wrong
			Error: &Error{Message: fmt.Sprintf("failed to plan query: %s", err)},
		}
	}

	result, err := sdb.executor.Execute(plan)
	if err != nil {
		return &Result{
			Code:  "NG",
			Error: &Error{Message: fmt.Sprintf("failed to execute query: %s", err)},
		}
	}

	successResp := &Result{
		Code: "OK",
		RS: &ResultSet{
			Message: result.RS.Message,
		},
	}

	// switch stmt.Typ {
	// case parser.SELECT_STMT:
	// 	successResp.RS.Columns = result.Columns
	// 	successResp.RS.Values = result.Values
	// 	successResp.RS.Count = result.Count
	// case parser.UPDATE_STMT, parser.DELETE_STMT:
	// 	successResp.RS.Count = result.Count
	// }

	return successResp
}

func (sdb *SDB) Shutdown() error {
	// TODO: should retry on error?
	if err := sdb.catalog.Persist(); err != nil {
		return err
	}
	return sdb.engine.Shutdown()
}
