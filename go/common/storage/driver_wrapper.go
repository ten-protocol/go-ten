package storage

import (
	"context"
	"database/sql/driver"
	"fmt"

	gethlog "github.com/ethereum/go-ethereum/log"
)

// NewPanicOnDBErrorDriver returns either a wrapped DriverContext or a Driver that panics on database errors
// When the database is not available, we prefer to fail instead of entering an inconsistency state.
func NewPanicOnDBErrorDriver(d driver.Driver, logger gethlog.Logger, isCritical func(error) bool, onNewConn func(conn driver.Conn)) driver.Driver {
	driverContext, isDriverContext := d.(driver.DriverContext)
	if isDriverContext {
		return &panicOnDBErrorDriverContext{
			driver:     driverContext,
			logger:     logger,
			isCritical: isCritical,
			onNewConn:  onNewConn,
		}
	}

	return &panicOnDBErrorDriver{
		driver:     d,
		logger:     logger,
		isCritical: isCritical,
		onNewConn:  onNewConn,
	}
}

// panicOnDBErrorDriver
type panicOnDBErrorDriver struct {
	driver     driver.Driver
	logger     gethlog.Logger
	isCritical func(error) bool
	onNewConn  func(conn driver.Conn)
}

func (d *panicOnDBErrorDriver) Open(dsn string) (driver.Conn, error) {
	conn, err := d.driver.Open(dsn)
	if err != nil {
		return nil, err
	}
	d.onNewConn(conn)
	return &panicOnErrorConnection{conn: conn, logger: d.logger, isCritical: d.isCritical}, nil
}

// panicOnDBErrorDriverContext
type panicOnDBErrorDriverContext struct {
	driver     driver.DriverContext
	logger     gethlog.Logger
	isCritical func(error) bool
	onNewConn  func(conn driver.Conn)
}

func (d *panicOnDBErrorDriverContext) OpenConnector(name string) (driver.Connector, error) {
	connector, err := d.driver.OpenConnector(name)
	if err != nil {
		return nil, err
	}

	return &panicOnErrorConnector{
		connector:  connector,
		logger:     d.logger,
		isCritical: d.isCritical,
		onNewConn:  d.onNewConn,
	}, nil
}

func (d *panicOnDBErrorDriverContext) Open(dsn string) (driver.Conn, error) {
	return nil, fmt.Errorf("not implemented")
}

// panicOnErrorConnector
type panicOnErrorConnector struct {
	connector  driver.Connector
	logger     gethlog.Logger
	isCritical func(error) bool
	onNewConn  func(conn driver.Conn)
}

func (conn *panicOnErrorConnector) Connect(ctx context.Context) (driver.Conn, error) {
	connection, err := conn.connector.Connect(ctx)
	if err != nil {
		return nil, err
	}
	conn.onNewConn(connection)
	return &panicOnErrorConnection{conn: connection, logger: conn.logger, isCritical: conn.isCritical}, nil
}

func (conn *panicOnErrorConnector) Driver() driver.Driver {
	return &panicOnDBErrorDriverContext{
		driver:     conn.connector.Driver().(driver.DriverContext),
		logger:     conn.logger,
		isCritical: conn.isCritical,
	}
}

// panicOnErrorConnection wraps a driver.Conn and panics on connection refused errors
type panicOnErrorConnection struct {
	conn       driver.Conn
	logger     gethlog.Logger
	isCritical func(error) bool
}

func (c *panicOnErrorConnection) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	execer, ok := c.conn.(driver.ExecerContext)
	if !ok {
		return nil, fmt.Errorf("underlying connection does not implement ExecContext")
	}

	result, err := execer.ExecContext(ctx, query, args)
	if err != nil && c.isCritical(err) {
		c.logger.Crit(fmt.Sprintf("Database operation failed with critical error: %v", err))
	}
	return result, err
}

func (c *panicOnErrorConnection) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	queryer, ok := c.conn.(driver.QueryerContext)
	if !ok {
		return nil, fmt.Errorf("underlying connection does not implement QueryContext")
	}

	rows, err := queryer.QueryContext(ctx, query, args)
	if err != nil && c.isCritical(err) {
		c.logger.Crit(fmt.Sprintf("Database query failed with critical error: %v", err))
	}
	return rows, err
}

func (c *panicOnErrorConnection) Exec(query string, args []driver.Value) (driver.Result, error) {
	execer, ok := c.conn.(driver.Execer)
	if !ok {
		return nil, fmt.Errorf("underlying connection does not implement Exec")
	}

	result, err := execer.Exec(query, args)
	if err != nil && c.isCritical(err) {
		c.logger.Crit(fmt.Sprintf("Database operation failed with critical error: %v", err))
	}
	return result, err
}

func (c *panicOnErrorConnection) Query(query string, args []driver.Value) (driver.Rows, error) {
	queryer, ok := c.conn.(driver.Queryer)
	if !ok {
		return nil, fmt.Errorf("underlying connection does not implement Query")
	}

	rows, err := queryer.Query(query, args)
	if err != nil && c.isCritical(err) {
		c.logger.Crit(fmt.Sprintf("Database query failed with critical error: %v", err))
	}
	return rows, err
}

func (c *panicOnErrorConnection) Ping(ctx context.Context) error {
	pinger, ok := c.conn.(driver.Pinger)
	if !ok {
		return fmt.Errorf("underlying connection does not implement Pinger")
	}

	err := pinger.Ping(ctx)
	if err != nil && c.isCritical(err) {
		c.logger.Crit(fmt.Sprintf("Database ping failed with critical error: %v", err))
	}
	return err
}

func (c *panicOnErrorConnection) Prepare(query string) (driver.Stmt, error) {
	stmt, err := c.conn.Prepare(query)
	if err != nil && c.isCritical(err) {
		c.logger.Crit(fmt.Sprintf("Database prepare failed with critical error: %v", err))
	}
	return stmt, err
}

func (c *panicOnErrorConnection) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	preparer, ok := c.conn.(driver.ConnPrepareContext)
	if !ok {
		return nil, fmt.Errorf("underlying connection does not implement ConnPrepareContext")
	}

	stmt, err := preparer.PrepareContext(ctx, query)
	if err != nil && c.isCritical(err) {
		c.logger.Crit(fmt.Sprintf("Database prepare context failed with critical error: %v", err))
	}
	return stmt, err
}

func (c *panicOnErrorConnection) Begin() (driver.Tx, error) {
	tx, err := c.conn.Begin()
	if err != nil && c.isCritical(err) {
		c.logger.Crit(fmt.Sprintf("Database begin transaction failed with critical error: %v", err))
	}
	return tx, err
}

func (c *panicOnErrorConnection) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	beginner, ok := c.conn.(driver.ConnBeginTx)
	if !ok {
		return nil, fmt.Errorf("underlying connection does not implement ConnBeginTx")
	}

	tx, err := beginner.BeginTx(ctx, opts)
	if err != nil && c.isCritical(err) {
		c.logger.Crit(fmt.Sprintf("Database begin transaction context failed with critical error: %v", err))
	}
	return tx, err
}

func (c *panicOnErrorConnection) Close() error {
	return c.conn.Close()
}

// ensure that panicOnErrorConnection implements all the required interfaces
var (
	_ driver.ExecerContext      = (*panicOnErrorConnection)(nil)
	_ driver.QueryerContext     = (*panicOnErrorConnection)(nil)
	_ driver.Execer             = (*panicOnErrorConnection)(nil)
	_ driver.Queryer            = (*panicOnErrorConnection)(nil)
	_ driver.Pinger             = (*panicOnErrorConnection)(nil)
	_ driver.ConnPrepareContext = (*panicOnErrorConnection)(nil)
	_ driver.ConnBeginTx        = (*panicOnErrorConnection)(nil)
)
