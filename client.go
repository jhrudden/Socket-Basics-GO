package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func main() {
	// grab command line arguments entered on package intialization (or throw error if invalid format is detected).
	port, isSecure, hostname, nuid := getInputs();
	// create either a TLS or non-TLS connection to hostname:port based on whether -s flag is given.
	conn, err := createConnection(isSecure, hostname, port)
	if err != nil {
		panic("failed to connect: " + err.Error())
	}
	defer conn.Close()
	// Once connection has been establish, begin by writing initial HELLO request to server.
	_, err = Write(conn, fmt.Sprintf("ex_string HELLO %s\n", nuid))
	if (err != nil) {
		panic("Error Writing :" + err.Error())
	}
	// Process FIND responses to COUNT requests until a BYE response is recieved or an error occurs.
	countAndTerminate(conn, hostname, port)
	return
}

// Finds the command line arguments input by user.
// Returns : 
// pFlag - port (number following -p), which can be optionally input (defaults to 27993, unless isSecure is given).
// sFlag - whether to use tls connection, boolean given via the -s command line argument (changes default port value to 27994, unless -p is given).
// hostname - name of server we are connecting to.
// nuid - nuid of user we are using for credentials on server.
func getInputs() (port int, isSecure bool, hostname string, nuid string) {
	pFlag := flag.Int("p", 0, "Server port");
	sFlag := flag.Bool("s", false, "If give socket should be TLS")
	flag.Parse()
	args := flag.Args();

	if *pFlag == 0 {
		if *sFlag {
			*pFlag = 27994
		} else {
			*pFlag = 27993
		}
	}

	if (len(args) != 2) {
		panic("You must input a hostname and nuid")
	}

	hostname = args[0];
	nuid = args[1];
	return *pFlag, *sFlag, hostname, nuid;
}


// creates either a tls or non-tls connection to a given hostname with a given port dependent on bool value of isSecure
// Returns: connection interface if connection was possible and an error otherwise.
func createConnection(isSecure bool, hostname string, port int) (net.Conn, error) {
	if (isSecure) {
	tlsCon, err := tls.Dial("tcp", hostname + ":" + strconv.Itoa(port), &tls.Config{})
	return net.Conn(tlsCon), err;

	} 
	return net.Dial("tcp", hostname + ":" + strconv.Itoa(port))
}

// Writes the given message to given connection.
// Returns: byte array of response or an error (if one occured).
func Write(conn net.Conn, message string) (int, error) {
	num, err := conn.Write([]byte(message))
	return num, err;
}

// Reads incoming messages from a either tls or non-tls connection.
// Returns: A completed string message from currently connected server or an error if there was a connection issue.
func Read(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	lines, err := reader.ReadString('\n')
	return lines, err;
}


// Count the occurence of a given focusString in another, then write the count to through a given connection as COUNT request.
// Returns: void
func handleCount(conn net.Conn, focusString string, areaOfSearch string) {
	focusCount := strings.Count(areaOfSearch, focusString);
	_, err := Write(conn, fmt.Sprintf("ex_string COUNT %d\n", focusCount))
	if (err != nil) {
		panic("Error Writing :" + err.Error())
	}
}

// Read FIND responses from connection, then write COUNT requests to connection until either an error occurs or a BYE responses is 
// recieved, then terminate.
// Returns: void
func countAndTerminate(conn net.Conn, hostname string, port int) {
	res, err := Read(conn)

	for {
		if err != nil {
			panic("Error reading from server: " + err.Error())
		}
			splitRes := strings.Split(res, " ")
		// if response from connection doesn't begin with ex_string, then it was in properly formatted. Throw an error.
		if (splitRes[0] != "ex_string") {
			panic("Invalid Response from server");
		}
		// if we have recieved a BYE response beginning with ex_string and with length of three, it was formatted correctly. So,
		// terminate the program.
		if (splitRes[1] == "BYE" && len(splitRes) == 3) {
			fmt.Println("secret key:", splitRes[2])
			break
		} else {
			// If we are given a FIND response with invalid formatting, throw an error.
			if ( splitRes[1] != "FIND" && len(splitRes) != 4)   {
				panic("Invalid Response from server");
			}
			// Otherwise, handle the valid FIND response and SEND a Count request.
			handleCount(conn, splitRes[2], splitRes[3]);
			// Read the next pending request from connection.
			res, err = Read(conn)
		}
	}
	return
}


