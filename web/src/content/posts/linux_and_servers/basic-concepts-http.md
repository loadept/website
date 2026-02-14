---
title: ¿Qué es HTTP?
date: 2025-09-21
keywords: [linux, http, servidores]
category: Linux and Servers
---

**HTTP** (*Hyper Text Transfer Protocol*) es un protocolo que se utiliza para la comunicación en la web entre dos partes: un **cliente** (como un navegador) y un **servidor**.
Es el fundamento de cómo se cargan y transfieren recursos como información texto, imágenes, videos y otro tipo de datos en internet.

# ¿Cómo funciona?
Cuando un cliente, como un navegador, quiere acceder a un recurso (por ejemplo, una página web, una imagen, o un archivo PDF), envía una solicitud **HTTP** al servidor correspondiente.

1. **Solicitud del cliente**:
	Esto ocurre cuando hacemos clic en un enlace, escribimos una URL en el navegador, o realizamos una búsqueda.
	La solicitud incluye información como el tipo de recurso solicitado y su ubicación.
  ![http-request](https://github.com/user-attachments/assets/d6646baa-6c7e-4d35-b358-6683decd28e8)

1. **Respuesta del servidor**:
	El servidor procesa la solicitud y devuelve una **respuesta HTTP** con el contenido solicitado por medio de una conexión **TCP**.
  ![processing](https://github.com/user-attachments/assets/2efa5fab-3e18-44bd-87e1-4b7509d442eb)

  Esta respuesta puede incluir el recurso deseado (como el **HTML** de una página web, una imagen, etc.). 
	O de lo contrario un mensaje de error si el recurso no está disponible (como un, "404 Not Found").
	![http-response](https://github.com/user-attachments/assets/11ac4286-d427-43fb-a33d-1a684da3adac)


De esta manera, el navegador utiliza el contenido recibido para mostrar la página web, el texto, o la imagen en la pantalla del usuario.
![http-content](https://github.com/user-attachments/assets/4e640421-f0a1-4183-84af-7da7c041e137)
