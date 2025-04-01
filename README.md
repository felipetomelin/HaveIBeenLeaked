# HaveIBeenLeaked
You need to pass a hash with the first 5 characters, and the backend will return a list of all suffix hashes (the remaining characters) that match your prefix. This way, your complete password hash never transits over the network.

When checking if your password has been leaked:

Your browser calculates the complete SHA-1 hash of your password locally

Only the first 5 characters of this hash are sent to our server

The server returns all hash suffixes (the remaining characters) that begin with those 5 characters

Your browser then checks locally if any of the returned suffixes, when combined with your prefix, match your complete password hash

This k-anonymity model ensures your actual password or complete hash is never transmitted, protecting your privacy while still allowing you to check if your password has appeared in known data breaches.

![image](https://github.com/user-attachments/assets/ead7b3ce-191a-490a-8fc1-0d789b4b62fc)
