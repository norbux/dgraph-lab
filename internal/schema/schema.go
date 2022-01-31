package schema

import (
	"context"

	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
)

func Setup(c *dgo.Dgraph) (err error) {
	// Drop all data and schema.
	err = c.Alter(context.Background(), &api.Operation{DropOp: api.Operation_ALL})
	if err != nil {
		return
	}

	err = c.Alter(context.Background(), &api.Operation{
		Schema: `
		type Email {
			address
			isdefault
		}

		address: string @index(term) .
		isdefault: bool .
		`,
	})

	if err != nil {
		return
	}

	err = c.Alter(context.Background(), &api.Operation{
		Schema: `
		type Phone {
			number
			isdefault
		}

		number: string @index(term) .
		isdefault: bool .
		`,
	})

	if err != nil {
		return
	}

	err = c.Alter(context.Background(), &api.Operation{
		Schema: `
		type Account {
			firstName
			lastName
			birthdate
			EmailList
			PhoneList
			status
		}

		firstName: string @index(term) .
		lastName: string @index(term) .
		birthdate: dateTime .
		EmailList: [uid] .
		PhoneList: [uid] .
		status: string .		
		`,
	})

	if err != nil {
		return
	}

	return
}
