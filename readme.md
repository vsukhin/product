# Overview

This is a dead simple server implementing RESTful API.
Please find below all methods description and examples.
All dependencies are stored in Godep vendor folders.
By default it's running at 3000 port, but you can change it using -port parameter.   

## Create Product

```sh
curl -X POST -H "Content-Type: application/json; charset=utf-8" -d '{"Name":"New Product", "Description": "Should be cool", "Price": 5.0, "Tags":["new", "cool"]}' "http://localhost:3000/products"
```

### REQUEST HEADERS

Content-Type: application/json; charset=utf-8

### REQUEST BODY
```
{
	"Name": "New Product", 

	"Description": "Should be cool", 

	"Price": 5.0, 

	"Tags": ["new", "cool"]
}
```

### RESPONSE HEADERS

content-type: application/json; charset=utf-8

status: 200 OK

### RESPONSE BODY
```
{
	"ID": 1
}
```

## Update Product

```sh
curl -X PUT -H "Content-Type: application/json; charset=utf-8" -d '{"Name":"New Product Updated", "Description": "Should be cool", "Price": 15.0, "Tags":["new", "cool", "updated"]}' "http://localhost:3000/products/1"
```

### REQUEST HEADERS

Content-Type: application/json; charset=utf-8

### REQUEST BODY
```
{
	"Name": "New Product Updated", 

	"Description": "Should be cool", 

	"Price": 15.0, 

	"Tags": ["new", "cool", "updated"]
}
```

### RESPONSE HEADERS

content-type: application/json; charset=utf-8

status: 200 OK

### RESPONSE BODY
```
```

## List Products

```sh
curl -X GET -H "Content-Type: application/json; charset=utf-8" "http://localhost:3000/products"
```

### REQUEST HEADERS

Content-Type: application/json; charset=utf-8

### REQUEST BODY
```
```

### RESPONSE HEADERS

content-type: application/json; charset=utf-8

status: 200 OK

### RESPONSE BODY
```
{
	[
		{
			"ID": 1,

			"Name": "New Product Updated",

			"Description": "Should be cool",

			"Price": 15, 

			"Tags": ["new", "cool", "updated"]
		}
	]
}
```

## Set Prices

```sh
curl -X PUT -H "Content-Type: application/json; charset=utf-8" -d '{"EUR": 20.15, "RUR":100.75}' "http://localhost:3000/products/1/prices"
```

### REQUEST HEADERS

Content-Type: application/json; charset=utf-8

### REQUEST BODY
```
{
	"EUR": 20.15, 

	"RUR": 1000.75, 
}
```

### RESPONSE HEADERS

content-type: application/json; charset=utf-8

status: 200 OK

### RESPONSE BODY
```
```

## Get Product

```sh
curl -X GET -H "Content-Type: application/json; charset=utf-8" "http://localhost:3000/products/1"
```

### REQUEST HEADERS

Content-Type: application/json; charset=utf-8

### REQUEST BODY
```
```

### RESPONSE HEADERS

content-type: application/json; charset=utf-8

status: 200 OK

### RESPONSE BODY
```
{
	"ID": 1,

	"Name": "New Product Updated",

	"Description": "Should be cool",

	"Price": 15, 

	"Tags": ["new", "cool", "updated"],

	"Prices":{"EUR": 20.15, "RUR": 100.75}
}
```

## Delete Product

```sh
curl -X DELETE -H "Content-Type: application/json; charset=utf-8" "http://localhost:3000/products/1"
```

### REQUEST HEADERS

Content-Type: application/json; charset=utf-8

### REQUEST BODY
```
```

### RESPONSE HEADERS

content-type: application/json; charset=utf-8

status: 200 OK

### RESPONSE BODY
```
```

## Status codes
* 200 - OK, success
* 400 - Bad request with wrong product data
* 404 - Product is not found
* 500 - Internal server error
