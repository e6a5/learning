# Playing on terminal with go

## Learning Rules

### **TDD Rules (Red-Green-Refactor)**

#### **The Three Laws:**
1. **Write ONLY enough of a failing test** to demonstrate the problem
2. **Write ONLY enough code** to make that test pass
3. **Refactor** the code while keeping tests green

#### **The Cycle:**
- **Red**: Write a failing test
- **Green**: Write minimal code to pass
- **Refactor**: Clean up without changing behavior

---

### **Unix Philosophy Rules**

#### **Core Principles:**
1. **Do one thing well** - Each program has a single, clear purpose
2. **Work with others** - Programs communicate via standard interfaces
3. **Handle text streams** - Use stdin/stdout/stderr appropriately
4. **Fail fast and clearly** - Exit codes and error messages matter
5. **No unnecessary output** - Silent success, verbose only when needed

#### **CLI Best Practices:**
- Exit code 0 = success, non-zero = failure
- Use flags consistently (`-h` for help, `-v` for verbose)
- Read from stdin when no file specified
- Write output to stdout, errors to stderr

---

### **Terminal Drawing Fundamentals**

#### **ANSI Escape Sequences:**
- Start with `\033[` or `\x1b[`
- End with a letter command
- Examples: `\033[2J` (clear screen), `\033[1;1H` (move to top-left)

#### **Coordinate System:**
- Top-left is (1,1) not (0,0)
- X = column, Y = row
- Most terminals are 80x24 by default