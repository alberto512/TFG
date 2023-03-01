import requests

url = 'http://localhost:8080/'

# Create user
body = """
    mutation {
        createUser(username: "user1", password: "123", rol: ADMIN) {
            id
            username
            password
            rol
            operations{
                id
            }
        }
    }
"""

response = requests.post(url=url, json={'query': body})

id = response.json()['data']['createUser']['id']

if response.status_code == 200:
    print('Success creating user')

# Login
body = """
    mutation {
        login(username: "user1", password: "123")
    }
"""

response = requests.post(url=url, json={'query': body})

if response.status_code == 200:
    print('Success login')

jwt = response.json()['data']['login']

# Get all users
body = """
    query {
        users {
            id
            username
            password
            rol
            operations {
                id
            }
        }
    }
"""

response = requests.post(url=url, json={'query': body}, headers={'Authorization': jwt})

if response.status_code == 200:
    print('Success get all users')

# Get user by id
body = """
    query {
        userById(id: "%s") {
            id
            username
            password
            rol
            operations {
                id
            }
        }
    }
""" % id

response = requests.post(url=url, json={'query': body}, headers={'Authorization': jwt})

if response.status_code == 200:
    print('Success get user by id')

# Get user by token
body = """
    query {
        userByToken {
            id
            username
            password
            rol
            operations {
                id
                description
                date
                amount
                category
                userId
            }
        }
    }
"""

response = requests.post(url=url, json={'query': body}, headers={'Authorization': jwt})

if response.status_code == 200:
    print('Success get user by token')

# Update user
body = """
    mutation {
        updateUser(id: "%s", password: "1234") {
            id
            username
            password
            rol
            operations {
                id
                description
                date
                amount
                category
            }
        }
    }
""" % id

response = requests.post(url=url, json={'query': body}, headers={'Authorization': jwt})

if response.status_code == 200:
    print('Success update user')

# Delete user
body = """
    mutation {
        deleteUser(id: "%s")
    }
""" % id

response = requests.post(url=url, json={'query': body}, headers={'Authorization': jwt})

if response.status_code == 200:
    print('Success delete user')