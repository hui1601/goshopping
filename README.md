# Simple Shopping Mall
## Quick Start
```shell
$ git clone https://github.com/hui1601/goshopping
```
To use the Toss Payments, you need to sign up for the Toss Payments and get the client key and secret key.(API 개별 연동 키)

Edit the `.env` and `api/.env` files to set the environment variables.
Also, You might edit the `web/static/js/payment.js` file to edit the `clientKey` for Toss Payments.
```shell
$ docker-compose up -d --build
```
Open your browser and navigate to `http://localhost:8000/` to view the app.