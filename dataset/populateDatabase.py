import pandas as pd
from pymongo import MongoClient
import random
from datetime import datetime

'''
https://www.kaggle.com/datasets/kartik2112/fraud-detection

   Unnamed: 0 trans_date_trans_time            cc_num                              merchant        category    amt  ...         dob                         trans_num   unix_time  merch_lat  merch_long is_fraud
0           0   2020-06-21 12:14:25  2291163933867244                 fraud_Kirlin and Sons   personal_care   2.86  ...  1968-03-19  2da90c7d74bd46a0caf3777415b3ebd3  1371816865  33.986391  -81.200714        0
1           1   2020-06-21 12:14:33  3573030041201292                  fraud_Sporer-Keebler   personal_care  29.84  ...  1990-01-17  324cc204407e99f51b0d6ca0055005e7  1371816873  39.450498 -109.960431        0
2           2   2020-06-21 12:14:53  3598215285024754  fraud_Swaniawski, Nitzsche and Welch  health_fitness  41.28  ...  1970-10-21  c81755dbbbea9d5c77f094348a7579be  1371816893  40.495810  -74.196111        0
3           3   2020-06-21 12:15:15  3591919803438423                     fraud_Haley Group        misc_pos  60.05  ...  1987-07-25  2159175b9efe66dc301f149d3d5abf8c  1371816915  28.812398  -80.883061        0
4           4   2020-06-21 12:15:17  3526826139003047                 fraud_Johnston-Casper          travel   3.19  ...  1955-07-06  57ff021bd3f328f8738bb535c302a31b  1371816917  44.959148  -85.884734        0

Users with two accounts
['jeffrey.smith', 'jennifer.scott', 'john.nichols', 'justin.bell', 'linda.davis', 'robert.james', 'scott.martin']
'''

# Read csv
df = pd.read_csv('./dataset.csv')

# Init connection to mongodb
client = MongoClient('mongodb://root:root@localhost:27017/?authMechanism=DEFAULT')
db = client.tfg
users = db.users
categories = db.categories
accounts = db.accounts
transactions = db.transactions

# Set variables
bodyUser = {
    'username': '',
    'password': '$2a$14$Wqs/DJ7HaFCSKzluQc4t..sovbPxpPjpA3h3QfztMmX7CT2iOSOKm',
    'role': 'USER'
}

bodyAccount = {
    'iban': '',
    'name': 'DBIT',
    'currency': 'EUR',
    'amount': 0,
    'bank': 'Prueba',
    'updateDate': datetime.now(),
}

bodyCategory = {
    'name': '',
}

bodyTransaction = {}

# Create column name to combine first and last
df['name'] = df['first'].str.lower() + '.' + df['last'].str.lower()
names = len(df['name'].unique())

# Iterate through all the users
for name, index in zip(df['name'].unique(), range(1, names + 1)):
    print('User ', index, '/', names, ': ', name)
    # Create user
    bodyUser['username'] = name
    responseUser = users.insert_one(bodyUser.copy())

    # Get all categories of user
    categoriesData = df.loc[df['name'] == name, 'category'].unique()
    dbCategories = {}
    # Iterate through all the categories
    for category in categoriesData:
        # Parse category
        words = category.split('_')
        capitalizedWords = [word.capitalize() for word in words]
        result = ' '.join(capitalizedWords)
        # Create category and save the id to use it later
        bodyCategory['name'] = result
        bodyCategory['userId'] = responseUser.inserted_id
        responseCategory = categories.insert_one(bodyCategory.copy())
        dbCategories[result] = responseCategory.inserted_id

    # Create category salary
    bodyCategory['name'] = 'Salary'
    bodyCategory['userId'] = responseUser.inserted_id
    responseCategory = categories.insert_one(bodyCategory.copy())
    dbCategories['Salary'] = responseCategory.inserted_id

    # Get all accounts of user
    accountsData = df.loc[df['name'] == name, 'cc_num'].unique()
    # Iterate through all the accounts
    for account in accountsData:
        print('--Account: ', account)
        # Create account
        bodyAccount['iban'] = str(account)
        amount = random.uniform(-2000, 2000)
        bodyAccount['amount'] = amount
        bodyAccount['userId'] = responseUser.inserted_id
        responseAccount = accounts.insert_one(bodyAccount.copy())

        # Get all transactions of account
        transactionsData = df[df['cc_num'] == account]
        # Iterate through all the transactions
        for _, transaction in transactionsData.iterrows():
            # Update amount
            amount = amount + transaction['amt'] * -1
            # Parse category
            words = transaction['category'].split('_')
            capitalizedWords = [word.capitalize() for word in words]
            result = ' '.join(capitalizedWords)
            # Create transaction
            bodyTransaction['description'] = transaction['merchant'].replace('fraud_', '')
            bodyTransaction['date'] = datetime.fromisoformat(transaction['trans_date_trans_time'])
            bodyTransaction['amount'] = transaction['amt'] * -1
            bodyTransaction['category'] = dbCategories[result]
            bodyTransaction['accountId'] = responseAccount.inserted_id
            transactions.insert_one(bodyTransaction.copy())

        # Add salary transactions
        for month in range(6, 13):
            # Update amount
            salary = random.uniform(5000, 10000)
            amount = amount + salary
            bodyTransaction['description'] = 'Salary ' + str(month)
            bodyTransaction['date'] = datetime(2020, month, 1)
            bodyTransaction['amount'] = salary
            bodyTransaction['category'] = dbCategories['Salary']
            bodyTransaction['accountId'] = responseAccount.inserted_id
            transactions.insert_one(bodyTransaction.copy())

        # Update amount after all the transactions
        accounts.update_one({'_id': responseAccount.inserted_id}, {'$set': {'amount': amount}})
