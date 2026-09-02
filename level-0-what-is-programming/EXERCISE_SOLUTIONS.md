# Level 0: Exercise Solutions (Pseudocode)

## Exercise 1: Simple Calculation

```
ALGORITHM Calculate(num1, num2):
    sum = num1 + num2
    difference = num1 - num2
    product = num1 × num2
    
    IF num2 = 0 THEN
        PRINT "Cannot divide by zero"
        quotient = undefined
    ELSE
        quotient = num1 ÷ num2
    END IF
    
    RETURN sum, difference, product, quotient
END ALGORITHM
```

**Usage:**
```
results = Calculate(10, 5)
// Output:
// Sum: 15
// Difference: 5
// Product: 50
// Quotient: 2
```

---

## Exercise 2: Number Classification

```
ALGORITHM ClassifyNumber(num):
    IF num > 0 THEN
        PRINT num, "is positive"
    ELSE IF num < 0 THEN
        PRINT num, "is negative"
    ELSE
        PRINT "The number is zero"
    END IF
END ALGORITHM
```

**Usage:**
```
ClassifyNumber(5)    // Output: 5 is positive
ClassifyNumber(-3)   // Output: -3 is negative
ClassifyNumber(0)    // Output: The number is zero
```

---

## Exercise 3: Grade Calculator

```
ALGORITHM CalculateGrade(score):
    IF score >= 90 AND score <= 100 THEN
        PRINT "Grade: A"
    ELSE IF score >= 80 AND score < 90 THEN
        PRINT "Grade: B"
    ELSE IF score >= 70 AND score < 80 THEN
        PRINT "Grade: C"
    ELSE IF score >= 60 AND score < 70 THEN
        PRINT "Grade: D"
    ELSE IF score < 60 THEN
        PRINT "Grade: F"
    ELSE
        PRINT "Invalid score"
    END IF
END ALGORITHM
```

**Usage:**
```
CalculateGrade(95)   // Output: Grade: A
CalculateGrade(75)   // Output: Grade: C
CalculateGrade(55)   // Output: Grade: F
```

---

## Exercise 4: Fibonacci Sequence

```
ALGORITHM Fibonacci(n):
    IF n <= 0 THEN
        PRINT "n must be positive"
        RETURN
    END IF
    
    IF n >= 1 THEN
        PRINT 0
    END IF
    
    IF n >= 2 THEN
        PRINT 1
    END IF
    
    IF n > 2 THEN
        prev2 = 0
        prev1 = 1
        
        FOR i = 3 TO n DO
            current = prev1 + prev2
            PRINT current
            prev2 = prev1
            prev1 = current
        END FOR
    END IF
END ALGORITHM
```

**Usage:**
```
Fibonacci(8)
// Output: 0, 1, 1, 2, 3, 5, 8, 13
```

**Explanation:**
- Start with 0 and 1
- Each next number = sum of previous two
- Repeat until we have n numbers

---

## Exercise 5: Prime Number Checker

```
ALGORITHM IsPrime(num):
    IF num < 2 THEN
        PRINT num, "is not prime"
        RETURN false
    END IF
    
    IF num = 2 THEN
        PRINT num, "is prime"
        RETURN true
    END IF
    
    IF num MOD 2 = 0 THEN
        PRINT num, "is not prime"
        RETURN false
    END IF
    
    FOR i = 3 TO sqrt(num) STEP 2 DO
        IF num MOD i = 0 THEN
            PRINT num, "is not prime"
            RETURN false
        END IF
    END FOR
    
    PRINT num, "is prime"
    RETURN true
END ALGORITHM
```

**Usage:**
```
IsPrime(17)   // Output: 17 is prime, Return: true
IsPrime(10)   // Output: 10 is not prime, Return: false
IsPrime(2)    // Output: 2 is prime, Return: true
```

**Why sqrt(num)?**
- If num has a divisor > sqrt(num), it must also have one < sqrt(num)
- This reduces checks from n to sqrt(n)

---

## Exercise 6: Reverse a String

```
ALGORITHM ReverseString(str):
    reversed = ""
    length = length(str)
    
    FOR i = length - 1 TO 0 STEP -1 DO
        reversed = reversed + str[i]
    END FOR
    
    RETURN reversed
END ALGORITHM
```

**Usage:**
```
ReverseString("Hello")        // Output: "olleH"
ReverseString("Programming")  // Output: "gnimmargorP"
```

**Explanation:**
- Start from the last character
- Move backwards through the string
- Build new string by concatenating characters

---

## Exercise 7: Palindrome Checker

```
ALGORITHM IsPalindrome(str):
    // Convert to lowercase for comparison
    str = lowercase(str)
    
    left = 0
    right = length(str) - 1
    
    WHILE left < right DO
        IF str[left] ≠ str[right] THEN
            PRINT "Not a palindrome"
            RETURN false
        END IF
        left = left + 1
        right = right - 1
    END WHILE
    
    PRINT "Is a palindrome"
    RETURN true
END ALGORITHM
```

**Usage:**
```
IsPalindrome("racecar")   // Output: Is a palindrome, Return: true
IsPalindrome("hello")     // Output: Not a palindrome, Return: false
IsPalindrome("madam")     // Output: Is a palindrome, Return: true
```

**Two-Pointer Approach:**
- Compare characters from both ends moving inward
- If any mismatch found, it's not a palindrome
- Time complexity: O(n), Space complexity: O(1)

---

## Exercise 8: Count Occurrences

```
ALGORITHM CountCharacter(str, char):
    count = 0
    length = length(str)
    
    FOR i = 0 TO length - 1 DO
        IF str[i] = char THEN
            count = count + 1
        END IF
    END FOR
    
    PRINT "Character '", char, "' appears", count, "times"
    RETURN count
END ALGORITHM
```

**Usage:**
```
CountCharacter("programming", 'r')  // Output: 2
CountCharacter("hello world", 'l')   // Output: 3
CountCharacter("mississippi", 's')   // Output: 4
```

**Explanation:**
- Iterate through each character in string
- Compare with target character
- Increment counter on match
- Return final count

---

## 🧠 Key Takeaways

1. **Algorithm Design:** Think step-by-step before implementing
2. **Edge Cases:** Always consider boundary conditions (n=0, n=1, etc.)
3. **Efficiency:** Some approaches are better than others (e.g., checking up to sqrt(n) for primes)
4. **Reusability:** These algorithms work in any language
5. **Testing:** Test with multiple inputs including edge cases

---

## 🎯 Challenge Exercises

### Challenge 1: Reverse Fibonacci
Find all Fibonacci numbers less than 1000, then reverse the list.

### Challenge 2: Prime Range
Find all prime numbers between two given numbers.

### Challenge 3: String Manipulation
Given a string, remove all vowels and return the result.

### Challenge 4: Nested Loops
Generate a multiplication table for numbers 1-10.

### Challenge 5: Complex Logic
Check if a string contains balanced parentheses: "((()))" = valid, "((())" = invalid.
