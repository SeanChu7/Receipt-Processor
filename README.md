# Receipt-Processor

Simple web service built in GO that receives receipts in the form of JSON POST requests and calculates the amount of points awarded to the receipt.
Can also return the amount of points awarded to a certain receipt that has been submitted to the program.

To run with Docker, simply build the docker image using the provided Dockerfile, "docker build --tag docker-receipt-processor ." . 
Then simply run the image and expose PORT 8080: "docker run --publish 8080:8080 docker-receipt-processor"

**Functionality**
To send a receipt to the web server, make a POST request to the endpoint "/receipts/process" with the receipt attached as a JSON. On a successful submission, the webserver will return a JSON object containing the ID associated with the receipt. 

To check the number of points a receipt is awarded, make a GET request to the endpoint "/receipts/:id/points" where :id is the ID receieved from the POST request. On successful submission, the webserver will return a JSON object containing the number of points. Otherwise a JSON object containing an error message will be sent. Note that receipts do not persist between application restarts, once the application is closed all stored receipts will be lost. 

