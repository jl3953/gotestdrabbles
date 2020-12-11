package gotestdrabbles

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/cockroachdb/cockroach/pkg/workload"
	"github.com/cockroachdb/cockroach/pkg/workload/histogram"
)

type kvOp struct {
	readStmt        workload.StmtHandle
	writeStmt       workload.StmtHandle
	db              *sql.DB
	mcp             *workload.MultiConnPool
}

func (o *kvOp) read(ctx context.Context) {

	tx, err := o.mcp.Get().BeginEx(ctx, &pgx.TxOptions {
		IsoLevel: pgx.Serializable,
		AccessMode: pgx.ReadOnly,
	})

	crdb.ExecuteInTx(ctx, (*workload.PgxTx)(tx), func() error {
		rows, err := o.readStmt.QueryTx(ctx, tx, args...)
		if err != nil {
			return err
		}
		for rows.Next() {
			val, err := rows.Values()
			fmt.Printf("jenndebug read val:[%+v], err:[%+v]\n", val, err)
		}
		if rowErr := rows.Err(); rowErr != nil {
			return rowErr
		}
		rows.Close()
		return nil
	})
}