This project is supposed to be a quick tcp chatroom in golang, named olc.
Welcome to olc server README file !
Right now, it implements :
-Renaming members (and handling two same-named members, by adding a 'bis' after the nick)
-Handling a number of chatters in the room
-Automatically removing connections when a member leaves
And in the future, the project might include :
-Blacklisting or Whitelisting IPs
-Comunicating the name of the server to the client
-Encryption of the chat (probably only encrypt it with a simple protocol)
