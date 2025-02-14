package renderer

import (
	"fmt"

	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/kubeshop/testkube/pkg/ui"
)

func TestSuiteRenderer(ui *ui.UI, obj interface{}) error {
	ts, ok := obj.(testkube.TestSuite)
	if !ok {
		return fmt.Errorf("can't use '%T' as testkube.TestSuite in RenderObj for test suite", obj)
	}

	ui.Warn("Name:     ", ts.Name)
	ui.Warn("Namespace:", ts.Name)
	ui.Warn("Labels:   ", fmt.Sprintf("%v", ts.Labels))

	steps := append(ts.Before, ts.Steps...)
	steps = append(steps, ts.After...)

	ui.Warn("Test steps:   ", fmt.Sprintf("%d", len(steps)))
	d := [][]string{{"Name", "Stop on failure", "Type"}}
	for _, step := range steps {
		d = append(d, []string{
			step.FullName(),
			fmt.Sprintf("%v", step.StopTestOnFailure),
			string(*step.Type()),
		})
	}

	ui.Table(ui.NewArrayTable(d), ui.Writer)

	return nil

}
