# Pasos para ejecutar el proyecto
##1. Crear el entorno virtual en python
```
python -m venv venv
#Para Windows
.\venv\Scripts\activate
pip install -r requirements.txt
```
##2. Crear nuestro modulo en go
```
go mod init project
```
##3. Ejecutar los archivos go en distintas terminales
```
go run calculator/cmd/main.go
go run api.go
```
