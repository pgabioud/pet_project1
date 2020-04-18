# CS523-Project1
This is the repository for the project 1 of the CS-523 course for Justinas Sukaitis and Pierre Gabioud


The parameters of the main are: [Party ID] [Input] [CircuitID] [BFV Beaver Generation]
where "CircuitID" is the ID of the circuit to compute and "BFV Beaver Generation" is a boolean which is false if we want to run Part 1 and true if we want to run Part 2.

For example, to run the executable on inputs xi for i ∈ {0, 1, 2} on circuit7 for part 2, type:
./mpc 0 x0 1 true & ./mpc 1 x1 1 true &./mpc 2 x2 1 true


To run all the test, type:
go test -v

To run the test "TestEval" on a specific circuit (for Part 2), type:
go test -v -run=^TestEval$/^circuit7$




