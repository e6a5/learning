// Package basics demonstrates Go fundamentals
// This package contains examples of core Go concepts used in the HTTP server

package basics

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// RunAllExamples executes all Go fundamentals demonstrations
func RunAllExamples() {
	fmt.Println("🎓 Go Fundamentals Examples")
	fmt.Println(strings.Repeat("=", 40))

	// 1. Variables and Types
	fmt.Println("\n1️⃣ Variables and Types")
	DemonstrateVariables()

	// 2. Control Structures
	fmt.Println("\n2️⃣ Control Structures")
	DemonstrateControlStructures()

	// 3. Functions
	fmt.Println("\n3️⃣ Functions")
	DemonstrateFunctions()

	// 4. Structs and Methods
	fmt.Println("\n4️⃣ Structs and Methods")
	DemonstrateStructs()

	// 5. Slices and Maps
	fmt.Println("\n5️⃣ Slices and Maps")
	DemonstrateCollections()

	// 6. Pointers
	fmt.Println("\n6️⃣ Pointers")
	DemonstratePointers()

	// 7. Error Handling
	fmt.Println("\n7️⃣ Error Handling")
	DemonstrateErrorHandling()
}

// DemonstrateVariables shows different ways to declare and use variables
func DemonstrateVariables() {
	// Different ways to declare variables
	var name string = "Alice" // Explicit type
	age := 30                 // Type inference
	const maxUsers = 100      // Constant
	var isActive bool = true  // Boolean
	var score float64 = 95.5  // Float

	fmt.Printf("Name: %s, Age: %d, Max Users: %d\n", name, age, maxUsers)
	fmt.Printf("Active: %t, Score: %.1f\n", isActive, score)

	// Zero values
	var count int      // 0
	var message string // ""
	var ready bool     // false
	fmt.Printf("Zero values - Count: %d, Message: '%s', Ready: %t\n", count, message, ready)
}

// DemonstrateControlStructures shows if/else, loops, switch statements
func DemonstrateControlStructures() {
	// If/else
	age := 25
	if age >= 18 {
		fmt.Println("Adult")
	} else {
		fmt.Println("Minor")
	}

	// For loop (only loop in Go!)
	fmt.Print("Numbers: ")
	for i := 1; i <= 5; i++ {
		fmt.Printf("%d ", i)
	}
	fmt.Println()

	// Range over slice
	fruits := []string{"apple", "banana", "orange"}
	fmt.Print("Fruits: ")
	for index, fruit := range fruits {
		fmt.Printf("%d:%s ", index, fruit)
	}
	fmt.Println()

	// Switch statement
	day := "Monday"
	switch day {
	case "Monday", "Tuesday":
		fmt.Println("Start of work week")
	case "Friday":
		fmt.Println("TGIF!")
	default:
		fmt.Println("Regular day")
	}
}

// DemonstrateFunctions shows function patterns including multiple returns and error handling
func DemonstrateFunctions() {
	// Function with multiple return values
	result, err := Divide(10, 2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("10 / 2 = %d\n", result)
	}

	// Function with error
	result, err = Divide(10, 0)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Variadic function
	sum := AddNumbers(1, 2, 3, 4, 5)
	fmt.Printf("Sum of 1,2,3,4,5 = %d\n", sum)
}

// Divide demonstrates function with multiple return values (common Go pattern)
func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("cannot divide by zero")
	}
	return a / b, nil
}

// AddNumbers demonstrates variadic function (accepts multiple arguments)
func AddNumbers(numbers ...int) int {
	total := 0
	for _, num := range numbers {
		total += num
	}
	return total
}

// DemonstrateStructs shows struct definition, creation, and usage
func DemonstrateStructs() {
	// Define and use struct
	type Person struct {
		Name    string
		Age     int
		Email   string
		IsAdmin bool
	}

	// Create struct instances
	person1 := Person{
		Name:    "Bob",
		Age:     35,
		Email:   "bob@example.com",
		IsAdmin: false,
	}

	person2 := Person{"Carol", 28, "carol@example.com", true}

	fmt.Printf("Person 1: %+v\n", person1)
	fmt.Printf("Person 2: %+v\n", person2)

	// Access and modify fields
	person1.Age = 36
	fmt.Printf("Updated age: %d\n", person1.Age)
}

// DemonstrateCollections shows slices and maps usage
func DemonstrateCollections() {
	// Slices (dynamic arrays)
	var users []string                    // Empty slice
	users = append(users, "Alice", "Bob") // Add elements
	users = append(users, "Carol")        // Add more

	fmt.Printf("Users slice: %v\n", users)
	fmt.Printf("Length: %d, Capacity: %d\n", len(users), cap(users))

	// Slice operations
	fmt.Printf("First user: %s\n", users[0])
	fmt.Printf("Last user: %s\n", users[len(users)-1])
	fmt.Printf("First two: %v\n", users[:2])

	// Maps (key-value pairs)
	userAges := make(map[string]int)
	userAges["Alice"] = 30
	userAges["Bob"] = 35
	userAges["Carol"] = 28

	fmt.Printf("User ages: %v\n", userAges)

	// Check if key exists
	if age, exists := userAges["Alice"]; exists {
		fmt.Printf("Alice is %d years old\n", age)
	}

	// Iterate over map
	fmt.Print("All ages: ")
	for name, age := range userAges {
		fmt.Printf("%s:%d ", name, age)
	}
	fmt.Println()
}

// DemonstratePointers shows pointer usage and memory management
func DemonstratePointers() {
	// Pointers are addresses to memory locations
	age := 25
	agePtr := &age // Get address of age

	fmt.Printf("Age value: %d\n", age)
	fmt.Printf("Age address: %p\n", agePtr)
	fmt.Printf("Value at address: %d\n", *agePtr)

	// Modify through pointer
	*agePtr = 26
	fmt.Printf("Modified age: %d\n", age)

	// Pointers with structs
	type User struct {
		Name string
		Age  int
	}

	user := User{"Dave", 40}
	userPtr := &user

	// Access struct fields through pointer
	fmt.Printf("User name: %s\n", userPtr.Name) // Go automatically dereferences
	userPtr.Age = 41
	fmt.Printf("Updated user: %+v\n", user)
}

// DemonstrateErrorHandling shows Go's explicit error handling patterns
func DemonstrateErrorHandling() {
	// Go uses explicit error handling (no exceptions)

	// Example 1: String to int conversion
	numStr := "123"
	if num, err := strconv.Atoi(numStr); err != nil {
		fmt.Printf("Error converting '%s': %v\n", numStr, err)
	} else {
		fmt.Printf("Converted '%s' to %d\n", numStr, num)
	}

	// Example 2: Invalid conversion
	invalidStr := "abc"
	if _, err := strconv.Atoi(invalidStr); err != nil {
		fmt.Printf("Error converting '%s': %v\n", invalidStr, err)
	}

	// Example 3: Working with time
	timeStr := "2023-12-25"
	layout := "2006-01-02" // Go's reference time
	if parsedTime, err := time.Parse(layout, timeStr); err != nil {
		fmt.Printf("Error parsing time: %v\n", err)
	} else {
		fmt.Printf("Parsed time: %v\n", parsedTime.Format("January 2, 2006"))
	}
}

/*
Key Go Concepts Demonstrated:

1. **Package Structure**: Exported functions (capital letters) vs unexported
2. **Variables**: var, :=, const, zero values
3. **Types**: string, int, bool, float64
4. **Control Flow**: if/else, for, range, switch
5. **Functions**: multiple returns, error handling, variadic
6. **Structs**: definition, initialization, field access
7. **Slices**: dynamic arrays, append, indexing
8. **Maps**: key-value pairs, existence checking
9. **Pointers**: addresses, dereferencing, struct pointers
10. **Error Handling**: explicit errors, nil checks

These concepts are the foundation of the HTTP server in main.go!
*/
