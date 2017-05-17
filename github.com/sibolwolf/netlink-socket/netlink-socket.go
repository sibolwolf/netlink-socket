package main

import (
    "fmt"
    "golang.org/x/sys/unix"
)

func main() {
    // To communicate with netlink, a netlink socket must be opened. This is done using the socket() system call:
    fd, err := unix.Socket(
        // Always used when opening netlink sockets.
        unix.AF_NETLINK,
        // Seemingly used interchangeably with SOCK_DGRAM,
        // but it appears not to matter which is used.
        unix.SOCK_RAW,
        // The netlink family that the socket will communicate
        // with, such as NETLINK_ROUTE or NETLINK_GENERIC.
        unix.NETLINK_KOBJECT_UEVENT,
    )

    // Once the socket is created, bind() must be called to prepare it to send and receive messages.
    err := unix.Bind(fd, &unix.SockaddrNetlink{
        // Always used when binding netlink sockets.
        Family: unix.AF_NETLINK,
        // A bitmask of multicast groups to join on bind.
        // Typically set to zero.
        Groups: 0,
        // If you'd like, you can assign a PID for this socket
        // here, but in my experience, it's easier to leave
        // this set to zero and let netlink assign and manage
        // PIDs on its own.
        Pid: 0,
    })
    for {
        b := make([]byte, os.Getpagesize())
        for {
            // Peek at the buffer to see how many bytes are available.
            n, _, _ := unix.Recvfrom(fd, b, unix.MSG_PEEK)
            fmt.Println(n)
            // Break when we can read all messages.
            if n < len(b) {
                break
            }
            // Double in size if not enough bytes.
            b = make([]byte, len(b)*2)
        }
        // Read out all available messages.
        n, _, _ := unix.Recvfrom(fd, b, 0)
        fmt.Println(n)
    }

}
