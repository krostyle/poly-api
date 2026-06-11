# SPEC-11 — Rediseño del modelo de estados y subestados

## Contexto

El modelo de 8 estados original (LLAMADA → … → CIERRE/TERMINADO) no refleja el flujo real
del proceso Ley 20.009. La experta de negocio indicó correcciones estructurales:

- El banco **siempre** pasa por medida precautoria (prejudicial) antes de la demanda.
- **CIERRE viene DESPUÉS de TERMINADO**, no antes.
- Los motivos de término son un **enum fijo** de 12 valores, no texto libre.
- El proceso judicial tiene etapas propias: audiencia → sentencia → apelación (opcional)
  → sentencia 2ª instancia → cumplimiento.

## Nuevo grafo de estados (12 estados)

```
INGRESO → REVISION → PREJUDICIAL → PAGO_NORMATIVO → JUDICIAL
                ↓                        ↓                ↓
            TERMINADO               TERMINADO         AUDIENCIA
                                                          ↓
                                                      SENTENCIA → APELACION → SENTENCIA_SEGUNDA
                                                          ↓                           ↓
                                                      CUMPLIMIENTO ←─────────────────┘
                                                          ↓
                                                      TERMINADO → CIERRE
```

Tabla completa de transiciones:

| Estado actual    | Estados destino permitidos                       |
|------------------|--------------------------------------------------|
| INGRESO          | REVISION, TERMINADO                              |
| REVISION         | PREJUDICIAL, TERMINADO                           |
| PREJUDICIAL      | PAGO_NORMATIVO, JUDICIAL, TERMINADO              |
| PAGO_NORMATIVO   | JUDICIAL, TERMINADO                              |
| JUDICIAL         | AUDIENCIA, TERMINADO                             |
| AUDIENCIA        | SENTENCIA, TERMINADO                             |
| SENTENCIA        | APELACION, CUMPLIMIENTO, TERMINADO               |
| APELACION        | SENTENCIA_SEGUNDA, TERMINADO                     |
| SENTENCIA_SEGUNDA| CUMPLIMIENTO, TERMINADO                          |
| CUMPLIMIENTO     | TERMINADO                                        |
| TERMINADO        | CIERRE                                           |
| CIERRE           | (ninguno — estado final)                         |

## Semántica de estados

| Estado            | Significado de negocio                                              |
|-------------------|---------------------------------------------------------------------|
| INGRESO           | Banco registra el reclamo del cliente                               |
| REVISION          | Banco evalúa si la denuncia es válida                               |
| PREJUDICIAL       | Medida precautoria solicitada al tribunal                           |
| PAGO_NORMATIVO    | Tribunal acoge MP y ordena restituir el abono normativo             |
| JUDICIAL          | Demanda presentada (plazo: 10 días hábiles desde notificación MP)   |
| AUDIENCIA         | Audiencia fijada en tribunal                                        |
| SENTENCIA         | Tribunal dicta sentencia de primera instancia                       |
| APELACION         | Recurso de apelación interpuesto                                    |
| SENTENCIA_SEGUNDA | Sentencia de segunda instancia                                      |
| CUMPLIMIENTO      | Fase de cumplimiento o ejecución de la sentencia                    |
| TERMINADO         | Caso terminado por motivo específico (requiere `motivo_termino`)    |
| CIERRE            | Cierre administrativo — no se realizarán más gestiones              |

## Enum motivo_termino (12 valores)

| Valor                         | Descripción                                    |
|-------------------------------|------------------------------------------------|
| IMPROCEDENTE                  | Reclamo improcedente                           |
| EXTEMPORANEO                  | Presentado fuera de plazo                      |
| BUSQUEDAS_NEGATIVAS           | Sin resultados en búsquedas                    |
| DEUDOR_FALLECIDO              | Deudor fallecido                               |
| DESISTIMIENTO_CLIENTE         | Cliente se desiste del reclamo                 |
| DESISTIMIENTO_BANCO           | Banco se desiste de acciones                   |
| DESISTIMIENTO_DENUNCIA_INVALIDA| Desistimiento por denuncia inválida           |
| DESISTIMIENTO_SIN_DENUNCIA    | Desistimiento sin denuncia                     |
| SENTENCIA_FAVORABLE_BANCO     | Sentencia favorable al banco                   |
| SENTENCIA_DESFAVORABLE_BANCO  | Sentencia desfavorable al banco                |
| AVENIMIENTO                   | Acuerdo entre partes (avenimiento)             |
| ABANDONO_PROCEDIMIENTO        | Abandono del procedimiento                     |

## Plazos automáticos creados en cada transición

| Transición a      | Plazo creado                                      |
|-------------------|---------------------------------------------------|
| PREJUDICIAL       | PRECAUTELAR — 13 días hábiles                     |
| PAGO_NORMATIVO    | DEMANDA — 10 días hábiles, RESTITUCION_RECHAZO — 3 días |
| JUDICIAL          | DEMANDA — 10 días hábiles (seguimiento)           |

## Migración de datos existentes

| Estado antiguo       | Estado nuevo      | Nota                                    |
|----------------------|-------------------|-----------------------------------------|
| LLAMADA              | INGRESO           |                                         |
| REVISION             | REVISION          | sin cambio                              |
| SUSPENSION           | PREJUDICIAL       |                                         |
| PRE_JUDICIALIZACION  | PREJUDICIAL       |                                         |
| RESTITUCION          | PAGO_NORMATIVO    |                                         |
| JUDICIALIZACION      | JUDICIAL          |                                         |
| CIERRE               | TERMINADO         | El viejo cierre equivale al nuevo término|
| TERMINADO            | TERMINADO         | sin cambio                              |

Los `motivo_termino` en texto libre que no coincidan con el nuevo enum son limpiados (NULL).
