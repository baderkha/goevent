# goevent

<p align="center"><img width="512px"src="./assets/logo.png"></p>

goevent is an event emitter similar to javascript's custom events.

*Note* :

This is not a pattern you should have inside your code all over the place . You should use this package for specific use cases. Some of those use cases are listed below.


Use Case :
- Init / Shutdown events for your rest api
- Async Processing
- PUB->SUB Local events
- etc ...

## Contents
- [Installation](#installation)
- [Usage](#usage)
    - [Global](#global)
    - [Local](#local)

## Installation

```bash
go get github.com/baderkha/goevent
```

## Usage
This section will cover exemplar scenarios for this package. This package can be used to emit global events you can subscribe to anywhere in the code , or local events if you want to manage / orchestrate the events.

### Global 

For Global Events you first need to init the global emitter as such

1) Init
    ```go
    // somewhere in your main func
    func main() {
        // IF YOU DO NOT DO THIS , THIS WILL CRASH
        handlePanics := true // choose to have your panics handled gracefully , if you plan on having risky event listeners
        goevent.InitGlobal(handlePanics)
    }
    ```
2) Add Event Listener
    ```go
    // somewhere else
    listenerHash := goevent.
            Global().
            AddListener("slim_shady", func(d interface{}) {
                name := d.(string)
                fmt.Println("BAND NAME => " + name)
            })
    ```

3) Emit Event 
    ```go
    // somewhere else
    goevent.
            Global().
            Emit("slim_shady","d12")

    // with the event emitted above this will  output
    // "BAND NAME => d12"
    ```
4) Remove Event Listener
     ```go
    // somewhere else
    hasBeenRemoved := goevent.
            Global().
            RemoveListener(listenerHash)
    // will be false if it doesn't exist
    fmt.Println(hasBeenRemoved)
    ```

### Local

For Local events you can construct them from any where and attach them to structs

1) Init
    ``` go
    // from anywhere


    handlePanics := true // choose to have your panics handled 

    // this object is local ie everytime you call new you get a new one
    ev := goevent.New(handlePanics) 
    ```

2) Add Event Listener
    ``` go
    listenerHash:= ev.AddListener("slim_shady", func(d interface{}) {
                    name := d.(string)
                    fmt.Println("BAND NAME => " + name)
                })
    ```

3) Emit Event
    ``` go
    ev.Emit("slim_shady","d12")
    ```

4) Remove Event Listener

    ``` go
    ev.RemoveListener(listenerHash)
    ```
