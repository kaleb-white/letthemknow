package sqlite

const CONTACT_TABLEDEF string = `
CREATE TABLE contacts (
	Id INTEGER PRIMARY KEY,
	Phone INTEGER,
	Phone2 INTEGER,
	Phone3 INTEGER,
	FirstName TEXT,
	LastName TEXT,
	FullName TEXT,
	ListMembership BLOB,
	CreatedAtDateTime BLOB,
	CreatedBy TEXT,
	LastUpdatedAtDateTime BLOB,
	LastUpdatedBy TEXT,
	Org TEXT	
);
CREATE INDEX pk_idx on contacts(Id);
`

const CONTACT_WRITE_NEW string = `
INSERT INTO contacts(Phone, Phone2, Phone3, FirstName, LastName, FullName, ListMembership, CreatedAtDateTime, CreatedBy, LastUpdatedAtDateTime, LastUpdatedBy, Org)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	RETURNING Id;
`

const CONTACT_WRITE_EXISTING string = `
INSERT INTO contacts(Id, Phone, Phone2, Phone3, FirstName, LastName, FullName, ListMembership, CreatedAtDateTime, CreatedBy, LastUpdatedAtDateTime, LastUpdatedBy, Org)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(Id) DO UPDATE SET
		Id=excluded.Id,
		Phone=excluded.Phone,
		Phone2=excluded.Phone2,
		Phone3=excluded.Phone3,
		FirstName=excluded.FirstName,
		LastName=excluded.LastName,
		ListMembership=excluded.ListMembership,
		FullName=excluded.FullName,
		CreatedAtDateTime=excluded.CreatedAtDateTime,
		CreatedBy=excluded.CreatedBy,
		LastUpdatedAtDateTime=excluded.LastUpdatedAtDateTime,
		LastUpdatedBy=excluded.LastUpdatedBy,
		Org=excluded.Org
	RETURNING Id;
`

const CONTACT_READ string = `
SELECT * FROM contacts  
	WHERE Id = ?;
`

const CONTACT_DELETE string = `
DELETE FROM contacts
	WHERE Id = ?
	RETURNING Id;
`
