"""
Toyota Chile ETL - Sample Project Structure

Este archivo demuestra la estructura de un proyecto ETL completo.

Para crear un nuevo proyecto:
1. Copiar toda la carpeta toyota-chile-etl-template
2. Renombrar <nombre_proyecto> al nombre del proyecto
3. Configurar config/dev.yaml y config/prd.yaml
4. Personalizar los scripts en estudios/<nombre_proyecto>/
"""

# Ejemplo de estructura de archivo ETL para Glue
# Usar skill: data-engineer-etl-s3 / data-engineer-etl-glue / data-engineer-etl-sharepoint

# Ejemplo: etl_from_s3.py pattern
S3_PATTERN = """
# Glue job: ETL_mi_tabla_carga
# source: S3

import sys
import boto3
import datetime
import re
import pandas as pd
import unicodedata
import awswrangler as wr
import pyspark.sql.functions as F
from awsglue.context import GlueContext
from awsglue.job import Job
from awsglue.utils import getResolvedOptions
from pyspark.context import SparkContext

# 1. CONFIGURACIÓN
S3_FILE_PATTERN = "MI_ARCHIVO_"
S3_DESTINATION_PATH = "mi_tabla"
STG_BASE_PREFIX = "panel_experiencia_tcl"

# 2. FUNCIONES
def clean_column_name(name):
    if not name:
        return "unnamed_column"
    name = unicodedata.normalize("NFKD", name).encode("ascii", "ignore").decode("utf-8")
    name = name.lower().strip()
    name = re.sub(r"[^a-z0-9]", "_", name)
    name = re.sub(r"_+", "_", name).strip("_")
    return name

def get_last_load_timestamp(database, table):
    # Obtiene timestamp máximo de carga previa
    pass

def get_s3_whitelist(bucket, prefix, pattern, last_load_timestamp):
    # Lista archivos nuevos según watermark
    pass

# 3. INICIALIZACIÓN GLUE
sc = SparkContext.getOrCreate()
glueContext = GlueContext(sc)
spark = glueContext.spark_session

args = getResolvedOptions(sys.argv, ["JOB_NAME", "BUCKET_RAW", "BUCKET_STG", "TOPIC_SNS"])
job = Job(glueContext)
job.init(args["JOB_NAME"], args)

# 4. LECTURA CON WATERMARK
last_load_timestamp = get_last_load_timestamp(GLUE_DATABASE, GLUE_TABLE)
s3_files = get_s3_whitelist(BUCKET_RAW, "Output/", S3_FILE_PATTERN, last_load_timestamp)

df = spark.read.csv(s3_files, header=True, sep=";", encoding="latin-1")

# 5. TRANSFORMACIÓN
# ... transformaciones ...

# 6. ESCRITURA (PARA FASE 3 - DESCOMENTAR)
# df_final.write.mode("append").partitionBy("year", "month", "day").parquet(S3_OUTPUT_PATH)
# wr.athena.repair_table(table=GLUE_TABLE, database=GLUE_DATABASE)

# Para Fase 1 (extracción de esquema):
# df_final.printSchema()

job.commit()
"""

# Ejemplo: etl_from_glue.py pattern
GLUE_PATTERN = """
# Glue job: ETL_mi_tabla_carga
# source: Glue Catalog

from awsglue.context import GlueContext
from awsglue.job import Job
from awsglue.utils import getResolvedOptions
from pyspark.context import SparkContext

sc = SparkContext.getOrCreate()
glueContext = GlueContext(sc)
spark = glueContext.spark_session

args = getResolvedOptions(sys.argv, ["JOB_NAME"])
job = Job(glueContext)
job.init(args["JOB_NAME"], args)

# Lectura desde Glue Catalog
CATALOG_TABLES = [
    ("db_origen", "tabla_origen_1"),
    ("db_origen", "tabla_origen_2"),
]

for db_name, table_name in CATALOG_TABLES:
    dynamic_frame = glueContext.create_dynamic_frame.from_catalog(
        database=db_name,
        table_name=table_name,
        transformation_ctx=f"lectura_{table_name}",
    )
    df = dynamic_frame.toDF()
    df.createOrReplaceTempView(table_name)

# SQL con CTEs
query = \"\"\"
WITH paso_1 AS (
    SELECT * FROM tabla_origen_1 WHERE condicion = 1
),
paso_2 AS (
    SELECT * FROM tabla_origen_2
)
SELECT * FROM paso_2
\"\"\"

df_final = spark.sql(query)

# Para Fase 1 (extracción de esquema):
# df_final.printSchema()

# Para Fase 3 (descomentar escritura):
# s3_output_path = f"s3://{BUCKET}/mi_tabla/mi_tabla.parquet"
# df_final.coalesce(1).write.mode("overwrite").parquet(s3_output_path)

job.commit()
"""

# Ejemplo: etl_from_sharepoint.py pattern
SHAREPOINT_PATTERN = """
# Lambda function: ETL_mi_tabla_sharepoint

import os
import json
import msal
import pandas as pd
import awswrangler as wr
import boto3
from datetime import datetime

TENANT_ID = os.getenv("TENANT_ID")
CLIENT_ID = os.getenv("CLIENT_ID")
CLIENT_SECRET = os.getenv("CLIENT_SECRET")
SHAREPOINT_HOSTNAME = os.getenv("SHAREPOINT_HOSTNAME")
SITE_NAME = os.getenv("SITE_NAME")
FOLDER_PATH = os.getenv("FOLDER_PATH")
BUCKET_DESTINO = os.getenv("BUCKETSTGPROCESS")

def get_access_token():
    app = msal.ConfidentialClientApplication(
        client_id=CLIENT_ID,
        authority=f"https://login.microsoftonline.com/{TENANT_ID}",
        client_credential=CLIENT_SECRET
    )
    result = app.acquire_token_for_client(scopes=["https://graph.microsoft.com/.default"])
    token = result.get("access_token")
    if not token:
        raise RuntimeError(f"MSAL auth failed: {result.get('error_description', result.get('error', 'unknown'))}")
    return token

def lambda_handler(event, context):
    access_token = get_access_token()
    headers = {"Authorization": f"Bearer {access_token}"}

    # Descargar archivo de SharePoint
    # ...

    df = pd.read_csv(download_url, sep=";", header=None)
    df.columns = ["col1", "col2", "col3"]
    df["fecha_proceso"] = pd.to_datetime("today").date()

    # Para Fase 1: usar df.dtypes para extraer esquema
    # Para Fase 3:
    path_destino = f"s3://{BUCKET_DESTINO}/mi_tabla"
    wr.s3.to_parquet(df=df, path=path_destino, dataset=True, mode="overwrite")

    return {"statusCode": 200, "body": json.dumps({"mensaje": "OK"})}
"""