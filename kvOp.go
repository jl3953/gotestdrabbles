package gotestdrabbles

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/cockroachdb/cockroach-go/crdb"
	"github.com/cockroachdb/cockroach/pkg/workload"
	"github.com/jackc/pgx"
	"strings"
)

type kvOp struct {
	sr workload.SQLRunner
	readStmt        workload.StmtHandle
	writeStmt       workload.StmtHandle
	db              *sql.DB
	mcp             *workload.MultiConnPool
}


func Read(ctx context.Context) error {

	urls := make([]string, 1)
	urls[0] = "localhost"
	// jenndebug this^ line will come back to haunt you

	db, err := sql.Open(`cockroach`, strings.Join(urls, ` `))
	if err != nil {
		return err
	}

	cfg := workload.MultiConnPoolCfg{
		MaxTotalConnections: 1,
	}
	mcp, err := workload.NewMultiConnPool(cfg, urls...)
	if err != nil {
		return err
	}

	o := &kvOp{
		db:	db,
		mcp:	mcp,
	}

	o.readStmt = o.sr.Define(`SELECT k, v FROM kv WHERE k IN (0)`)
	if err := o.sr.Init(ctx, "kv", mcp, nil); err != nil {
		return err
	}

	tx, err := o.mcp.Get().BeginEx(ctx, &pgx.TxOptions {
		IsoLevel: pgx.Serializable,
		AccessMode: pgx.ReadOnly})

	if err != nil {
		return err
	}

	args := make([]interface{}, 1)
	args[0] = 1

	err = crdb.ExecuteInTx(ctx, (*(workload.PgxTx))(tx), func() error {
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

	return err
}

