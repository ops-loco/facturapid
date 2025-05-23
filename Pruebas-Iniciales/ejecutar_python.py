#!/usr/bin/env python
# -*- coding: utf-8 -*-

import os
import sys
import subprocess
import argparse

def verificar_python():
    """Verifica la instalación de Python y devuelve la ruta del ejecutable."""
    try:
        # Intentar usar python como comando
        result = subprocess.run(['python', '--version'], 
                               stdout=subprocess.PIPE, 
                               stderr=subprocess.PIPE,
                               text=True,
                               check=False)
        
        if result.returncode == 0:
            return 'python'
        
        # Intentar usar py como comando (Windows Python Launcher)
        result = subprocess.run(['py', '--version'], 
                               stdout=subprocess.PIPE, 
                               stderr=subprocess.PIPE,
                               text=True,
                               check=False)
        
        if result.returncode == 0:
            return 'py'
            
        # Intentar con python3
        result = subprocess.run(['python3', '--version'], 
                               stdout=subprocess.PIPE, 
                               stderr=subprocess.PIPE,
                               text=True,
                               check=False)
        
        if result.returncode == 0:
            return 'python3'
            
        print("No se pudo encontrar una instalación de Python. Por favor, instala Python y asegúrate de que esté en el PATH.")
        sys.exit(1)
    except Exception as e:
        print(f"Error al verificar Python: {e}")
        sys.exit(1)

def ejecutar_script(python_cmd, script_path, args=None):
    """Ejecuta un script Python con los argumentos proporcionados."""
    cmd = [python_cmd, script_path]
    
    if args:
        cmd.extend(args)
        
    try:
        print(f"Ejecutando: {' '.join(cmd)}")
        proceso = subprocess.run(cmd, check=True)
        return proceso.returncode
    except subprocess.CalledProcessError as e:
        print(f"Error al ejecutar el script: {e}")
        return e.returncode
    except Exception as e:
        print(f"Error inesperado: {e}")
        return 1

def main():
    parser = argparse.ArgumentParser(description='Ejecuta scripts Python en Windows de manera confiable.')
    parser.add_argument('script', help='Ruta al script Python que deseas ejecutar')
    parser.add_argument('args', nargs='*', help='Argumentos para pasar al script')
    parser.add_argument('-v', '--virtual-env', help='Ruta al entorno virtual a utilizar')
    
    args = parser.parse_args()
    
    # Verifica si el script existe
    if not os.path.exists(args.script):
        print(f"Error: El script '{args.script}' no existe.")
        sys.exit(1)
    
    # Manejo de entorno virtual
    if args.virtual_env:
        if os.path.exists(args.virtual_env):
            # Activar entorno virtual en Windows
            activate_script = os.path.join(args.virtual_env, 'Scripts', 'activate.bat')
            if os.path.exists(activate_script):
                # Ejecutar en un nuevo proceso cmd con el entorno activado
                cmd = f'cmd /c "{activate_script} && python {args.script}'
                if args.args:
                    cmd += f' {" ".join(args.args)}"'
                else:
                    cmd += '"'
                    
                return subprocess.call(cmd, shell=True)
            else:
                print(f"Error: No se pudo encontrar el script de activación en {activate_script}")
                sys.exit(1)
        else:
            print(f"Error: El entorno virtual '{args.virtual_env}' no existe.")
            sys.exit(1)
    
    # Ejecución normal (sin entorno virtual)
    python_cmd = verificar_python()
    return ejecutar_script(python_cmd, args.script, args.args)

if __name__ == "__main__":
    sys.exit(main()) 