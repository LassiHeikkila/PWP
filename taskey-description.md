# taskey

## idea
Web service where admin can schedule and trigger tasks to be performed on enrolled devices.
Similar to cron, but remotely configurable + reporting results to web service.

## Constraints
Daemon should support Linux only, I don't care about Windows machines.

Server should be containerized, with all dependencies included in the image.
Database(s) may run as separate container.
Daemon could also run in a container, but naturally that will limit what it can do.

## Task
Task can be anything:
- read a file
- run a single command
- do a web request
    - control headers, methods, body
- run a script
- etc.

## API
API should have endpoints for (human or scripted user):
- creating and managing tasks
- creating and managing user accounts
- creating and managing user groups
- creating and managing device groups
- enrolling and unenrolling devices
- scheduling tasks to run at certain time(s)
- triggering task to run immediately on machine or group of machines

API should also have endpoints for (daemon running on devices):
- getting tasks for device
- reporting task result from device

There should be a lightweight real-time communication API as well, used to trigger actions immediately.
Maybe that should be pub-sub? MQTT? Could also be used as presence indicator.
Publishing should be heavily restricted.

## Definitions
### Task
```graphql
type Task {
    id: ID!
    description: String
    actions: [Action]!
}

enum ActionType {
    HTTP_REQUEST
    SCRIPT
    COMMAND
}

interface Action {
    type: ActionType!
}

type ActionResult {
    Code: Integer
    Output: String
    Error: String
}

type HttpAction implements Action {
    type: ActionType!
    requestMethod: HttpMethod!
    requestUrl: String!
    headers: [String!]
    form: [KeyValuePair!]
    body: String
}

type ScriptAction implements Action {
    type: ActionType!
    content: String
}

type CommandAction implements Action {
    type: ActionType!
    command: String!
    arguments: [String!]
    environment: [KeyValuePair!]
}

enum HttpMethod {
    GET
    POST
    PUT
    DELETE
}

type KeyValuePair {
    Key: String
    Value: String
}
```
### User
```json
{
    "id": "i099j309dj3j40j3",
    "username": "myusername",
    "organization": "3f943j0fj943j0f"
}
```
### Organization
```json
{
    "id": "3f943j0fj943j0f",
    "name": "www.example.com",
    "members": [
        "i099j309dj3j40j3"
    ],
    "machines": [
        "i902j390cj32909c3"
    ]
}
```
### Machine
```json
{
    "id": "i902j390cj32909c3"
}
```
