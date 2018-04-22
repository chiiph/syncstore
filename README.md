# SyncStore

SyncStore is a simple object store sync implementation. It's an abstraction
over simple object stores to add versions to objects and synchronize data
from two object stores. Kind of like rsync for object stores.

The first iteration will be the simplest way to accomplish that definition above.
The end goal would be to have an online service to create accounts and use this
as a backend for an application.

The idea is that stores can be local or remote. Remote stores can go through an
authentication proxy or not. The algorithm remains, what changes is how complex
is it to fetch an object, whether it's read a file, or it's connecting to S3 and
downloading data, doesn't matter.

# Disclaimer

Aren't there systems like this already? I don't know, it doesn't matter, the goal
of this project is to develop something to have some fun when I'm bored.

The goal is to also do it in the best possible way I can, but it'll take time to
get there.

Use of this code in production is discouraged unless you understand perfectly how
it works or you're willing to risk it.
