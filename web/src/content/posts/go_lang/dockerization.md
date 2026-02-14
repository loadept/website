---
title: Dockerización en Go
date: 2025-06-09
keywords: [binario, go, golang, programacion, development, docker]
category: Go-lang
---

Conocer cómo dockerizar de manera correcta en **Go** es muy importante, ya que puede tener un impacto positivo en el rendimiento, el consumo de recursos y el almacenamiento. Entender aspectos como **multi-stage builds** o el uso de versiones ligeras de las imágenes puede ser muy útil.

## ¿Qué imagen usar?
Para la dockerización con buenas prácticas en **Go**, siempre es recomendable usar versiones `alpine` de las imágenes de **golang**. Dependiendo de la versión de **Go** que estemos utilizando, podemos elegir imágenes como `golang:1.23.4-alpine3.21`. Las versiones dependen de nuestra necesidad: si queremos una versión específica, una versión flexible con respecto al parche, etc.

A continuación, dockerizaremos una pequeña API usando una versión flexible de Alpine.
```diff
+ FROM golang:1.23-alpine3.21

+ RUN apk add --no-cache git

+ WORKDIR /app
```
La versión actual con la que se creó la API es **1.23.0**, pero no hay una versión específica del parche disponible. Por lo tanto, usaremos la imagen sin especificar el parche. También descargamos `git` usando `apk` (*Gestor de paquetes en alpine*) sin almacenar en caché para aligerar la imagen, por ultimo establecemos el entorno de trabajo (`WORKDIR`), que será el espacio o `$HOME` dentro del contenedor, donde se almacenaran los archivos o recursos copiados o descargados.

## Copiando de dependencias por separado
A continuación, copiamos los archivos de dependencias de **Go** (`go.mod` y `go.sum`) y descargamos las dependencias en la imagen:
```diff
FROM golang:1.23-alpine3.21

RUN apk add --no-cache git

WORKDIR /app

+ COPY go.mod go.sum ./
```

>**¿Por qué copiamos `go.mod` y `go.sum` de forma separada?** 
>Esto lo hacemos para aprovechar la caché de **Docker**. Al copiar estos archivos por separado, si se modifica algún archivo de código excepto las dependencias, **Docker** no volverá a descargarlas. Esto acelera la construcción de la imagen, ya que **Docker** utiliza un sistema de capas y cada instrucción en el `Dockerfile` corresponde a una capa independiente, si no se modifican estos archivos **Docker** podrá reutilizar la capa sin necesidad de descargar nuevamente todas las dependencias.
>Y, si hacemos un `COPY . .` para copiar todo el código fuente, los archivos de dependencias e instalarlas en una sola capa, cualquier cambio en los archivos reiniciará el proceso completo, lo que incluye volver a descargar las dependencias. Al separar estas acciones en diferentes capas, evitamos este problema.

Ahora descargamos las dependencias:
```diff
FROM golang:1.23-alpine3.21

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

+ RUN go mod download
```

>**¿Por qué descargamos dependencias?**
>Esto lo hacemos para que las dependencias estén en la caché del sistema, así al compilar, **Go** se saltará el proceso de descargarlo por si solo, porque internamente este, descarga las dependencias al compilar por lo que puede tardar un tiempo, así que nosotros lo hacemos manualmente en una capa diferente para ahorrar tiempo y aprovechar la caché de **Docker**. 

Después, copiamos el resto del código fuente al área de trabajo:
```diff
FROM golang:1.23-alpine3.21

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

+ COPY . .
```
>Con este paso, garantizamos que cualquier cambio en el código fuente no afecte las capas anteriores.

## Compilar el código
En este paso, convertimos todo el código en un binario ejecutable:
```diff
FROM golang:1.23-alpine3.21

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

+ RUN go build cmd/app/main.go
```
Este comando genera un binario dentro del directorio de trabajo (`/app`). El archivo resultante será `/app/main`.

Podríamos finalizar aquí la dockerización, agregando un `CMD` o `ENTRYPOINT` para ejecutar el binario, pero nosotros queremos exprimir al máximo la optimización, así que vamos a hacer uso de **multi-stage builds**

## Uso de Multi-Stage Builds.
Los **multi-stage builds** permiten tener múltiples etapas dentro de un `Dockerfile`. Esto ayuda a reducir el tamaño final de la imagen, ya que los archivos temporales o innecesarios de la etapa de construcción serán "limpiados" por lo que no se incluyen en la imagen final.
Resultando así en un contenedor más pequeño y con un menor consumo.

## Creando el entorno de construcción
Para hacer uso de **multi-stage builds** necesitamos primero cambiar la primera línea, haremos uso de **aliases** para nombrar nuestras etapas, por ejemplo vamos a nombrar nuestra etapa actual como `build`, porque es el área de descarga y compilación.
```diff
- FROM golang:1.23-alpine3.21
+ FROM golang:1.23-alpine3.21 AS build
```

## Creando el entorno de ejecución
Cada etapa necesita una nueva imagen, en nuestra etapa de construcción hicimos uso de la imagen de `golang` en la que descargamos dependencias y compilamos. Ahora, nuestra segunda etapa, necesita de una imagen también, podríamos hacer uso de la imagen de `golang` nuevamente, pero, hay que tener en cuenta que esta imagen puede ser pesada, ya que trae todo lo necesario para compilar o ejecutar nuestro código `golang`, trae el compilador, formateador, caché o gestor de dependencias, lo que hace que esta imagen sea pesada.

Ya que actualmente solo necesitamos de un pequeño entorno donde ejecutar nuestro binario, no necesitamos mucho, así que haremos uso de una imagen pequeña de `alpine` en su versión `3.21` que es la misma versión `alpine` correspondiente a la usada con `golang`.
```diff
FROM golang:1.23-alpine3.21 AS build

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build cmd/app/main.go

+ FROM alpine:3.21
```

Listo, así creamos nuestra etapa, ahora vamos a repetir ligeramente los pasos, como creación de un espacio de trabajo, y vamos a copiar nuestros archivos, pero lo haremos de diferente manera.
```diff
FROM golang:1.23-alpine3.21 AS build

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build cmd/app/main.go

FROM alpine:3.21

+ WORKDIR /app

+ COPY --from=build /app/main .
```

Aplicamos un `COPY --from` el cual nos permite copiar recursos entre diferentes etapas, en este caso desde la etapa `build`, donde se copia el binario que se encuentra en espacio de trabajo `/app` (`/app/main`) el cual es el binario resultante después de la compilación.

Hecho esto, solo quedaría dictar el `CMD` o `ENTRYPOINT`para ejecutar nuestro binario.
```diff
FROM golang:1.23-alpine3.21 AS build

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build cmd/app/main.go

FROM alpine:3.21

WORKDIR /app

COPY --from=build /app/main .

+ ENTRYPOINT [ "./main" ]
```

Finalmente logramos dockerizar nuestra aplicación en **Go**, es una dockerización que sigue las mejores prácticas, ligera y muy rápida, haciendo uso correcto de la caché de **Docker** y múltiples etapas.

# Comparación
Vamos a hacer una pequeña comparación de una dockerización bien optimizada versus una que no, y vamos a recordar el por qué debemos hacer cada cosa.

## Malas prácticas
```dockerfile
FROM golang:1.23

WORKDIR /app

COPY . .

RUN go build cmd/app/main.go

ENTRYPOINT [ "./main" ]
```

### Resultado
![Pasted image 20250108015147](https://github.com/user-attachments/assets/40081688-fe47-460f-8f6b-abbf3f103615)

>En este ejemplo se muestra un archivo `Dockerfile` mal optimizado, no hace uso de [versiones ligeras como alpine](#qué-imagen-usar) o [copiado de archivos de dependencias por separado](#copiando-de-dependencias-por-separado), ya que hace uso de un solo `COPY . .` para copiar directamente los archivos de dependencias y el código fuente y la falta de uso múltiples etapas de construcción, lo que podría traer consecuencias negativas a la hora de la construcción como tiempos muy exagerados o imágenes muy pesadas, como en este caso, donde la imagen llega a pesar **1.08GB**.

## Buenas prácticas
```dockerfile
FROM golang:1.23-alpine3.21 AS build

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build cmd/app/main.go

FROM alpine:3.21

WORKDIR /app

COPY --from=build /app/main .

ENTRYPOINT [ "./main" ]
```

### Resultado
![Pasted image 20250108015742](https://github.com/user-attachments/assets/5f7cfc33-4719-4057-9f50-2d3e832e321d)

>Este es el ejemplo que realizamos utilizando las mejores prácticas de optimización y construcción, haciendo uso correcto de la caché de **Docker**, multiples etapas de construcción y versiones ligeras de imágenes, lo que se ve reflejado en el tamaño de imagen y tiempo de construcción.

# Consideraciones finales
Algunas dependencias pueden hacer uso de una característica de **Go** que es `CGO` que permite ejecutar módulos en **C** desde nuestras apps en **Go**. Para poder hacer uso de esta característica basta con habilitar la variable `CGO_ENABLED=1` a la hora de compilar nuestras aplicaciones. No hacer esto podría traer fallas en la ejecución o compilación de nuestros programas.
```diff
FROM golang:1.23-alpine3.21 AS build

+ ENV CGO_ENABLED=1

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build cmd/app/main.go

FROM alpine:3.21

WORKDIR /app

COPY --from=build /app/main .

ENTRYPOINT [ "./main" ]
```
Hecho esto podremos hacer uso de los módulos en **C** para **Go** como `librdkafka` para trabajar con **Apache Kafka** u otros módulos.
