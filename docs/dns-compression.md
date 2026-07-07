# DNS Message Compression (RFC 1035)

## The Problem
A standard DNS UDP packet has a strict limit of 512 bytes. If a DNS server responds with multiple records for the same domain (e.g., `www.google.com`, `mail.google.com`), spelling out `google` and `com` repeatedly wastes a massive amount of space. 

To solve this, DNS uses a **compression scheme**: instead of repeating strings, the message says, "I already wrote this string earlier in the packet. Just go look at byte offset X."

## The Label Restriction (Max 63 Characters)
To understand how compression pointers work, you first have to understand normal domain encoding. A domain like `www.google.com` is broken down into "labels".
It is encoded as a sequence of length-prefixed strings: `[3] w w w [6] g o o g l e [3] c o m [0]`.

According to RFC 1035, **a single label can never exceed 63 characters**. 
Because 63 in binary is `00111111`, any valid length byte will **always** start with two zero bits (`00`). 
It is mathematically impossible for a normal length byte to start with a `1`. 

## The "11" Trap Card (The Pointer)
Because normal length bytes always start with `00`, the designers of DNS created a rule:
If the parser reads a byte and the top two bits are `11`, it is a "trap card." It is **not** a length byte; it is the first half of a **2-byte pointer**.

In hexadecimal, `11000000` is `0xC0` (or `192` in decimal). So, if `byte >= 192` (or `byte & 0xC0 == 0xC0`), you know you've hit a pointer.

## Why a 2-Byte Sequence?
Why not just use a 1-byte pointer? 
If a pointer was only 1 byte long, the maximum value it could hold would be 255. That means it would be physically impossible to point to any string located past byte 255 in the packet. 

By making the pointer **2 bytes (16 bits)**:
1. The first 2 bits are consumed by the `11` marker.
2. We are left with exactly **14 bits** for the actual offset number.
3. The maximum number you can store in 14 bits is `16,383`.

Since standard DNS UDP packets max out at 512 bytes, `16,383` is more than enough space to point to any byte index inside any valid DNS packet!

### Anatomy of a Pointer
If the offset we want to point to is byte `20` (binary `00000000010100`), the two pointer bytes look like this:
```text
[11] [000000]   [00010100]
 ▔▔   ▔▔▔▔▔▔     ▔▔▔▔▔▔▔▔
Marker  Offset (High)   Offset (Low)
```
- Byte 1: `11000000` (`0xC0`)
- Byte 2: `00010100` (`0x14`)

## The Core Algorithm (How to Parse)

When looping through bytes to decode a domain name, your algorithm must handle three scenarios for every byte it inspects:

1. **If the byte is `0x00`**: 
   You have reached the end of the domain name. Stop parsing.

2. **If the byte is `< 192` (e.g. `0x03`)**: 
   It is a standard length byte. 
   - Read the number `N`.
   - Read the next `N` bytes and append them as a string.
   - Advance your overall packet offset by `N + 1`.

3. **If the byte is `>= 192` (starts with `11`)**: 
   You have hit a compression pointer!
   - Grab the current byte and the *next* byte to form a 16-bit integer (using `binary.BigEndian.Uint16`).
   - Erase the top two `11` bits by doing a bitwise AND: `pointerBytes & 0x3FFF`.
   - The resulting number is the exact byte offset from the start of the packet where the rest of the string lives.
   - Jump to that offset in the original byte array and recursively parse the rest of the string from there.
   - *Note: A pointer always signifies the absolute end of the current domain representation. You do not keep reading bytes after the 2-byte pointer.*

### Important Concept: Pointers Always Point to a Suffix
A very common misconception is that a pointer can grab the "middle" of a string and then return. **This is false.** 
There is no "return" after jumping to a pointer. A pointer always points to a **suffix** that naturally terminates (by hitting a `0` byte). 

For example, if the original string at byte 20 is `[6]google[3]com[0]`, and you want to encode `api.google.com`:
- You write: `[3] a p i [Pointer to 20]`
- The parser reads `api`, hits the pointer, jumps to byte 20, and reads `google.com.`. It terminates successfully at the `0` at the end of `com`. 

If you wanted to encode `api.google.org` but only `google.com` was in the packet, you **cannot** use a pointer to just grab the `google` part. You would have to spell out `api.google.org` manually. A pointer always represents the **entire rest of the domain name**.
