# Go keyword

The go keyword is used to start a new goroutine, which is a lightweight thread managed by the Go runtime.
A goroutine is a function or method that runs concurrently with other goroutines in the same program.

### 
```go
    go s.StartServer()
```

#### Inline func - w/wo params
```go
    param:= "test param"
	
    go func(param string) {
        s.CalculateSomething
        s.StartServer(param)	
    }(param)
	
//------

    go func() {
        s.CalculateSomethin()
        s.StartServer(param)
    }()
```

# Waiting on channel
We can block execution by waiting until channel receiving a value.
In the following example the app is waiting on the given line until context is cancelled.
We will use this for waiting in the main goroutine until server runs and exits only when 
it is stopped by SIGINT or SIGTERM signals, like ctrl+c.

```
<-ctx.Done()
```
For example the Execute function in the next lecture will use the following structure, which 
starts a goroutine for monitoring the aforementioned signals and waits until signal comes.
One it has been received the context will be cancelled. The context cancellation causes
the webserver shutdown.
```go
func Execute() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig
		slog.Info("Shutdown signal received")
		cancel()
	}()


	return executeStg(ctx)
}
```

# Inline Struct
Inline struct can be defined and instantiated immediately. 
These structs don't have name. 

```go
functionWaithSomeStructs( struct{ Status string }{"UP"})

func functionWaithSomeStructs(s any){}
```
