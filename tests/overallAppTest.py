
import requests

url = 'http://localhost:53072'

# test all endpoints not using authentication headers must to return missing authentication headers

print("ALL ENDPOINTS NOT USING HEADERS AUTHENTICATION")

def noAuthenticationTest(response):
    print(
        "----> " + response.request.method + " " + response.request.path_url + " " + 
        "STATUS_CODE:" + ("OK" if response.status_code == 401 else "NOK") + " " + 
        "BODY:" + ("OK" if response.text == "authorization header missing\n" else "NOK")
        )

headers = {
    "Content-Type": "application/json"
}
data = {
    "asd":"asd"
}

noAuthenticationTest(requests.post(url + "/con/ads", json=data, headers=headers))
noAuthenticationTest(requests.put(url + "/con/ads", json=data, headers=headers))
noAuthenticationTest(requests.get(url + "/con/ads", json=data, headers=headers))
noAuthenticationTest(requests.delete(url + "/con/ads", json=data, headers=headers))
noAuthenticationTest(requests.post(url + "/adm/token", json=data, headers=headers))
noAuthenticationTest(requests.delete(url + "/adm/token", json=data, headers=headers))
noAuthenticationTest(requests.post(url + "/adm/context", json=data, headers=headers))
noAuthenticationTest(requests.get(url + "/adm/context", json=data, headers=headers))
noAuthenticationTest(requests.delete(url + "/adm/context", json=data, headers=headers))

print()

# test all endpoints using the Authorization headers with "Authorization: bearer WRONG_TOKEN" must return unauthorized
print("ALL ENDPOINTS THE WRONG TOKEN ON THE HEADER")
def wrongAuthenticationToken(response):
    print(
        "----> " + response.request.method + " " + response.request.path_url + " " + 
        "STATUS_CODE:" + ("OK" if response.status_code == 401 else "NOK") + " " + 
        "BODY:" + ("OK" if response.text == "not authorized\n" else "NOK")
        )

headers = {
    "Content-Type": "application/json",
    "Authorization": "bearer ASDASD"
}
data = {
    "asd":"asd"
}

wrongAuthenticationToken(requests.post(url + "/con/ads", json=data, headers=headers))
wrongAuthenticationToken(requests.put(url + "/con/ads", json=data, headers=headers))
wrongAuthenticationToken(requests.get(url + "/con/ads", json=data, headers=headers))
wrongAuthenticationToken(requests.delete(url + "/con/ads", json=data, headers=headers))
wrongAuthenticationToken(requests.post(url + "/adm/token", json=data, headers=headers))
wrongAuthenticationToken(requests.delete(url + "/adm/token", json=data, headers=headers))
wrongAuthenticationToken(requests.post(url + "/adm/context", json=data, headers=headers))
wrongAuthenticationToken(requests.get(url + "/adm/context", json=data, headers=headers))
wrongAuthenticationToken(requests.delete(url + "/adm/context", json=data, headers=headers))


# test all endpoints using the authentication headers with "Authentication: WRONG_TOKEN" must authentication headers not well written, something like that

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
