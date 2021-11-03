test_timeout := 30s

# =============================
# 			TESTS 
# =============================

run-package-tests: 
	go test -timeout $(test_timeout) ${package} -v

run-tests: 
	$(MAKE) run-package-tests package="."

run-benchmarks: 
	go test -bench=. -benchmem -v 
