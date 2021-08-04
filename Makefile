test_timeout := 30s

# =================================
#  			 SANDBOX 
# =================================

run-sandbox:
	cd .sandbox/ && go run .

build-sandbox: 
	cd .sandbox/ && go build .

build-run: 
	$(MAKE) build-sandbox 
	$(MAKE) run-exe args=""

run-exe: 
	cd .sandbox/ && ./.sandbox ${args}
		
# =============================
# 			TESTS 
# =============================

run-package-tests: 
	go test -timeout $(test_timeout) ${package} -v

run-tests: 
	$(MAKE) run-package-tests package="."
