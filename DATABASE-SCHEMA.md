# 📊 Esquema Completo de Base de Datos - Plot Center

## Tablas Identificadas

### 1. **usuarios**
```sql
- id (PK)
- nombre
- password_hash
- rol (enum: 'administracion', 'taller', 'mostrador')
- last_seen
```

### 2. **ordenes_trabajo** (Tabla Principal)
```sql
- id (PK)
- numero_op
- cliente
- descripcion
- fecha_entrega
- estado (varchar)
- prioridad (varchar: 'Alta', 'Normal', 'Baja')
- fecha_creacion
- fecha_ingreso
- operario_asignado
- complejidad (enum: 'Baja', 'Media', 'Alta')
- sector (enum: 'Taller Gráfico', 'Mostrador')
- hora_estimada_entrega
- hora_entrega_efectiva
- id_usuario_creador (FK -> usuarios)
- usuario_trabajando_id
- usuario_trabajando_nombre
- timestamp_inicio_trabajo
```

### 3. **materiales**
```sql
- id (PK)
- codigo
- descripcion
```

### 4. **orden_materiales** (Relación Many-to-Many)
```sql
- id (PK)
- id_orden (FK -> ordenes_trabajo)
- id_material (FK -> materiales)
- cantidad (decimal)
```

### 5. **sectores**
```sql
- id (PK)
- nombre
- color
- activo (boolean)
- orden_visualizacion
- created_at
```

**Sectores predefinidos:**
- Diseño Gráfico
- En Espera
- Imprenta (Área de Impresión)
- Taller de Imprenta
- Taller Gráfico
- Instalaciones
- Metalúrgica
- Finalizado en Taller
- Almacén de Entrega
- Entregado o Instalado
- Mostrador

### 6. **orden_sectores** (Relación Many-to-Many)
```sql
- id (PK)
- id_orden (FK -> ordenes_trabajo)
- id_sector (FK -> sectores)
- fecha_asignacion
```

### 7. **archivos_adjuntos**
```sql
- id (PK)
- id_orden (FK -> ordenes_trabajo)
- nombre_archivo
- nombre_original
- fecha_subida
```

### 8. **historial_movimientos**
```sql
- id (PK)
- id_orden (FK -> ordenes_trabajo)
- id_usuario (FK -> usuarios)
- nombre_usuario
- estado_anterior
- estado_nuevo
- duracion_estado_anterior_seg
- timestamp
- comentario
```

### 9. **chat_rooms**
```sql
- id (PK)
- nombre
- tipo (enum: 'publico', 'privado')
- created_at
```

### 10. **chat_messages**
```sql
- id (PK)
- room_id (FK -> chat_rooms)
- id_usuario (FK -> usuarios)
- nombre_usuario
- mensaje
- timestamp
- message_type (varchar, opcional)
```

### 11. **notificaciones**
```sql
- id (PK)
- usuario_destino
- tipo
- mensaje
- id_orden (FK -> ordenes_trabajo, nullable)
- leida (boolean)
- timestamp
```

### 12. **notificaciones_vistas**
```sql
- id (PK)
- id_usuario (FK -> usuarios)
- id_historial (FK -> historial_movimientos)
- timestamp
```

### 13. **user_notifications**
```sql
- id (PK)
- user_id (FK -> usuarios)
- title
- description
- type
- orden_id (FK -> ordenes_trabajo, nullable)
- is_read (boolean)
- created_at
```

### 14. **alertas_enviadas**
```sql
- id (PK)
- id_orden (FK -> ordenes_trabajo)
- tipo_alerta (enum: 'estancada', 'plazo')
- timestamp
```

### 15. **smart_alerts**
```sql
- id (PK)
- tipo_alerta (enum: 'retraso_predicho', 'sobrecarga_operario', 'cuello_botella', 'eficiencia_baja')
- prioridad (enum: 'baja', 'media', 'alta', 'critica')
- titulo
- descripcion
- datos_contexto (JSON)
- fecha_creacion
- fecha_resuelto
- resuelto (boolean)
```

### 16. **tareas**
```sql
- id (PK)
- id_orden (FK -> ordenes_trabajo)
- descripcion_tarea
- estado_kanban (enum: 'Pendiente', 'En Proceso', 'Finalizado')
```

### 17. **online_users**
```sql
- user_id (FK -> usuarios)
- user_nombre
- last_seen
```

### 18. **stats_cache**
```sql
- id (PK)
- cache_key
- cache_value (JSON)
- expires_at
- created_at
```

### 19. **trending_metrics**
```sql
- id (PK)
- fecha
- metrica
- valor
- categoria
- subcategoria
- created_at
```

### 20. **prediction_metrics**
```sql
- id (PK)
- orden_id (FK -> ordenes_trabajo)
- numero_op
- tiempo_predicho_horas
- tiempo_real_horas
- confianza_prediccion
- error_absoluto
- error_porcentual
- factores_aplicados
- fecha_prediccion
- fecha_completado
```

### 21. **comentarios_orden**
```sql
- id (PK)
- id_orden (FK -> ordenes_trabajo)
- id_usuario (FK -> usuarios)
- comentario
- timestamp
```

### 22. **enlaces_adjuntos**
```sql
- id (PK)
- id_orden (FK -> ordenes_trabajo)
- url
- descripcion
- timestamp
```

### 23. **v_ordenes_stats** (Vista)
Vista materializada para estadísticas de órdenes.

## Relaciones Principales

```
usuarios (1) ──< (N) ordenes_trabajo (id_usuario_creador)
usuarios (1) ──< (N) historial_movimientos
usuarios (1) ──< (N) chat_messages
usuarios (1) ──< (N) user_notifications

ordenes_trabajo (1) ──< (N) orden_materiales
materiales (1) ──< (N) orden_materiales
ordenes_trabajo (N) >──< (N) materiales (vía orden_materiales)

ordenes_trabajo (1) ──< (N) orden_sectores
sectores (1) ──< (N) orden_sectores
ordenes_trabajo (N) >──< (N) sectores (vía orden_sectores)

ordenes_trabajo (1) ──< (N) archivos_adjuntos
ordenes_trabajo (1) ──< (N) historial_movimientos
ordenes_trabajo (1) ──< (N) tareas
ordenes_trabajo (1) ──< (N) alertas_enviadas
ordenes_trabajo (1) ──< (N) comentarios_orden
ordenes_trabajo (1) ──< (N) enlaces_adjuntos

chat_rooms (1) ──< (N) chat_messages
```

## Estados de Órdenes

Los estados se almacenan como strings, pero los principales son:
- `Pendiente`
- `Diseño Gráfico`
- `Diseño en Proceso`
- `En Espera`
- `Imprenta (Área de Impresión)`
- `Taller de Imprenta`
- `Taller Gráfico`
- `Instalaciones`
- `Metalúrgica`
- `Finalizado en Taller`
- `Almacén de Entrega`
- `Entregado o Instalado`
- `Mostrador`

## Prioridades

- `Alta`
- `Normal`
- `Baja`

## Complejidad

- `Baja`
- `Media`
- `Alta`

## Roles de Usuario

- `administracion` / `administración`
- `taller`
- `mostrador`

