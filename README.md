# Project 1: Socket Basics

## My Approach

For this class I decided to learn a new language, Golang, because I have heard it is pretty good for projects similar to what we will being doing in this class. With that in mind, I approached this project in two steps, creating a tcp connection via socket, then reading and writing to connections. I solved the first problem using net and crypto packets for go, which I was able to use in order to implement both tls and non-tls connections. Once I had established a connection, I tackled writing to a connection in order to create the HELLO request. This was a little finicky, but with a little elbow grease I got it to work. With that done, I built out the reading functionality, which is where I spent most my time on this project (I will talk about this in challenges section). With that I was able to easily fetch FIND and BYE responses. Since I already had writing done, COUNT was very to implement. I finished by making sure my program wouldn't break when given garbage inputs or responses from the server and now here we are.

## Challenges

As I said before, I had never used Golang before this project, so there were a ton of challenges in implementing this project. One of the largest was having to deal with allocating space for buffers which I would to capture responses from the server. Originally I set the size of the buffer to be a very large amount of bytes, but when I learned that the messages vary in length, I decided that it would be better to use a method which would either extend or squish the byte size of buffer dependent on the size of response message. This issue actually took the most amount of my time, but once I solved it (after a few hours of google), the project went smoothly.

## Testings Process

The majority of my testings was done through print statements and trial and error through using the servers responses. This is not a very good practice as there must be better ways to test my code. I could create a mock server and test giving my client bogus responses etc. While I probably should have done this, I didn't feel I had the time, so with strenous fmt.println testings, we have arrived at my finished piece of code.
