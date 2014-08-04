Feed a sheep, rule the world.

# Message format:

````
      message ::== "(" <message-part> ")"
 message-part ::== <mouse-message>
                 | <sheep-message>
mouse-message ::== "mouse" xcoord ycoord
sheep-message ::== "sheep" xcoord ycoord
````

