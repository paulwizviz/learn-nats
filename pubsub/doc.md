# Pub Sub model

A pub sub communication model distribute messages on a one-to-many model. A publisher sends a message via a subject to consumers listening or subscribed to a subject on NATS broker.

![Pub Sub Model](../assets/img/pub-sub.png)

## Examples

### Example 1 - Simple one-to-many

This [example](./ex1/main.go) demonstrates a simple one producer to two consumers arrangement via a subject name `my-topic`.