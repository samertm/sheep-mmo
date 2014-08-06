Feed a sheep, rule the world.

# Message format:

````
             message ::== "(" <message-part> ")"
        message-part ::== <mouse-message>
                        | <sheep-message>
                        | <tick-message>
        tick-message ::== "tick"
       mouse-message ::== "mouse " xcoord " " ycoord
server-mouse-message ::== "mouse " id " " xcoord " " ycoord
       sheep-message ::== "sheep " id " " xcoord " " ycoord
  id, xcoord, ycoord ::== non-negative integer
````

Server-to-client messages: sheep-message, server-mouse-message

Client-to-server messages: mouse-message

