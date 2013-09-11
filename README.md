appleapns
=========

go implement apple apns



$ openssl pkcs12 -clcerts -nokeys -out cert.pem -in apn-teemuikonen-cert.p12 
Enter Import Password:
MAC verified OK

$ openssl pkcs12 -nocerts -out key.pem -in apn-teemuikonen-key.p12 
Enter Import Password:
MAC verified OK
Enter PEM pass phrase:
Verifying - Enter PEM pass phrase:

$ openssl rsa -in key.pem -out key-noenc.pem
Enter pass phrase for key.pem:
writing RSA key
