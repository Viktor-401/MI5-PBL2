/*
Participant Service

Each participant handles prepare, commit, and abort requests using MongoDB to ensure atomicity within their scope.
*/
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Account struct {
	ID      string  `bson:"_id"`
	Balance float64 `bson:"balance"`
}

type PendingTransaction struct {
	TxID      string    `bson:"tx_id"`
	AccountID string    `bson:"account_id"`
	Amount    float64   `bson:"amount"`
	Status    string    `bson:"status"`
	CreatedAt time.Time `bson:"created_at"`
}

type PrepareRequest struct {
	TxID      string  `json:"tx_id"`
	AccountID string  `json:"account_id"`
	Amount    float64 `json:"amount"`
}

var client *mongo.Client
var accountsCollection *mongo.Collection
var pendingCollection *mongo.Collection

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	accountsCollection = client.Database("bank").Collection("accounts")
	pendingCollection = client.Database("bank").Collection("pending_transactions")

	http.HandleFunc("/prepare", prepareHandler)
	http.HandleFunc("/commit", commitHandler)
	http.HandleFunc("/abort", abortHandler)

	fmt.Println("Participant service starting on :8081")
	http.ListenAndServe(":8081", nil)
}

// prepareHandler is an HTTP handler function that processes a "prepare" request.
// It validates the request payload, checks account balance, ensures the transaction
// is not already prepared, and inserts a new pending transaction into the database.
func prepareHandler(w http.ResponseWriter, r *http.Request) {
    // Define a variable to hold the request payload.
    var req PrepareRequest

    // Decode the JSON body of the request into the req variable.
    // If decoding fails, return a 400 Bad Request error to the client.
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Create a context with a timeout of 5 seconds to ensure the operation does not hang indefinitely.
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel() // Ensure the context is canceled to release resources.

    // Start a MongoDB session for the transaction.
    session, err := client.StartSession()
    if err != nil {
        // If starting the session fails, return a 500 Internal Server Error to the client.
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer session.EndSession(ctx) // Ensure the session is ended after the operation.

    // Execute the transaction using the session.
    _, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
        // Fetch the account document from the accounts collection using the provided AccountID.
        var account Account
        if err := accountsCollection.FindOne(sessCtx, bson.M{"_id": req.AccountID}).Decode(&account); err != nil {
            // If the account is not found, return an error indicating the account does not exist.
            return nil, fmt.Errorf("account not found")
        }

        // Check if the account has sufficient balance for the requested amount.
        if account.Balance < req.Amount {
            // If the balance is insufficient, return an error.
            return nil, fmt.Errorf("insufficient balance")
        }

        // Check if a transaction with the same TxID already exists in the pending transactions collection.
        var existingTx PendingTransaction
        if err := pendingCollection.FindOne(sessCtx, bson.M{"tx_id": req.TxID}).Decode(&existingTx); err == nil {
            // If a transaction with the same TxID exists, return an error indicating it is already prepared.
            return nil, fmt.Errorf("transaction already prepared")
        } else if err != mongo.ErrNoDocuments {
            // If an error occurs other than "no documents found," return the error.
            return nil, err
        }

        // Create a new pending transaction document.
        pendingTx := PendingTransaction{
            TxID:      req.TxID,          // Transaction ID from the request.
            AccountID: req.AccountID,    // Account ID from the request.
            Amount:    req.Amount,       // Amount from the request.
            Status:    "prepared",       // Set the status to "prepared."
            CreatedAt: time.Now(),       // Set the creation timestamp to the current time.
        }

        // Insert the new pending transaction into the pending transactions collection.
        if _, err := pendingCollection.InsertOne(sessCtx, pendingTx); err != nil {
            // If the insertion fails, return the error.
            return nil, err
        }

        // Return nil to indicate the transaction was successful.
        return nil, nil
    })

    // If the transaction fails, return a 400 Bad Request error with the error message.
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // If everything succeeds, return a 200 OK status and a "prepared" message to the client.
    w.WriteHeader(http.StatusOK)
    fmt.Fprint(w, "prepared")
}

func commitHandler(w http.ResponseWriter, r *http.Request) {
	var req struct{ TxID string }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	session, err := client.StartSession()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		var pendingTx PendingTransaction
		if err := pendingCollection.FindOne(sessCtx, bson.M{"tx_id": req.TxID, "status": "prepared"}).Decode(&pendingTx); err != nil {
			return nil, fmt.Errorf("transaction not found")
		}

		filter := bson.M{"_id": pendingTx.AccountID}
		update := bson.M{"$inc": bson.M{"balance": -pendingTx.Amount}}
		if res, err := accountsCollection.UpdateOne(sessCtx, filter, update); err != nil || res.MatchedCount == 0 {
			return nil, fmt.Errorf("update failed")
		}

		if _, err := pendingCollection.UpdateOne(sessCtx, bson.M{"tx_id": req.TxID}, bson.M{"$set": bson.M{"status": "committed"}}); err != nil {
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "committed")
}

func abortHandler(w http.ResponseWriter, r *http.Request) {
	var req struct{ TxID string }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := pendingCollection.DeleteOne(ctx, bson.M{"tx_id": req.TxID}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "aborted")
}