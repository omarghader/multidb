package multidb

var (
	EXCEPTION_CONNECTION_FAILED = "Cannot connect to the database"
	EXCEPTION_NO_DATAFOUND      = "No data found"
	EXCEPTION_QUERY             = "Cannot execute query"

	EXCEPTION_DB_ALREADY_EXISTS = "Database already exists"
	EXCEPTION_DB_NOTFOUND       = "Database not found"
	EXCEPTION_DB_CREATE         = "Cannot create db"
	EXCEPTION_DB_DROP           = "Cannot drop db"
	EXCEPTION_DB_INSERT_ERROR   = "Cannot insert in db"
	EXCEPTION_DB_DELETE_ERROR   = "Cannot delete from db"
	EXCEPTION_DB_UPDATE_ERROR   = "Cannot update in db"
	EXCEPTION_DB_FIND_ERROR     = "Cannot search in db"

	EXCEPTION_GRAPH_ALREADY_EXISTS = "Graph already exists"
	EXCEPTION_GRAPH_NOTFOUND       = "Graph not found"
	EXCEPTION_GRAPH_CREATE         = "Cannot create graph"
	EXCEPTION_GRAPH_DROP           = "Cannot drop graph"
	EXCEPTION_GRAPH_INSERT_ERROR   = "Cannot insert in graph"
	EXCEPTION_GRAPH_DELETE_ERROR   = "Cannot delete from graph"
	EXCEPTION_GRAPH_UPDATE_ERROR   = "Cannot update in graph"
	EXCEPTION_GRAPH_FIND_ERROR     = "Cannot search in graph"

	EXCEPTION_TABLE_ALREADY_EXISTS = "Table already exists"
	EXCEPTION_TABLE_NOTFOUND       = "Table not found"
	EXCEPTION_TABLE_CREATE         = "Cannot create table"
	EXCEPTION_TABLE_DROP           = "Cannot drop table"
	EXCEPTION_TABLE_INSERT_ERROR   = "Cannot insert in table"
	EXCEPTION_TABLE_DELETE_ERROR   = "Cannot delete from table"
	EXCEPTION_TABLE_UPDATE_ERROR   = "Cannot update in table"
	EXCEPTION_TABLE_FIND_ERROR     = "Cannot search in table"

	EXCEPTION_RELATION_ALREADY_EXISTS = "Relation already exists"
	EXCEPTION_RELATION_CREATE         = "Cannot create relation"
	EXCEPTION_RELATION_DROP           = "Cannot drop relation"
	EXCEPTION_RELATION_INSERT_ERROR   = "Cannot insert in relation"
	EXCEPTION_RELATION_DELETE_ERROR   = "Cannot delete from relation"
	EXCEPTION_RELATION_UPDATE_ERROR   = "Cannot update in relation"
	EXCEPTION_RELATION_FIND_ERROR     = "Cannot search in relation"

	//----------------------------------------------------------------
	EXCEPTION_JSON_MARSHAL   = "Cannot marshal struct"
	EXCEPTION_JSON_UNMARSHAL = "Cannot unmarshal string"
)
