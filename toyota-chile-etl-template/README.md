# Toyota Chile ETL Project Template

Template para proyectos ETL de Toyota Chile siguiendo las convenciones del equipo de datos.

## Estructura del Proyecto

```
toyota-chile-etl-template/
├── Makefile                    # Comandos de deployment y ejecución
├── README.md                   # Este archivo
├── requirements.txt            # Dependencias Python
├── .gitignore
├── estudios/                   # Análisis y scripts ETL
│   ├── config/                 # Configuraciones dev/prd
│   │   ├── dev.yaml
│   │   └── prd.yaml
│   ├── <nombre_proyecto>/      # Proyecto específico
│   │   ├── dev_*.py           # Scripts para desarrollo
│   │   ├── prd_*.py           # Scripts para producción
│   │   ├── input/             # Archivos fuente (raw)
│   │   ├── output/            # Archivos procesados
│   │   │   ├── Quicksight/    # Outputs para QuickSight
│   │   │   └── dev/          # Outputs de desarrollo
│   │   ├── versions/          # Versiones históricas
│   │   └── esquema.md         # Documentación del esquema
│   └── docs/                  # Documentación
│       ├── arquitectura/
│       ├── datos/
│       └── migracion/
├── migracion/PAP/              # SAM templates (referencia)
│   ├── carga-datos-stg-encuestas/   # Template para Glue Jobs
│   ├── infra-datos-stg-encuestas/    # Template para Glue Tables
│   └── docs/
├── src/
│   └── artefactos/            # ETL scripts productivos
└── scripts/
    └── pyproject.toml
```

## Workflow de Uso

### 1. Nuevo Proyecto

1. Copiar este template
2. Renombrar `<nombre_proyecto>` al nombre del proyecto
3. Configurar `estudios/config/dev.yaml` y `estudios/config/prd.yaml`
4. Colocar archivos fuente en `estudios/<nombre_proyecto>/input/`

### 2. Análisis de Fuente (etl_from_*)

Usar skill `data-engineer-study-file` para analizar archivos:
- Encoding (UTF-8, Latin-1, ISO-8859-1)
- Delimitador (comma, semicolon, pipe, tab)
- Formato de fechas (ISO8601, DD/MM/YYYY, etc.)
- Caracteres anormales a limpiar

### 3. Generación de SQL (si no existe)

Usar skill `data-engineer-sql-from-logic` para generar SQL desde la lógica de negocio.
Guardar en `artefactos/<nombre_tabla>.sql` (FUERA de los repos CARGA/INFRA).

### 4. Crear Tabla (Workflow 3 Fases)

Usar skill `data-engineer-create-table`:

**Phase 1**: ETL en CARGA repo (comentar writes)
**Phase 2**: Tabla en INFRA repo (crear YAML + deploy)
**Phase 3**: Habilitar writes en CARGA repo

### 5. Integration Completa

Usar skill `data-engineer-integrate` para orquestar todo el workflow.

## Commands

```bash
make install      # Instalar dependencias
make dev          # Ejecutar ETL en desarrollo
make prd          # Ejecutar ETL en producción
make test         # Correr tests
make clean        # Limpiar archivos generados

# Deployment
make deploy-carga  # Deploy CARGA stack (Glue Jobs)
make deploy-infra  # Deploy INFRA stack (Glue Tables)
```

## Convenciones

- **Encoding Toyota Chile**: Latin-1 / ISO-8859-1 para archivos legacy
- **Delimiter**: Semicolon (;) para archivos CSV Toyota
- **Particiones**: year/month/day para tablas con volumen
- **Artefactos**: SQL y docs辅助ires en `artefactos/` (FUERA de CARGA/INFRA)

## Skills Relacionados

| Skill | Propósito |
|-------|-----------|
| `data-engineer-study-file` | Análisis de archivos (CSV/Excel) |
| `data-engineer-sql-from-logic` | Generación SQL desde lógica |
| `data-engineer-etl-s3` | ETL desde S3 con watermark |
| `data-engineer-etl-glue` | ETL desde Glue Catalog |
| `data-engineer-etl-sharepoint` | ETL desde SharePoint |
| `data-engineer-create-table` | Workflow 3 fases para crear tablas |
| `data-engineer-integrate` | Orquestación completa del proyecto |

## Repos Oficiales

- **CARGA**: Repositorio de ETL (Glue Jobs)
- **INFRA**: Repositorio de tablas (Glue Tables)
- **Artefactos**: SQL y docs auxiliares (fuera de CARGA/INFRA)