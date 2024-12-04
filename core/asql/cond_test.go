package asql_test

import (
	"github.com/aarioai/airis/asql"
	"github.com/aarioai/airis/core/aenum"
	"testing"
)

func TestCond(t *testing.T) {
	var cond = &asql.Cond{}
	cond.And("t.id", "100")
	cond.Write("AND", aenum.StsInvalid("t.status"))
	cond.Try("t.ranking_woman DESC, t.vuid", 0, 20)

	if cond.Stmt() != " WHERE `t`.`id`=\"100\" AND t.status<0 ORDER BY t.ranking_woman DESC, t.vuid LIMIT 0,20" {
		t.Errorf("test cond failed `%s`", cond.Stmt())
	}
}
