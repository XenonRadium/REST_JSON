# REST_JSON

A copy of the test database is included as messages.sql to be imported into sql for demo purposes.

The program is able to read REST commands from an external program such as Postman or Advanced REST Client.
Depending on whether GET, POST, PATCH, DELETE is used, different outcome will be produced.

/getMessage/{version} GET protocol is used to get specific messages from database with the specified version.
/addMessage is used to add new messages to the database according to the JSON body.
/deleteMessage/{version} is used to delete messages with the specified version.
/updateMessage/{version} is used to update messages with the specified version according to the JSON body.
