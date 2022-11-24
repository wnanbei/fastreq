.PHONY: test jsonBenchmark

test:
	go test -v -cover

jsonBenchmark:
	go test -v -benchmem -run=^$$ -bench ^BenchmarkJson