package instructions

import (
	"testing"

	"github.com/matryer/is"
	"github.com/veggiemonk/backlog/internal/core"
)

func TestDocCLICreate(t *testing.T) {
	for _, e := range Create.Examples {
		t.Run(e.Name, func(t *testing.T) {
			is := is.New(t)
			res, err := execute[core.CreateTaskParams, core.Task]("create", e.Params)
			is.NoErr(err)

			is.Equal(res, e.Expected)
		})
	}
}

func execute[In, Out any](name string, Params In) (Out, error) {
	var o Out
	return o, nil
}
