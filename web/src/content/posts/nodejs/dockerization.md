---
title: Dockerización en Node.js
date: 2025-12-14
keywords: [Docker, Node.js, Dockerfile, Multi-stage builds, npm ci, TypeScript, Vite, Next.js]
category: Node.js
---

Al igual que en **Go** u otros lenguajes, una dockerización con buenas prácticas en **Node.js** puede tener impactos muy positivos. Dockerizar aplicaciones en **Node.js** es un tanto diferente a como lo sería en **Go**, ya que podemos dockerizar aplicaciones que solo usen **JavaScript** (lo cual es muy fácil de hacer) o aplicaciones que necesiten **transpilarse**, como las desarrolladas con **TypeScript**. También podemos dockerizar aplicaciones frontend creadas con herramientas como **Vite**, **Next.js**, **Angular** o **Webpack**, donde es necesario un paso de construcción que puede beneficiarse del uso de **multi-stage builds**.

## ¿Qué imagen usar?
Así como en **Go**, haremos uso de versiones `alpine` por su gran ligereza. Las versiones de **Node.js** que vamos a usar dependerán de la versión en la que hayamos construido nuestras aplicaciones, como `node:22.10.0` u otras. Siempre debemos optar por versiones `alpine` para optimizar el tamaño de las imágenes.

A continuación, se presentan ejemplos de dockerización para aplicaciones con **TypeScript**, **Next.js** y **JavaScript**. En cualquier caso, utilizaremos las siguientes instrucciones como base para cualquier aplicación en **Node.js**:
```diff
+ FROM node:22-alpine3.21

+ WORKDIR /app
```

La versión utilizada será `node:22-alpine3.21`, ya que es la misma empleada para construir las aplicación de ejemplo. También definimos el entorno de trabajo (`WORKDIR`), que será el directorio raíz dentro del contenedor donde se almacenarán los archivos y recursos copiados o descargados.

## Copiando de dependencias por separado
Al igual que en **Go**, copiaremos los archivos de dependencias de **Node.js** (`package.json` y `package-lock.json`) de forma independiente:
```diff
FROM node:22-alpine3.21

WORKDIR /app

+ COPY package*.json .
```

>¿Por qué copiamos las dependencias por separado?
Como ya se hizo mención en multiples ocaciones al dockerizar una aplicación en **Go**. Copiar los archivos de dependencias de forma independiente permite aprovechar al máximo la caché de **Docker**. Si copiamos todos los archivos con `COPY . .` y modificamos cualquier archivo distinto de las dependencias (`package.json` o `package-lock.json`), **Docker** reinstalará todas las dependencias una y otra vez, lo que resulta ineficiente. Al copiar las dependencias por separado, este paso se almacena en una capa distinta y no se repetirá a menos que agreguemos nuevas dependencias.

A continuación, instalamos las dependencias:
```diff
FROM node:22-alpine3.21

WORKDIR /app

COPY package*.json .

+ RUN npm install
```
>¿Se puede optimar este paso?
Sí. En lugar de `npm install`, podemos usar `npm ci` para una instalación más rápida y limpia.

## Instalación optimizada
Utilizaremos el comando `npm ci` para instalar las dependencias de forma eficiente:
```diff
FROM node:22-alpine3.21

WORKDIR /app

COPY package*.json .

- RUN npm install
+ RUN npm ci
```

>¿Qué es `npm ci`? Es un comando de **Node.js** que instala dependencias de manera estricta a partir del archivo `package-lock.json`, garantizando que se usen versiones exactas de las dependencias. Además, limpia completamente el directorio `node_modules` antes de instalar las dependencias.

Luego, copiamos el resto del código de nuestra aplicación:
```diff
FROM node:22-alpine3.21

WORKDIR /app

COPY package*.json .

RUN npm ci

+ COPY . .
```

## Transpilando la aplicación
Si nuestra aplicación utiliza **TypeScript** o está construida con herramientas como **Vite** o **Next.js**, debemos incluir un paso adicional de construcción:
```diff
FROM node:22-alpine3.21

WORKDIR /app

COPY package*.json .

RUN npm ci

COPY . .

+ RUN npm run build
```

Esto generará un directorio, comúnmente llamado `dist/`, que contiene la aplicación transpilada a **JavaScript**. Posteriormente, podemos ejecutar la aplicación con `npm start` o `node dist/index.js`. Sin embargo, podemos optimizar aún más el proceso mediante **multi-stage builds**.

## Uso de Multi-Stage Builds
Ya mencionamos esto en la dockerizacion en **Go**. Los **multi-stage builds** nos permiten separar las etapas de construcción y ejecución, logrando imágenes más ligeras y optimizadas.

## Creando entorno de construcción
Para crear la primer etapa vamos a modificar la primera linea de nuestro `Dockerfile`, asignando un **alias** para que **Docker** entienda que será una etapa.
```diff
- FROM node:22-alpine3.21
+ FROM node:22-alpine3.21 AS build
```

## Creando entorno de ejecución
Una vez creada la primera etapa vamos a seguir con la segunda, que será la etapa de ejecución, donde copiaremos solo los archivos transpilados para ser ejecutados, elimando el resto de recursos de la primera etapa.

Usaremos nuevamenta la instruccion `FROM` junto con la imagen de **Node.js** con la misma versión para crear la segunda etapa.
```diff
FROM node:22-alpine3.21 AS build

WORKDIR /app

COPY package*.json .

RUN npm ci

COPY . .

RUN npm run build

+ FROM node:22-alpine3.21
```

Ahora, vamos a establecer nuevamente el entorno de trabajo y copiaremos el directorio `dist/` que contiene el código transpilado y `node_modules` que contiene todas las dependencias necesarias para la ejecución.
```diff
FROM node:22-alpine3.21 AS build

WORKDIR /app

COPY package*.json .

RUN npm ci

COPY . .

RUN npm run build

FROM node:22-alpine3.21

+ WORKDIR /app

+ COPY --from=build /app/dist .
+ COPY --from=build /app/node_modules node_modules
```
>Aplicamos `COPY --from` porque como lo explicamos en la dockerización en **Go**. `COPY --from` permite transferir recursos entre etapas, copiando únicamente el código transpilado y las dependencias necesarias para ejecutar la aplicación.

Hecho esto, solo quedaría dictar el `CMD` o `ENTRYPOINT`para ejecutar nuestra aplicación.
```diff
FROM node:22-alpine3.21 AS build

WORKDIR /app

COPY package*.json .

RUN npm ci

COPY . .

RUN npm run build

FROM node:22-alpine3.21

+ WORKDIR /app

COPY --from=build /app/dist .
COPY --from=build /app/node_modules node_modules

+ ENTRYPOINT [ "node", "index.js" ]
```

## ¿Y si solo uso JavaScript?
Si nuestra aplicación está escrita únicamente en **JavaScript**, no es necesario utilizar múltiples etapas ni incluir un paso de transpilación. En este caso, el proceso es mucho más simple.
Volveremos al pasado hasta las instrucciones de `RUN npm ci` y `COPY . .` y quitaremos el alias `build` y solo ejecutaremos la aplicación:

```diff
FROM node:22-alpine3.21

WORKDIR /app

COPY package*.json .

RUN npm ci

COPY . .

+ ENTRYPOINT [ "node", "index.js" ]
```
Hemos dockerizado nuestra aplicación en **Node.js** siguiendo las mejores prácticas, logrando una imagen ligera, rápida y eficiente gracias al uso adecuado de la caché de **Docker** y de los **multi-stage builds**.

# Comparación
Vamos a hacer una pequeña comparación de una dockerización bien optimizada versus una que no, y vamos a recordar el por qué debemos hacer cada cosa.

## Malas prácticas
```dockerfile
FROM node:lts-iron

WORKDIR /app
COPY . .
RUN npm install
RUN npm run build

CMD [ "node", "dist/index.js" ]
```
>En este ejemplo se muestra un archivo `Dockerfile` mal optimizado, no hace uso de [versiones ligeras como alpine](#qué-imagen-usar) o [copiado de archivos de dependencias por separado](#copiando-de-dependencias-por-separado), ya que hace uso de un solo `COPY . .` para copiar directamente los archivos de dependencias y el código fuente y la falta de uso múltiples etapas de construcción, lo que podría traer consecuencias negativas a la hora de la construcción como tiempos muy exagerados o imágenes muy pesadas.

## Buenas prácticas
```dockerfile
FROM node:22-alpine3.21 AS build

WORKDIR /app

COPY package*.json .

RUN npm ci

COPY . .

RUN npm run build

FROM node:22-alpine3.21

WORKDIR /app

COPY --from=build /app/dist .
COPY --from=build /app/node_modules node_modules

ENTRYPOINT [ "node", "index.js" ]
```
>Este es el ejemplo que realizamos utilizando las mejores prácticas de optimización y construcción, haciendo uso correcto de la caché de **Docker**, multiples etapas de construcción y versiones ligeras de imágenes, lo que se ve reflejado en el tamaño de imagen y tiempo de construcción.
