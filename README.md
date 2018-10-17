## Dumb! The Chatbot
This is minimum functionality chatbot written in **Go**.

It can give responses like, `current weather`, `time`,`date`, `day` and `month`.
The `JSON` request and response objects are shown below.

### JSON Request Object

```json
{
  "messages": [
    {
      "type": "string",
      "unstructured": {
        "id": "string",
        "text": "weather",
        "timestamp": "string"
      }
    }
  ]
}
```

### JSON Response Object

```json
{
    "messages": [
        {
            "type": "POST",
            "unstructured": {
                "id": "6ea3b24c-d223-11e8-b5ea-557478d78d46",
                "text": "Temperature in New York is 9°C and it's, clear sky.",
                "timestamp": "2018-10-17 15:43:45.833831992 +0000 UTC m=+2558.925851738"
            }
        }
    ]
}
```

It uses an **API Gateway** on Amazon with the following specification.

#### API specification
```yaml
swagger: '2.0'
info:
  title: AI Customer Service API
  description: 'AI Customer Service application, built during the Cloud and Big Data course at Columbia University.'
  version: 1.0.0
schemes:
  - https
basePath: /v1
produces:
  - application/json
paths:
  /chatbot:
    post:
      summary: The endpoint for the Natural Language Understanding API.
      description: |
        This API takes in one or more messages from the client and returns
        one or more messages as a response. The API leverages the NLP
        backend functionality, paired with state and profile information
        and returns a context-aware reply.
      tags:
        - NLU
      operationId: sendMessage
      produces:
        - application/json
      parameters:
        - name: body
          in: body
          required: true
          schema:
            $ref: '#/definitions/BotRequest'
      responses:
        '200':
          description: A Chatbot response
          schema:
            $ref: '#/definitions/BotResponse'
        '403':
          description: Unauthorized
          schema:
            $ref: '#/definitions/Error'
        '500':
          description: Unexpected error
          schema:
            $ref: '#/definitions/Error'
definitions:
  BotRequest:
    type: object
    properties:
      messages:
        type: array
        items:
          $ref: '#/definitions/Message'
  BotResponse:
    type: object
    properties:
      messages:
        type: array
        items:
          $ref: '#/definitions/Message'
  Message:
    type: object
    properties:
      type:
        type: string
      unstructured:
        $ref: '#/definitions/UnstructuredMessage'
  UnstructuredMessage:
    type: object
    properties:
      id:
        type: string
      text:
        type: string
      timestamp:
        type: string
        format: datetime
  Error:
    type: object
    properties:
      code:
        type: integer
        format: int32
      message:
        type: string
``` 


The API is also uses an **Access Key**, which must be included in the header as:

`x-api-key : <Access Key>`

The bot is currently hosted on **Amazon S3** and can be accessed at 

https://s3.amazonaws.com/k8s.procrypt.xyz/index.html 

OR

http://k8s.procrypt.xyz/