import pyodbc
import random
import os

def access_to_sql(access_file_path, output_sql_file):
    # Conectar a la base de datos Access
    conn_str = f'DRIVER={{Microsoft Access Driver (*.mdb, *.accdb)}};DBQ={access_file_path}'
    conn = pyodbc.connect(conn_str)
    cursor = conn.cursor()

    # Abrir archivo para escritura
    with open(output_sql_file, 'w', encoding='utf-8') as sql_file:
        # Obtener todas las tablas
        tables = []
        for row in cursor.tables():
            if row.table_type == 'TABLE':
                tables.append(row.table_name)
        
        for table in tables:
            # Escribir información del esquema
            sql_file.write(f"-- Estructura de la tabla `{table}`\n")
            sql_file.write(f"CREATE TABLE `{table}` (\n")
            
            # Obtener información de columnas
            columns = []
            primary_keys = []
            
            # Obtener columnas
            column_info = cursor.columns(table=table)
            for column in column_info:
                column_name = column.column_name
                data_type = column.type_name
                column_size = column.column_size
                nullable = "NULL" if column.nullable else "NOT NULL"
                
                # Mapear tipos de datos de Access a SQL
                sql_type = map_access_type_to_sql(data_type, column_size)
                
                columns.append(f"  `{column_name}` {sql_type} {nullable}")
            
            # Obtener claves primarias
            try:
                pk_info = cursor.primaryKeys(table=table)
                for pk in pk_info:
                    primary_keys.append(pk.column_name)
                
                if primary_keys:
                    columns.append(f"  PRIMARY KEY ({', '.join([f'`{pk}`' for pk in primary_keys])})")
            except:
                pass  # Algunas tablas pueden no tener clave primaria
            
            sql_file.write(',\n'.join(columns))
            sql_file.write("\n);\n\n")
            
            # Exportar aproximadamente el 10% de los datos
            cursor.execute(f"SELECT * FROM [{table}]")
            rows = cursor.fetchall()
            if rows:
                # Calcular cuántos registros son el 10%
                sample_size = max(1, int(len(rows) * 0.1))
                # Tomar una muestra aleatoria
                sample_rows = random.sample(rows, sample_size)
                
                # Obtener nombres de columnas
                column_names = [column[0] for column in cursor.description]
                
                # Escribir sentencias INSERT
                sql_file.write(f"-- Datos de ejemplo para la tabla `{table}` (10% aproximado)\n")
                for row in sample_rows:
                    values = []
                    for i, value in enumerate(row):
                        if value is None:
                            values.append("NULL")
                        elif isinstance(value, (int, float)):
                            values.append(str(value))
                        else:
                            # Escapar comillas simples en cadenas
                            value_str = str(value).replace("'", "''")
                            values.append(f"'{value_str}'")
                    
                    sql_file.write(f"INSERT INTO `{table}` ({', '.join([f'`{col}`' for col in column_names])}) VALUES ({', '.join(values)});\n")
                
                sql_file.write("\n")
    
    conn.close()
    print(f"Exportación completada. Archivo SQL generado: {output_sql_file}")

def map_access_type_to_sql(access_type, size):
    """Mapear tipos de datos de Access a tipos SQL genéricos"""
    access_type = access_type.upper()
    
    if access_type == 'COUNTER':
        return 'INTEGER AUTO_INCREMENT'
    elif access_type == 'INTEGER':
        return 'INTEGER'
    elif access_type == 'SMALLINT':
        return 'SMALLINT'
    elif access_type == 'BYTE':
        return 'TINYINT'
    elif access_type == 'DECIMAL' or access_type == 'CURRENCY' or access_type == 'MONEY':
        return 'DECIMAL(19,4)'
    elif access_type == 'REAL' or access_type == 'FLOAT' or access_type == 'DOUBLE':
        return 'FLOAT'
    elif access_type == 'DATETIME' or access_type == 'DATE':
        return 'DATETIME'
    elif access_type == 'BIT' or access_type == 'YESNO' or access_type == 'BOOLEAN':
        return 'BOOLEAN'
    elif access_type == 'LONGBINARY' or access_type == 'BINARY' or access_type == 'VARBINARY':
        return 'BLOB'
    elif access_type == 'MEMO' or access_type == 'LONGTEXT':
        return 'TEXT'
    elif access_type == 'CHAR' or access_type == 'VARCHAR' or access_type == 'TEXT':
        return f'VARCHAR({size})'
    else:
        return f'VARCHAR({size})'  # Tipo predeterminado

if __name__ == "__main__":
    # Ruta al archivo de Access
    access_file = input("Ingrese la ruta completa al archivo de Access (.mdb o .accdb): ")
    
    # Verificar que el archivo existe
    if not os.path.exists(access_file):
        print(f"Error: El archivo {access_file} no existe.")
        exit(1)
    
    # Ruta para el archivo SQL de salida
    output_file = input("Ingrese la ruta para guardar el archivo SQL (por defecto: export.sql): ") or "export.sql"
    
    # Ejecutar la exportación
    access_to_sql(access_file, output_file)
    