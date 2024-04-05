package adapters

import (
	"context"

	"github.com/kndndrj/nvim-dbee/dbee/core"
	"github.com/kndndrj/nvim-dbee/dbee/core/builders"
)

var _ core.Driver = (*spannerDriver)(nil)

type spannerDriver struct {
	c *builders.Client
}

func (c *spannerDriver) Query(ctx context.Context, query string) (core.ResultStream, error) {
	// run query, fallback to affected rows
	// return c.c.QueryUntilNotEmpty(ctx, query, "select changes() as 'Rows Affected'")
	return c.c.QueryUntilNotEmpty(ctx, query)
}

func (c *spannerDriver) Columns(opts *core.TableOptions) ([]*core.Column, error) {
	return c.c.ColumnsFromQuery("SELECT column_name, spanner_type FROM information_schema.columns WHERE table_name = '%s'", opts.Table)
}

func (c *spannerDriver) Structure() ([]*core.Structure, error) {
	query := `SELECT table_name FROM information_schema.tables WHERE table_schema = ''`

	rows, err := c.Query(context.TODO(), query)
	if err != nil {
		return nil, err
	}

	var schema []*core.Structure
	for rows.HasNext() {
		row, err := rows.Next()
		if err != nil {
			return nil, err
		}

		// We know for a fact there is only one string field (see query above)
		table := row[0].(string)
		schema = append(schema, &core.Structure{
			Name:   table,
			Schema: "",
			Type:   core.StructureTypeTable,
		})
	}

	return schema, nil
}

func (c *spannerDriver) Close() {
	c.c.Close()
}
