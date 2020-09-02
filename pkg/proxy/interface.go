package proxy

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/pingcap/parser/auth"
	"github.com/pingcap/tidb/sessionctx/variable"
	"github.com/pingcap/tidb/types"
	"github.com/pingcap/tidb/util"
	"github.com/pingcap/tidb/util/chunk"
)

// IDriver opens IContext.
type IDriver interface {
	// OpenCtx opens an IContext with connection id, client capability, collation, dbname and optionally the tls state.
	OpenCtx(connID uint64, capability uint32, collation uint8, dbname string, tlsState *tls.ConnectionState) (QueryCtx, error)
}

// QueryCtx is the interface to execute command.
type QueryCtx interface {
	// Status returns server status code.
	Status() uint16

	// LastInsertID returns last inserted ID.
	LastInsertID() uint64

	// LastMessage returns last info message generated by some commands
	LastMessage() string

	// AffectedRows returns affected rows of last executed command.
	AffectedRows() uint64

	// Value returns the value associated with this context for key.
	Value(key fmt.Stringer) interface{}

	// SetValue saves a value associated with this context for key.
	SetValue(key fmt.Stringer, value interface{})

	SetProcessInfo(sql string, t time.Time, command byte, maxExecutionTime uint64)

	// CommitTxn commits the transaction operations.
	CommitTxn(ctx context.Context) error

	// RollbackTxn undoes the transaction operations.
	RollbackTxn()

	// WarningCount returns warning count of last executed command.
	WarningCount() uint16

	// CurrentDB returns current DB.
	CurrentDB() string

	// Execute executes a SQL statement.
	Execute(ctx context.Context, db string, sql string) ([]ResultSet, error)

	// ExecuteInternal executes a internal SQL statement.
	ExecuteInternal(ctx context.Context, sql string) ([]ResultSet, error)

	// SetClientCapability sets client capability flags
	SetClientCapability(uint32)

	// Prepare prepares a statement.
	Prepare(sql string) (statement PreparedStatement, columns, params []*ColumnInfo, err error)

	// GetStatement gets PreparedStatement by statement ID.
	GetStatement(stmtID int) PreparedStatement

	// FieldList returns columns of a table.
	FieldList(tableName string) (columns []*ColumnInfo, err error)

	// Close closes the QueryCtx.
	Close() error

	// Auth verifies user's authentication.
	Auth(user *auth.UserIdentity, auth []byte, salt []byte) bool

	// ShowProcess shows the information about the session.
	ShowProcess() *util.ProcessInfo

	// GetSessionVars return SessionVars.
	GetSessionVars() *variable.SessionVars

	SetCommandValue(command byte)

	SetSessionManager(util.SessionManager)
}

// PreparedStatement is the interface to use a prepared statement.
type PreparedStatement interface {
	// ID returns statement ID
	ID() int

	// Execute executes the statement.
	Execute(context.Context, []types.Datum) (ResultSet, error)

	// AppendParam appends parameter to the statement.
	AppendParam(paramID int, data []byte) error

	// NumParams returns number of parameters.
	NumParams() int

	// BoundParams returns bound parameters.
	BoundParams() [][]byte

	// SetParamsType sets type for parameters.
	SetParamsType([]byte)

	// GetParamsType returns the type for parameters.
	GetParamsType() []byte

	// StoreResultSet stores ResultSet for subsequent stmt fetching
	StoreResultSet(rs ResultSet)

	// GetResultSet gets ResultSet associated this statement
	GetResultSet() ResultSet

	// Reset removes all bound parameters.
	Reset()

	// Close closes the statement.
	Close() error
}

// ResultSet is the result set of an query.
type ResultSet interface {
	Columns() []*ColumnInfo
	NewChunk() *chunk.Chunk
	Next(context.Context, *chunk.Chunk) error
	StoreFetchedRows(rows []chunk.Row)
	GetFetchedRows() []chunk.Row
	Close() error
}

// fetchNotifier represents notifier will be called in COM_FETCH.
type fetchNotifier interface {
	// OnFetchReturned be called when COM_FETCH returns.
	// it will be used in server-side cursor.
	OnFetchReturned()
}
