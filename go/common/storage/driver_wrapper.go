package storage

import (
	"context"
	"database/sql/driver"
	"fmt"

	gethlog "github.com/ethereum/go-ethereum/log"
)

// PanicOnDBErrorDriver wraps the MySQL driver and panics on connection refused errors
type PanicOnDBErrorDriver struct {
	driver.Driver
	logger     gethlog.Logger
	isCritical func(error) bool
}

// NewPanicOnDBErrorDriver is a constructor for PanicOnDBErrorDriver
func NewPanicOnDBErrorDriver(driver driver.Driver, logger gethlog.Logger, isCritical func(error) bool) *PanicOnDBErrorDriver {
	return &PanicOnDBErrorDriver{
		Driver:     driver,
		logger:     logger,
		isCritical: isCritical,
	}
}

// Open implements the driver.Driver interface
func (d *PanicOnDBErrorDriver) Open(dsn string) (driver.Conn, error) {
	conn, err := d.Driver.Open(dsn)
	if err != nil {
		return nil, err
	}
	return &panicOnErrorConn{Conn: conn, logger: d.logger, isCritical: d.isCritical}, nil
}

// panicOnErrorConn wraps a driver.Conn and panics on connection refused errors
type panicOnErrorConn struct {
	driver.Conn
	driver.ExecerContext
	driver.QueryerContext
	logger     gethlog.Logger
	isCritical func(error) bool
}

func (c *panicOnErrorConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	execer, ok := c.Conn.(driver.ExecerContext)
	if !ok {
		return nil, driver.ErrSkip
	}

	result, err := execer.ExecContext(ctx, query, args)
	if err != nil && c.isCritical(err) {
		c.logger.Crit(fmt.Sprintf("Database operation failed with connection refused: %v", err))
	}
	return result, err
}

func (c *panicOnErrorConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	queryer, ok := c.Conn.(driver.QueryerContext)
	if !ok {
		return nil, driver.ErrSkip
	}

	rows, err := queryer.QueryContext(ctx, query, args)
	if err != nil && c.isCritical(err) {
		c.logger.Crit(fmt.Sprintf("Database query failed with connection refused: %v", err))
	}
	return rows, err
}
