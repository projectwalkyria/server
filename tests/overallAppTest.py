import requests

# test all endpoints not using auhtentication must to return unauthorized

# try create a token without the barear token and with a wrong token, must to return unauthorized

# try create a token with the correct bearer token, must retun created

# try create a context without the barear token and with a wrong token, must to return unauthorized

# try create a context with the correct bearer token must return created

# try delete a token without the barear token and with a wrong token, must to return unauthorized

# try delete a context without the barear token and with a wrong token, must to return unauthorized

# try to use the token created to insert entries on the context must to return unauthorized

# grant the privilege to write on the context and test create read write and delete on the context, everything but write must to return unauthorized

# revoke the privilege to write and give the token the privilege to read, everything but read must to return unauthorized

# revoke the privilege to read and give the token the privilege to update, everything but update must to return unauthorized

# revoke the privilege to write and give the token the privilege to delete, everything but delete must to return unauthorized

# try delete a token with the correct bearer token, must retun deleted

# try delete a context with the correct bearer token must return deleted
