#!/bin/bash
echo "¿Seguro que quieres dropear las tablas?"

echo -n "y/n: "
read -r respuesta

y="y"

if [ $respuesta = $y ]; then
    echo "DROPEANDO..."
    echo "El pwd es user"
    mysql -h 127.0.0.1 -P 6033 -u user -p app_db < Drop_tables.sql        
    exit 1
fi

echo "No hay dropeo (su respuesta: $respuesta)"

