package main

import (
    "fmt"
    "os"
    "golang.org/x/sys/unix"
)



func main() {
    // To communicate with netlink, a netlink socket must be opened. This is done using the socket() system call:
    fmt.Println("Hello, this is netlink socket")
    fd, socket_err := unix.Socket(
        // Always used when opening netlink sockets.
        unix.AF_NETLINK,
        // Seemingly used interchangeably with SOCK_DGRAM,
        // but it appears not to matter which is used.
        unix.SOCK_RAW,
        // The netlink family that the socket will communicate
        // with, such as NETLINK_ROUTE or NETLINK_GENERIC.
        unix.NETLINK_KOBJECT_UEVENT,
    )

    fmt.Println("unix.NETLINK_KOBJECT_UEVENT:", unix.NETLINK_KOBJECT_UEVENT)
    fmt.Println("fd is:", fd)

    if socket_err != nil {
        fmt.Println("Socket_err is:" + socket_err.Error())
    }

    // Once the socket is created, bind() must be called to prepare it to send and receive messages.
    bind_err := unix.Bind(fd, &unix.SockaddrNetlink{
        // Always used when binding netlink sockets.
        Family: unix.AF_NETLINK,
        // A bitmask of multicast groups to join on bind.
        // Typically set to zero.
        Groups: 0xffffffff,
        // If you'd like, you can assign a PID for this socket
        // here, but in my experience, it's easier to leave
        // this set to zero and let netlink assign and manage
        // PIDs on its own.
        Pid: uint32(os.Getpid()),
    })

    if bind_err != nil {
        fmt.Println("Bind_err is:" + bind_err.Error())
    }

    bstore := make([]byte, os.Getpagesize())
    fmt.Println("length b is:", len(bstore))
    for {
        fmt.Println("Start reading2 ...")
        for {
            // Peek at the buffer to see how many bytes are available.
            n, _, _ := unix.Recvfrom(fd, bstore, 4096)
            fmt.Println("Length of data is:", n)
            // http://stackoverflow.com/questions/14230145/what-is-the-best-way-to-convert-byte-array-to-string
            fmt.Println(string(bstore[:n]))
            // Break when we can read all messages.
            if n < len(bstore) {
                break
            }
            // Double in size if not enough bytes.
            bstore = make([]byte, len(bstore)*2)
        }
        fmt.Println("Start reading3 ...")
    }

}
