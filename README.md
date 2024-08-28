# Phone Book API

A simple phone book Web Server API.

## Features
- Add a contact
- Get contact 
- Search for contacts
- Update a contact
- Delete a contact

## Design Decisions

### File Structure
The project follows a structured file organization to maintain readability and manageability. The `cmd` folder contains the `main.go` file, which serves as the entry point for the application. This structure is familiar to me and helps in keeping the code organized.

### Database Choice
SQLite is used for its simplicity and ease of installation, making it suitable for this project’s development phase. For a production environment, MySQL would be preferred due to its robustness. The choice of a structured database (SQL) over a NoSQL database is based on the clearly structured nature of the data. Structured data benefits from the relational model of SQL databases, which provides clear schemas and relationships between data entities. In contrast, NoSQL databases like MongoDB are more suited for unstructured or semi-structured data and scenarios requiring flexible schema designs. For this project, where the data structure is well-defined and consistent, a structured database aligns better with the project’s needs.


### Pagination limit 10 contacts per page.
Prev and next are links to previous and next pages.
Example response:
{
   "contacts": [
    {
           "id": "545ba702-f79b-49c7-900a-2752d0b1fe6d",
           "firstName": "shayna",
           "lastName": "segal",
           "phone": "0580000",
           "address": "my address"
       }...
   ],
   "pagination": {
       "count": 23,
       "next": "/contacts?limit=10&offset=10&count=23",
       "prev": ""
   }
}


### Scaling
Horizontal Pod Autoscaling (HPA) is used to dynamically adjust the number of application instances based on CPU usage. This approach ensures that the application can scale efficiently under varying load conditions. The configuration sets a minimum of 2 replicas to maintain robustness, with additional replicas automatically added as CPU utilization increases.


### API Documentation
Swagger is utilized to represent the API's endpoints, including all optional requests and responses. This documentation aids in understanding and integrating with the API by providing a clear overview of its functionalities.

### Caching
Added a redis layer last minute bonus.
Implemented the get contact by id should retrieve from redis if it exists.

### Future Improvements
Improve caching and Implement saving to cache for search method where the entire response would be saved in redis to reduce all the calls for next pages.
User management 
Security





