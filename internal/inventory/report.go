package inventory

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type ReportEntry struct {
	Subject string
	Count   int
}

func Report(writer io.Writer) {
	var buff bytes.Buffer
	var list EntryList

	defer func() {
		buff.WriteString("\n")
		writer.Write(buff.Bytes())
	}()

	if err := list.Read(); err != nil {
		buff.WriteString(err.Error())
		return
	}

	var report []ReportEntry
	var colSize int

	rows := list.Distinct()

	for subject, entries := range rows {
		report = append(report, ReportEntry{
			Subject: subject,
			Count:   entries.Sum(),
		})

		if colSize < len(subject) {
			colSize = len(subject)
		}
	}

	for _, item := range report {
		pad := strings.Repeat(" ", colSize-len(item.Subject))
		fmt.Fprintf(&buff, "%s:%s %d\n", item.Subject, pad, item.Count)
	}
}
