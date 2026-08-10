package utilidades

import (
	"encoding/json"
	"fmt"
	"os"
)

// AsegurarDirectorio verifica si la carpeta indicada existe. Si no existe, la crea con permisos estándar (0755).
func AsegurarDirectorio(rutaCarpeta string) error {
	if _, err := os.Stat(rutaCarpeta); os.IsNotExist(err) {
		err := os.MkdirAll(rutaCarpeta, 0755)
		if err != nil {
			return fmt.Errorf("no se pudo crear la carpeta '%s': %v", rutaCarpeta, err)
		}
	}
	return nil
}

// ExisteArchivo comprueba si un archivo existe en el sistema de archivos local.
func ExisteArchivo(ruta string) bool {
	info, err := os.Stat(ruta)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// GuardarJSON convierte cualquier estructura o slice de Go a formato JSON indentado
// y lo guarda en la ruta de archivo especificada.
func GuardarJSON(ruta string, datos interface{}) error {
	// Asegurar que el directorio contenedor exista antes de guardar
	err := AsegurarDirectorio("datos")
	if err != nil {
		return err
	}

	// Convertir la estructura de datos a bytes en formato JSON legible con sangría
	contenido, err := json.MarshalIndent(datos, "", "  ")
	if err != nil {
		return fmt.Errorf("error al convertir datos a JSON: %v", err)
	}

	// Escribir los bytes en el archivo con permisos de lectura y escritura (0644)
	err = os.WriteFile(ruta, contenido, 0644)
	if err != nil {
		return fmt.Errorf("error al escribir el archivo '%s': %v", ruta, err)
	}

	return nil
}

// LeerJSON lee el contenido de un archivo JSON y lo convierte a la estructura destino pasada por puntero.
// Si el archivo no existe o está vacío, crea un archivo por defecto con un array vacío "[]".
func LeerJSON(ruta string, destino interface{}) error {
	// Si el archivo no existe o está vacío, inicializamos un archivo JSON base con "[]"
	if !ExisteArchivo(ruta) {
		err := GuardarJSON(ruta, json.RawMessage("[]"))
		if err != nil {
			return err
		}
	}

	contenido, err := os.ReadFile(ruta)
	if err != nil {
		return fmt.Errorf("error al leer el archivo '%s': %v", ruta, err)
	}

	// Si el archivo tiene 0 bytes, guardar array vacío por defecto
	if len(contenido) == 0 {
		contenido = []byte("[]")
		_ = GuardarJSON(ruta, json.RawMessage("[]"))
	}

	// Convertir el JSON a la estructura de Go correspondiente
	err = json.Unmarshal(contenido, destino)
	if err != nil {
		return fmt.Errorf("error al deserializar JSON desde '%s': %v", ruta, err)
	}

	return nil
}
