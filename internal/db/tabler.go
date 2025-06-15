package db

import (
	"bytes"
	"fmt"
	"iter"
	"text/tabwriter"
)

type Tabler interface {
	getHeader() string
	getRows() []string

}

func MakeTable(t Tabler) string {

	buf := &bytes.Buffer{}
	tb := tabwriter.NewWriter(buf, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tb, t.getHeader())
	for _, r := range t.getRows() {
		fmt.Fprintln(tb, r)
	}
	tb.Flush()
	return buf.String()
}

func (r *Donation) getRows() []string {
	return []string{
		fmt.Sprintf("%s\t%d\t%s\t%s\n", r.Channel, r.Amount, r.Startingdate, r.Endingdate),
	}
}

func (d *Donation) getHeader() string {
	return "Channel\tAmount\tStartingDate\tEndingDate"
}