# Esquema de Datos - Proyecto Ejemplo

## Objetivo

Documentar el esquema de datos para el proyecto ETL.

## Fuentes de Datos

| Fuente | Tipo | Descripción |
|--------|------|-------------|
| `SERVICIO_EXPORT_S3.csv` | S3 | Encuestas de servicio CSI |
| `TOYOTA_SERVICIOS_HISTORICO.csv` | S3 (Histórico) | Historial de encuestas |
| `ipsos_servicio_s3/` | S3 | Archivos de Ipsos |

## Estructura de Archivos

### CSV Encoding Standard

| Archivo | Encoding | Delimiter | Header |
|---------|----------|-----------|--------|
| `SERVICIO_EXPORT_S3.csv` | UTF-8 | `\|` | Sí |
| `TOYOTA_SERVICIOS_HISTORICO.csv` | Latin-1 | `;` | Sí |

### Date Formats

| Campo | Formato | Locale | Ejemplo |
|-------|---------|--------|---------|
| `fecha_carga_proveedor` | `yyyy-MM-dd HH:mm:ss` | ISO8601 | `2024-01-15 10:30:00` |
| `fecha_proceso` | `yyyy-MM-dd HH:mm:ss` | ISO8601 | `2024-01-15 14:22:00` |

## Columnas del DataFrame Final

| Columna | Tipo | Descripción | nullable |
|---------|------|-------------|----------|
| `respondent_serial` | string | ID único de respondente | no |
| `respondent_id` | string | ID del respondente | no |
| `datacollection_status` | string | Estado de la encuesta | yes |
| `cod_sap` | string | Código SAP del dealer | no |
| `fecha` | string | Fecha de la encuesta | no |
| `year` | string | Partición: año | no |
| `month` | string | Partición: mes | no |
| `day` | string | Partición: día | no |
| `fecha_carga_proveedor` | timestamp | Fecha de carga | no |
| `fecha_proceso` | timestamp | Fecha de procesamiento | no |

## Cleaning Rules

1. **Encoding**: Normalizar a UTF-8
2. **Delimitador de fecha**: Asegurar formato ISO8601
3. **Whitespace**: Trim en todos los strings
4. **Caracteres anormales**: Remover control chars excepto TAB/LF/CR

## Metadatos

| Atributo | Valor |
|----------|-------|
| Particionado | Sí (year/month/day) |
| Formato | Parquet |
| Compresión | Snappy |
| Tipo de tabla | EXTERNAL_TABLE |

## Glosario

- **respondent_serial**: Identificador único de cada encuesta completada
- **respondent_id**: ID del sistema de encuestas externo
- **cod_sap**: Código de dealer SAP (6 dígitos)
- **fecha_carga_proveedor**: Timestamp del proveedor de datos
- **fecha_proceso**: Timestamp de procesamiento interno