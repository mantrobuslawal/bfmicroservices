# Borough Furniture (BF) Online Store 

## About
This is the monorepo for Borough Furniture ACME, a fictional online furniture retailer for golang gophers.
It's designed for educational purposes, to create a somewhat complete project to use in demonstrations of cloud-based
deployments.

## BF Functional Requirements
**IN PROGRESS -- IN COMPLETE**
The following are a list of functional requirements specified by the product owner:

1. Customers should be able to search for products by product category/sub category.
2. Customers should be able to search for products by brand.
3. New customers should be able to sign-up for an account.
4. Existing customers should be able to login to their accounts.
5. Customers should be able to purchase products from the website.
6. Customers who don't have accounts should be able to perform a _guest checkout_.
7. Customers should be able to track the progress of their order(s).
8. Customers should be able to see a purchase history covering the last past 6 months.


## List of Microservices
- catalog service: provides endpoints for users to get information about products
- order service: processes orders.
- payment service: processes payments.
- notification service
- shipping service
- login service: covers login and authorization
- loyalty service: manages user loyalty points
- inventory service
 
