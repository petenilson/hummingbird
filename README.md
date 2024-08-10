## About

Hummingbird is a double entry ledger written in Go with a basic HTTP API on top in order to expose the functionality of the ledger. Hummingbird is meant as an experiment in a ledger design and it's implementation using Go. The goal of this project is to test my ability to build a ledger that is stable under a high level of concurrency, measure the performance, then optimize the speed of the system for Transactions per second. 

As the goal of the project was less focused on the API, and more on the underlying ledger, there are some notable things Hummingbird currently lacks such as authentication and proper validation of user data on endpoints. Hummingbird is thus not fit for production use in it's current state.

The ledger provides primitives for building on top of in the form of Accounts, Transactions and Entry's. A Transaction represents a complete movement of money between any number of Accounts. An Entry is a record of one of those movements. As this is a double entry ledger, every Transaction will have at least two associated Entry's. A Transaction may for example be the movement of funds between two Accounts or a payment for a service and the fee's associated with that service.

## Load Testing

This repo contains a load test with the goal of testing how many transactions per second can be created. Not transactions in the Database sense, but rather business transactions. Therefore the resulting requests per second of the load test accurately measure TPS. Below are the results of the initial load test where new payments are created between 1000 different accounts randomly.



Running 30s test @ http://localhost:8000
  12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    48.80ms    7.81ms 176.01ms   93.80%
    Req/Sec   409.59    121.03   760.00     73.08%
  146874 requests in 30.03s, 55.89MB read
  Socket errors: connect 157, read 109, write 0, timeout 0
Requests/sec:   4891.36
Transfer/sec:      1.86MB
