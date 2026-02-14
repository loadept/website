---
title: Concurrencia en Go - Canales y Select en Acción
date: 2025-03-10
keywords: [Go, Concurrencia, Canales, Select]
category: Go-lang
---

Un canal es un tipo de dato en **Go**, que se usa para comunicar **goroutines**. Es un tipo de dato primitivo del lenguaje, al igual que `slice`, `array` o `map`.

Su propio nombre lo dice, es un canal o "tubería" por donde la información pasa de una **goroutine** a otra. Se pueden enviar o recibir valores de un tipo específico, porque los canales necesitan ser de un tipo de dato, ya sea `int`, `string`, etc.
![go-channel](https://github.com/user-attachments/assets/2546f49a-144f-4076-b4f8-a4626cf51b13)

## Introducción al uso de canales
Para trabajar con canales existe el nombre `chan`, que representa este tipo de dato.

## ¿Cómo se crean?
Al igual que las estructuras conocidas, como `map` o `slice`, para crear un canal usamos la función `make` y asignamos un tipo de dato.
```go
package main

import "fmt"

func main() {
	myChannel := make(chan int)

	fmt.Println(myChannel)
}
```
>¿Qué hará Println?
>De momento, **fmt.Println(myChannel)** solo imprimirá la dirección de memoria donde se encuentra el canal, algo similar a: `0xc000020060`

Acabamos de crear un canal (`myChannel`) de tipo entero (`int`). Este canal está diseñado para enviar y recibir datos de inmediato, ya que no tiene **buffer**.
## ¿Qué es un canal con o sin buffer?
### Canales sin buffer
Un canal sin buffer necesita que un receptor esté listo para recibir los datos que se envían, y viceversa. Si no hay un receptor o emisor disponible, la goroutine se bloqueará, impidiendo que el programa avance. En casos donde haya una goroutine (`main`) o más, que donde estén intentado enviar/recibir datos sin emisor/receptor disponible se producirá un `panic deadlock`. Veamos un ejemplo:
```go
package main

import "fmt"

func main() {
	myChannel := make(chan int)

	myChannel <- 5

	fmt.Println(myChannel)
}
```
El código anterior causará un error como este:
```plaintext
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan send]:
main.main()
	/home/user/Golang/channels.go:10 +0x2d
exit status 2
```
![go-no-receptor-channel](https://github.com/user-attachments/assets/3442da60-5e67-4462-be8b-240dccde1f02)
![arrow](https://github.com/user-attachments/assets/b38c86b1-9238-47aa-bf9f-3d166e931bb0)
![go-lock-channel](https://github.com/user-attachments/assets/ded68c16-94ef-4afd-9f39-c1aa12231a0f)
Por defecto, los canales se crean sin buffer. Esto significa que necesitan que el receptor y el emisor estén sincronizados.
### Canales con buffer
Si queremos que un canal pueda almacenar valores temporalmente sin bloquearse, podemos asignarle un buffer:
```go
package main

import "fmt"

func main() {
	myChannel := make(chan int, 4)

	myChannel <- 5

	fmt.Println(myChannel)
}
```
Aquí, usamos `make` para crear un canal con un buffer con capacidad para 4 valores. Ahora el canal puede guardar los datos temporalmente hasta que alguien los reciba.
![go-capacity-channel](https://github.com/user-attachments/assets/56e97802-ccca-4f70-b7d4-53bf89864246)
## ¿Cómo empezar a trabajar con canales?
Para trabajar con canales, usamos el operador `<-`. Este operador es intuitivo:
- Si el canal está **antes** del operador, como en `<-myChannel`, significa que estamos extrayendo datos del canal.
- Si el canal está **después**, como en `myChannel <- 3`, significa que estamos enviando datos al canal.

Ejemplo:
```go
package main

import "fmt"

func main() {
	myChannel := make(chan int, 4)

	go func() {
		// Enviar datos por el canal
		myChannel <- 3
	}()

	// Recibir datos enviados por el canal
	fmt.Println(<-myChannel)
}
```

## Canales unidireccionales
En **Go**, también podemos usar canales unidireccionales mediante el operador `<-`. Sin embargo, hay que tener cuidado, ya que es necesario comprender bien su uso para evitar confusiones.

Imagina que necesitas un canal unidireccional, ya sea de solo envío o solo recepción de datos. Aquí es importante aplicar criterio, ya que declarar directamente un canal como unidireccional puede ser problemático. Veamos un ejemplo incorrecto:
### Ejemplo incorrecto
* **Solo envío**
```go
package main

import "fmt"

func main() {
	myChannel := make(chan<- int)
}
```

* **Solo recepción**
```go
package main

import "fmt"

func main() {
	myChannel := make(<-chan int)
}
```

En estos casos, los canales son unidireccionales desde su declaración. Esto es inútil porque un canal necesita tanto un emisor como un receptor para funcionar. Si no hay uno de ellos, el programa no podrá utilizar el canal y se producirá un bloqueo o será simplemente inservible. Por lo tanto, esta forma de declarar canales unidireccionales no tiene utilidad.

### Forma Correcta de Usar Canales Unidireccionales
La manera adecuada de trabajar con canales unidireccionales es convertir canales bidireccionales a unidireccionales según sea necesario. Esto se logra utilizando funciones que especifiquen si el canal será de solo envío o solo recepción.
* **Solo envío**:
	Un canal de solo envío permite transmitir datos, pero no recibirlos:
```go
package main

import "fmt"

func sendMessage(ch chan<- string) {
	message := "Hello World"
	ch <- message
}

func main() {
	myChannel := make(chan string)

	// Convertimos el canal a uno de solo envío en la función
	go sendMessage(myChannel)

	// Recibimos el dato desde el canal
	fmt.Println(<-myChannel)
	// > Hello World
}
```

* **Solo recepción**
	Un canal de solo recepción permite recibir datos, pero no enviarlos:
```go
package main

import "fmt"

func sendMessage(ch chan<- string) {
	message := "Hello World"
	ch <- message
}

func getMessage(ch <-chan string) {
	// Recibimos el dato del canal y lo imprimimos
	fmt.Println(<-ch)
}

func main() {
	myChannel := make(chan string)

	go sendMessage(myChannel)

	// Convertimos el canal a uno de solo recepción en la función
	getMessage(myChannel)
	// > Hello World
}
```
>Conversión de canales
>Hay que tener en cuenta que podemos convertir canales bidireccionales a unidirrecionales, pero no podemos hacerlo al revés, convertir canales unidireccionales a bidireccionales no es posible.

>Seguridad con canales unidireccionales
>Usar canales unidireccionales aumenta la seguridad y la claridad del código, ya que especifica exactamente cómo debe utilizarse el canal en cada función.

>Dato importante sobre el operador `<-`
>Cuando usamos `<-myChannel`, siempre extraemos datos del canal. Esto ocurre incluso si no los almacenamos en una variable:
>```go
>func main() {
>	ch := make(chan int)
>	go func() {
>		ch <- 5
>	}()
>	<-ch
>}
>```
Aquí se extrae los valores del canal, no importa si se almacena en una variable o no.

## Cierre de canales (close)
Una parte importante del uso de canales es el cierre, el cual se realiza con la función `close(channel)`.  
Al cerrar un canal, hacemos que los emisores dejen de enviar datos. Además, notificamos a los receptores que ya no se enviarán más datos por ese canal.

Ejemplo:
```go
package main

import "fmt"

func sendMessage(ch chan<- string) {
	message := "Hello World"
	ch <- message
	close(ch)
}

func getMessage(ch <-chan string) {
	fmt.Println(<-ch)
}

func main() {
	myChannel := make(chan string)

	go sendMessage(myChannel)

	getMessage(myChannel)
	// > Hello World
}
```

**¿Qué sucede si un receptor intenta extraer datos de un canal cerrado?**
En este caso, el receptor obtendrá un valor nulo dependiendo del tipo de dato del canal. Por ejemplo:
- Si es `int`, obtendrá `0`.
- Si es `string`, obtendrá `""`.
- Si es `bool`, obtendrá `false`.

Ejemplo:
```go
package main

import "fmt"

func sendNumber(ch chan<- int) {
	number := 5
	ch <- number
	close(ch)
}

func getNumber(ch <-chan int) {
	for i := 0; i < 10; i++ {
		fmt.Println(<-ch)
	}
}

func main() {
	myChannel := make(chan int)

	go sendNumber(myChannel)

	getNumber(myChannel)
	// > 5
	// > 0
	// > 0
	// ...
}
```

**¿Qué sucede si un emisor intenta enviar datos a un canal cerrado?**
Intentar enviar datos a un canal cerrado genera un `panic: send on closed channel`. Esto provoca que el programa falle de forma abrupta.
```go
package main

import "fmt"

func main() {
	myChannel := make(chan int, 4)
	close(myChannel)

	myChannel <- 5

	fmt.Println(myChannel)
}
```
El código anterior generará el siguiente error:
```plaintext
panic: send on closed channel

goroutine 1 [running]:
main.main()
	/home/user/Golang/channels.go:9 +0x3a
exit status 2
```
## Gestión de errores  e iteraciones
En la sección anterior [(Cierre de canales (close))](#cierre-de-canales-close), vimos cómo detectar si un canal está cerrado o no. Esto se puede lograr mediante un segundo parámetro, denominado `ok`, que se obtiene al extraer datos del canal. Este parámetro nos indica si el canal sigue abierto.

Ejemplo:
```go
package main

import "fmt"

func sendMessage(ch chan<- string) {
	message := "Hello World"
	ch <- message
}

func getMessage(ch <-chan string) {
	for {
		message, ok := <-ch
		if !ok {
			break
		}
		fmt.Println(message)
	}
}

func main() {
	myChannel := make(chan string)

	go sendMessage(myChannel)

	getMessage(myChannel)
	// > Hello World
}
```
En este caso, `ok` devuelve un valor booleano:
- `true`: si el canal está abierto.
- `false`: si el canal está cerrado.
### Uso de bucle `for`
Además de verificar si un canal está abierto, podemos iterar directamente sobre él. Esto es útil para enviar y recibir múltiples datos. Veamos un ejemplo donde usamos un bucle para enviar y recibir datos.

Ejemplo con `ok`:
```go
package main

import (
	"fmt"
	"time"
)

func sendNumbers(ch chan<- int) {
	for i := 1; i <= 10; i++ {
		ch <- i
		time.Sleep(time.Second)
	}
}

func getNumbers(ch <-chan int) {
	for {
		if num, ok := <-ch; ok {
			fmt.Println(num)		
		} else {
			fmt.Println("closed channel")
			break
		}
	}
}

func main() {
	myChannel := make(chan int)

	go sendNumbers(myChannel)

	getNumbers(myChannel)
	// > 1
	// > 2
	// > 3
	// ...
}
```
Este es un ejemplo muy claro de la utilidad de los canales, pudimos enviar muchos datos cuanto sean posibles haciendo uso de un bucle y también pudimos recibirlos conforme van llegando usando otro bucle.

Dijimos que podemos iterar sobre los canales, pero en el ejemplo anterior no lo hicimos. ¿Es mentira acaso? No, en realidad podemos iterar sobre canales.

En **Go** existe la palabra reservada `range`, que se usa mucho para iterar sobre elementos, ya sea un `map` o un `array`. Pero también podemos usarla para iterar sobre canales. De hecho, iterar canales con ayuda de `range` es mucho más fácil. Veamos el código:
```go
package main

import (
	"fmt"
	"time"
)

func sendNumbers(ch chan<- int) {
	for i := 1; i <= 10; i++ {
		ch <- i
		time.Sleep(time.Second)
	}
}

func getNumbers(ch <-chan int) {
	for num := range ch {
		fmt.Println(num)
	}
}

func main() {
	myChannel := make(chan int)

	go sendNumbers(myChannel)

	getNumbers(myChannel)
	// > 1
	// > 2
	// > 3
	// ...
}
```
Como te podrás dar cuenta, la función es mucho más corta y ya no se hace uso de un `if` para verificar si el canal está abierto o no. Esto se debe a que, al usar `for` y `range` para la iteración, el mecanismo interno de `range` se encarga de evaluar si el canal está cerrado. Si un canal se cierra y deja de tener valores disponibles, `range` rompe automáticamente el ciclo y deja de iterar sobre el canal. Es como hacer un `break` dentro de un `if`, pero de una manera más sencilla.
### Uso de `select`
El uso de `select` también es importante al manejar canales. Una gran ventaja de `select` es que podemos escuchar múltiples canales. Pero no solo eso, `select` también nos permite enviar datos a varios canales. Veamos cómo se utiliza:
* **Recibir datos**
```go
package main

import (
	"fmt"
	"time"
)

func sendNumber(ch chan<- int) {
	ch <- 10
}

func sendMessage(ch chan<- string) {
	ch <- "Hello"
}

func getValue(ch <-chan int, ch2 <-chan string) {
	for {
		select {
		case num := <-ch:
			fmt.Println("a number was obtained:", num)
			return
		case message, ok := <-ch2:
			if !ok {
				return
			}
			fmt.Println("a message was obtained:", num)
			return
		default:
			fmt.Println("waiting for some signal on the channels")
			time.Sleep(time.Second)
		}
	}
}

func main() {
	myChannel := make(chan int)
	myOtherChannel := make(chan string)

	go sendNumbers(myChannel)
	go sendMessage(myOtherChannel)

	getNumbers(myChannel, myOtherChannel)
	// > a number was obtained: 10
	// > a message was obtained: Hello
}
```

* **Enviar datos**
```go
package main

import (
	"fmt"
	"time"
)  

func getValues(ch chan int, ch2 chan string) {
	fmt.Println(<-ch)
	fmt.Println(<-ch2)

}

func sendValues(ch chan<- int, ch2 chan<- string) {
	for i := 0; i < 2; i++ {
		select {
		case ch <- 2:
			fmt.Println("the number was sent successfully")	
		case ch2 <- "Hello":
			fmt.Println("the message was sent successfully")	
		default:
			fmt.Println("the channels are closed or full")
		}
		time.Sleep(time.Second)
	}
}

func main() {
	myChannel := make(chan int, 1)
	myOtherChannel := make(chan string, 1)

	go getValues(myChannel, myOtherChannel)

	sendValues(myChannel, myOtherChannel)
	// the message was sent successfully
	// the number was sent successfully
	// 2
	// Hello
}
```

En estos ejemplos vimos cómo, con `select`, pudimos enviar o recibir datos a múltiples canales. `Select` es una gran herramienta porque nos ofrece mucha más flexibilidad y control. Además, no solo podemos usarlo con un `for`; también es posible utilizar `select` de forma independiente, como en este caso:
```go
func sendValues(ch chan<- int, ch2 <-chan string) {
	select {
	case ch <- 2:
		fmt.Println("the number was sent successfully")	
	case message := <-ch2:
		fmt.Println("the message was get successfully:", message)	
	default:
		fmt.Println("the channels are closed or full")
	}
}
```
Otra ventaja es que también podemos combinar `select` con las distintas estructuras de control disponibles en **Go**, como `if`, `switch`, etc. Esto aumenta aún más su versatilidad al manejar canales.
>Espera de bajo consumo
>Cuando se usa `select` junto con un `for loop` infinito, el uso de `select` puede llegar a ser muy eficiente, esto porque cuando `select` está esperando que algún canal tenga datos, el runtime de >**Go** pone la goroutine en un estado de espera (**waiting**). Lo cual hace que que la goroutine no consuma ciclos de CPU mientras está esperando, esto hace que su uso sea muy eficiente.
