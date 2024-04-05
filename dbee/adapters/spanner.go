//go:build (darwin && (amd64 || arm64)) || (freebsd && (386 || amd64 || arm || arm64)) || (linux && (386 || amd64 || arm || arm64 || ppc64le || riscv64 || s390x)) || (netbsd && amd64) || (openbsd && (amd64 || arm64)) || (windows && (amd64 || arm64))

package adapters

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/googleapis/go-sql-spanner"
	"github.com/kndndrj/nvim-dbee/dbee/core"
	"github.com/kndndrj/nvim-dbee/dbee/core/builders"
)

// Register client
func init() {
	_ = register(&Spanner{}, "spanner")
}

var _ core.Adapter = (*Spanner)(nil)

type Spanner struct{}

func (s *Spanner) Connect(url string) (core.Driver, error) {
	db, err := sql.Open("spanner", url)
	if err != nil {
		log.Fatal(err)
	}

	return &spannerDriver{
		c: builders.NewClient(db),
	}, nil
}

func (*Spanner) GetHelpers(opts *core.TableOptions) map[string]string {
	return map[string]string{
		"List":    fmt.Sprintf("SELECT * FROM %q LIMIT 500", opts.Table),
		"Columns": fmt.Sprintf("SELECT column_name, spanner_type FROM information_schema.columns WHERE table_name = %q", opts.Table),
		"Indexes": fmt.Sprintf("SELECT * FROM information_schema.indexes WHERE table_name=%q", opts.Table),
		// "Foreign Keys": fmt.Sprintf("%s WHERE constraint_type = 'FOREIGN KEY' AND tc.table_name = '%s' AND tc.table_schema = '%s'",
		// 	basicConstraintQuery,
		// 	opts.Table,
		// 	opts.Schema,
		// ),
		// "References": fmt.Sprintf("%s WHERE constraint_type = 'FOREIGN KEY' AND ccu.table_name = '%s' AND tc.table_schema = '%s'",
		// 	basicConstraintQuery,
		// 	opts.Table,
		// 	opts.Schema,
		// ),
		// "Primary Keys": fmt.Sprintf("%s WHERE constraint_type = 'PRIMARY KEY' AND tc.table_name = '%s' AND tc.table_schema = '%s'",
		// 	basicConstraintQuery,
		// 	opts.Table,
		// 	opts.Schema,
		// ),
	}
}
