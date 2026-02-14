---
title: ¿Qué es HTTPS?
date: 2025-09-21
keywords: [linux, http, servidores]
category: Linux and Servers
---

**HTTPS** es una versión segura del protocolo **HTTP**.
Utiliza el protocolo [TLS](04_tls.md) (*Transport Layer Security*) para cifrar la comunicación, lo que garantiza una mayor seguridad al navegar.

El mecanismo de **HTTPS** funciona de manera similar a **HTTP**, pero al encriptar los datos transmitidos, protege la información contra posibles lecturas o manipulaciones por parte de terceros.
Además, este proporciona autenticación, lo que asegura  que el servidor con el que el cliente interactúra es legítimo.

La aunticación es posible gracias a que el servidor está certificado por una autoridad certificadora (**CA**), como **Let's Encypt**. Estos certificados validan la identidad del y servidor.

# ¿Cómo funciona?
1. **Cifrado inicial**:
	Cuando se establece una conexión HTTPS, el cliente y el servidor realizan un **intercambio de claves** utilizando el protocolo [TLS](04_tls.md).
	Este proceso asegura que ambas partes acuerden un método de cifrado, permitiendo que los datos transmitidos estén protegidos contra interceptaciones.
  ![https-encrypt](https://github.com/user-attachments/assets/ce7fd1d1-76bd-4e31-b216-53deba42ded7)

1. **Solicitud del cliente**:
	Después de establecer el cifrado, el cliente realiza una solicitud al servidor, como ocurre en **HTTP**.
  ![https-request](https://github.com/user-attachments/assets/e6382b32-684f-4fe4-9c12-310f2491c326)

	Sin embargo, en **HTTPS**, tanto la solicitud como la respuesta están cifradas, lo que garantiza que los datos sean **ilegibles** e **inmodificables** para cualquier atacante o entidad no autorizada.
  ![https-response](https://github.com/user-attachments/assets/f2fdefa7-99a2-457b-95be-c2040ed1492f)

Gracias a esta capa adicional de seguridad, **HTTPS** protege la privacidad del usuario y ofrece una experiencia de navegación más confiable.
Esto también puede mejorar el posicionamiento en motores de búsqueda, ya que es un factor valorado en [SEO](04_ceo).
