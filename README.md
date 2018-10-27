# gofish

A simple file sharing server with basich authentication and read-only/admin users written in go (golang)

If you ever (like me) searched for a simple tool to install on your own server a secure file sharing service, this little tool might be what you need. gofish establishes a server that provides acces to files in a configurable directory tree on the server for authorised users using basic authentication.

Two different user roles are supported:

* a standard user with read-only access to the files

* an admin user that has also the right to add/remove files and/or folders

Please be aware that basic authentication makes only sense in combination with SSL/TLS and certified keys, e.g as provided by [Let's encrypt](https://letsencrypt.org/) not directly supported by gofish.

gofish may be configured using command line arguments or by a configuration file. The latter is highly recommend.
