---
title: Command Pattern en Go
date: 2025-03-12
keywords: [design-pattern, command-pattern, go, golang, development, diseño, programación, patrón]
category: Go-lang
---

# Definición
El patrón de diseño **Command** es un patrón de **comportamiento** que encapsula una solicitud como un objeto, lo que permite desacoplar el solicitante de la ejecución de la acción. Este patrón facilita que las acciones sean ejecutadas, deshechas o programadas, haciendo más flexible y extensible la implementación de comandos.

El patrón **Command** es una solución efectiva para encapsular funcionalidades específicas que se ejecutan según determinadas condiciones. En lugar de depender de estructuras como **if-else** o **switch-case**, este patrón es ideal cuando se busca extender o modificar la funcionalidad de un sistema sin alterar su código existente.

En **Go**, podemos implementar este patrón usando **interfaces** para definir los comandos y las estructuras correspondientes, de la siguiente manera:

### Implementación paso a paso.
1. Definimos una **interfaz** con los métodos `Execute` y `Undo`, que no implementamos directamente. Por convención, estos son los métodos comunes en un comando.
```go
type Command interface {
	Execute()
	Undo()
}
```

2. Creamos los **receptores** (**Receiver**) que contienen la lógica que se ejecutará cuando se invoquen los comandos.
```go
import "fmt"

type Car struct{}

func (c *Car) Advance() {
	fmt.Println("The car is moving forward")
}

func (c *Car) Stop() {
	fmt.Println("The car is stopped")
}
```

```go
import "fmt"

type Car2 struct{}

func (c *Car2) Reload() {
	fmt.Println("The car is reloading")
}

func (c *Car2) RetakePath() {
	fmt.Println("The car is retaking path")
}

```

3.  Creamos los **concretadores** (**ConcreteCommand**) que implementan los métodos de la interfaz `Command` y llaman a los métodos apropiados de los receptores.
```go
type CarCommand struct {
	car *Car
}

func (o *CarCommand) Execute() {
	o.car.Advance()
}

func (o *CarCommand) Undo() {
	o.car.Stop()
}
```

```go
type CarCommand2 struct {
	car *Car2
}

func (o *CarCommand2) Execute() {
	o.car.Reload()
}

func (o *CarCommand2) Undo() {
	o.car.RetakePath()
}
```

4. Finalmente, creamos el **invocador** (**Invoker**) que es quien ejecuta o deshace las acciones de los comandos.
```go
type Invoker struct{
	command Command
}

func (i *Invoker) SetCommand(command Command) {
	i.command = command
}

func (i *Invoker) Start() {
	i.command.Execute()
}

func (i *Invoker) Stop() {
	i.command.Undo()
}
```


<span id="implementación-invocación" style="color: transparent">implementacion-invocacion</span>

5. Inicializamos los componentes en el `main` y demostramos la invocación de los comandos.
```go
func main() {
	car1 := &Car{}
	car2 := &Car2{}

	advanceCar := &CarCommand{car: car1}
	reloadCar := &CarCommand{car: car2}

	invoker := &Invoker{}
	invoker.SetCommand(advanceCar)
	invoker.Start()
	invoker.Stop()

	invoker.SetCommand(advanceCar)
	invoker.Start()
	invoker.Stop()
}
```
>Este ejemplo es una implementación básica del patrón **command**, no es el único, ni el definitivo

# Una propuesta más cómoda.
La implementación anterior [(Implementación paso a paso)](#implementación-paso-a-paso) proporciona una base sólida para comprender el patrón Command, pero presenta algunas limitaciones cuando los comandos y receptores crecen. Si tenemos más comandos, como se muestra en la implementación, el invocador está restringido a aceptar un solo comando a la vez. **¿Qué sucede si necesitamos ejecutar múltiples comandos de forma flexible y eficiente?**

Si observamos bien el bloque de [Implementación e invocación (líneas 9-13)](#implementación-invocación), al usar `SetCommand()`, estamos modificando el invocador con cada comando. Esto puede no ser tan flexible, ya que nos obliga a establecer solo un comando a la vez. **¿Qué pasa si necesitamos configurar varios comandos?** Podríamos recurrir a un `if` o a una lógica condicional para asignar diferentes comandos según ciertas condiciones, pero… **¿No era precisamente la intención del patrón Command evitar ese tipo de estructuras condicionales?**

En la propuesta siguiente, abordaremos cómo solucionar este problema utilizando un enfoque más flexible, que nos permita manejar múltiples comandos de manera eficiente, sin necesidad de modificar el invocador constantemente o depender de condicionales.

## Uso de Map `map`.
Para mejorar la flexibilidad del invocador, utilizaremos un `map` en lugar de un solo comando. Este mapa nos permitirá asociar cada comando con una clave única, de modo que podamos registrar varios comandos y ejecutarlos dinámicamente sin recurrir a `Start()` como en el ejemplo anterior de [Implementación e invocación (línea 10)](#implementación-invocación).

### Implementación paso a paso (con `map`).
1. Definimos una **interfaz** para los comandos, al igual que en el ejemplo anterior, con los métodos `Execute` y `Undo`.
```go
type PersonCommand interface {
	Execute()
	Undo()
}
```

2. Creamos los **receptores** (acciones) que se ejecutarán cuando se invoquen los comandos, de forma similar a los ejemplos anteriores.
```go
import "fmt"

type Drink struct{}

func (d *Drink) StartDriking() {
	fmt.Println("Start drinking")
}

func (d *Drink) StopDriking() {
	fmt.Println("Stop drinking")
}
```

```go
import "fmt"

type Eat struct{}

func (d *Eat) StartEating() {
	fmt.Println("Start eating")
}

func (d *Eat) StopEating() {
	fmt.Println("Stop eating")
}
```

3. Creamos los **concretadores** de comandos, implementando los métodos de la interfaz `PersonCommand` y llamando a los métodos de los receptores.
```go
type DrinkCommand struct {
	drink *Drink
}

func (d *DrinkCommand) Execute() {
	d.drink.StartDriking()
}

func (d *DrinkCommand) Undo() {
	d.drink.StopDriking()
}
```

```go
type EatCommand struct {
	eat *Eat
}

func (e *EatCommand) Execute() {
	e.eat.StartEating()
}

func (e *EatCommand) Undo() {
	e.eat.StopEating()
}
```

4. Ahora, modificamos el **invocador** para manejar múltiples comandos utilizando un `map` que asocie comandos a claves.
```go
type Invoker struct{
	commands map[string]PersonCommand
}

func (i *Invoker) SetCommand(action string, command Command) {
	if i.commands == nil {
		i.commands = make(map[string]Command)
	}
	i.commands[action] = command
}

func (i *Invoker) StartCommand(action string) {
	if cmd, exists := i.commands[action]; exists {
		cmd.Excute()
	} else {
		fmt.Println("Unknown command", action)
	}
}

func (i *Invoker) StopCommand() {
	if cmd, exists := i.commands[action]; exists {
		cmd.Undo()
	} else {
		fmt.Println("Unknown command", action)
	}
}
```

5. Finalmente, inicializamos los receptores e invocador, y demostramos cómo se pueden ejecutar y deshacer múltiples comandos con un solo invocador.
```go
func main() {
	action := "drink" // eat

	eat := &Eat{}
	drink := &Drink{}

	eatCommand := &EatCommand{eat: eat}
	drinkCommand := &DrinkCommand{drink: drink}

	invoker := &Invoker{}
	invoker.SetCommand("eat", eatCommand)
	invoker.SetCommand("drink", drinkCommand)

	invoker.StartCommand(action)
	invoker.StopCommand(action)
}
```

>Con este enfoque, el invocador se convierte en una herramienta mucho más flexible y escalable,
>permitiéndonos manejar múltiples comandos sin depender de condicionales ni de modificaciones constantes del invocador.
