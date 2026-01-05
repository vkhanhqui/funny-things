## Definition of %
The remainder after integer division

There are 2 ways to demonstrate.

### 1. Modulo:
```
a % b = a − int(a / b) * b
```
Note: int(a / b) means the fractional part is discarded

Example:
```
48 % 18
= 48 - int(48 / 18) * 18
= 48 - 2 * 18
= 48 - 36
= 12
```

### 2. Subtraction:
Modulo is just a compressed version of repeated subtraction
```
a % b = a - b - ... - b = c
```

Example:
```
48 % 18
= 48 - 18 - 18
= 12
```
