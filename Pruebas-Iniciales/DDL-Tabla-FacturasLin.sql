-- `FacturasLin` definition

-- Drop table

-- DROP TABLE `FacturasLin`;

CREATE TABLE `FacturasLin` (
	`CodigoFactura` INTEGER NOT NULL,
	`UnidadesOld` SMALLINT,
	`Subtotal` DECIMAL(100,4),
	`CodigoProducto` VARCHAR(15),
	`Producto` VARCHAR(200) NOT NULL,
	`IvaAplicado` NUMERIC(100,7),
	`Linea` INTEGER NOT NULL,
	`Unidades` DECIMAL(100,4),
	`CombinadoCon` VARCHAR(15) NOT NULL,
	`LigaSiguiente` VARCHAR(1),
	`Serie` VARCHAR(1),
	CONSTRAINT SYS_PK_10895 PRIMARY KEY (`CodigoFactura`,`Producto`,`Linea`)
);
CREATE UNIQUE INDEX SYS_IDX_SYS_PK_10895_10896 ON `FacturasLin` (`CodigoFactura`,`Producto`,`Linea`);