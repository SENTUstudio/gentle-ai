# Glue Jobs

Place Glue ETL Python scripts here.

Example: `ETL_stg_encuestas_carga.py` — follows the pattern from `data-engineer-etl-*` skills.

## Creating a New Glue Job

1. Create your ETL script following the pattern from `data-engineer-etl-*` skills
2. Add a `AWS::Glue::Job` resource to `template.yaml`
3. Deploy: `aws cloudformation deploy --template-file template.yaml --stack-name carga-stg-encuestas`