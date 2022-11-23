.PHONY: jsonBenchmark

jsonBenchmark:
	go test -benchmem -run=^$$ -bench ^BenchmarkJson