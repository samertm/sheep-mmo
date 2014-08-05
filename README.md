Feed a sheep, rule the world.

# Message format:

````
           message ::== "(" <message-part> ")"
      message-part ::== <mouse-message>
                      | <sheep-message>
     mouse-message ::== "mouse" id xcoord ycoord
     sheep-message ::== "sheep" xcoord ycoord
id, xcoord, ycoord ::== non-negative integer
````

Server-to-client messages: sheep-message

Client-to-server messages: mouse-message

