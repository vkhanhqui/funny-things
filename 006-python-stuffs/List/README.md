# List in Python


## So what is list?
- A collection which is changeable and allows duplicate members.
- To store multiple items in a single variable.


## Create a list
```python
new_list = ["Qui", "Steve"]
print(new_list)
## ["Qui", "Steve"]
```


## Some advanced things
- Generate 2D lists
```python
# 1
nums = [[] for _ in range((3))]
nums[1] = 1
print(nums)
## [[], 1, []]

# 2
nums = [[]] * 3
nums[1] = 1
print(nums)
## [[], 1, []]
```

- Pass by value VS pass by reference
```python
# Pass by value
nums = [[] for _ in range((3))]
nums[1].append(1)
print(nums)
## [[], [1], []]

# Pass by reference
nums = [[]] * 3
nums[1].append(1)
print(nums)
## [[1], [1], [1]]
```