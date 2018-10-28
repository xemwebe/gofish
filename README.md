# gofish

A simple file sharing server with basic authentication and read-only/admin users written in go (golang)

If you ever (like me) searched for a simple tool to install on your own server a secure file sharing service, this little tool might be what you need. gofish establishes a server that provides acces to files in a configurable directory tree on the server for authorised users using basic authentication.

Two different user roles are supported:

* a standard user with read-only access to the files

* an admin user that has also the right to add/remove files and/or folders

Please be aware that basic authentication makes only sense in combination with SSL/TLS and certified keys, e.g as provided by [Let's encrypt](https://letsencrypt.org/) not directly supported by gofish.

gofish may be configured using command line arguments or by a configuration file. The latter is highly recommend.

`gofish -help` gives a list of all comman line arguments (note: if a configuration file is used, the command line arguments will be overwritten).

## Setup a gofish server

1. If not already done, install [go](https://golang.org/) on your server or some other binary compatible computer and compile gofish like any other go executable

2. Please make sure that the folder in which you execute the gofish application must contain the views folder (the gofish binary may be in any othe folder). The images folder should be in the same folder as the views folder, so this is not a requirement.

3. Define a folder for where you wish to put the shared files.

4. Generate password hashes for the read-only user (user name is the empty string "") and the admin user (user name is "admin") by the command:

```bash
gofish -gen-pwd <password>
```

5. Make the appropriate changes in the config file:

```json
{
  "allow_admin": false,
	"file_path":   "/var/www/My_File_Storage",
	"ip_address":  "127.0.0.1",
	"port":        "7356",
	"title":       "My Private Web File Sharing Site",
	"favicon":     "/images/sgws.ico",
	"realm": 	     "GoFiSh",
	"author":      "<Page's owner name>",
	"email":       "<Page's owner email address>",
	"subject":     "MyPrivateWebSpace",
	"userpwhash":	 "",
	"adminpwhash": "",
	"colors": {
		"title": "Navy",
		"buttonfg": "OrangeRed",
		"buttonbg": "Gold"
	}
```

The user password hash and the admin password hash must be set to the hashes generated in the previous step. Also, make sure the file_path is set to an absolute path on your server. The user, under which the gofish is running, must have read and write acces (if run in admin mode) for this folder. If `allow_admin` is set, the server runs in admin mode (i.e. the admin may upload/delete files and folders), otherwise the server runs in read-only mode. 

6. After configuration is complete, start the server:

```bash
gofish -admin -config=./config.json
```

In this example, the additional `-admin` argument overrides the flag `allow_admin` in the configuration file

