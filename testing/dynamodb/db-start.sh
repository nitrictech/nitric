#!/bin/bash

java -Djava.library.path=/usr/local/dynamodb/DynamoDBLocal_lib -jar /usr/local/dynamodb/DynamoDBLocal.jar -sharedDb &
