package cl

const CONTACT_LIST_TABLEDEF string = `
CREATE TABLE contact_lists (
	Id INTEGER PRIMARY KEY,
	Name TEXT,
	Description TEXT,
	ContactList BLOB,
	CreatedAtDateTime BLOB,
	CreatedBy TEXT,
	LastUpdatedAtDateTime BLOB,
	LastUpdatedBy TEXT
);
CREATE INDEX pk_idx on contact_lists(Id);
`

const CONTACT_LIST_WRITE_NEW string = `
INSERT INTO contact_lists(Name, Description, ContactList, CreatedAtDateTime, CreatedBy, LastUpdatedAtDateTime, LastUpdatedBy)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	RETURNING Id;
`

const CONTACT_LIST_WRITE_EXISTING string = `
INSERT INTO contact_lists(Id, Name, Description, ContactList, CreatedAtDateTime, CreatedBy, LastUpdatedAtDateTime, LastUpdatedBy)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(Id) DO UPDATE SET
		Id=excluded.Id,
		Name=excluded.Name,
		Description=excluded.Description,
		ContactList=excluded.ContactList,
		CreatedAtDateTime=excluded.CreatedAtDateTime,
		CreatedBy=excluded.CreatedBy,
		LastUpdatedAtDateTime=excluded.LastUpdatedAtDateTime,
		LastUpdatedBy=excluded.LastUpdatedBy
	RETURNING Id;
`

const CONTACT_LIST_READ string = `
SELECT * FROM contact_lists  
	WHERE Id = ?;
`

const CONTACT_LIST_DELETE string = `
DELETE FROM contact_lists
	WHERE Id = ?
	RETURNING Id;
`
