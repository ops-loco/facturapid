# Prompt para Agente-IA-Developer-01: Desarrollo del Backoffice de FactuRapid

## Contexto del Proyecto
Estás a cargo del desarrollo del backoffice para el sistema FactuRapid, una plataforma que permite generar facturas oficiales en PDF a partir de tickets de TPV de hostelería. Tu responsabilidad es crear una aplicación web completa con React y Supabase que gestione todo el proceso de facturación, almacenamiento de datos de clientes y sincronización con el programa cliente.

## Objetivos Específicos del Backoffice
1. Desarrollar una interfaz web responsive para clientes que solicitan facturas
2. Implementar sistema de autenticación vía email y Gmail
3. Crear API para comunicación con el sincronizador cliente
4. Gestionar almacenamiento y recuperación de datos fiscales
5. Generar facturas oficiales en PDF con diseño profesional
6. Permitir consulta de histórico de facturas por cliente

## Requisitos Técnicos

### Frontend (React)
- Desarrollar con React (versión moderna con hooks)
- Crear formularios para captura de datos fiscales
- Implementar sistema de autenticación
- Diseñar interfaces para:
  - Solicitud de facturas por número de ticket
  - Gestión de perfil y datos fiscales
  - Consulta de historial de facturas
  - Descarga y reenvío de facturas
- Asegurar diseño responsive para móviles

### Backend (Supabase)
- Configurar base de datos en Supabase con tablas para:
  - Clientes y sus datos fiscales
  - Tickets sincronizados
  - Facturas generadas
  - Relaciones entre tickets y facturas
- Implementar funciones para:
  - Autenticación segura
  - Almacenamiento de datos fiscales
  - Generación de facturas
  - Envío de emails con facturas

### API
- Desarrollar endpoints RESTful para:
  - Recepción de datos de tickets desde el sincronizador
  - Verificación de tickets para facturación
  - Generación y recuperación de facturas
  - Gestión de datos de clientes
- Implementar validaciones y manejo de errores
- Asegurar la documentación completa de la API

### Sistema de Facturas
- Implementar generación de PDFs conforme a normativa fiscal
- Incluir todos los campos obligatorios:
  - Datos fiscales del emisor
  - Datos fiscales del receptor
  - Número de factura
  - Fecha de emisión
  - Desglose de conceptos, cantidades y precios
  - Cálculo de impuestos (IVA)
  - Total facturado
- Archivar facturas de forma segura y vinculadas al ID del cliente
- Permitir envío por email y descarga directa

## Flujo de Usuario
1. Cliente escanea QR en ticket físico
2. Accede a la web y se autentica (email o Gmail)
3. Introduce número de ticket
4. Si es cliente nuevo, completa formulario con datos fiscales
5. Sistema genera factura con datos del ticket y cliente
6. Cliente recibe factura por email y/o descarga directa
7. En visitas posteriores, puede consultar histórico de facturas

## Integración con Sincronizador
- Recibir datos de tickets enviados por el programa sincronizador
- Verificar y validar datos recibidos
- Almacenar información de tickets para posterior facturación
- Proporcionar endpoint seguro para la comunicación

## Requisitos de Seguridad
- Implementar autenticación segura (JWT, OAuth)
- Proteger datos fiscales sensibles
- Validar todas las entradas de usuario
- Asegurar que cada cliente solo acceda a sus propias facturas
- Implementar registro de actividad para auditoría

## Consideraciones de Diseño
- Interfaz intuitiva y sencilla para usuarios no técnicos
- Tiempos de carga rápidos
- Feedback claro durante procesos
- Manejo adecuado de errores con mensajes informativos
- Confirmaciones de acciones importantes (envío de email, etc.)

## Entregables Esperados
1. Código fuente completo del backoffice (React)
2. Estructura de base de datos en Supabase
3. Documentación de API para el sincronizador
4. Manual de despliegue y configuración
5. Pruebas unitarias y de integración

Puedes seleccionar las librerías y frameworks adicionales que mejor conozcas para React, gestión de estado, generación de PDFs, validación de formularios y cualquier otra funcionalidad necesaria, siempre justificando tu elección en base a los requisitos del proyecto.