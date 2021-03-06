import contracts
from service import Service

if __name__ == "__main__":
    account = {
        'id': 'example_id',
        'key': 'example_key',
        'permissions': []
    }

    auth_serivce = Service('http://localhost:5258')

    contracts.validate_basic(auth_serivce, account)

    account['permissions'] = ['test123']
    contracts.validate_basic(auth_serivce, account, ['test123'])

    print('contracts accepted')