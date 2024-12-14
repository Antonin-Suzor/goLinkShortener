# Link shortener in go

Small personal project written during summer 2024.
It served as an introduction to go and the gin framework.

### Functionality

The project simply allows:
- creation of a user
- login with JWT cookie
- creation of redirect links related to the user
- usage of user-defined redirect links

### Code structure

No database was set up / used, all the data is stored in the files in `data/accounts`.
All back-end code is located in the `innards` folder.
The `exposableFiles` folder contains the files served by the server.
