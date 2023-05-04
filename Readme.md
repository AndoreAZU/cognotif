
# Sample Project Ecommerce

This project contains 3 main file. The first one is core api which have api for ecommerce purpose. The second one is scheduler logic which have responsible to sending email for pending order for each user. This scheduler is customizable based on our requirement. The last main is for generate csv order.

For sending email, since i dont have any smtp server, so in the logic only printing the data that should be send through email.

## Tech Stack

**Backend:** Golang, Echo Framework, Postgress

## Deployment

We need initiate database on first run
```bash
  make initiate-database
```

Then we can running core api
```bash
  make run-core
```

If we want to generate report order
```bash
  make generate-report
```

If we want to trigger scheduler
```bash
  make run-scheduler
```

After that, if we want to destroy database
```bash
  make destroy-database
```


## Test api
I also provide json for postman collection. You can use that one if you want to testing the api.

## Badges

[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://choosealicense.com/licenses/mit/)
[![GPLv3 License](https://img.shields.io/badge/License-GPL%20v3-yellow.svg)](https://opensource.org/licenses/)
[![AGPL License](https://img.shields.io/badge/license-AGPL-blue.svg)](http://www.gnu.org/licenses/agpl-3.0)
