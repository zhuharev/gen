package gen

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPrepareProject(t *testing.T) {

	Convey("test prepare", t, func() {
		cases := map[string]struct {
		}{
			"relation title": {},
		}
		for name, _ := range cases {
			Convey(name, func() {
				//So(c.exp.Models, ShouldResemble, prepareProject(c.inp).Models)
			})
		}
	})
}
