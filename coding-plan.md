# Plan incremental para un agente CLI CodeAct en Go

## Objetivo

Construir un **prototipo/sandbox** de un agente CLI en Go que:

1. Reciba un archivo de entrada y una intención en lenguaje natural.
2. Inspeccione el archivo para inferir estructura básica.
3. Genere un plan de transformación.
4. Produzca y ejecute código Go de forma controlada.
5. Capture errores de compilación o ejecución y los use para corregirse en iteraciones limitadas.
6. Deje trazabilidad para que otro agente pueda retomar el trabajo sin perder contexto.

## Alcance y restricciones

- **No es producción**.
- **Minimizar código**: evitar abstracciones innecesarias.
- **Sandbox primero**: cualquier ejecución de código generado debe ir aislada.
- **Progreso por etapas**: no avanzar si la etapa actual no pasa sus pruebas.
- **Pruebas unitarias estándar**: suficientes para asegurar robustez básica y manejo de errores, sin sobreingeniería.
- **Tareas paralelizables**: cuando dos partes sean independientes, se pueden resolver en paralelo.
 - **Prohibición de palabras**: las palabras prohibidas se listan en el archivo `.words` (no incluido en el repositorio).

---

## Stack mínimo recomendado

- `cobra` para la CLI.
- `os/exec` para invocar procesos y capturar `stdout/stderr`.
- `encoding/json`, `encoding/xml`, `encoding/csv`.
- `go.yaml.in/yaml/v4` para YAML.
- `excelize` para Excel.
- `goldmark` para Markdown.
- `pdfcpu` para lectura/apoyo con PDF si el caso lo permite.
- `log/slog` o logging simple.
- `testing` + `t.Run` para pruebas unitarias.
- Sandbox por proceso externo o contenedor simple.

---

## Flujo general del agente

1. **Entrada**: prompt + archivo.
2. **Inspección**: identificar tipo, estructura y señales del contenido.
3. **Planificación**: decidir qué transformación aplicar.
4. **Generación**: producir código Go temporal.
5. **Ejecución**: compilar/ejecutar en sandbox.
6. **Feedback loop**: si falla, corregir con base en el error.
7. **Salida**: artefacto transformado + registro de la corrida.

---

# Etapas de implementación

Cada etapa tiene una regla: **solo se avanza si las pruebas de esa etapa pasan**.

## Etapa 0 — Esqueleto del proyecto

### Objetivo
Tener una base mínima de CLI, estructura de carpetas y comando inicial.

### Entregables
- Proyecto Go inicializado.
- Comando `agente`.
- Subcomando `run`.
- Configuración mínima de flags.
- Estructura de carpetas establecida.

### Checklist
- [ ] `main.go` arranca la CLI.
- [ ] `agente --help` funciona.
- [ ] `agente run --help` funciona.
- [ ] `--input`, `--task`, `--output`, `--target` definidos.
- [ ] Config inicial cargable desde archivo opcional.
- [ ] Salida de errores clara y predecible.

### Pruebas mínimas
- [ ] Test del comando raíz.
- [ ] Test de ayuda de `run`.
- [ ] Test de validación de flags obligatorios.
- [ ] Test de error cuando falta `--input`.
- [ ] Test de error cuando falta `--task`.

### Criterio de salida
La CLI arranca, valida entrada y falla de forma limpia ante uso incorrecto.

---

## Etapa 1 — Inspección del archivo

### Objetivo
Detectar el tipo de archivo y extraer una estructura inicial útil para planificación.

### Entregables
- Detector básico de formato.
- Inspector por tipo: JSON, XML, CSV, YAML, TXT, Markdown, Excel, PDF.
- Modelo de esquema o resumen estructural.

### Checklist
- [ ] Detecta extensión o tipo por contenido.
- [ ] Para JSON: detecta llaves, arrays, tipos simples.
- [ ] Para XML: detecta nodos y atributos principales.
- [ ] Para CSV: detecta separador, encabezados y columnas.
- [ ] Para YAML: detecta mapas/listas relevantes.
- [ ] Para TXT/Markdown: detecta bloques y patrones.
- [ ] Para Excel: detecta hojas, filas y columnas.
- [ ] Para PDF: extrae texto si es posible y genera resumen estructural mínimo.
- [ ] Reporta errores cuando el archivo no existe o está corrupto.

### Pruebas mínimas
- [ ] Test de detección por extensión.
- [ ] Test de JSON válido.
- [ ] Test de JSON inválido.
- [ ] Test de CSV con encabezado.
- [ ] Test de archivo inexistente.
- [ ] Test de salida vacía o no legible.
- [ ] Test de ruta con tipo no soportado.

### Criterio de salida
El agente puede describir la estructura del archivo de forma confiable y sin romperse ante entradas malas.

### Tareas paralelizables
- Inspector de JSON/XML/CSV/YAML en paralelo.
- Detección de tipo y validación de entrada en paralelo con el diseño del modelo de esquema.

---

## Etapa 2 — Planificación de transformación

### Objetivo
Convertir la inspección en un plan corto y ejecutable.

### Entregables
- Estructura `Plan`.
- Generación de instrucciones de transformación.
- Traducción del prompt del usuario a una intención técnica mínima.

### Checklist
- [ ] El plan describe objetivo, formato origen, formato destino y pasos.
- [ ] El plan identifica si la operación es conversión, filtrado, extracción o normalización.
- [ ] El plan rechaza tareas ambiguas sin contexto suficiente.
- [ ] El plan conserva el resumen de esquema.
- [ ] El plan es serializable a JSON o Markdown.

### Pruebas mínimas
- [ ] Test de planificación desde esquema JSON.
- [ ] Test de planificación desde CSV.
- [ ] Test de manejo de prompt ambiguo.
- [ ] Test de serialización del plan.
- [ ] Test de estructura vacía.

### Criterio de salida
El sistema puede pasar de “qué hay” a “qué hacer” sin generar código todavía.

### Tareas paralelizables
- Mapeo de intención natural a tipo de transformación.
- Serialización y persistencia del plan.

---

## Etapa 3 — Generación de código Go

### Objetivo
Emitir código Go temporal para la transformación concreta.

### Entregables
- Plantilla base de programa Go.
- Generación de código según plan.
- Archivo temporal `main.go` por corrida.

### Checklist
- [ ] El código generado compila en el caso feliz.
- [ ] El código usa un formato simple y legible.
- [ ] El código incluye manejo básico de errores.
- [ ] El código no depende de lógica innecesaria.
- [ ] El código solo implementa lo que pide el plan.
- [ ] El código evita usar reflexión salvo que sea imprescindible.

### Pruebas mínimas
- [ ] Test de generación con un plan conocido.
- [ ] Test de que el código no queda vacío.
- [ ] Test de que contiene imports válidos.
- [ ] Test de que no produce funciones duplicadas.
- [ ] Test de plantilla con datos faltantes.

### Criterio de salida
El agente puede producir código Go temporal plausible y coherente.

### Tareas paralelizables
- Plantilla base de generación.
- Validación estática del código emitido.
- Construcción del sistema de prompts.

---

## Etapa 4 — Ejecución en sandbox

### Objetivo
Compilar y ejecutar el código generado con control de tiempo y salida.

### Entregables
- Runner de sandbox.
- Captura de `stdout`, `stderr` y código de salida.
- Timeout por ejecución.

### Checklist
- [ ] Ejecuta el código generado.
- [ ] Captura errores de compilación.
- [ ] Captura errores de ejecución.
- [ ] Aplica timeout.
- [ ] Aísla archivos temporales.
- [ ] Limpia artefactos intermedios.

### Pruebas mínimas
- [ ] Test de ejecución exitosa.
- [ ] Test de error de compilación.
- [ ] Test de error de runtime.
- [ ] Test de timeout.
- [ ] Test de limpieza de temporales.

### Criterio de salida
El agente puede correr el código generado sin dejar procesos colgados ni basura temporal.

### Tareas paralelizables
- Captura de salida y manejo de error.
- Gestión de temporales y cleanup.
- Timeout y cancelación.

---

## Etapa 5 — Feedback loop y autocorrección

### Objetivo
Usar errores de compilación o runtime para reintentar con corrección incremental.

### Entregables
- Bucle de iteración limitado.
- Estrategia simple de reparación por error.
- Historial de intentos.

### Checklist
- [ ] Lee stderr y clasifica el error.
- [ ] Regenera el código con base en el fallo.
- [ ] Limita el número de intentos.
- [ ] Detiene el loop si el error es repetitivo.
- [ ] Conserva trazabilidad de cada intento.

### Pruebas mínimas
- [ ] Test de corrección tras error de compilación.
- [ ] Test de abandono tras N intentos.
- [ ] Test de error repetido.
- [ ] Test de historial de intentos.

### Criterio de salida
El agente puede corregirse solo en escenarios simples y no entra en loops infinitos.

### Tareas paralelizables
- Clasificador de errores.
- Mecanismo de reintento.
- Registro de intentos.

---

## Etapa 6 — Persistencia de contexto y trazabilidad

### Objetivo
Dejar evidencia suficiente para retomar una corrida desde el último estado útil.

### Entregables
- Directorio de corrida.
- Archivos de inspección, plan, código generado, stderr y salida final.
- Identificador de ejecución.

### Checklist
- [ ] Guarda esquema.
- [ ] Guarda plan.
- [ ] Guarda código generado.
- [ ] Guarda stdout/stderr.
- [ ] Guarda decisión de éxito o fracaso.
- [ ] Guarda etapa actual y número de iteración.

### Pruebas mínimas
- [ ] Test de persistencia de artefactos.
- [ ] Test de lectura de corrida previa.
- [ ] Test de recuperación parcial.

### Criterio de salida
Cualquier agente puede leer el estado y continuar el trabajo.

### Tareas paralelizables
- Persistencia de artefactos.
- Carga de estado.
- Formato del registro de ejecución.

---

# División sugerida del trabajo por bloques

## Bloque A — Base y CLI
- Etapa 0 completa.
- Parte inicial de Etapa 6 para guardar corridas.

## Bloque B — Inspección
- Etapa 1 completa.
- Tests de errores y formatos.

## Bloque C — Planificación
- Etapa 2 completa.
- Serialización del plan.

## Bloque D — Generación y sandbox
- Etapa 3 y Etapa 4 completas.
- Ajuste del timeout y cleanup.

## Bloque E — Feedback loop
- Etapa 5 completa.
- Reintentos limitados.

## Bloque F — Trazabilidad final
- Etapa 6 completa.
- Registro consistente para retomar corridas.

---

# Definición de comandos del CLI

## Comando principal
```bash
agente run --input ./archivo.pdf --task "extrae y convierte a xml" --target xml --output ./out
```

## Comandos auxiliares
```bash
agente init
agente inspect --input ./archivo.pdf
agente plan --input ./archivo.pdf --task "..."
agente execute --code ./tmp/generated/main.go
agente replay --run-id 2026-06-01-001
```

---

# Criterios de calidad mínimos

## Debe cumplir
- Falla limpia ante archivos inexistentes o inválidos.
- Detecta tipos soportados sin romperse.
- Genera plan antes de generar código.
- Ejecuta en sandbox con timeout.
- Registra errores y decisiones.
- No avanza etapa sin tests verdes.

## No debe cumplir
- No hace falta arquitectura hexagonal completa.
- No hace falta plugin system.
- No hace falta observabilidad pesada.
- No hace falta orquestación distribuida.
- No hace falta un framework de agentes complejo.

---

# Registro de progreso para IA

Este bloque sirve para que cualquier agente pueda leer el archivo, ubicar el estado y continuar.

## Estado actual
- Etapa actual: `ETAPA_0`
- Subtarea actual: `pendiente`
- Última acción realizada: `ninguna`
- Último error relevante: `ninguno`
- Próximo paso: `definir e implementar la siguiente tarea pequeña`

## Formato de actualización

Cada vez que se avance, actualizar este bloque con información concreta.

```text
[IA-STEP]
fecha: 2026-06-03
etapa: FINAL
subtarea: integración completa
estado: completado
resultado: agente funcional de principio a fin integrado en el CLI con soporte de telemetría
siguiente: proyecto finalizado
bloqueos: ninguno
[/IA-STEP]
```

## Reglas para la IA que retome el trabajo
1. Leer el estado actual.
2. Verificar la última etapa completada.
3. Ejecutar o revisar los tests de esa etapa.
4. Solo avanzar si la etapa actual está estable.
5. Escribir una nueva entrada `[IA-STEP]` al finalizar cada avance.

---

# Lista de verificación global

- [ ] CLI base funcionando.
- [ ] Inspección por tipo de archivo.
- [ ] Planificación desde esquema.
- [ ] Generación de código Go temporal.
- [ ] Ejecución en sandbox.
- [ ] Feedback loop con corrección.
- [ ] Persistencia de corridas.
- [ ] Tests mínimos por etapa.
- [ ] Registro de progreso para retomar trabajo.
