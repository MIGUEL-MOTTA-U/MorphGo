# Guía de Ejecución Completa - MorphGo con Datasets Reales

Este documento describe cómo probar todas las funcionalidades de MorphGo utilizando los archivos proporcionados en `test_data`.

## Preparación

Primero, asegúrate de compilar el binario:

```bash
go build -o morphgo main.go
```

---

## 1. Transformaciones de Datos (Casos de Éxito)

### JSON a CSV
```bash
./morphgo run \
  --input test_data/01_inspection_valid/employees01.json \
  --task "convert employees to csv" \
  --target csv
```

### CSV a JSON
```bash
./morphgo run \
  --input test_data/01_inspection_valid/employees.csv \
  --task "convert to json preserving all fields" \
  --target json
```

### XML a JSON
```bash
./morphgo run \
  --input test_data/01_inspection_valid/employees.xml \
  --task "extract all employees to json" \
  --target json
```

### YAML a JSON
```bash
./morphgo run \
  --input test_data/01_inspection_valid/employees.yaml \
  --task "convert to json" \
  --target json
```

### Excel (XLSX) a JSON
```bash
./morphgo run \
  --input test_data/01_inspection_valid/employees.xlsx \
  --task "convert excel sheets to json" \
  --target json
```

### CSV con Punto y Coma (Detectar Separador)
```bash
./morphgo run \
  --input test_data/01_inspection_valid/employees_semicolon.csv \
  --task "convert semicolon csv to json" \
  --target json
```

---

## 2. Manejo de Errores e Invalidaciones

### Archivo JSON Corrupto (Broken)
```bash
./morphgo run \
  --input test_data/02_inspection_invalid/broken.json \
  --task "convert to csv" \
  --target csv
```
*Resultado esperado:* Falla en etapa de inspección.

### Archivo XML Mal Formado
```bash
./morphgo run \
  --input test_data/02_inspection_invalid/broken.xml \
  --task "convert to json" \
  --target json
```
*Resultado esperado:* Error de parsing XML.

### Prompt Ambiguo (Falla de Planificación)
```bash
./morphgo run \
  --input test_data/01_inspection_valid/employees.csv \
  --task "haz algo con esto" \
  --target json
```
*Resultado esperado:* Error: "prompt is ambiguous".

---

## 3. Pruebas de Trazabilidad y Telemetría

Después de ejecutar cualquier comando anterior, verifica la carpeta `runs/`:

1. Identifica la última carpeta creada (formato `YYYYMMDD-HHMMSS-...`).
2. Abre `run.json` para ver:
   - El esquema detectado.
   - El plan técnico generado.
   - El código Go que se intentó ejecutar.
   - El resultado final (`success` o `failure`).

---

## 4. Pruebas de Integridad de Datos (E2E)

Para validar que los datos se mantienen íntegros, ejecuta una conversión circular:

1. **JSON -> CSV**:
   ```bash
   ./morphgo run --input test_data/01_inspection_valid/employees01.json --task "convert" --target csv
   ```
2. **CSV -> JSON** (usando la salida del paso anterior):
   ```bash
   ./morphgo run --input output.csv --task "convert back to json" --target json
   ```
3. Compara `output.json` con el original `test_data/01_inspection_valid/employees01.json`.

---

## 5. Pruebas Automáticas (Suite de Tests)

Para ejecutar la validación técnica completa incluyendo los nuevos tests de integridad:

```bash
go test -v ./internal/agent/...
```
