package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"dgraph-lab/internal/schema"

	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"google.golang.org/grpc"
)

type Account struct {
	Uid       string    `json:"uid,omitempty"`
	FirstName string    `json:"firstName,omitempty"`
	LastName  string    `json:"lastName,omitempty"`
	EmailList []Email   `json:"emailList,omitempty"`
	PhoneList []Phone   `json:"phoneList,omitempty"`
	Birthdate time.Time `json:"birthdate,omitempty"`
	Status    string    `json:"status,omitempty"`
	DType     []string  `json:"dgraph.type,omitempty"`
}

type Email struct {
	Address   string   `json:"address,omitempty"`
	IsDefault bool     `json:"isDefault,omitempty"`
	DType     []string `json:"dgraph.type,omitempty"`
}

type Phone struct {
	Number    string   `json:"phone,omitempty"`
	IsDefault bool     `json:"isDefault,omitempty"`
	DType     []string `json:"dgraph.type,omitempty"`
}

func newClient() *dgo.Dgraph {
	dgraphClient, err := grpc.Dial("localhost:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	return dgo.NewDgraphClient(
		api.NewDgraphClient(dgraphClient),
	)
}

func main() {
	dgraphClient := newClient()
	err := schema.Setup(dgraphClient)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Schema created.")

	ctx := context.Background()

	// Input data
	name := "Juan"
	lastName := "Perez"
	birthdate := time.Date(1980, 01, 01, 22, 0, 0, 0, time.UTC)
	email := "jperez@example.com"

	// Insert 2 nodes and edge at once
	account := &Account{
		FirstName: name,
		LastName:  lastName,
		Birthdate: birthdate,
		Status:    "GUEST",
		DType:     []string{"Account"},
		EmailList: []Email{
			{Address: email, IsDefault: true, DType: []string{"Email"}},
		},
	}

	mut := &api.Mutation{
		CommitNow: true,
	}

	accountBytes, err := json.Marshal(account)
	if err != nil {
		log.Fatal(err)
	}

	mut.SetJson = accountBytes
	response, err := dgraphClient.NewTxn().Mutate(ctx, mut)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response Uids: %+v", response.Uids)

	// Query a node and add an edge to a new node on it
	query := `{
		user(func: eq(firstName, "Juan")) {
			uid
		}
	}`

	queryResp, err := dgraphClient.NewTxn().Query(ctx, query)
	if err != nil {
		log.Fatal(err)
	}

	type uid map[string]string
	type qResp struct {
		User []uid `json:"user"`
	}
	var u qResp
	err = json.Unmarshal(queryResp.Json, &u)
	if err != nil {
		log.Fatal()
	}

	id := u.User[0]["uid"]
	log.Printf("Uid: %v", id)

	updateAccount := &Account{
		Uid:   id,
		DType: []string{"Account"},
		EmailList: []Email{
			{Address: "newemail@gmail.com", IsDefault: false, DType: []string{"Email"}},
		},
	}

	update := &api.Mutation{
		CommitNow: true,
	}

	accountBytes, err = json.Marshal(updateAccount)
	if err != nil {
		log.Fatal(err)
	}

	update.SetJson = accountBytes
	updateResponse, err := dgraphClient.NewTxn().Mutate(ctx, update)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response Uids: %+v", updateResponse.Uids)
}
