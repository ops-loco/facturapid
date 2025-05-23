# Prompt para Agente-IA-Developer-02: Desarrollo del Sincronizador Cliente para FactuRapid

## Contexto del Proyecto
Estás a cargo del desarrollo del programa sincronizador para el sistema FactuRapid, que operará en el lado del cliente. Tu misión es crear una aplicación en Go que extraiga información de tickets desde la base de datos de un TPV (Terminal Punto de Venta) de hostelería y la sincronice con el backend del sistema, todo ello sin interrumpir las operaciones normales del negocio.

## Objetivos Específicos del Sincronizador
1. Desarrollar una aplicación en Go que se ejecute como servicio en Windows 7
2. Conectar y extraer datos de una base de datos Microsoft Access
3. Sincronizar información de tickets con el backend a intervalos regulares
4. Gestionar errores, reintentos y registros de actividad
5. Minimizar la huella en el sistema cliente

## Requisitos Técnicos

### Lenguaje y Compilación
- Desarrollar íntegramente en Go (Golang)
- Compilación estática para minimizar dependencias externas
- Objetivo: instalación simple con mínimo impacto en sistema cliente
- Compatibilidad garantizada con Windows 7

### Conexión a Base de Datos
- Conectar a base de datos Microsoft Access (archivo "c:\tpv\tpv.mdb")
- Utilizar password de acceso: "mcdqn"
- Implementar manejo seguro de credenciales
- Adaptar a posibles limitaciones de drivers ODBC/OLE DB en Windows 7

### Funcionalidad de Sincronización
- Consultar periódicamente la base de datos (cada 5 minutos)
- Identificar tickets nuevos no sincronizados previamente
- Extraer datos completos de tickets:
  - Número de ticket
  - Fecha y hora
  - Detalles de productos/servicios
  - Cantidades, precios e impuestos
  - Totales
  - Otros campos relevantes para facturación
- Enviar datos al backend mediante API REST
- Mantener registro de tickets ya sincronizados

### Gestión de Errores y Reintentos
- Implementar estrategia robusta de manejo de errores
- Establecer política de reintentos con backoff exponencial
- Registro detallado de actividad (logging)
- Recuperación automática tras fallos
- Notificación de errores persistentes

### Instalación y Configuración
- Crear proceso de instalación simple (idealmente un solo ejecutable)
- Configuración mediante archivo externo o flags
- Parámetros configurables:
  - URL del backend
  - Credenciales de API
  - Ruta a base de datos
  - Intervalo de sincronización
  - Niveles de logging

### Operación como Servicio
- Implementar como servicio de Windows
- Arranque automático con el sistema
- Ejecución en segundo plano
- Mecanismo para iniciar/detener/reiniciar el servicio
- Monitorización de estado

## Integración con Backend
- Comunicación con API REST del backend
- Autenticación segura para comunicación con API
- Formato estandarizado para transmisión de datos de tickets
- Confirmaciones de recepción correcta de datos
- Gestión de versiones de API

## Consideraciones de Rendimiento
- Mínimo consumo de recursos (CPU, memoria)
- Operación sin interferir con el TPV
- Optimización de consultas a base de datos
- Eficiencia en transmisión de datos
- Gestión de concurrencia

## Requisitos de Seguridad
- Protección de credenciales de acceso
- Comunicación cifrada con backend (HTTPS)
- Validación de certificados SSL
- No exposición de datos sensibles en logs
- Verificación de integridad de datos transmitidos

## Entregables Esperados
1. Código fuente completo de la aplicación en Go
2. Ejecutable compilado para Windows 7 (32/64 bits según necesidad)
3. Documentación de instalación y configuración
4. Manual de solución de problemas
5. Pruebas unitarias y de integración

Puedes seleccionar las librerías de Go que mejor conozcas para conectar con bases de datos Access, comunicación HTTP, logging, y cualquier otra funcionalidad necesaria, siempre justificando tu elección en base a los requisitos del proyecto y asegurando que sean compatibles con la compilación estática.