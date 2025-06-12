# Balances Backend

[![codecov](https://codecov.io/gh/kerti/balances/graph/badge.svg?token=FTH30BY4MN)](https://codecov.io/gh/kerti/balances)

## How To Run

1. Clone the repository.
2. Prepare the database.
   - Create a new database
   - Run the SQL scripts located in `/migrations`
   - If you would like to reset the data, run the SQL scripts in `/migrations/demo`
3. Prepare the server
   - Make a copy of `.env.example` and name it `.env`
   - Setup your environment variables to point to the database you've setup in step 2 above
4. Run the service.
   ```bash
   make local
   ```