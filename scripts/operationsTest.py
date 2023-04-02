import requests

url = 'http://localhost:8080/'

# Create user
body = """
    mutation {
        createUser(username: "user1", password: "123", role: ADMIN) {
            id
            username
            password
            role
            operations{
                id
            }
        }
    }
"""

response = requests.post(url=url, json={'query': body})

userId = response.json()['data']['createUser']['id']

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

# Create operation
body = """
    mutation {
        createOperation(
            input: {
                description: "Descripción de prueba 3"
                date: 1609542000000
                amount: 10
                category: "Sports"
            }
        ) {
            id
            description
            date
            amount
            category
            userId
        }
    }
"""

response = requests.post(url=url, json={'query': body}, headers={'Authorization': jwt})

id = response.json()['data']['createOperation']['id']

if response.status_code == 200:
    print('Success create operation')

# Get all operations
body = """
    query {
        operations {
            id
            description
            date
            amount
            category
            userId
        }
    }
"""

response = requests.post(url=url, json={'query': body}, headers={'Authorization': jwt})

if response.status_code == 200:
    print('Success get all operations')

# Get all operations by date
body = """
    query {
        operationsByDate(initDate: 1640991600000 endDate: 1641078000000) {
            id
            description
            date
            amount
            category
            userId
        }
    }
"""

response = requests.post(url=url, json={'query': body}, headers={'Authorization': jwt})

if response.status_code == 200:
    print('Success get all operations by date')

# Get all operations by category
body = """
    query {
        operationsByCategory(category: "Sports") {
            id
            description
            date
            amount
            category
            userId
        }
    }
"""

response = requests.post(url=url, json={'query': body}, headers={'Authorization': jwt})

if response.status_code == 200:
    print('Success get all operations by category')

# Get operation by id
body = """
    query {
        operationById(id: "%s") {
            id
            description
            date
            amount
            category
            userId
        }
    }
""" % id

response = requests.post(url=url, json={'query': body}, headers={'Authorization': jwt})

if response.status_code == 200:
    print('Success get operation by id')

# Update operation
body = """
    mutation {
        updateOperation(input: {
            id: "%s"
            description: "Descripción de prueba 2"
        }) {
            id
            description
            date
            amount
            category
            userId
        }
    }
""" % id

response = requests.post(url=url, json={'query': body}, headers={'Authorization': jwt})

if response.status_code == 200:
    print('Success update operation')

# Delete operation
body = """
    mutation {
        deleteOperation(id: "%s")
    }
""" % id

response = requests.post(url=url, json={'query': body}, headers={'Authorization': jwt})

if response.status_code == 200:
    print('Success delete operation')

# Delete user
body = """
    mutation {
        deleteUser(id: "%s")
    }
""" % userId

response = requests.post(url=url, json={'query': body}, headers={'Authorization': jwt})

if response.status_code == 200:
    print('Success delete user')