package inventory

import (
	"time"

	"strconv"

	"github.com/crit/inventory/internal/storage/models"
)

const EntryType = "inventory.Entry"

type Entry struct {
	ID      string    `json:"id"`
	Subject string    `json:"subject"`
	Change  int       `json:"change"`
	Date    time.Time `json:"date"`
	Editor  string    `json:"editor"`
}

func NewEntry(change, subject, editor string) error {
	i, err := strconv.Atoi(change)

	if err != nil {
		return err
	}

	e := Entry{
		ID:      models.NextID(),
		Subject: subject,
		Change:  i,
		Date:    time.Now().UTC(),
		Editor:  editor,
	}

	return e.Write()
}

func (e *Entry) Read() error {
	return models.Read(e)
}

func (e *Entry) Write() error {
	return models.Write(e)
}

func (e *Entry) Delete() error {
	return models.Delete(e)
}

func (e *Entry) Key() string {
	return e.ID
}

func (e *Entry) Type() string {
	return EntryType
}

func (e *Entry) Data() []byte {
	return models.ToBytes(e)
}

func (e *Entry) SetData(data []byte) {
	models.FromBytes(data, e)
}

type EntryList struct {
	Entries []Entry `json:"entries"`
}

func (l *EntryList) Type() string {
	return EntryType
}

func (l *EntryList) SetCount(count int64) {
	// do nothing
}

func (l *EntryList) Append(id, typ string, data []byte) {
	var entry Entry
	models.FromBytes(data, &entry)
	l.Entries = append(l.Entries, entry)
}

func (l *EntryList) Read() error {
	return models.List(l)
}

func (l *EntryList) Sum() (total int) {
	for _, entry := range l.Entries {
		total += entry.Change
	}

	return total
}

func (l *EntryList) Distinct() map[string]*EntryList {
	res := map[string]*EntryList{}

	for _, entry := range l.Entries {
		if _, ok := res[entry.Subject]; !ok {
			res[entry.Subject] = &EntryList{}
		}

		res[entry.Subject].Entries = append(res[entry.Subject].Entries, entry)
	}

	return res
}
