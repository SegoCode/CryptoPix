# CryptoPix

CryptoPix is concept for image sharing service fully anonymous, database-less and client side encryption these data is encrypted and decrypted in the browser using 256bit AES using [crypto-js](https://github.com/brix/crypto-js) lib.

The images are converted to b64, at this point they are treated as text (Pastebin-like system), encrypted and sent to server.

Easy and minimalistic code easy to setup and modified with zero configuration


## Functionality 
<del>*in very questionable diagram*</del>
## Database-less

All sessions in the web generate a UID these uses to name the local file and uses to referer the file when it is shared, there is no any database configuration but can be easily modified to connect to a database such as Firebase. The server never knows what files are uploaded and it does not retrieve any user-related in UID data.
