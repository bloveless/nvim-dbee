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
	basicConstraintQuery := `
	SELECT tc.constraint_name, tc.table_name, kcu.column_name, ccu.table_name AS foreign_table_name, ccu.column_name AS foreign_column_name, rc.update_rule, rc.delete_rule
	FROM
		information_schema.table_constraints AS tc
		JOIN information_schema.key_column_usage AS kcu
			ON tc.constraint_name = kcu.constraint_name
		JOIN information_schema.referential_constraints as rc
			ON tc.constraint_name = rc.constraint_name
		JOIN information_schema.constraint_column_usage AS ccu
			ON ccu.constraint_name = tc.constraint_name
	`

	return map[string]string{
		"List":         fmt.Sprintf("SELECT * FROM %s LIMIT 500", opts.Table),
		"Columns":      fmt.Sprintf("SELECT column_name, spanner_type FROM information_schema.columns WHERE table_name = '%s'", opts.Table),
		"Indexes":      fmt.Sprintf("SELECT * FROM information_schema.indexes WHERE table_name='%s'", opts.Table),
		"Foreign Keys": fmt.Sprintf("SELECT * FROM information_schema.table_contraints WHERE table_name=%q AND constraint_type='FOREIGN KEY'", opts.Table),
		"Primary Keys": fmt.Sprintf("SELECT * FROM information_schema.table_contraints WHERE table_name=%q AND constraint_type='PRIMARY KEY'", opts.Table),
		"References":   "",
		"Check Constraints": fmt.Sprintf(`
			SELECT
			  cc.constraint_name,
			  tc.table_name,
			  tc.constraint_type,
			  tc.enforced,
			  cc.constraint_name,
			  cc.check_clause,
			  cc.spanner_state
			FROM
			  information_schema.table_constraints as tc
			  JOIN information_schema.check_constraints as cc on cc.constraint_name = tc.constraint_name
			WHERE table_name = %q`,
			opts.Table,
		),
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
