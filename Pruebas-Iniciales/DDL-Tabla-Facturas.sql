-- `Facturas` definition

-- Drop table

-- DROP TABLE `Facturas`;

CREATE TABLE `Facturas` (
	`Codigo` INTEGER NOT NULL,
	`Cuenta` VARCHAR(20),
	`Fecha` TIMESTAMP,
	`Hora` TIMESTAMP,
	`Total` DECIMAL(100,4),
	`TipoCobro` VARCHAR(50),
	`Vendedor` VARCHAR(50),
	`CuotaIVA` DECIMAL(100,4),
	`Abonado` DECIMAL(100,4),
	`Terminal` VARCHAR(15),
	`Traspasada` VARCHAR(1),
	`Tarifa` VARCHAR(10) NOT NULL,
	`Base1` DECIMAL(100,4),
	`Base2` DECIMAL(100,4),
	`Base3` DECIMAL(100,4),
	`Iva1` DECIMAL(100,4),
	`Iva2` DECIMAL(100,4),
	`Iva3` DECIMAL(100,4),
	`CuotaIva1` DECIMAL(100,4),
	`CuotaIva2` DECIMAL(100,4),
	`CuotaIva3` DECIMAL(100,4),
	`Serie` VARCHAR(1) NOT NULL,
	`Cliente1` VARCHAR(30),
	`Cliente2` VARCHAR(30),
	`Cliente3` VARCHAR(30),
	`Cliente4` VARCHAR(30),
	`Revisable` VARCHAR(1),
	`Impresa` VARCHAR(1),
	`CobroMixto` DECIMAL(100,4),
	`EfectivoMixto` DECIMAL(100,4),
	`TipoCobroMixto` VARCHAR(50),
	`TipoCobroMixto2` VARCHAR(50),
	`Comensales` INTEGER,
	`CodigoDeFactura` INTEGER,
	`FechaDeFactura` TIMESTAMP,
	`HoraDeFactura` TIMESTAMP,
	`CobroMixto2` DECIMAL(100,4),
	`Base4` DECIMAL(100,4),
	`Base5` DECIMAL(100,4),
	`Base6` DECIMAL(100,4),
	`Iva4` DECIMAL(100,4),
	`Iva5` DECIMAL(100,4),
	`Iva6` DECIMAL(100,4),
	`CuotaIva4` DECIMAL(100,4),
	`CuotaIva5` DECIMAL(100,4),
	`CuotaIva6` DECIMAL(100,4),
	CONSTRAINT SYS_PK_10885 PRIMARY KEY (`Codigo`)
);
CREATE INDEX FACTURAS_KEYCLIENTE ON `Facturas` (`Cliente4`);
CREATE INDEX FACTURAS_KEYCLIENTE1 ON `Facturas` (`Cliente1`);
CREATE UNIQUE INDEX SYS_IDX_SYS_PK_10885_10886 ON `Facturas` (`Codigo`);

