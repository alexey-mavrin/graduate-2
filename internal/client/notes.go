package client

import (
	"github.com/alexey-mavrin/graduate-2/internal/common"
)

// ListNotes lists account for the current user
func (c *Client) ListNotes() (common.Notes, error) {
	records, err := c.listRecords(common.NoteRecord)
	if err != nil {
		return common.Notes{}, err
	}
	return *records.Notes, nil
}

// DeleteNote returns account record with the given id
func (c *Client) DeleteNote(id int64) error {
	return c.deleteRecord(id, common.NoteRecord)
}

// GetNote returns account record with the given id
func (c *Client) GetNote(id int64) (common.Note, error) {
	record, err := c.getRecord(id, common.NoteRecord)
	if err != nil {
		return common.Note{}, err
	}
	return *record.Note, nil
}

// UpdateNote updates account record with the given id
func (c *Client) UpdateNote(id int64, acc common.Note) error {
	record := common.Record{
		Type: common.NoteRecord,
		Note: &acc,
	}
	return c.updateRecord(id, record)
}

// StoreNote stores account record
func (c *Client) StoreNote(acc common.Note) (int64, error) {
	record := common.Record{
		Type: common.NoteRecord,
		Note: &acc,
	}
	return c.storeRecord(record)
}
