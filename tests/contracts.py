
def validate_basic(service, account, client_permissions = []):
    client_key = service.register_client(client_permissions)['key']
    assert(len(client_key) > 0)

    register_response = service.register_resource(
        client_key,
        account
    )

    assert(register_response['valid'] == True)

    dispatch_response = service.dispatch_token(account)
    identity_token = dispatch_response['identity_token']

    assert(len(identity_token) > 0)
    
    validate_response = service.validate(client_key, identity_token)
    assert(validate_response['valid'] == True)
