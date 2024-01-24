all := core cloud/aws cloud/gcp cloud/azure cloud/common
providers := cloud/aws cloud/gcp cloud/azure

install-tools:
	$(MAKE) install-tools -C core

binaries: $(providers)
	for dir in $(providers); do \
		$(MAKE) -C $$dir || exit 1; \
	done

sec:
	for dir in $(all); do \
		$(MAKE) sec -C $$dir || exit 1; \
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
