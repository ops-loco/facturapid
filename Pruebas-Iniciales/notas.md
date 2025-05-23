Filename: C:\Tpv\tpv.mdb
The password is: mcdqn

Filename: C:\Tpv\data.mdb
The password is: mcdqn

Filename: C:\Tpv\master.mdb
The password is: mcetlqsm

Filename: C:\Tpv\profile.mdb
The password is: mcdqn

# Prompting

Me puedes mejorar y este prompt?, necesito tres prompts, uno para el "Agente IA-Arquitecto" tiene que tener toda la informacion general del proyecto, otro prompt que especifique el desarrollo del backoffice para el "Agente-IA-Developer-01", y para el "Agente-IA-Developer-02" que encarga del desarrollo del progrma sincronizador del lado cliente", deja que los agentes escojan las librerias que conozacan mejor y se adapten a las necesidades del proyecto. 

# Objetivo FactuRapid

Generar facturas oficiales en pdf a partir de los tickets que genera el tpv, usando un qr con la url del backoffice para que el cliente introduza el email, su nif, y el numero de ticket, si es nuevo recoger sus datos fiscales y guardarlos en la base de datos.
Dar la opcion de descargar o enviar por email.
Archivar la factura emitida enlazada con el id del cliente.
Autenticar a los clientes con su email, gmail.
El cliente tiene que tener la posibilidad de entrar posteriormente y poder consultar y descargarse sus facturas anteriores.
El sistema en lado cliente es un sistema con windows 7 y un programa de gestion de hosteleria (TpvFacil).
La base de datos del programa esta en un archivo Microsoft acces "tpv.mdb" protegida con la password "mcdqn"
Hay que acceder a los tickets para sincronizar los tickets con el backend sin interrumpir la operaciones del programa tpv.

El backoffice tiene que usar react, supabase y una api para sincronizar los tickets.

El programa sincronizador en el lado cliente, tiene que estar compilado en "go" con todas las librerias necesarias, compilado  estaticamente, el objetivo es instalar el minimo posible en el sistema. 
El programa sincronizador tiene que acceder a los tickets para sincronizar los tickets con el api backend sin interrumpir la operaciones del programa tpv, cada 5 mins.


****************

Poyecto Facturapid

Problema: 
Necesitamos que el  cliente se pueda descargar la factura en pdf rellenando el sus datos fiscales.
El sistema en lado cliente es un sistema con windows 7 y un programa de gestion de hosteleria (TpvFacil).
La base de datos del programa esta en un archivo Microsoft acces "tpv.mdb" protegida con la password "mcdqn"
el archivo es “c:\tpv\tpv.mdb”
queremos que cada vez que haya una factura con el cliente “QR” y se haya enviado a la impresora, suba los todos los datos de esa factura al backend a través de una api.
Genere un “qr” con el enlace al backend para rellenar los datos fiscales y lo imprima a la impresora de tickets que esta en el “lpt1:” la impersora es una Epson tm-t88
 Este sincronizador tiene que estar compilado en go y con las mínimas depencias del sistema operativo, seria ideal ejecutarlo como servicio.

Por otra parte el el frontend tiene que estar hecho con react y el hosting estara en un servidor plesk.

Aquí hay un ejemplo de la tabla factura:


Codigo	Cuenta	Fecha	Hora	Total	TipoCobro	Vendedor	CuotaIVA	Abonado	Terminal	Traspasada	Tarifa	Base1	Base2	Base3	Iva1	Iva2	Iva3	CuotaIva1	CuotaIva2	CuotaIva3	Serie	Cliente1	Cliente2	Cliente3	Cliente4	Revisable	Impresa	CobroMixto	EfectivoMixto	TipoCobroMixto	TipoCobroMixto2	Comensales	CodigoDeFactura	FechaDeFactura	HoraDeFactura	CobroMixto2	Base4	Base5	Base6	Iva4	Iva5	Iva6	CuotaIva4	CuotaIva5	CuotaIva6
40296	BAR1	07/03/2006	15:34:06	17,05 €	_EFECTIVO	_ADMIN	1,12 €	17,05 €	Term1	N	_NORMAL	15,93 €			7,00 €			1,12 €			A	QR	AVDA VILA DE TOSSA 56	17310 LLORET DE MAR	GIRONA	N	S																		

El campo que contiene el cliente es "Cliente1", y el campo que indica si se ha impreso es "Impresa"

y este es el ejemplo de la tabla que contiene las lineas de las facturas "FacturasLin":

CodigoFactura	UnidadesOld	Subtotal	CodigoProducto	Producto	IvaAplicado	Linea	Unidades	CombinadoCon	LigaSiguiente	Serie
40296		0,84 €	CAFESOLO	CAFES-CAFE SOLO	7	2	1,00 €			
40296		1,07 €	TEE	CAFES-TEE O TILA	7	1	1,00 €			
40296		14,02 €	MENU DEL DIA	PLATOS COMB-MENU DEL DIA	7	0	2,00 €			

