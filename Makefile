all := core provider/aws provider/gcp provider/azure
providers := provider/aws provider/gcp provider/azure

install-tools:
	$(MAKE) install-tools -C core

binaries: $(providers)
	for dir in $(providers); do \
		$(MAKE) -C $$dir; \
	done

check: lint test

fmt: $(all)
	for dir in $(all); do \
		$(MAKE) fmt -C $$dir; \
	done

lint: $(all)
	for dir in $(all); do \
		$(MAKE) lint -C $$dir; \
	done

test-integration:
	@echo Running integration tests
	@cd ./e2e && make

test: $(all)
	for dir in $(all); do \
		$(MAKE) test -C $$dir; \
	done

test-coverage: $(all)
	for dir in $(all); do \
		$(MAKE) test-coverage -C $$dir; \
	done

license-check: $(providers)
	for dir in $(providers); do \
		$(MAKE) license-check -C $$dir; \
	done

generate-sources: $(all)
	for dir in $(all); do \
		$(MAKE) generate-sources -C $$dir; \
	done
