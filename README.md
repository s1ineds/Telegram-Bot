The bot may seem not finished. I had written it without using any libraries specifically to learn how to use postgre database, learn how does HTTP works and how to work with API.

The bot have a three tables, Users, Income and Expenses. The Users table stores all of the users that will use the bot. It contains chat_id column and first_name that I take from Telegram API. The table like this I need to make the bot to process a multiple users simultaneously. Thus the tables Income and Expenses store entries about your income and expenses.

In fact, to make this bot to run all you have to do is create a database -> clone repo - and run the bot. This repo contains sql backup that named "gobot_bak.sql" that you should import with "psql" tool.
