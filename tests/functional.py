import requests

URL = "http://localhost:5258"

# register client
resp = requests.post(URL + "/register-client", json={
    "permissions": 'all'
}).json()

client_key = resp['key']

assert(len(client_key) > 0)

TEST_ACC = "test@test.com"
TEST_PWD = "password123"

# register resource
resp = requests.post(URL + "/register-resource", json={
    "client_key": client_key,
    "resource_id": TEST_ACC,
    "resource_key": TEST_PWD
}).json()

assert(resp['valid'] == True)

# dispatch token
resp = requests.post(URL + "/dispatch", json={
    "resource_id": TEST_ACC,
    "resource_key": TEST_PWD,
}).json()

identity_token = resp['identity_token']

assert(len(identity_token) > 0)

# validate
resp = requests.post(URL + "/validate", json={
    "client_key": client_key,
    "identity_token": identity_token,
}).json()

assert(resp['valid'] == True)

print('Functional Test Success')
