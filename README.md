# CryptoPix

CryptoPix is concept for image sharing service fully anonymous, database-less and client side encryption these data is encrypted and decrypted in the browser using 256bit AES using [crypto-js](https://github.com/brix/crypto-js) lib.

The images are converted to b64, at this point they are treated as text (Pastebin-like system), encrypted and sent to server.

Easy and minimalistic code, easy to setup and modified with zero configuration


## Functionality 
<del>*In very questionable diagram*</del>

<img src="https://github.com/SegoCode/CryptoPix/blob/main/documents/CryptoPix%20Diagram.png">

## Database-less

All sessions in the web generate a UID these uses to name the local file and uses to referer the file when it is shared, there is no any database configuration but can be easily modified to connect to a database such as Firebase. The server never knows what files are uploaded and it does not retrieve any user-related in UID data.

## Share link 

The sahre link contains id for encrypted file and contains a # character is a fragment URL. The portion of the URL to the right of the # is the key for decrypt, If you try using fragment URLs in an HTTP sniffer like HttpWatch, you’ll never see the fragment IDs in the requested URL or Referer header. The reason is that the fragment identifier is only used by the browser, fragments Are not Sent in HTTP Request Messages!
