# Nuance Retrieval Service 

## Overview 

The Nuance Retrieval Service is exposes a RAG application over a tcp server

agenda:


[07/17/2024]

implemente horizontal skaling w/ k8s(2 or more)

get client to connect to one of the servers and then the responding server broadcasts to all clients

put a load balancing application in front of the two servers

eventing service for handling the messages 

can speed up implementation by leveraging k8s from docker desktop 

test coverage as well 

NewServer() should start the listener and can then return itself plus the error so you can still check the error in main

pass the config to new server as well so you can parameterize port of service

we should have a distinct client package

have a start function(will be your handle client thing) and then call that with a goroutine in main, and then the start function will have your main for loop in it

go s.HandleClient(conn) should be in Start or Begin function

simplify server

implement concurrent broadcast to clients, needs to handle for arbitray scaling 

fence off open API stuff

proper commenting, documentation, and testing 
