import requests

class Service:
    def __init__(self, url):
        self.url = url
    
    def register_client(self, permissions):
        return requests.post(f'{self.url}/register-client', json={
            "permissions": permissions
        }).json()

    def register_resource(self, client_key, account):
        return requests.post(f'{self.url}/register-resource', json={
            "client_key": client_key,
            "resource_id": account['id'],
            "resource_key": account['key'],
            "permissions": account['permissions']
        }).json()

    def dispatch_token(self, resource):
        return requests.post(f'{self.url}/dispatch', json={
            "resource_id": resource['id'],
            "resource_key": resource['key'],
        }).json()

    def validate(self, client_key, identity_token):
        return requests.post(f'{self.url}/validate', json={
            "client_key": client_key,
            "identity_token": identity_token,
        }).json()