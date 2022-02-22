# Intro

Room Bookings web app powered by Go

- [chi router](https://github.com/go-chi/chi/v5)
- Alex Edwards's sessions manager [SCS](https://github.com/alexedwards/scs/v2)
- [nosurf](https://github.com/justinas/nosurf)

# Soda usage

From the project root directory,
```bash
    soda generate fizz <name-of-table>

    # eg
    soda generate fizz CreateUserTable
```

this creates a directory called `migrations` in the project root dir. We can then populate it using the `fizz` syntax.

Once the `.fizz` files have been populated, run:
`soda migrate` to migrate the schema to the database.

To run all the down migrations and then the up migrations,

```bash
    soda reset
```

Note for this to work, the database must not be currently accessed by any other user.

