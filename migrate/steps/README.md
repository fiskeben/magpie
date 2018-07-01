# Database migrations

## Naming convention

For each migration create two files, 
one for creating/altering stuff (up),
one for reverting stuff (down).

Each file should be named as follows:

`{timestamp}_hyphen-separated-descriptive-name.{up|down}.sql`
