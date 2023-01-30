all := core cloud/aws cloud/gcp cloud/azure
providers := cloud/aws cloud/gcp cloud/azure

install-tools:
	$(MAKE) install-tools -C core

binaries: $(providers)
	for dir in $(providers); do \
		$(MAKE) -C $$dir || exit 1; \
	done

check: lint test

fmt: $(all)
	for dir in $(all); do \
		$(MAKE) fmt -C $$dir || exit 1; \
	done

lint: $(all)
	for dir in $(all); do \
		$(MAKE) lint -C $$dir || exit 1; \
	done

test-integration:
	@echo Running integration tests
	@cd ./e2e && make

test: $(all)
	for dir in $(all); do \
		$(MAKE) test -C $$dir || exit 1; \
	done

test-coverage: $(all)
	for dir in $(all); do \
		$(MAKE) test-coverage -C $$dir || exit 1; \
	done

license-check: $(providers)
	for dir in $(providers); do \
		$(MAKE) license-check -C $$dir || exit 1; \
	done

generate-sources: $(all)
	for dir in $(all); do \
		$(MAKE) generate-sources -C $$dir || exit 1; \
	done
