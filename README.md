Para deplegar el proyecto en local se recomienda usar el dataset disponible en https://www.kaggle.com/datasets/kartik2112/fraud-detection junto al script de python disponible en la carpeta dataset.

Como ajustes adicionales habría que configurar las variables de entorno.

El archivo para el docker compose contendría las siguientes variables:
PORT
SECRET_KEY_JWT
MONGO_PATH
MONGO_DATABASE
MONGO_USERNAME
MONGO_PASSWORD
MONGO_INITDB_ROOT_USERNAME
MONGO_INITDB_ROOT_PASSWORD

En el frontend habría que cambiar el .env.development por el .env.production

En el backend las variables son las siguientes:
PORT
SECRET_KEY_JWT
MONGO_PATH
MONGO_DATABASE
MONGO_USERNAME
MONGO_PASSWORD
SANTANDER_ID
SANTANDER_SECRET

Por último habría que ajustar las fechas de las estadísticas, ya que los datos de prueba son del 2020
