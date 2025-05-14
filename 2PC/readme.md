# Explanation:

## Coordinator Service:

- Initiates transactions via /start-transaction.

- Sends prepare requests to all participants. If all succeed, sends commit; otherwise, abort.

## Participant Service:

- Prepare: Checks account balance, records a pending transaction in MongoDB within a transaction to ensure atomicity.

- Commit: Deducts the balance and updates the transaction status within a MongoDB transaction.

- Abort: Deletes the pending transaction.

# Usage:

- Run the coordinator on port 8080.

- Run participant(s) on different ports (e.g., 8081).

- Use a REST client to POST a transaction request to the coordinator's /start-transaction endpoint with participant details.

This example demonstrates the 2PC flow, ensuring atomicity across distributed services using HTTP and MongoDB transactions.