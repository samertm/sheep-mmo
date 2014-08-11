Feed a sheep, rule the world.

# Message format:

````
             message ::== "(" <message-part> ")"
        message-part ::== <*-message>
        tick-message ::== "tick"
       mouse-message ::== "mouse " xcoord " " ycoord
server-mouse-message ::== "mouse " id " " xcoord " " ycoord
       sheep-message ::== "sheep " id " " xcoord " " ycoord " "
                          sheep-name " " state
      rename-message ::== "rename " id " " sheep-name
   gen-sheep-message ::== "gen-sheep"
       fence-message ::== "fence " xcoord " " ycoord " " width " " height
  id, xcoord, ycoord,
       height, width ::== non-negative integer
   sheep-name, state ::== string (can be delimited with double quotes)
````

Server-to-client messages: sheep-message, server-mouse-message

Client-to-server messages: mouse-message, rename-message, gen-sheep-message

