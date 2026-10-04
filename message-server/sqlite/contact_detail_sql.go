package sqlite

const CONTACT_DETAIL_READ string = `
SELECT Id, Phone, FirstName, LastName, LastUpdatedAtDateTime, LastUpdatedBy, Org FROM contacts  
	WHERE Id = ?;
`

const CONTACT_DETAIL_DELETE string = `
DELETE FROM contacts
	WHERE Id = ?
	RETURNING Id;
`
