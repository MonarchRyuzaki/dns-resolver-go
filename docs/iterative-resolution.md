# Iterative DNS Resolution

## Recursive vs. Iterative
When your computer needs to resolve a domain (like `dns.google.com`), it typically sends a **Recursive** query (`RD=1`) to a public DNS resolver like `8.8.8.8` or your ISP's router. 
A recursive query means: *"Go find the IP for me, do all the hard work of traversing the internet, and just hand me the final result."*

If you are building a true DNS Resolver, you don't want to rely on `8.8.8.8`. You want to do the traversal yourself! This is called **Iterative Resolution**.

To do this, you set the Recursion Desired bit to 0 (`RD=0`) in your DNS Header. This tells the server you query: *"Just give me a hint on who to ask next."*

## The Root Name Servers
An iterative resolution always begins at the very top of the DNS hierarchy: **The Root Servers**. 
These servers don't know the IP of `dns.google.com`, but they *do* know the IP of the servers that manage the `.com` Top-Level Domain (TLD). 

There are 13 Root Servers (A-M). For example, `198.41.0.4` is the A Root Server.

## The Iterative Algorithm

To build a full resolver, you use a `for` loop that implements the following logic:

### 1. Set the Initial Target
You always start by aiming your query at a Root Server:
`server = "198.41.0.4"`

### 2. Check for an Answer (Type 1)
When you receive the `DNSMessage` back, check the `Answers` array.
If you find a record where `Type == 1` (A Record), you have successfully found the IP address! You can print it and exit the loop.

### 3. Check for Glue Records (Additionals)
If `Answers` is empty, it means the server is referring you down the chain. Look at the `Additionals` array.
DNS servers will often attach "Glue Records" here. They are essentially saying: *"I don't know the IP you want, but you should ask this NS server, and by the way, here is that NS server's IP."*
If you find a `Type == 1` record here, extract the IP, set `server = newIP`, and let the loop run again!

### 4. The "Missing IP" Edge Case (Authorities)
Sometimes, a server will refer you to the next nameserver, but will be "lazy" and only give you the *name* of the nameserver in the `Authorities` section, without providing the IP in the `Additionals` section.
If `Additionals` contains no IPs, you must:
1. Find the `NS` record (`Type 2`) in the `Authorities` array.
2. Decode the domain name of that nameserver (e.g., `ns1.google.com`).
3. **Recursively** invoke your entire Iterative Algorithm to resolve the IP of `ns1.google.com`.
4. Once you get that IP, set `server = resolvedIP` and continue your original loop.

### Example Flow
1. Query `198.41.0.4` for `dns.google.com` 
   - *Result*: Refers you to `.com` TLD servers. (Found IP `192.12.94.30` in Additionals).
2. Query `192.12.94.30` for `dns.google.com`
   - *Result*: Refers you to Google's authoritative servers. (Found IP `216.239.34.10` in Additionals).
3. Query `216.239.34.10` for `dns.google.com`
   - *Result*: Returns `8.8.8.8` in Answers!
