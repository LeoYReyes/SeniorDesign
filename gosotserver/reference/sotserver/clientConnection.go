package main

import ("net"; "container/list"; "bytes")

type Client struct {
    // Device Id
    id string
    // Channel for incoming messages
    Incoming chan string
    // Channel for outgoing messages
    Outgoing chan string
    // Connection to client
    Conn net.Conn
    Quit chan bool
    // List of connected clients
    ClientList *list.List
}

// Reads in a message from client
func (c *Client) Read(buffer []byte) bool {
    bytesRead, error := c.Conn.Read(buffer)
    if error != nil {
        c.Close()
        Log(error)
        return false
    }
    Log("Read ", bytesRead, " bytes")
    return true
}
// Closes client connection
func (c *Client) Close() {
    c.Quit <- true
    c.Conn.Close()
    c.RemoveMe()
}

func (c *Client) Equal(other *Client) bool {
    if bytes.Equal([]byte(c.id), []byte(other.id)) {
        if c.Conn == other.Conn {
            return true
        }
    }
    return false
}

func (c *Client) RemoveMe() {
    for entry := c.ClientList.Front(); entry != nil; entry = entry.Next() {
        client := entry.Value.(Client)
        if c.Equal(&client) {
            Log("RemoveMe: ", c.id)
            c.ClientList.Remove(entry)
        }
    }
}