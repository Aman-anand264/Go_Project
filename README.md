A RESTful API for managing books in a bookstore, built with GoFr and MongoDB.

Add a Book
curl -X POST -H "Content-Type: application/json" -d '{"isbn": "123456", "title": "The Book Title", "author": "Author Name"}' http://localhost:your_port/book/add

List Books
curl http://localhost:your_port/books/list

Update a Book
curl -X PUT -H "Content-Type: application/json" -d '{"title": "New Title", "author": "New Author"}' http://localhost:your_port/books/list/123456

Remove a Book
curl -X DELETE http://localhost:your_port/book/remove/123456

Running the Project
git clone https://github.com/your-username/bookstore-api.git
cd bookstore-api
go mod tidy
go run main.go

Testing the API
Postman to send HTTP requests to the endpoints.

License
This project is licensed under the MIT License.




