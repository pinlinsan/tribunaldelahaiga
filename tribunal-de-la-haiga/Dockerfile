# Usa una imagen base de Go
FROM golang:1.20

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia los archivos del proyecto al contenedor
COPY . .

# Descarga las dependencias
RUN go mod download

# Expone el puerto 8080
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["go", "run", "cmd/server/main.go"]
