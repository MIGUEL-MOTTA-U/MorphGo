# MorphGo

## Descripción

MorphGo es una arquitectura de referencia para agentes CodeAct (Code-as-Action) desarrollada en Go. El proyecto implementa un sistema autónomo capaz de recibir archivos estructurados y descripciones de tareas en lenguaje natural para transformarlos mediante la generación y ejecución dinámica de código Go.

El sistema resuelve el problema de la transformación de datos (conversión, filtrado, extracción y normalización) operando en un ciclo cerrado que incluye inspección de archivos, planificación de tareas, generación de código temporal, ejecución en un entorno controlado (sandbox) y un bucle de retroalimentación para la autocorrección basada en errores de compilación o ejecución.

---

## Características

- **Inspección de Archivos:** Identificación automática de formato y extracción de metadatos estructurales (columnas, tipos básicos, conteo de registros).
- **Formatos Soportados:** Procesamiento verificado para JSON, XML, CSV, YAML, XLSX (Excel), Markdown y PDF.
- **Planificación Dinámica:** Inferencia de planes de transformación técnica a partir de prompts de usuario.
- **Ciclo CodeAct:** Generación automática de código fuente Go basado en el plan y el esquema detectado.
- **Sandbox Controlado:** Ejecución aislada de código generado con límites de tiempo (timeout) y captura de salida estándar/errores.
- **Autocorrección (Feedback Loop):** Bucle de reintentos parametrizable que analiza errores para intentar corregir el código generado.
- **Trazabilidad Completa:** Persistencia sistemática de cada ejecución, incluyendo el esquema inicial, el plan de acción, el historial de intentos de código y los resultados finales.

---

## Arquitectura

El proyecto está organizado en componentes desacoplados que gestionan el ciclo de vida del agente:

- **Orquestador (`internal/agent`)**: Coordina el flujo completo desde la entrada hasta el resultado final.
- **Inspector (`internal/parsers`)**: Módulos especializados en la lectura y análisis de diferentes formatos de archivo.
- **Planificador (`internal/schema`)**: Traduce la intención del usuario y el esquema del archivo en una estructura de datos `Plan`.
- **Generador de Código (`internal/codeact`)**: Sistema de plantillas que produce programas Go autónomos para ejecutar la transformación.
- **Sandbox (`internal/sandbox`)**: Motor de ejecución basado en procesos externos que aísla la corrida del código generado.
- **Telemetría (`internal/telemetry`)**: Subsistema encargado de la persistencia y auditoría de las corridas (Run Logs).

---

## Tecnologías Utilizadas

| Categoría | Tecnología |
|------------|------------|
| Lenguaje | Go 1.26.3 |
| CLI Framework | Cobra |
| Procesamiento de Excel | Excelize v2 |
| Procesamiento de YAML | yaml.v3 |
| Sandbox/Ejecución | os/exec (Go Standard Library) |
| Gestión de Estados | encoding/json |

---

## Requisitos

- Go 1.26.3 o superior.
- Sistemas compatibles con la ejecución de procesos externos vía `os/exec` (Windows, Linux, macOS).

---

## Instalación

Para compilar el binario del proyecto desde el código fuente:

```bash
go build -o morphgo main.go
```

---

## Ejecución

El agente se opera principalmente a través del subcomando `run`.

### Comandos CLI

```bash
# Ejecutar el flujo completo del agente
./morphgo run \
  --input <archivo_origen> \
  --task "<descripcion_tarea>" \
  --target <formato_destino> \
  [--output <ruta_trazas>]
```

---

## Uso

### Ejemplo de transformación de JSON a CSV

```bash
./morphgo run \
  --input test.json \
  --task "convert to csv" \
  --target csv
```

Este comando realizará las siguientes acciones:

1. Inspeccionará `test.json` para entender su estructura.
2. Generará un plan de transformación interna.
3. Producirá código Go para realizar la conversión.
4. Ejecutará el código en el sandbox.
5. Guardará un log detallado en el directorio `runs/`.

---

## Pruebas

El proyecto incluye una suite de pruebas unitarias y de integración que cubren los parsers, la lógica de planificación y el bucle CodeAct.

### Ejecutar todas las pruebas

```bash
go test ./...
```

### Ejecutar pruebas con salida detallada

```bash
go test -v ./...
```

---

## Estructura del Proyecto

```text
MorphGo/
├── cmd/                # Comandos de la CLI (root, run)
├── internal/
│   ├── agent/          # Orquestación e inspección centralizada
│   ├── codeact/        # Generación de código y bucle de reintentos
│   ├── parsers/        # Implementación de inspectores por formato
│   ├── sandbox/        # Ejecución controlada de código generado
│   ├── schema/         # Modelos de datos y lógica de planificación
│   └── telemetry/      # Sistema de trazabilidad y persistencia de corridas
├── main.go             # Punto de entrada de la aplicación
├── go.mod              # Definición del módulo y dependencias
└── coding-plan.md      # Plan de implementación incremental
```

---

## Observabilidad

MorphGo implementa un sistema de telemetría basado en archivos que registra cada ejecución en directorios específicos dentro de la carpeta `runs/` (o la ruta especificada mediante `--output`).

Cada directorio de corrida contiene un archivo `run.json` con la siguiente evidencia:

- ID único y marca de tiempo.
- Plan técnico generado.
- Historial de intentos, incluyendo código generado, salida estándar (`stdout`), errores (`stderr`) y códigos de salida de cada iteración del sandbox.
- Estado final de la ejecución (`success` / `failure`).