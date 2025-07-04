OUTPUT := src/sdk/gen
VENV := .venv
PYTHON := $(VENV)/bin/python
CONTRACTS := ../../proto

.PHONY: generate-proto install build clean

PACKAGE_NAME := sdk
BUILD_DIR := dist

generate-proto:
	@echo "Ensuring virtual environment exists..."
	@test -d $(VENV) || python3 -m venv $(VENV)
	@$(PYTHON) -m pip install -q -U pip grpcio-tools
	@$(VENV)/bin/python -m ensurepip --upgrade

	@echo "Cleaning old generated files..."
	@rm -rf $(OUTPUT)

	@echo "Generating Proto Sources into $(OUTPUT)..."
	@mkdir -p $(OUTPUT)
	@$(PYTHON) -m grpc_tools.protoc \
		--proto_path=$(CONTRACTS) \
		--python_out=$(OUTPUT) \
		--grpc_python_out=$(OUTPUT) \
		$(shell find $(CONTRACTS) -name '*.proto' ! -path '*/google/protobuf/*')

	@echo "Patching imports with Python script..."
	@$(PYTHON) ./patch/patch_imports.py
	
	@echo "Adding __init__.py files to make generated code importable..."
	@find $(OUTPUT) -type d -exec touch {}/__init__.py \;


install:
	@echo "Installing module in editable mode..."
	@$(PYTHON) -m pip install -e .

build:
	@echo "Building Python package..."
	@$(PYTHON) -m pip install -q build
	@$(PYTHON) -m build

clean:
	@echo "Cleaning up build artifacts..."
	@rm -rf build/ $(BUILD_DIR)/ *.egg-info $(OUTPUT)


publish-test: build
	@echo "Publishing to Test PyPI..."
	@$(PYTHON) -m pip install -q twine
	@$(PYTHON) -m twine upload --repository-url https://test.pypi.org/legacy/ dist/*.tar.gz dist/*.whl

publish-prod: build
	@echo "Publishing to PyPI..."
	@$(PYTHON) -m pip install -q twine
	@$(PYTHON) -m twine upload dist/*.tar.gz dist/*.whl

