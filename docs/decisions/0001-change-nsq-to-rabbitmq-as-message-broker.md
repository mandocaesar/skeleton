# Change NSQ to RabbitMQ as Message Broker

* Status: accepted
* Deciders: Armand, Caesario
* Date: 2023-08-14

Technical Story: Replace NSQ with Rabbit MQ in terms of assurance of sequence

## Context and Problem Statement

NSQ don't have feature to ensure queue sequences, and we need to make sure message re-queue when consumer not ack or down

## Decision Drivers

* Order processing need processed at most one and need in processed sequence
* NSQ messed up the sequence of message which impacted to business revenue
* Queue on NSQ don't have DLX when consumer down

## Considered Options

* change message broker to Rabbit MQ
* implement DLX on NSQ

## Decision Outcome

Chosen option: "change message broker to Rabbit MQ", because comes out best.

## Pros and Cons of the Options

### change message broker to Rabbit MQ

Rabbit MQ comply all needs from business needs

* Good, because cater all scenario in bussiness need
* Good, because use by wider community
* Good, because battle tested
* Bad, because we need imlemented it one by one in each service
* Bad, because we don't have monitor queue yet for RMQ

### implement DLX on NSQ

implement virtual queue in NSQ to mimic how DLX work

* Good, because we dont need to re-implement on services
* Bad, because we need to have extensive test since this feature is not NSQ nature as message broker
