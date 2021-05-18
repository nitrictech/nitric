#!/bin/bash

# Create Standard Table (customer)
echo 'Create table: application'

aws dynamodb delete-table \
    --table-name customer \
    --endpoint-url http://localhost:8000

aws dynamodb create-table \
    --table-name customer \
    --attribute-definitions AttributeName=key,AttributeType=S \
    --key-schema AttributeName=key,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1 \
    --endpoint-url http://localhost:8000

aws dynamodb put-item \
    --table-name customer \
    --item '{
        "key": {
            "S": "jane.smith@server.com"
        }, 
        "value": {
            "M": { 
                "firstName": {"S": "Jane"}, 
                "lastName": {"S": "Smith"}, 
                "email": {"S": "jane.smith@server.com"},
                "mobile": {"S": "0482847293"}
            }
        }
    }' \
    --endpoint-url http://localhost:8000

aws dynamodb put-item \
    --table-name customer \
    --item '{
        "key": {
            "S": "paul.davis@server.com"
        }, 
        "value": {
            "M": { 
                "firstName": {"S": "Paul"}, 
                "lastName": {"S": "Davis"}, 
                "email": {"S": "paul.davis@server.com"},
                "mobile": {"S": "041231234"}
            }
        }
    }' \
    --endpoint-url http://localhost:8000


aws dynamodb scan \
    --table-name customer \
    --endpoint-url http://localhost:8000

# Create Single Table Design Table (application)
echo 'Create table: application'    

aws dynamodb delete-table \
    --table-name application \
    --endpoint-url http://localhost:8000

aws dynamodb create-table \
    --table-name application \
    --attribute-definitions AttributeName=pk,AttributeType=S AttributeName=sk,AttributeType=S \
    --key-schema AttributeName=pk,KeyType=HASH AttributeName=sk,KeyType=RANGE \
    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1 \
    --endpoint-url http://localhost:8000

aws dynamodb put-item \
    --table-name application \
    --item '{
        "pk": {
            "S": "Customer#1000"
        }, 
        "sk": {
            "S": "Customer#1000"
        }, 
        "value": {
            "M": { 
                "firstName": {"S": "Jane"}, 
                "lastName": {"S": "Smith"},
                "email": {"S": "jane.smith@gmail.com"},
                "mobile": {"S": "0482847293"} 
            } 
        }
    }' \
    --endpoint-url http://localhost:8000

aws dynamodb put-item \
    --table-name application \
    --item '{
        "pk": {
            "S": "Customer#1000"
        }, 
        "sk": {
            "S": "Order#501"
        }, 
        "value": {
            "M": {
                "order": {"S": "501"}, 
                "sku": {"S": "83-292-236"}, 
                "number": {"S": "2"}
            } 
        }
    }' \
    --endpoint-url http://localhost:8000

aws dynamodb put-item \
    --table-name application \
    --item '{
        "pk": {
            "S": "Customer#1000"
        }, 
        "sk": {
            "S": "Order#502"
        }, 
        "value": {
            "M": { 
                "order": {"S": "502"}, 
                "sku": {"S": "171-823-623"}, 
                "number": {"S": "12"}
            } 
        }
    }' \
    --endpoint-url http://localhost:8000

aws dynamodb put-item \
    --table-name application \
    --item '{
        "pk": {
            "S": "Customer#1000"
        }, 
        "sk": {
            "S": "Order#503"
        }, 
        "value": {
            "M": { 
                "order": {"S": "503"}, 
                "sku": {"S": "171-823-623"}, 
                "number": {"S": "1"}
            } 
        }
    }' \
    --endpoint-url http://localhost:8000

aws dynamodb put-item \
    --table-name application \
    --item '{
        "pk": {
            "S": "Product#170-592-923"
        }, 
        "sk": {
            "S": "Product#170-592-923"
        }, 
        "value": {
            "M": { 
                "sku": {"S": "171-823-623"}, 
                "color": {"S": "red"},
                "weight": {"S": "34.2"}
            } 
        }
    }' \
    --endpoint-url http://localhost:8000

aws dynamodb scan \
    --table-name application \
    --endpoint-url http://localhost:8000

