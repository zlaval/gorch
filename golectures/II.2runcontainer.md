# Context

The context package is used to efficiently manage long-running
or interruptible operations. With context, we can pass timing,
deadlines, and cancellation signals between functions.
This is especially useful in scenarios like network operations or
concurrent programs. For example, commands sent to Docker engine
can be interrupted if the service takes too long to
respond (e.g., due to an error), preventing the thread from waiting indefinitely.
In GO, context is usually the first input parameter of the functions.
We are going to use context to stop background jobs
and for docker method invocations
as it is a required param there.

Creating new context is easy, here is some example:

```go
ctx := context.Background()
ctx, cancel := context.WithCancel(context.Background())

context.TODO()
context.WithTimeout(context.Background(), 1 * time.Seconds)
```

# Defer keyword

The defer keyword is used to delay the execution of a function
or operation until the surrounding function completes.
Functions marked with defer are typically used for resource cleanup
(e.g., closing files, releasing locks).
It is similar to the try-finally construct in other programming languages,
so it runs even the method return early.

```
func MyFunc() int {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    //do some operation
    
    //if condition then early return
    if cond == true{
        return 2
    }
    
    //do some other operation
    
    return 1
}
//the cancel() function executes regardless of 
//which branch the function returns on.
```

# Error handling

Error handling in GO is straightforward and emphasizes clarity.
Unlike exceptions in other languages, Go uses explicit error values
that functions return alongside regular results.
This approach requires developers to check for and handle errors directly.
While this approach may seem unusual at first, it has many advantages.

```
func ReadLineFromFile(path string) (string,error){
    return ...
}

line, err := ReadLineFromFile("path/to/file")
if err != nil{
    return err
}
//no error line can be processed

-----or-----

func sendMessage(msg Message) error{
    return ...
}

if err:=sendMessage(msg); err != nil{
    return err
}
//no error, continue execution
```

# Struct

Struct is typed collection of fields. We will discuss it in detail later.
For now, you only need to know that it can be declared inline as `struct{}`.
If no fields are defined, it is an empty struct.
We can create an instance using the `{}` notation.
So the expression `struct{}{}` creates an instance of an empty struct,
which can be used to represent empty values.
```
nat.PortSet{
    nat.Port("8000"): struct{}{},
},
```
# Channels basic
A channel is a feature used to communicate between goroutines (lightweight threads).
Channels allow you to send and receive values, 
enabling safe and synchronized data exchange in concurrent programs.
They can be thought of as pipes that connect different goroutines, 
helping to coordinate and synchronize their execution. 
We will use channels to implement different jobs, and we will discuss them in detail.
In this lecture we need to use a function, which waits until the container is ready.
This function uses error channel to signal any error. The only thing we need to know is
we can wait until the async function finish or until error occurs using the ```<-errCh```
expression, which also reads the error.

So the following code waits until something presents on the errCh or the docker engine signal
us that the container is running. In case of error, `log.Fatal` writes the error to the console and exits.
```
errCh := cli.ContainerWait()
if err := <-errCh; err != nil {
	log.Fatal(err)
}
//continue processing if no error
```

