[183. Customers Who Never Order](https://leetcode.com/problems/customers-who-never-order)

```sql
SELECT Name as Customers
FROM Customers
LEFT JOIN Orders
ON Orders.customerId  = Customers.Id
WHERE Orders.customerId is null

```
