
import requests
import json
import uuid

hostname = input('Type the server hostname > ').strip()
port = input('Type the server port > ').strip()
url = f'http://{hostname}:{port}'
ADM_TOKEN = input("Type ADM Token > ").strip()

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
            response.status_code
        )
        print(
            response.text
        )
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
    return token

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
def getContext(response):
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

getContext(requests.get(url + "/adm/context", json=data, headers=headers))

print()

# grant POST on token on context
print("GRANT POST ON TOKEN")
def grantPost(response):
    print(
        "----> " + response.request.method + " " + response.request.path_url + " " + 
        "STATUS_CODE:" + ("OK" if response.status_code == 201 else "NOK") + " " + 
        "BODY:" + ("OK" if response.text == "" else "NOK")
        )

headers = {
    "Content-Type": "application/json",
    "Authorization": f"bearer {ADM_TOKEN}"
}
data = {
    "token":f"{token}",
    "grant":"POST",
    "context":f"{context}"
}

grantPost(requests.post(url + "/adm/token/grant", json=data, headers=headers))

print()
# POST an entry with the token on the context
print("POST ENTRY WITH TOKEN ON CONTEXT")
def postEntry(response):
    responseJSON= json.loads(response.text)
    print(
        "----> " + response.request.method + " " + response.request.path_url + " " + 
        "STATUS_CODE:" + ("OK" if response.status_code == 201 else "NOK") + " " + 
        "BODY:" + ("OK" if responseJSON['context'] == context 
                        and responseJSON['key'] == entry_key 
                        and responseJSON['value'] == entry_value 
                        and responseJSON['status'] == 'created' 
                        else "NOK")
        )

headers = {
    "Content-Type": "application/json",
    "Authorization": f"bearer {token}"
}
entry_key = str(uuid.uuid4())
entry_value = str(uuid.uuid4())
data = {
    entry_key:entry_value
}

postEntry(requests.post(url + f"/con/{context}", json=data, headers=headers))

print()

# grant GET on token on context and check the entry if it is correct
print("GRANT GET ON TOKEN")
def grantGet(response):
    print(
        "----> " + response.request.method + " " + response.request.path_url + " " + 
        "STATUS_CODE:" + ("OK" if response.status_code == 201 else "NOK") + " " + 
        "BODY:" + ("OK" if response.text == "" else "NOK")
        )

headers = {
    "Content-Type": "application/json",
    "Authorization": f"bearer {ADM_TOKEN}"
}
data = {
    "token":f"{token}",
    "grant":"GET",
    "context":f"{context}"
}

grantGet(requests.post(url + "/adm/token/grant", json=data, headers=headers))

print()

# check if the entry is stored the right way
print("GET THE ENTRY CREATED")
def getEntry(response):
    responseJSON= json.loads(response.text)
    print(
        "----> " + response.request.method + " " + response.request.path_url + " " + 
        "STATUS_CODE:" + ("OK" if response.status_code == 200 else "NOK") + " " + 
        "BODY:" + ("OK" if responseJSON['key'] == entry_key 
                        and responseJSON['value'] == entry_value 
                        and responseJSON['context'] == context 
                        else "NOK")
        )

headers = {
    "Content-Type": "application/json",
    "Authorization": f"bearer {token}"
}
data = {
    "key":entry_key
}

getEntry(requests.get(url + f"/con/{context}", json=data, headers=headers))

print()

# revoke GET and POST and grant UPDATE on token on context
print("REVOKE GET ON TOKEN")
def revokeGet(response):
    print(
        "----> " + response.request.method + " " + response.request.path_url + " " + 
        "STATUS_CODE:" + ("OK" if response.status_code == 200 else "NOK") + " " + 
        "BODY:" + ("OK" if response.text == "" else "NOK")
        )

headers = {
    "Content-Type": "application/json",
    "Authorization": f"bearer {ADM_TOKEN}"
}
data = {
    "token":f"{token}",
    "grant":"GET",
    "context":f"{context}"
}

revokeGet(requests.delete(url + "/adm/token/revoke", json=data, headers=headers))

print()

print("GRANT PUT ON TOKEN")
def grantPut(response):
    print(
        "----> " + response.request.method + " " + response.request.path_url + " " + 
        "STATUS_CODE:" + ("OK" if response.status_code == 201 else "NOK") + " " + 
        "BODY:" + ("OK" if response.text == "" else "NOK")
        )

headers = {
    "Content-Type": "application/json",
    "Authorization": f"bearer {ADM_TOKEN}"
}
data = {
    "token":f"{token}",
    "grant":"PUT",
    "context":f"{context}"
}

grantPut(requests.post(url + "/adm/token/grant", json=data, headers=headers))

print()

# UPDATE an entry with the token on the context
print("UPDATE ENTRY")
def updateEntry(response):
    responseJSON= json.loads(response.text)
    print(
        "----> " + response.request.method + " " + response.request.path_url + " " + 
        "STATUS_CODE:" + ("OK" if response.status_code == 200 else "NOK") + " " + 
        "BODY:" + ("OK" if responseJSON['context'] == context 
                        and responseJSON['key'] == entry_key 
                        and responseJSON['value'] == entry_new_value 
                        and responseJSON['status'] == 'updated' 
                        else "NOK")
        )

headers = {
    "Content-Type": "application/json",
    "Authorization": f"bearer {token}"
}
entry_new_value = str(uuid.uuid4())
data = {
    entry_key:entry_new_value
}

updateEntry(requests.put(url + f"/con/{context}", json=data, headers=headers))

print()
# grant GET on token on context and check the entry if it is correct
headers = {
    "Content-Type": "application/json",
    "Authorization": f"bearer {ADM_TOKEN}"
}
data = {
    "token":f"{token}",
    "grant":"GET",
    "context":f"{context}"
}

requests.post(url + "/adm/token/grant", json=data, headers=headers)

print()

# revoke GET, UPDATE and grant DELETE on token on context
headers = {
    "Content-Type": "application/json",
    "Authorization": f"bearer {ADM_TOKEN}"
}
data = {
    "token":f"{token}",
    "grant":"GET",
    "context":f"{context}"
}

revokeGet(requests.delete(url + "/adm/token/revoke", json=data, headers=headers))

print("GRANT DELETE ON TOKEN")
def grantDelete(response):
    print(
        "----> " + response.request.method + " " + response.request.path_url + " " + 
        "STATUS_CODE:" + ("OK" if response.status_code == 201 else "NOK") + " " + 
        "BODY:" + ("OK" if response.text == "" else "NOK")
        )

headers = {
    "Content-Type": "application/json",
    "Authorization": f"bearer {ADM_TOKEN}"
}
data = {
    "token":f"{token}",
    "grant":"DELETE",
    "context":f"{context}"
}

grantDelete(requests.post(url + "/adm/token/grant", json=data, headers=headers))

print()
# DELETE an entry with the token on the context
print("DELETE ENTRY")
def deleteEntry(response):
    print(
        "----> " + response.request.method + " " + response.request.path_url + " " + 
        "STATUS_CODE:" + ("OK" if response.status_code == 204 else "NOK") + " " + 
        "BODY:" + ("OK" if response.text == "" else "NOK")
        )

headers = {
    "Content-Type": "application/json",
    "Authorization": f"bearer {token}"
}
data = {
    "key":entry_key
}

deleteEntry(requests.delete(url + f"/con/{context}", json=data, headers=headers))

print()

# grant GET on token on context and check the entry if it was deleted
headers = {
    "Content-Type": "application/json",
    "Authorization": f"bearer {ADM_TOKEN}"
}
data = {
    "token":f"{token}",
    "grant":"GET",
    "context":f"{context}"
}

requests.post(url + "/adm/token/grant", json=data, headers=headers)

print("GET THE ENTRY DELETED")
def getEntry(response):
    print(
        "----> " + response.request.method + " " + response.request.path_url + " " + 
        "STATUS_CODE:" + ("OK" if response.status_code == 404 else "NOK") + " " + 
        "BODY:" + ("OK" if response.text == "\n" else "NOK")
        )

headers = {
    "Content-Type": "application/json",
    "Authorization": f"bearer {token}"
}
data = {
    "key":entry_key
}

getEntry(requests.get(url + f"/con/{context}", json=data, headers=headers))

print()

# revoke GET of the token on the context
headers = {
    "Content-Type": "application/json",
    "Authorization": f"bearer {ADM_TOKEN}"
}
data = {
    "token":f"{token}",
    "grant":"GET",
    "context":f"{context}"
}

requests.delete(url + "/adm/token/revoke", json=data, headers=headers)

# post token grant with a non valid token to grant
print("TRY TO GRANT A NON EXISTANT TOKEN")
def grantNonExistantToken(response):
    print(
        "----> " + response.request.method + " " + response.request.path_url + " " + 
        "STATUS_CODE:" + ("OK" if response.status_code == 400 else "NOK") + " " + 
        "BODY:" + ("OK" if response.text == "token does not exist\n" else "NOK")
        )

headers = {
    "Content-Type": "application/json",
    "Authorization": f"bearer {ADM_TOKEN}"
}
data = {
    "token":f"asdasd",
    "grant":"POST",
    "context":f"{context}"
}

grantNonExistantToken(requests.post(url + "/adm/token/grant", json=data, headers=headers))

print()

# delete token grant with a non valid token to grant
print("TRY TO REVOKE A NON EXISTANT TOKEN")
def grantNonExistantToken(response):
    print(
        "----> " + response.request.method + " " + response.request.path_url + " " + 
        "STATUS_CODE:" + ("OK" if response.status_code == 400 else "NOK") + " " + 
        "BODY:" + ("OK" if response.text == "token does not exist\n" else "NOK")
        )

headers = {
    "Content-Type": "application/json",
    "Authorization": f"bearer {ADM_TOKEN}"
}
data = {
    "token":f"asdasd",
    "grant":"POST",
    "context":f"{context}"
}

grantNonExistantToken(requests.delete(url + "/adm/token/revoke", json=data, headers=headers))

print()

# delete context
print("DELETE CONTEXT")
def deleteContext(response):
    print(
        "----> " + response.request.method + " " + response.request.path_url + " " + 
        "STATUS_CODE:" + ("OK" if response.status_code == 200 else "NOK") + " " + 
        "BODY:" + ("OK" if response.text == "" else "NOK")
        )

headers = {
    "Content-Type": "application/json",
    "Authorization": f"bearer {ADM_TOKEN}"
}
data = {
    "context":f"{context}"
}

deleteContext(requests.delete(url + "/adm/context", json=data, headers=headers))

print()

# check if the context was deleted
print("GET CONTEXT")
def getContext(response):
    print(
        "----> " + response.request.method + " " + response.request.path_url + " " + 
        "STATUS_CODE:" + ("OK" if response.status_code == 404 else "NOK") + " " + 
        "BODY:" + ("OK" if response.text == "\n" else "NOK")
        )

headers = {
    "Content-Type": "application/json",
    "Authorization": f"bearer {ADM_TOKEN}"
}
data = {
    "context":f"{context}"
}

getContext(requests.get(url + "/adm/context", json=data, headers=headers))

print()

# delete token
print("DELETE TOKEN")
def deleteContext(response):
    print(
        "----> " + response.request.method + " " + response.request.path_url + " " + 
        "STATUS_CODE:" + ("OK" if response.status_code == 200 else "NOK") + " " + 
        "BODY:" + ("OK" if response.text == "" else "NOK")
        )

headers = {
    "Content-Type": "application/json",
    "Authorization": f"bearer {ADM_TOKEN}"
}
data = {
    "token":f"{token}"
}

deleteContext(requests.delete(url + "/adm/token", json=data, headers=headers))

print()

# check if the token was deleted - functionality does not exists yet
