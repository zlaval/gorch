# Struct

Struct is a composite data type that groups together fields to create custom data types. 
Structs are similar to classes in other programming languages but without the concept of inheritance.


```go
type Person struct {
    Name string
    Age  int
}

p := Person{Name: "Alice", Age: 25}
p.Name = "Charlie"
```

# Receiver function
Receiver functions are methods associated with a specific type. 
They allow us to add functionality to structs without cluttering your
main program logic. Receiver functions are defined by specifying a receiver parameter in the method signature,
which can be either a value or a pointer. 
We won't dive deep into the difference between values and pointers in this course, just simply use
pointers only when we need to modify a value of the struct in the receiver functions.

```go
func (p Person) IsAdult() bool {
    return p.Age >= 18
}

func (p *Person) IncrementAge() {
    p.Age += 1
}

if p.IsAdult(){
    p.IncrementAge()
}
```

# Embedding
Embedding allows one struct to include another struct as a field.
This enables a form of composition, letting us create hierarchical and reusable data models. 
The embedded struct’s fields and methods are promoted, meaning they can be accessed 
directly from the embedding struct without explicitly referencing the embedded field.
A struct can be embedded into another struct using its type only, without field name

```go
type Address struct {
    Street string
    City   string
    Zip    int
}

type Person struct {
    Address
    
    Name    string
    Age     int 
}
```

# Visibility
In Go, visibility (public or private) is determined by the case of the first letter of an identifier 
(e.g., variables, functions, types, methods, or struct fields).

*   Public: If the identifier starts with an uppercase letter, it is exported (visible outside its defining package).
*   Private: If the identifier starts with a lowercase letter, it is unexported (visible only within its defining package).

```go 
type Person struct {
    Name  string // Public (exported)
    age   int    // Private (unexported)
}
```

# Sprintf
fmt.Sprintf is a function is used to create formatted strings.

```go
msg := fmt.Sprintf("Name: %s, Age: %d", name, age)
```

#  For Loop
We can use the range keyword to retrieve both the index and value of each element in a slice.
```go
for i, p := range people {
  fmt.Printf("Person %d: Name: %s, Age: %d\n", i, people[i].Name, p.Age)
}
```
We can use slices.Values function to omit the index.
```go
for p := range slices.Values(people) {
  fmt.Printf("Person: Name: %s, Age: %d\n", p.Name, p.Age)
}
```
For loop on map retrieves the key (and value)
```go
for key, value := range map {
    
}
```

First value can be skipped using _.
```go
for _, p := range people {
fmt.Printf("Person: Name: %s, Age: %d\n", p.Name, p.Age)
}
```

Second value can be omitted.
```go
for i := range people {
fmt.Printf("Person: Name: %s, Age: %d\n", people[i].Name, people[i].Age)
}
```


# Default folder structure
In Go, proper organization of code is crucial. 
Common patterns for structuring folders in a Go project are the use of internal and pkg directories. 

The internal directory is a special folder in Go used to restrict the visibility 
of code to the current module or submodule.
Code inside an internal directory cannot be imported by other modules.

The pkg directory is often used for code intended to be reused across multiple projects or modules. 
It’s a convention rather than a strict Go feature.