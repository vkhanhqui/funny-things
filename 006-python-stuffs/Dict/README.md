# Dict in Python


## So what is dict?
- Stand for dictionary
- Dictionaries are used to store data values in key:value pairs.
- A dictionary is a collection which is changeable and cannot have two items with the same key.
- There are many categories of dict in Python such as <b>normal dict, defaultdict, OrderedDict, Counter, ChainMap, UserDict.</b>


## Dictionaries are written with curly brackets, and have keys and values:
```python
new_dict = {
  "id": "abcd",
  "name": "Steve",  # this one is duplicated so that will be removed
  "name": "Qui"
}
print(new_dict)
# {'id': 'abcd', 'name': 'Qui'}
```


## Some advanced things
```python
# We cannot put a list as dict key but we can do with tuple
## 1 list
new_dict = {}
new_dict.update({
    ['a', 'b']: ['a', 'b']
})
##  new_dict.update({
## TypeError: unhashable type: 'list'

## 2 tuple
new_dict = {}
new_dict.update({
    ('c', 'd'): ['c', 'd']
})
## {('c', 'd'): ['c', 'd']}
```