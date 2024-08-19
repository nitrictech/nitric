all := core cloud/aws cloud/gcp cloud/azure cloud/common
providers := cloud/aws cloud/gcp cloud/azure

install-tools:
	$(MAKE) install-tools -C core

binaries: $(providers)
	for dir in $(providers); do \
		echo "Building $$dir"; \
		$(MAKE) -C $$dir || exit 1; \
	done

sec:
	for dir in $(all); do \
		echo "Running gosec on $$dir"; \
		$(MAKE) sec -C $$dir || exit 1; \
	done

check: lint test

fmt: $(all)
	for dir in $(all); do \
		echo "Formatting $$dir"; \
		$(MAKE) fmt -C $$dir || exit 1; \
	done

lint: $(all)
	for dir in $(all); do \
		echo "Linting $$dir"; \
		$(MAKE) lint -C $$dir || exit 1; \
	done

test: $(all)
	for dir in $(all); do \
		echo "Testing $$dir"; \
		$(MAKE) test -C $$dir || exit 1; \
	done

test-coverage: $(all)
	for dir in $(all); do \
		echo "Testing $$dir"; \
		$(MAKE) test-coverage -C $$dir || exit 1; \
	done

license-check: $(providers)
	for dir in $(providers); do \
		echo "Checking licenses for $$dir"; \
		$(MAKE) license-check -C $$dir || exit 1; \
	done

generate-sources: $(all)
	for dir in $(all); do \
		echo "Generating sources for $$dir"; \
		$(MAKE) generate-sources -C $$dir || exit 1; \
	done

tidy: $(all)
	@go work sync
	for dir in $(all); do \
		echo "Tidying $$dir"; \
		$(MAKE) tidy -C $$dir || exit 1; \
	done

.PHONY: install-tools binaries sec check fmt lint test test-coverage license-check generate-sources tidy