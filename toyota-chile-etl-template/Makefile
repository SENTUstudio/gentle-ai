# Toyota Chile ETL Project Template
# Structure for new ETL projects following Toyota Chile conventions

.PHONY: help install test dev prd clean

PYTHON_ENV := python3
PIP := pip3
VENV_DIR := .venv

help:
	@echo "Toyota Chile ETL Project - Available commands:"
	@echo ""
	@echo "  make install    - Install dependencies"
	@echo "  make dev        - Run development ETL"
	@echo "  make prd        - Run production ETL"
	@echo "  make test       - Run tests"
	@echo "  make clean      - Remove generated files"
	@echo ""

install:
	$(PIP) install -r requirements.txt

dev:
	@echo "Running development ETL..."
	@cd estudios && python3 -m <nombre_proyecto>.dev_<proyecto>

prd:
	@echo "Running production ETL..."
	@cd estudios && python3 -m <nombre_proyecto>.prd_<proyecto>

test:
	@echo "Running tests..."
	@cd estudios && python3 -m pytest

clean:
	@echo "Cleaning generated files..."
	@find . -type d -name "__pycache__" -exec rm -rf {} + 2>/dev/null || true
	@find . -type d -name "*.egg-info" -exec rm -rf {} + 2>/dev/null || true
	@find . -type f -name "*.pyc" -delete 2>/dev/null || true
	@rm -rf .pytest_cache 2>/dev/null || true

# ETL deployment commands (for Glue Jobs)
deploy-carga:
	@echo "Deploying CARGA stack..."
	aws cloudformation deploy \
		--template-file migracion/PAP/carga-datos-stg-encuestas/template.yaml \
		--stack-name $(STACK_NAME)-carga \
		--parameter-overrides Environment=prd

deploy-infra:
	@echo "Deploying INFRA stack..."
	aws cloudformation deploy \
		--template-file migracion/PAP/infra-datos-stg-encuestas/template.yaml \
		--stack-name $(STACK_NAME)-infra \
		--parameter-overrides Environment=prd

# S3 operations
sync-input:
	@echo "Syncing input files from S3..."
	aws s3 sync s3://$(BUCKET_RAW)/Input/ estudios/<nombre_proyecto>/input/ --exclude "HISTORICO/*"

sync-output:
	@echo "Syncing output files to S3..."
	aws s3 sync estudios/<nombre_proyecto>/output/ s3://$(BUCKET_STG)/<proyecto>/ --delete