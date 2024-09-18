
import requests
import json


url = 'http://localhost:53072'
ADM_TOKEN = ''

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

print()

# test all endpoints using the authentication headers with "Authentication: WRONG_TOKEN" must authentication headers not well written, something like that
print("ALL ENDPOINTS WITH THE WRONG AUTHENTICATION TOKEN STRUCTURE")
def wrongAuthenticationStructure(response):
    print(
        "----> " + response.request.method + " " + response.request.path_url + " " + 
        "STATUS_CODE:" + ("OK" if response.status_code == 401 else "NOK") + " " + 
        "BODY:" + ("OK" if response.text == "invalid Authorization header format\n" else "NOK")
        )

headers = {
    "Content-Type": "application/json",
    "Authorization": "ASDASD"
}
data = {
    "asd":"asd"
}

wrongAuthenticationStructure(requests.post(url + "/con/ads", json=data, headers=headers))
wrongAuthenticationStructure(requests.put(url + "/con/ads", json=data, headers=headers))
wrongAuthenticationStructure(requests.get(url + "/con/ads", json=data, headers=headers))
wrongAuthenticationStructure(requests.delete(url + "/con/ads", json=data, headers=headers))
wrongAuthenticationStructure(requests.post(url + "/adm/token", json=data, headers=headers))
wrongAuthenticationStructure(requests.delete(url + "/adm/token", json=data, headers=headers))
wrongAuthenticationStructure(requests.post(url + "/adm/context", json=data, headers=headers))
wrongAuthenticationStructure(requests.get(url + "/adm/context", json=data, headers=headers))
wrongAuthenticationStructure(requests.delete(url + "/adm/context", json=data, headers=headers))

print()

# create token
print("CREATE TOKEN")
def createToken(response):    
    try:
        token = json.loads(response.text)['token']
    except:
        print(
            "----> " + response.request.method + " " + response.request.path_url + " " + 
            "STATUS_CODE:" + ("OK" if response.status_code == 201 else "NOK") + " " + 
            "BODY:NOK1"
            )
    else:
        print(
            "----> " + response.request.method + " " + response.request.path_url + " " + 
            "STATUS_CODE:" + ("OK" if response.status_code == 201 else "NOK") + " " + 
            "BODY:" + ("OK" if response.text == '{"token":"' + f'{token}'+  '"}\n' else "NOK2")
            )

headers = {
    "Content-Type": "application/json",
    "Authorization": f"bearer {ADM_TOKEN}"
}
data = {
    "asd":"asd"
}

token = createToken(requests.post(url + "/adm/token", json=data, headers=headers))

print()

# create context
import random

print("CREATE CONTEXT")
def createContext(response):
    print(
        "----> " + response.request.method + " " + response.request.path_url + " " + 
        "STATUS_CODE:" + ("OK" if response.status_code == 201 else "NOK") + " " + 
        "BODY:" + ("OK" if response.text == "" else "NOK")
        )

headers = {
    "Content-Type": "application/json",
    "Authorization": f"bearer {ADM_TOKEN}"
}
context = f'context{random.randint(1,100)}'
data = {
    "context":f"{context}"
}

createContext(requests.post(url + "/adm/context", json=data, headers=headers))

print()

# check if the context was created
print("GET CONTEXT")
def createContext(response):
    print(
        "----> " + response.request.method + " " + response.request.path_url + " " + 
        "STATUS_CODE:" + ("OK" if response.status_code == 200 else "NOK") + " " + 
        "BODY:" + ("OK" if response.text == '{"context":"' + context + '"}' + "\n" else "NOK")
        )

headers = {
    "Content-Type": "application/json",
    "Authorization": f"bearer {ADM_TOKEN}"
}
data = {
    "context":f"{context}"
}

createContext(requests.get(url + "/adm/context", json=data, headers=headers))

print()
# grant POST on token on context
# POST an entry with the token on the context
# grant GET on token on context and check the entry if it is correct

# revoke GET and POST and grant UPDATE on token on context
# UPDATE an entry with the token on the context
# grant GET on token on context and check the entry if it is correct

# revoke GET, UPDATE and grant DELETE on token on context
# DELETE an entry with the token on the context
# grant GET on token on context and check the entry if it was deleted

# revoke GET of the token on the context

# delete context
# check if the context was deleted

# delete token
# check if the token was deleted