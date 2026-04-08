package model

import "context"

func (m *Model) Ping(ctx context.Context) string {
	db := m.db()
	_, e := db.QueryRow(ctx, `select 1`)
	if e != nil {
		return "[error] " + e.Msg
	}
	return "PONG"
}
