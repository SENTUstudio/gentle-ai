# Glue Job Deployment Instructions

## Overview

This directory contains reference SAM templates for AWS Glue Jobs (CARGA repo pattern).

## Structure

```
carga-datos-stg-encuestas/
├── template.yaml          # SAM template for Glue Jobs
└── glue-jobs/
    └── README.md          # Instructions for creating new Glue Jobs
```

## Creating a New Glue Job

### 1. Create the ETL Script

Create your ETL script following the pattern from `data-engineer-etl-*` skills:

```python
# Glue job: ETL_<table_name>_carga
# source: {s3|glue_catalog|sharepoint}
# author: gentle-ai

import sys
from awsglue.context import GlueContext
from awsglue.job import Job
from awsglue.utils import getResolvedOptions
from pyspark.context import SparkContext

# ... ETL code ...

# For schema extraction phase: COMMENT OUT these lines
# df_final.write.mode("overwrite").parquet(s3_output_path)
# job.commit()
```

### 2. Add Glue Job to template.yaml

```yaml
  EtlNombreTablaCarga:
    Type: AWS::Glue::Job
    Properties:
      Name: !Sub 'ETL_stg_encuestas_nombre_tabla_carga${DBDiscriminator}'
      Description: "ETL Nombre Tabla (Source: ...)"
      Role:
        Fn::ImportValue:
          !Sub "${InfraDatalakeStackName}-GlueGGExecutionRoleArn"
      GlueVersion: "5.0"
      WorkerType: G.4X
      NumberOfWorkers: 10
      Timeout: 2880
      Command:
        Name: glueetl
        PythonVersion: "3"
        ScriptLocation: ./glue-jobs/ETL_nombre_tabla_carga.py
      DefaultArguments:
        "--job-language": python
        "--enable-glue-datacatalog": "true"
        "--additional-python-modules": "awswrangler"
        "--job-bookmark-option": "job-bookmark-disable"
        "--TempDir": !Sub "s3://${BucketName}/temporary/"
        "--BUCKET_STG":
          Fn::ImportValue:
            !Sub '${InfraDatalakeStgStackName}-EncuestasName'
        "--BUCKET_RAW": !Sub "infra-contac-${Entorno}-contactabilidad-usuario"
        "--TOPIC_SNS":
          Fn::ImportValue:
            !Sub '${InfraWorkflowSNS}-topico1SnsDatalakeArn'
      Tags:
        Entorno: !Ref Entorno
        Project: "carga-datos-stg-encuestas-nombre-tabla"
```

### 3. Deploy

```bash
aws cloudformation deploy \
  --template-file template.yaml \
  --stack-name carga-datos-stg-encuestas \
  --parameter-overrides Environment=prd \
  --capabilities CAPABILITY_IAM
```

## Environment Variables Required

For Glue Jobs running with Python shell or PySpark:

- `AWS_DEFAULT_REGION`
- `AWS_ACCESS_KEY_ID` (via IAM role)
- `AWS_SECRET_ACCESS_KEY` (via IAM role)

For SharePoint jobs:
- `TENANT_ID`
- `CLIENT_ID`
- `CLIENT_SECRET`
- `SHAREPOINT_HOSTNAME`
- `SITE_NAME`
- `FOLDER_PATH`

## Schema Extraction Mode

When creating a new table:

1. **Phase 1**: Create ETL with write block COMMENTED OUT
2. Run locally or in dev to extract schema: `df_final.printSchema()`
3. Document columns for INFRA repo

## Schedule Triggers

Add to template.yaml:

```yaml
  ScheduleNombreTabla:
    Type: AWS::Glue::Trigger
    Properties:
      Type: SCHEDULED
      StartOnCreation: True
      Schedule: cron(0 3,10,13,18 * * ? *)
      Actions:
        - JobName: !Ref EtlNombreTablaCarga
      Name: !Sub 'job_ScheduleNombreTabla${DBDiscriminator}'
```