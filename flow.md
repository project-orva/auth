
### Identity Token Dispatch
#### input:
```json
{
    "resource_id": "",
    "resource_key": "",
    "resource_type": "",
}
```
Everything is a resource, from people to devices.. e.g for auth a user would be to pass the username as the `resource_id`, the password as the `resource_key` & `resource_type` of user. 

#### (valid resource details) output:
```json
{
    "identity_token": "",
    "iat": "<issued at>",
    "eat": "<expires at>",
}
```

### Validation
#### input:
```json
{
    "client_key": "<key of the client making the request>",
    "identity_token": "<token output from dispatch step>",
}
```

#### output:
```json
{
    "valid": "<boolean>",
}
```

### Register Resource
#### input:
```json
{
    "resource_id": "<id of the resource>",
    "resource_key": "<key of the resource>",
    "resource_type": "<type of resource>",
}
```
Note: during this step, all provided keys get bcrypt, if key is not given one will be assigned to the resource automatically.

#### output:
```json
{
    "valid": "<boolean>",
    "resource_key": "<initial key input / generated key (pre-bcrypt)>",
}
```

### Register Resource Client
This endpoint is restricted to whitelisted adresses.
#### input: 
```json
{
    "permissions": "<array of permissions>",
    "created_by": "<id of creator>",
    "expires_on": "<stamp of time expires>"    
}
```

#### output:
```json
{
    "key": "<resource key>",
    "iat": "<issued on>"
}
```